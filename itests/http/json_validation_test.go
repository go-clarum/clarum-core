package http

import (
	"github.com/goclarum/clarum/http/message"
	"net/http"
	"testing"
)

func TestJson(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Put().
			Json().
			Payload(""))

	firstTestServer.In(t).Receive().
		Message(message.Put("myApp"))
	firstTestServer.In(t).Send().
		Message(message.Response(http.StatusCreated))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusCreated))
}
