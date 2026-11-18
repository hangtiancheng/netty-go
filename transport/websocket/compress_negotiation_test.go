package websocket

import (
	"bytes"
	"context"
	"io"
	"net"
	"sync"
	"testing"
)

// recordingConn observes the actual client frames after the HTTP handshake.
type recordingConn struct {
	net.Conn
	mu      sync.Mutex
	written bytes.Buffer
}

func (c *recordingConn) Write(p []byte) (int, error) {
	c.mu.Lock()
	_, _ = c.written.Write(p)
	c.mu.Unlock()
	return c.Conn.Write(p)
}

func (c *recordingConn) take() []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	p := bytes.Clone(c.written.Bytes())
	c.written.Reset()
	return p
}

func TestCompressionNegotiationAndThreshold(t *testing.T) {
	for _, negotiated := range []bool{false, true} {
		name := "not negotiated"
		if negotiated {
			name = "negotiated"
		}
		t.Run(name, func(t *testing.T) {
			serverOptions, clientOptions := testOptions(), testOptions()
			serverOptions.CompressEnabled = negotiated
			clientOptions.CompressEnabled = true
			clientOptions.CompressThreshold = 64
			var recorded *recordingConn
			clientOptions.Dialer.NetDialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
				conn, err := (&net.Dialer{}).DialContext(ctx, network, address)
				if err != nil {
					return nil, err
				}
				recorded = &recordingConn{Conn: conn}
				return recorded, nil
			}
			server, client := transportPair(t, serverOptions, clientOptions, false)
			_ = recorded.take() // Discard the handshake request.
			for _, size := range []int{63, 64, 1024, 10} {
				payload := bytes.Repeat([]byte("a"), size)
				original := bytes.Clone(payload)
				if _, err := client.Write(payload); err != nil {
					t.Fatal(err)
				}
				frame := recorded.take()
				if len(frame) < 2 {
					t.Fatal("no frame written")
				}
				wantCompressed := negotiated && int64(size) >= clientOptions.CompressThreshold
				if compressed := frame[0]&0x40 != 0; compressed != wantCompressed {
					t.Fatalf("size %d: RSV1 = %v; want %v", size, compressed, wantCompressed)
				}
				if frame[1]&0x80 == 0 {
					t.Fatal("client frame is not masked")
				}
				if !bytes.Equal(payload, original) {
					t.Fatal("Write modified caller's payload")
				}
				got, err := io.ReadAll(server)
				if err != nil || !bytes.Equal(got, original) {
					t.Fatalf("payload = %q, %v", got, err)
				}
			}
		})
	}
}
