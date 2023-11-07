package http

import (
	"github.com/goclarum/clarum/http/message"
	"net/http"
	"testing"
)

func TestGet(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Get())

	testServer.In(t).Receive().
		Message(message.Get())
	testServer.In(t).Send().
		Message(message.Response(http.StatusOK).Payload("my test"))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK).Payload("my test"))
}

func TestPost(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Post().
			QueryParam("myParam", "myValue1").
			Payload("my plain text payload"))

	testServer.In(t).Receive().
		Message(message.Post().
			QueryParam("myParam", "myValue1").
			Payload("my plain text payload2"))
	testServer.In(t).Send().
		Message(message.Response(http.StatusOK).ContentType("text/xml"))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK).ContentType("text/xml"))
}
