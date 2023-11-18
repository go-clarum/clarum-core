package errors

import (
	"github.com/goclarum/clarum/http/message"
	"net/http"
	"testing"
)

// The following tests check server validation errors.

// HTTP method validation error. Client sends HTTP GET & server expects POST
func TestMethodValidation(t *testing.T) {
	expectedError := "HTTP server errorsServer: validation error - HTTP method mismatch - expected [POST] but received [GET]"

	e1 := errorsClient.Send().Message(message.Get().BaseUrl("http://localhost:8083/myApp"))

	e2 := errorsServer.Receive().Message(message.Post())
	errorsServer.Send().
		Message(message.Response(http.StatusInternalServerError))

	e3 := errorsClient.Receive().
		Message(message.Response(http.StatusInternalServerError))

	checkErrors(t, expectedError, e1, e2, e3)
}
