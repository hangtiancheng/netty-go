package netty

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"testing"

	bufferpool "github.com/libp2p/go-buffer-pool"
)

// Queue actions so tests can change caller buffers before the transport sees them.
type poolTestExecutor struct {
	actions []Action
}

func (e *poolTestExecutor) Exec(action Action) { e.actions = append(e.actions, action) }

func (e *poolTestExecutor) run() {
	for len(e.actions) > 0 {
		action := e.actions[0]
		e.actions = e.actions[1:]
		action()
	}
}

type poolTestTransport struct {
	net.Conn
	data       bytes.Buffer
	writeErr   error
	writePanic bool
}

func (tr *poolTestTransport) Write(p []byte) (int, error) {
	if tr.writeErr != nil {
		return 0, tr.writeErr
	}
	return tr.data.Write(p)
}

func (tr *poolTestTransport) Writev(p net.Buffers) (int64, error) {
	if tr.writeErr != nil {
		// A transport may consume slice headers before returning an error.
		p[0] = nil
		if tr.writePanic {
			panic(tr.writeErr)
		}
		return 0, tr.writeErr
	}
	return p.WriteTo(&tr.data)
}

func (*poolTestTransport) Flush() error         { return nil }
func (*poolTestTransport) Close() error         { return nil }
func (tr *poolTestTransport) RawTransport() any { return tr }

func newPoolTestChannel(ctx context.Context, size int, untilWrite bool) (*channel, *poolTestTransport, *poolTestExecutor) {
	tr := &poolTestTransport{}
	executor := &poolTestExecutor{}
	c := newChannelWith(ctx, NewPipeline(), tr, executor, 1, size, untilWrite).(*channel)
	return c, tr, executor
}

// Churn the same bucket to expose buffers returned before their queued write.
func churnPool(size int) {
	for i := 0; i < 100; i++ {
		buf := bufferpool.Get(size)
		for j := range buf {
			buf[j] = 0xff
		}
		bufferpool.Put(buf)
	}
}

func TestChannelPoolReadFrom(t *testing.T) {
	for _, size := range []int{0, 16} {
		name := "sync"
		if size > 0 {
			name = "async"
		}
		t.Run(name, func(t *testing.T) {
			c, tr, executor := newPoolTestChannel(context.Background(), size, false)
			defer c.cancel()
			want := strings.Repeat("hello", 1000)
			n, err := c.ReadFrom(strings.NewReader(want))
			if n != int64(len(want)) || err != nil {
				t.Fatalf("ReadFrom = %d, %v", n, err)
			}
			churnPool(1024)
			executor.run()
			if tr.data.String() != want {
				t.Fatal("ReadFrom data changed while awaiting the transport")
			}
		})
	}
}

func TestChannelPoolAsyncWriteCopiesCallerData(t *testing.T) {
	c, tr, executor := newPoolTestChannel(context.Background(), 8, false)
	defer c.cancel()
	first := []byte("hello")
	parts := [][]byte{[]byte(" "), []byte("world")}
	if n, err := c.Write1(first); n != len(first) || err != nil {
		t.Fatalf("Write1 = %d, %v", n, err)
	}
	if n, err := c.Writev(parts); n != 6 || err != nil {
		t.Fatalf("Writev = %d, %v", n, err)
	}
	copy(first, "xxxxx")
	copy(parts[1], "xxxxx")
	churnPool(8)
	executor.run()
	if tr.data.String() != "hello world" {
		t.Fatalf("transport received %q", tr.data.String())
	}
}

func TestChannelPoolRejectedWrites(t *testing.T) {
	for _, untilWrite := range []bool{false, true} {
		c, tr, executor := newPoolTestChannel(context.Background(), 1, untilWrite)
		defer c.cancel()
		want := strings.Repeat("a", 1024)
		if _, err := c.Write1([]byte(want)); err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if _, err := c.CtxWrite1(ctx, make([]byte, 1024)); !errors.Is(err, context.Canceled) {
			t.Fatalf("canceled Write1: %v", err)
		}
		if _, err := c.CtxWritev(ctx, [][]byte{make([]byte, 1024)}); !errors.Is(err, context.Canceled) {
			t.Fatalf("canceled Writev: %v", err)
		}
		if !untilWrite {
			if _, err := c.Write1(make([]byte, 1024)); !errors.Is(err, ErrAsyncNoSpace) {
				t.Fatalf("full queue Write1: %v", err)
			}
			if _, err := c.Writev([][]byte{make([]byte, 1024)}); !errors.Is(err, ErrAsyncNoSpace) {
				t.Fatalf("full queue Writev: %v", err)
			}
			if _, err := c.ReadFrom(strings.NewReader(want)); !errors.Is(err, ErrAsyncNoSpace) {
				t.Fatalf("full queue ReadFrom: %v", err)
			}
		}
		churnPool(1024)
		executor.run()
		if tr.data.String() != want {
			t.Fatal("a rejected write changed the pending packet")
		}
	}
}

func TestChannelPoolParentCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c, _, _ := newPoolTestChannel(ctx, 1, true)
	defer c.cancel()
	if _, err := c.Write1([]byte("queued")); err != nil {
		t.Fatal(err)
	}
	cancel()
	if _, err := c.ReadFrom(strings.NewReader("rejected")); !errors.Is(err, context.Canceled) {
		t.Fatalf("ReadFrom on a canceled channel: %v", err)
	}
	bufferpool.Put(<-c.writeQueue)
}

func TestChannelPoolWriteBatchFailure(t *testing.T) {
	for _, panicWrite := range []bool{false, true} {
		c := &channel{transport: &poolTestTransport{writeErr: io.ErrClosedPipe, writePanic: panicWrite}}
		send := net.Buffers{bufferpool.Get(16), bufferpool.Get(32)}
		recycle := append(net.Buffers(nil), send...)
		var recovered any
		func() {
			defer func() { recovered = recover() }()
			c.writeBatch(send, recycle)
		}()
		if err, ok := recovered.(error); !ok || !errors.Is(err, io.ErrClosedPipe) {
			t.Fatalf("unexpected write failure: %v", recovered)
		}
		for i := range send {
			if send[i] != nil || recycle[i] != nil {
				t.Fatal("failed write retained buffer references")
			}
		}
	}
}
