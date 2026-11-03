package websocket

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	ws "github.com/gorilla/websocket"
	"github.com/yukino-js/netty-go/transport"
)

func testOptions() *Options {
	opts := *DefaultOptions
	opts.ServeMux = http.NewServeMux()
	return &opts
}

func transportPair(t *testing.T, serverOptions, clientOptions *Options, secure bool) (*websocketTransport, *websocketTransport) {
	t.Helper()
	serverOptions = serverOptions.Apply()
	type result struct {
		conn *websocketTransport
		err  error
	}
	accepted := make(chan result, 1)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := serverOptions.Upgrader.Upgrade(w, r, serverOptions.ResponseHeader)
		if err != nil {
			accepted <- result{err: err}
			return
		}
		tr, err := newWebsocketTransport(conn, serverOptions, r)
		if err != nil {
			_ = conn.Close()
		}
		accepted <- result{conn: tr, err: err}
	})
	var server *httptest.Server
	if secure {
		server = httptest.NewTLSServer(handler)
		clientCopy := *clientOptions
		clientOptions = &clientCopy
		clientOptions.TLS = server.Client().Transport.(*http.Transport).TLSClientConfig
	} else {
		server = httptest.NewServer(handler)
	}
	t.Cleanup(server.Close)
	opts, err := transport.ParseOptions(context.Background(), "ws"+strings.TrimPrefix(server.URL, "http")+"/echo?token=hello%20world", WithOptions(clientOptions))
	if err != nil {
		t.Fatal(err)
	}
	conn, err := New().Connect(opts)
	if err != nil {
		t.Fatal(err)
	}
	client := conn.(*websocketTransport)
	t.Cleanup(func() { _ = client.Close() })
	var peer result
	select {
	case peer = <-accepted:
	case <-time.After(3 * time.Second):
		t.Fatal("server did not accept connection")
	}
	if peer.err != nil {
		t.Fatal(peer.err)
	}
	t.Cleanup(func() { _ = peer.conn.Close() })
	deadline := time.Now().Add(5 * time.Second)
	if err := client.SetDeadline(deadline); err != nil {
		t.Fatal(err)
	}
	if err := peer.conn.SetDeadline(deadline); err != nil {
		t.Fatal(err)
	}
	return peer.conn, client
}

func TestMessageBoundaries(t *testing.T) {
	for _, opCode := range []int{ws.TextMessage, ws.BinaryMessage} {
		t.Run(string(rune('0'+opCode)), func(t *testing.T) {
			opts := testOptions()
			opts.OpCode = opCode
			server, client := transportPair(t, opts, opts, false)
			messages := transport.Buffers{[]byte("first"), {}, bytes.Repeat([]byte("large"), 2000), []byte("last")}
			if _, err := client.Writev(messages); err != nil {
				t.Fatal(err)
			}
			for _, want := range messages {
				got, err := io.ReadAll(server)
				if err != nil || !bytes.Equal(got, want) {
					t.Fatalf("message = %q, %v; want %q", got, err, want)
				}
				if _, err := server.Write(got); err != nil {
					t.Fatal(err)
				}
				got, err = io.ReadAll(client)
				if err != nil || !bytes.Equal(got, want) {
					t.Fatalf("echo = %q, %v; want %q", got, err, want)
				}
			}
		})
	}
}

func TestFragmentedUTF8AndControlFrames(t *testing.T) {
	opts := testOptions()
	opts.CheckUTF8 = true
	server, client := transportPair(t, opts, opts, false)
	writer, err := client.conn.NextWriter(ws.TextMessage)
	if err != nil {
		t.Fatal(err)
	}
	// Exceed Gorilla's write buffer to produce continuation frames and split
	// multi-byte runes across the transport's one-byte reads.
	want := bytes.Repeat([]byte("中文🙂"), 1000)
	if _, err := writer.Write(want[:6001]); err != nil {
		t.Fatal(err)
	}
	if err := client.conn.WriteControl(ws.PingMessage, []byte("ping"), time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write(want[6001:]); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	var got bytes.Buffer
	if _, err := io.CopyBuffer(struct{ io.Writer }{&got}, struct{ io.Reader }{server}, make([]byte, 1)); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got.Bytes(), want) {
		t.Fatal("fragmented UTF-8 message changed")
	}
}

func TestInvalidMessages(t *testing.T) {
	for _, test := range []struct {
		name        string
		opCode      int
		payload     []byte
		limit       int64
		compression bool
		wantCode    int
	}{
		{name: "invalid UTF-8", opCode: ws.TextMessage, payload: []byte{0xff}, wantCode: ws.CloseInvalidFramePayloadData},
		{name: "truncated UTF-8", opCode: ws.TextMessage, payload: []byte{0xe4, 0xb8}, wantCode: ws.CloseInvalidFramePayloadData},
		{name: "unsupported binary", opCode: ws.BinaryMessage, payload: []byte("binary"), wantCode: ws.CloseUnsupportedData},
		{name: "message limit", opCode: ws.TextMessage, payload: bytes.Repeat([]byte("a"), 10000), limit: 100, wantCode: ws.CloseMessageTooBig},
		{name: "decompressed limit", opCode: ws.TextMessage, payload: bytes.Repeat([]byte("a"), 10000), limit: 100, compression: true, wantCode: ws.CloseMessageTooBig},
	} {
		t.Run(test.name, func(t *testing.T) {
			opts := testOptions()
			opts.CheckUTF8 = true
			opts.MaxFrameSize = test.limit
			opts.CompressEnabled = test.compression
			opts.CompressThreshold = 0
			server, client := transportPair(t, opts, opts, false)
			if err := client.conn.WriteMessage(test.opCode, test.payload); err != nil {
				t.Fatal(err)
			}
			if _, err := io.ReadAll(server); err == nil {
				t.Fatal("invalid message accepted")
			}
			_, err := io.ReadAll(client)
			if !ws.IsCloseError(err, test.wantCode) {
				t.Fatalf("close error = %v; want code %d", err, test.wantCode)
			}
		})
	}
}

func TestHeadersAndTLS(t *testing.T) {
	serverOptions, clientOptions := testOptions(), testOptions()
	serverOptions.ResponseHeader = http.Header{"X-Server": {"response"}}
	clientOptions.RequestHeader = http.Header{"X-Client": {"request"}}
	server, client := transportPair(t, serverOptions, clientOptions, true)
	if server.Route() != "/echo" || server.Request().URL.Query().Get("token") != "hello world" || server.Header().Get("X-Client") != "request" {
		t.Fatalf("unexpected server request: %+v", server.Request())
	}
	if client.Header().Get("X-Server") != "response" {
		t.Fatalf("missing response header: %v", client.Header())
	}
	if _, ok := client.RawTransport().(*tls.Conn); !ok {
		t.Fatalf("raw transport = %T; want TLS connection", client.RawTransport())
	}
	if _, err := client.Write([]byte("secure")); err != nil {
		t.Fatal(err)
	}
	if got, err := io.ReadAll(server); err != nil || string(got) != "secure" {
		t.Fatalf("TLS message = %q, %v", got, err)
	}
}

func TestCloseAndWriteDeadline(t *testing.T) {
	t.Run("close reason", func(t *testing.T) {
		server, client := transportPair(t, testOptions(), testOptions(), false)
		if err := client.WriteClose(ws.CloseNormalClosure, "done"); err != nil {
			t.Fatal(err)
		}
		_, err := io.ReadAll(server)
		var closeErr *ws.CloseError
		if !errors.As(err, &closeErr) || closeErr.Code != ws.CloseNormalClosure || closeErr.Text != "done" {
			t.Fatalf("close error = %v", err)
		}
	})
	t.Run("expired deadline", func(t *testing.T) {
		_, client := transportPair(t, testOptions(), testOptions(), false)
		if err := client.SetWriteDeadline(time.Now().Add(-time.Second)); err != nil {
			t.Fatal(err)
		}
		_, err := client.Write([]byte("timeout"))
		var netErr net.Error
		if !errors.As(err, &netErr) || !netErr.Timeout() {
			t.Fatalf("write error = %v; want timeout", err)
		}
	})
}

func TestConcurrentWrites(t *testing.T) {
	_, client := transportPair(t, testOptions(), testOptions(), false)
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 10 {
				if _, err := client.Write([]byte("message")); err != nil {
					t.Error(err)
					return
				}
			}
		})
	}
	wg.Wait()
}

func TestFactoryAcceptor(t *testing.T) {
	opts := testOptions()
	opts.Upgrader.Subprotocols = []string{"netty"}
	opts.Dialer.Subprotocols = []string{"netty"}
	options, err := transport.ParseOptions(context.Background(), "ws://127.0.0.1:0/echo", WithOptions(opts))
	if err != nil {
		t.Fatal(err)
	}
	acceptor, err := New().Listen(options)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = acceptor.Close() })
	_, port, err := net.SplitHostPort(acceptor.(*wsAcceptor).httpServer.Addr)
	if err != nil {
		t.Fatal(err)
	}
	options, err = transport.ParseOptions(context.Background(), "ws://"+net.JoinHostPort("127.0.0.1", port)+"/echo", WithOptions(opts))
	if err != nil {
		t.Fatal(err)
	}
	client, err := New().Connect(options)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close() })
	server, err := acceptor.Accept()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = server.Close() })
	if client.(*websocketTransport).conn.Subprotocol() != "netty" || server.(*websocketTransport).conn.Subprotocol() != "netty" {
		t.Fatal("subprotocol was not negotiated")
	}
	_ = server.SetDeadline(time.Now().Add(3 * time.Second))
	if _, err := client.Write([]byte("acceptor")); err != nil {
		t.Fatal(err)
	}
	if got, err := io.ReadAll(server); err != nil || string(got) != "acceptor" {
		t.Fatalf("accepted payload = %q, %v", got, err)
	}
}
