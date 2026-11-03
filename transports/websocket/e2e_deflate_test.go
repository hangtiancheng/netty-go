package websocket

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestPerMessageDeflateInterop(t *testing.T) {
	opts := testOptions()
	opts.CompressEnabled = true
	opts.CompressThreshold = 1
	server, client := transportPair(t, opts, opts, false)
	extensions := client.Header().Get("Sec-WebSocket-Extensions")
	for _, parameter := range []string{"permessage-deflate", "server_no_context_takeover", "client_no_context_takeover"} {
		if !strings.Contains(extensions, parameter) {
			t.Fatalf("missing %q in negotiated extensions: %q", parameter, extensions)
		}
	}
	for range 3 {
		want := bytes.Repeat([]byte("hello-compress-中文"), 1000)
		if _, err := client.Write(want); err != nil {
			t.Fatal(err)
		}
		got, err := io.ReadAll(server)
		if err != nil || !bytes.Equal(got, want) {
			t.Fatalf("server read: %v, length %d", err, len(got))
		}
		if _, err := server.Write(got); err != nil {
			t.Fatal(err)
		}
		got, err = io.ReadAll(client)
		if err != nil || !bytes.Equal(got, want) {
			t.Fatalf("client read: %v, length %d", err, len(got))
		}
	}
}
