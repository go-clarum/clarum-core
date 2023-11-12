package http

import (
	"github.com/goclarum/clarum/http/message"
	"net/http"
	"testing"
)

// this test fails intentionally
// 1. test header validation inbound communication
// 2. test header validation response communication
// 3. test query param validation inbound only
// 4. test query param value validation inbound only

func TestGet(t *testing.T) {
	Client1.In(t).Send().Message(message.Get())

	Server1.In(t).Receive().Message(message.Get())
	Server1.Send().
		Message(message.Response(http.StatusOK).ContentType("text/xml"))

	Client1.In(t).Receive().Message(message.Response(http.StatusBadRequest).Payload("something"))
}

func TestGet2(t *testing.T) {
	Client1.In(t).Send().
		Message(message.Get().QueryParam("myParam", "myValue1"))

	Server1.In(t).Receive().
		Message(message.Get().QueryParam("myParam", "myValue2"))
	Server1.Send().
		Message(message.Response(http.StatusOK).ContentType("text/xml"))

	Client1.In(t).Receive().
		Message(message.Response(http.StatusOK).ContentType("text/xml"))
}
