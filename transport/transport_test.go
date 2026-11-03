package transport

import (
	"net/url"
	"testing"
)

func TestSchemes(t *testing.T) {

	var schemes = Schemes{}.Add("tcp").Add("tcp4").Add("tcp6")

	if !schemes.Valid("tcp") {
		t.Fatal("fail")
	}

	if schemes.Valid("udp") {
		t.Fatal("fail")
	}

	if !schemes.ValidURL("tcp://127.0.0.1:9527") {
		t.Fatal("fail")
	}

	if schemes.ValidURL("udp://127.0.0.1:9527") {
		t.Fatal("fail")
	}

	u, err := url.Parse("//127.0.0.1:9527")
	if nil != err {
		t.Fatal(err)
	}

	if err := schemes.FixScheme(u); nil != err || u.Scheme != "tcp" {
		t.Fatal(err, u.Scheme)
	}

}
