package http

import (
	"github.com/goclarum/clarum/http/message"
	"net/http"
	"testing"
)

// Method GET
// + single query param validation
// + URL from client
func TestGet(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Get().QueryParam("myParam", "myValue1"))

	firstTestServer.In(t).Receive().
		Message(message.Get("myApp").QueryParam("myParam", "myValue1"))
	firstTestServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}

// Test method HEAD
// + URL overwrite
func TestHead(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Head("myOtherApp").
			BaseUrl("http://localhost:8084"))

	secondTestServer.In(t).Receive().
		Message(message.Head("myOtherApp").BaseUrl("has no effect on server"))
	secondTestServer.In(t).Send().
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

	firstTestServer.In(t).Receive().
		Message(message.Post().
			QueryParam("myParam1", "myValue1").
			QueryParam("myParam2", "myValue1").
			Payload("my plain text payload"))
	firstTestServer.In(t).Send().
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

	firstTestServer.In(t).Receive().
		Message(message.Put().
			QueryParam("myParam1", "myValue1").
			Payload("my plain text payload"))
	firstTestServer.In(t).Send().
		Message(message.Response(http.StatusCreated))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusCreated))
}

// DELETE
func TestDelete(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Delete())

	firstTestServer.In(t).Receive().
		Message(message.Delete())
	firstTestServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}

// DELETE
func TestOptions(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Options())

	firstTestServer.In(t).Receive().
		Message(message.Options())
	firstTestServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}

// TRACE
func TestTrace(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Trace())

	firstTestServer.In(t).Receive().
		Message(message.Trace())
	firstTestServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}

// PATCH
func TestPatch(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Patch())

	firstTestServer.In(t).Receive().
		Message(message.Patch())
	firstTestServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}
