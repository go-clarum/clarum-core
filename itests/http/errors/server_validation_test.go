package errors

import (
	"github.com/goclarum/clarum/http/message"
	"net/http"
	"testing"
)

// The following tests check server validation errors.

// HTTP method validation error.
// Client sends HTTP GET & server expects POST
func TestMethodValidation(t *testing.T) {
	expectedError := "HTTP server errorsServer: validation error - HTTP method mismatch - expected [POST] but received [GET]"

	e1 := errorsClient.Send().Message(message.Get().BaseUrl("http://localhost:8083/myApp"))

	e2 := errorsServer.Receive().Message(message.Post())
	e3 := errorsServer.Send().
		Message(message.Response(http.StatusInternalServerError))

	e4 := errorsClient.Receive().
		Message(message.Response(http.StatusInternalServerError))

	checkErrors(t, expectedError, e1, e2, e3, e4)
}

// HTTP status code error.
// Server receives a message to send with an invalid HTTP status code
func TestInvalidStatusCode(t *testing.T) {
	expectedError := "HTTP server errorsServer: message to send is invalid - unsupported status code [99]"

	e1 := errorsClient.Send().Message(message.Get().BaseUrl("http://localhost:8083/myApp"))

	e2 := errorsServer.Receive().Message(message.Get())
	e3 := errorsServer.Send().
		Message(message.Response(99))

	e4 := errorsClient.Receive().
		Message(message.Response(http.StatusOK))

	checkErrors(t, expectedError, e1, e2, e3, e4)
}

// HTTP status code validation error.
// Server responds with 400 Bad Request & client expects 200 OK
func TestStatusCodeValidation(t *testing.T) {
	expectedError := "HTTP client errorsClient: validation error - HTTP status mismatch - expected [200] but received [400]"

	e1 := errorsClient.Send().Message(message.Get().BaseUrl("http://localhost:8083/myApp"))

	e2 := errorsServer.Receive().Message(message.Get())
	e3 := errorsServer.Send().
		Message(message.Response(http.StatusBadRequest))

	e4 := errorsClient.Receive().
		Message(message.Response(http.StatusOK))

	checkErrors(t, expectedError, e1, e2, e3, e4)
}
