package websocket

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"github.com/yukino-js/netty-go/transport"
)

type websocketTransport struct {
	net.Conn
	conn        *websocket.Conn
	options     *Options
	opCode      int
	request     *http.Request
	msgReader   io.Reader
	messageSize int64
	utf8Tail    []byte
	writeLocker sync.Mutex
}

func newWebsocketTransport(conn *websocket.Conn, options *Options, request *http.Request) (*websocketTransport, error) {
	raw := conn.NetConn()
	if tlsConn, ok := raw.(*tls.Conn); ok {
		raw = tlsConn.NetConn()
	}
	if tcpConn, ok := raw.(*net.TCPConn); ok {
		if err := tcpConn.SetNoDelay(options.NoDelay); err != nil {
			return nil, err
		}
	}
	if options.CompressEnabled {
		if err := conn.SetCompressionLevel(options.CompressLevel); err != nil {
			return nil, err
		}
	}
	conn.SetReadLimit(options.MaxFrameSize)
	opCode := websocket.TextMessage
	if options.OpCode&websocket.TextMessage == 0 {
		opCode = websocket.BinaryMessage
	}
	return &websocketTransport{
		Conn: conn.NetConn(), conn: conn, options: options, opCode: opCode, request: request,
	}, nil
}

func (t *websocketTransport) Route() string          { return t.request.URL.Path }
func (t *websocketTransport) Header() http.Header    { return t.request.Header }
func (t *websocketTransport) Request() *http.Request { return t.request }

// Read returns io.EOF at each message boundary, including fragmented messages.
func (t *websocketTransport) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if t.msgReader == nil {
		messageType, reader, err := t.conn.NextReader()
		if err != nil {
			return 0, err
		}
		if messageType&t.options.OpCode == 0 {
			return 0, t.readError(websocket.CloseUnsupportedData, "unsupported data type")
		}
		t.msgReader = reader
		t.messageSize = 0
		t.utf8Tail = nil
		if t.options.CheckUTF8 && messageType == websocket.TextMessage {
			t.utf8Tail = make([]byte, 0, utf8.UTFMax)
		}
	}
	// Also limit decompressed payloads; Gorilla's limit counts wire bytes.
	if limit := t.options.MaxFrameSize; limit > 0 {
		remaining := limit - t.messageSize
		if int64(len(p)) > remaining {
			p = p[:remaining+1]
		}
	}
	n, err := t.msgReader.Read(p)
	t.messageSize += int64(n)
	if limit := t.options.MaxFrameSize; limit > 0 && t.messageSize > limit {
		return 0, t.readError(websocket.CloseMessageTooBig, "message too large")
	}
	if t.utf8Tail != nil && !t.validUTF8(p[:n], err == io.EOF) {
		return 0, t.readError(websocket.CloseInvalidFramePayloadData, "invalid UTF-8")
	}
	if err == io.EOF {
		t.msgReader = nil
	}
	return n, err
}

// validUTF8 keeps only an incomplete rune between reads.
func (t *websocketTransport) validUTF8(p []byte, eof bool) bool {
	if len(t.utf8Tail) > 0 {
		for len(p) > 0 && !utf8.FullRune(t.utf8Tail) {
			t.utf8Tail = append(t.utf8Tail, p[0])
			p = p[1:]
		}
		if utf8.FullRune(t.utf8Tail) {
			if r, size := utf8.DecodeRune(t.utf8Tail); r == utf8.RuneError && size == 1 {
				return false
			}
			t.utf8Tail = t.utf8Tail[:0]
		}
	}
	for len(p) > 0 {
		if !utf8.FullRune(p) {
			t.utf8Tail = append(t.utf8Tail, p...)
			break
		}
		r, size := utf8.DecodeRune(p)
		if r == utf8.RuneError && size == 1 {
			return false
		}
		p = p[size:]
	}
	return !eof || len(t.utf8Tail) == 0
}

func (t *websocketTransport) readError(code int, reason string) error {
	_ = t.WriteClose(code, reason)
	return &websocket.CloseError{Code: code, Text: reason}
}

func (t *websocketTransport) Write(p []byte) (int, error) {
	t.writeLocker.Lock()
	defer t.writeLocker.Unlock()
	t.conn.EnableWriteCompression(t.options.CompressEnabled && int64(len(p)) >= t.options.CompressThreshold)
	if err := t.conn.WriteMessage(t.opCode, p); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (t *websocketTransport) Writev(buffers transport.Buffers) (int64, error) {
	var n int64
	for _, p := range buffers {
		written, err := t.Write(p)
		n += int64(written)
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (t *websocketTransport) WriteClose(code int, reason string) error {
	return t.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(code, reason), time.Now().Add(time.Second))
}

func (t *websocketTransport) SetWriteDeadline(deadline time.Time) error {
	// Interrupt a blocked socket write before waiting for Gorilla's writer lock.
	if err := t.Conn.SetWriteDeadline(deadline); err != nil {
		return err
	}
	t.writeLocker.Lock()
	defer t.writeLocker.Unlock()
	return t.conn.SetWriteDeadline(deadline)
}

func (t *websocketTransport) SetDeadline(deadline time.Time) error {
	if err := t.SetReadDeadline(deadline); err != nil {
		return err
	}
	return t.SetWriteDeadline(deadline)
}

func (t *websocketTransport) Flush() error              { return nil }
func (t *websocketTransport) Close() error              { return t.conn.Close() }
func (t *websocketTransport) RawTransport() interface{} { return t.Conn }

var _ transport.Transport = (*websocketTransport)(nil)
