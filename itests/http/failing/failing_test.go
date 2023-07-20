package http

import (
	clrm "github.com/goclarum/clarum/http"
	"net/http"
	"testing"
)

// this test fails intentionally
// 1. test header validation inbound communication
// 2. test header validation response communication
// 3. test query param validation inbound only
// 4. test query param value validation inbound only

func TestGet(t *testing.T) {
	Client1.Send(clrm.Get())

	Server1.Receive(t, clrm.Get())
	Server1.Send(clrm.Response(http.StatusOK).ContentType("text/xml"))

	Client1.Receive(t, clrm.Response(http.StatusOK))
}

func TestGet2(t *testing.T) {
	Client1.Send(clrm.Get().QueryParam("myParam", "myValue1"))

	Server1.Receive(t, clrm.Get().QueryParam("myParam", "myValue2"))
	Server1.Send(clrm.Response(http.StatusOK).ContentType("text/xml"))

	Client1.Receive(t, clrm.Response(http.StatusOK).ContentType("text/xml"))
}
