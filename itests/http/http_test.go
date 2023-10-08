package http

import (
	clrm "github.com/goclarum/clarum/http"
	"net/http"
	"testing"
)

func TestGet(t *testing.T) {
	testClient.Send(t, clrm.Get())

	testServer.Receive(t, clrm.Get())
	testServer.Send(clrm.Response(http.StatusOK))

	testClient.Receive(t, clrm.Response(http.StatusOK))
}

func TestPost(t *testing.T) {
	testClient.Send(t, clrm.Post().QueryParam("myParam", "myValue1").
		Payload("plain text payload"))

	testServer.Receive(t, clrm.Post().QueryParam("myParam", "myValue1").
		Payload("plain text payload"))
	testServer.Send(clrm.Response(http.StatusOK).ContentType("text/xml"))

	testClient.Receive(t, clrm.Response(http.StatusOK).ContentType("text/xml"))
}
