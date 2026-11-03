package netty_ws

import (
	"bytes"
	"compress/flate"
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

type testDialer struct{ called atomic.Bool }

func (d *testDialer) Dial(network, addr string) (net.Conn, error) {
	return nil, errors.New("DialContext should be used")
}

func (d *testDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	d.called.Store(true)
	return (&net.Dialer{}).DialContext(ctx, network, addr)
}

func TestHTTPUpgradeAndCloseError(t *testing.T) {
	engine := NewWebsocket(WithCompress(flate.BestSpeed, 1), WithServerHeader(http.Header{"X-Server": {"netty"}}))
	opened := make(chan Conn, 1)
	closed := make(chan error, 1)
	engine.OnOpen = func(conn Conn) { opened <- conn }
	engine.OnData = func(conn Conn, data []byte) {
		if err := conn.Write(data); err != nil {
			t.Error(err)
		}
	}
	engine.OnClose = func(_ Conn, err error) { closed <- err }
	server := httptest.NewServer(engine)
	t.Cleanup(server.Close)
	t.Cleanup(func() { _ = engine.Close() })
	dialer := *websocket.DefaultDialer
	dialer.EnableCompression = true
	client, response, err := dialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"/echo?token=value", http.Header{"X-Client": {"gorilla"}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close() })
	_ = client.SetReadDeadline(time.Now().Add(5 * time.Second))
	var accepted Conn
	select {
	case accepted = <-opened:
	case <-time.After(3 * time.Second):
		t.Fatal("OnOpen was not called")
	}
	if accepted.Header().Get("X-Client") != "gorilla" || response.Header.Get("X-Server") != "netty" || accepted.Request().URL.Query().Get("token") != "value" {
		t.Fatal("handshake headers or URL were lost")
	}
	want := bytes.Repeat([]byte("hello"), 1000)
	if err := client.WriteMessage(websocket.TextMessage, want); err != nil {
		t.Fatal(err)
	}
	_, got, err := client.ReadMessage()
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("echo = %q, %v", got, err)
	}
	if err := client.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done"), time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-closed:
		var closeErr ClosedError
		if !errors.As(err, &closeErr) || closeErr.Code != websocket.CloseNormalClosure || closeErr.Reason != "done" {
			t.Fatalf("OnClose error = %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("OnClose was not called")
	}
}

func TestOpenHeadersAndCustomDialer(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Client") != "netty" {
			t.Error("missing client header")
		}
		conn, err := upgrader.Upgrade(w, r, http.Header{"X-Server": {"gorilla"}})
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			t.Error(err)
			return
		}
		if err := conn.WriteMessage(messageType, payload); err != nil {
			t.Error(err)
		}
		// Keep reading until the netty client closes.
		_, _, _ = conn.ReadMessage()
	}))
	t.Cleanup(server.Close)
	dialer := &testDialer{}
	client := NewWebsocket(WithClientHeader(http.Header{"X-Client": {"netty"}}), WithDialTimeout(time.Second), WithDialer(dialer))
	received := make(chan []byte, 1)
	client.OnData = func(_ Conn, data []byte) { received <- bytes.Clone(data) }
	t.Cleanup(func() { _ = client.Close() })
	conn, err := client.Open("ws" + strings.TrimPrefix(server.URL, "http"))
	if err != nil {
		t.Fatal(err)
	}
	if conn.Header().Get("X-Server") != "gorilla" {
		t.Fatal("missing server header")
	}
	if !dialer.called.Load() {
		t.Fatal("custom dialer was not called")
	}
	if err := conn.Write([]byte("hello")); err != nil {
		t.Fatal(err)
	}
	select {
	case got := <-received:
		if string(got) != "hello" {
			t.Fatalf("echo = %q", got)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("OnData was not called")
	}
}

func TestHTTPUpgradeRejectsForeignOrigin(t *testing.T) {
	engine := NewWebsocket()
	t.Cleanup(func() { _ = engine.Close() })
	request := httptest.NewRequest(http.MethodGet, "http://example.com/echo", nil)
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Upgrade", "websocket")
	request.Header.Set("Sec-WebSocket-Version", "13")
	request.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	request.Header.Set("Origin", "https://other.example.com")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d; want 403", response.Code)
	}
	if strings.Contains(response.Body.String(), "websocket: request origin not allowed by Upgrader.CheckOrigin") {
		t.Fatal("ServeHTTP appended a second error response")
	}
}
