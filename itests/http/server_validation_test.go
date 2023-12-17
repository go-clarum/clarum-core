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
		Message(message.Get("/myApp/").QueryParam("myParam", "myValue1"))
	firstTestServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}

// Method HEAD
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

// Method POST
// + multiple query params
func TestPost(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Post().
			QueryParam("myParam1", "myValue1").
			QueryParam("myParam2", "myValue1").
			Payload("my plain text payload"))

	firstTestServer.In(t).Receive().
		Message(message.Post("myApp").
			QueryParam("myParam1", "myValue1").
			QueryParam("myParam2", "myValue1").
			Payload("my plain text payload"))
	firstTestServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}

// Method PUT
// + query param with multiple values
func TestPut(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Put().
			QueryParam("myParam1", "myValue1").
			Payload("my plain text payload"))

	firstTestServer.In(t).Receive().
		Message(message.Put("myApp").
			QueryParam("myParam1", "myValue1").
			Payload("my plain text payload"))
	firstTestServer.In(t).Send().
		Message(message.Response(http.StatusCreated))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusCreated))
}

// Method DELETE
// + path validation
func TestDelete(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Delete("my", "/", "resource", "", "1234"))

	firstTestServer.In(t).Receive().
		Message(message.Delete("myApp/my/resource/1234"))
	firstTestServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}

// Method OPTIONS
func TestOptions(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Options())

	firstTestServer.In(t).Receive().
		Message(message.Options("myApp"))
	firstTestServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}

// Method TRACE
func TestTrace(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Trace())

	firstTestServer.In(t).Receive().
		Message(message.Trace("myApp"))
	firstTestServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}

// Method PATCH
func TestPatch(t *testing.T) {
	testClient.In(t).Send().
		Message(message.Patch())

	firstTestServer.In(t).Receive().
		Message(message.Patch("myApp"))
	firstTestServer.In(t).Send().
		Message(message.Response(http.StatusOK))

	testClient.In(t).Receive().
		Message(message.Response(http.StatusOK))
}
