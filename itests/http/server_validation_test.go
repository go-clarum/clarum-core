package http

import (
	"github.com/goclarum/clarum/http/message"
	"net/http"
	"testing"
)

// Method GET + single query param validation
func TestGet(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Get().QueryParam("myParam", "myValue1"))

	testServer.In(t).Receive().
		Message(message.Get().QueryParam("myParam", "myValue1"))
	testServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}

func TestHead(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Head())

	testServer.In(t).Receive().
		Message(message.Head())
	testServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}

// POST + multiple query params
func TestPost(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Post().
			QueryParam("myParam1", "myValue1").
			QueryParam("myParam2", "myValue1").
			Payload("my plain text payload"))

	testServer.In(t).Receive().
		Message(message.Post().
			QueryParam("myParam1", "myValue1").
			QueryParam("myParam2", "myValue1").
			Payload("my plain text payload"))
	testServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}

// PUT + query param with multiple values
func TestPut(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Put().
			QueryParam("myParam1", "myValue1").
			Payload("my plain text payload"))

	testServer.In(t).Receive().
		Message(message.Put().
			QueryParam("myParam1", "myValue1").
			Payload("my plain text payload"))
	testServer.In(t).Send().
		Message(message.Response(http.StatusCreated))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusCreated))
}

// DELETE
func TestDelete(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Delete())

	testServer.In(t).Receive().
		Message(message.Delete())
	testServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}

// DELETE
func TestOptions(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Options())

	testServer.In(t).Receive().
		Message(message.Options())
	testServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}

// TRACE
func TestTrace(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Trace())

	testServer.In(t).Receive().
		Message(message.Trace())
	testServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}

// PATCH
func TestPatch(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Patch())

	testServer.In(t).Receive().
		Message(message.Patch())
	testServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}
