package http

import (
	clrm "github.com/goclarum/clarum/http"
	"net/http"
	"testing"
)

func TestGet(t *testing.T) {
	testClient.Send(t, clrm.Get())

	testServer.Receive(t, clrm.Get())
	testServer.Send(clrm.Response(http.StatusOK).
		Payload("my test"))

	testClient.Receive(t, clrm.Response(http.StatusOK).
		Payload("my test"))
}

func TestPost(t *testing.T) {
	testClient.Send(t, clrm.Post().QueryParam("myParam", "myValue1").
		Payload("my plain text payload"))

	testServer.Receive(t, clrm.Post().QueryParam("myParam", "myValue1").
		Payload("my plain text payload2"))
	testServer.Send(clrm.Response(http.StatusOK).ContentType("text/xml"))

	testClient.Receive(t, clrm.Response(http.StatusOK).ContentType("text/xml"))
}
