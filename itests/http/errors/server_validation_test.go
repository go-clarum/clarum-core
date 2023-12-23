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

	e2 := errorsServer.Receive().Message(message.Post("myApp"))
	e3 := errorsServer.Send().
		Message(message.Response(http.StatusInternalServerError))

	e4 := errorsClient.Receive().
		Message(message.Response(http.StatusInternalServerError))

	checkErrors(t, expectedError, e1, e2, e3, e4)
}

// HTTP status code error.
// Server receives a message to send with an invalid HTTP status code -> default error response because of the error
func TestInvalidStatusCode(t *testing.T) {
	expectedError := "HTTP server errorsServer: message to send is invalid - unsupported status code [99]\n" +
		"HTTP client errorsClient: validation error - HTTP status mismatch - expected [200] but received [500]"

	e1 := errorsClient.Send().Message(message.Get().BaseUrl("http://localhost:8083/myApp"))

	e2 := errorsServer.Receive().Message(message.Get("myApp"))
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

	e2 := errorsServer.Receive().Message(message.Get("myApp"))
	e3 := errorsServer.Send().
		Message(message.Response(http.StatusBadRequest))

	e4 := errorsClient.Receive().
		Message(message.Response(http.StatusOK))

	checkErrors(t, expectedError, e1, e2, e3, e4)
}

// HTTP path validation error.
// Server responds with 404 Bad Request & client expects 200 OK
func TestPathValidation(t *testing.T) {
	expectedError := "HTTP server errorsServer: validation error - HTTP path mismatch - expected [my/resource/5433] but received [my/resource/1234]"

	e1 := errorsClient.Send().
		Message(message.Get("my", "resource", "1234").
			BaseUrl("http://localhost:8083"))

	e2 := errorsServer.Receive().Message(message.Get("my", "resource", "5433"))
	e3 := errorsServer.Send().
		Message(message.Response(http.StatusNotFound))

	e4 := errorsClient.Receive().
		Message(message.Response(http.StatusNotFound))

	checkErrors(t, expectedError, e1, e2, e3, e4)
}

// HTTP header validation error: multiple headers, one missing
func TestHeaderMissingValidation(t *testing.T) {
	expectedError := "HTTP server errorsServer: validation error - header <traceid> missing"

	e1 := errorsClient.Send().
		Message(message.Get().
			BaseUrl("http://localhost:8083").
			Authorization("Bearer: 123152123123"))

	e2 := errorsServer.Receive().Message(message.Get().
		Authorization("Bearer: 123152123123").
		Header("traceid", "777777777"))
	e3 := errorsServer.Send().
		Message(message.Response(http.StatusInternalServerError))

	e4 := errorsClient.Receive().
		Message(message.Response(http.StatusInternalServerError))

	checkErrors(t, expectedError, e1, e2, e3, e4)
}

// HTTP header validation error: header value incorrect
func TestHeaderInvalidValidation(t *testing.T) {
	expectedError := "HTTP server errorsServer: validation error - header <authorization> mismatch - expected [Bearer: 234121] but received [[Bearer: 123152123123]]"

	e1 := errorsClient.Send().
		Message(message.Get().
			BaseUrl("http://localhost:8083").
			Authorization("Bearer: 123152123123"))

	e2 := errorsServer.Receive().Message(message.Get().
		Authorization("Bearer: 234121"))
	e3 := errorsServer.Send().
		Message(message.Response(http.StatusInternalServerError))

	e4 := errorsClient.Receive().
		Message(message.Response(http.StatusInternalServerError))

	checkErrors(t, expectedError, e1, e2, e3, e4)
}

// HTTP query params validation error: query param missing
func TestQueryParamMissingValidation(t *testing.T) {
	expectedError := "HTTP server errorsServer: validation error - query param <param2> missing"

	e1 := errorsClient.Send().
		Message(message.Get().
			BaseUrl("http://localhost:8083").
			QueryParam("param1", "value1"))

	e2 := errorsServer.Receive().Message(message.Get().
		QueryParam("param1", "value1").
		QueryParam("param2", "value2"))
	e3 := errorsServer.Send().
		Message(message.Response(http.StatusInternalServerError))

	e4 := errorsClient.Receive().
		Message(message.Response(http.StatusInternalServerError))

	checkErrors(t, expectedError, e1, e2, e3, e4)
}

// HTTP query params validation error: query param value mismatch
func TestQueryParamInvalidValueValidation(t *testing.T) {
	expectedError := "HTTP server errorsServer: validation error - query param <param2> values mismatch - expected [[value3]] but received [[value2]]"

	e1 := errorsClient.Send().
		Message(message.Get().
			BaseUrl("http://localhost:8083").
			QueryParam("param1", "value1").
			QueryParam("param2", "value2"))

	e2 := errorsServer.Receive().Message(message.Get().
		QueryParam("param1", "value1").
		QueryParam("param2", "value3"))
	e3 := errorsServer.Send().
		Message(message.Response(http.StatusInternalServerError))

	e4 := errorsClient.Receive().
		Message(message.Response(http.StatusInternalServerError))

	checkErrors(t, expectedError, e1, e2, e3, e4)
}

// HTTP query params validation error: query param multi value mismatch
func TestQueryParamInvalidMultiValueValidation(t *testing.T) {
	expectedError := "HTTP server errorsServer: validation error - query param <param2> values mismatch - expected [[value2 value3]] but received [[value2 value4]]"

	e1 := errorsClient.Send().
		Message(message.Get().
			BaseUrl("http://localhost:8083").
			QueryParam("param1", "value1").
			QueryParam("param2", "value2", "value4"))

	e2 := errorsServer.Receive().Message(message.Get().
		QueryParam("param1", "value1").
		QueryParam("param2", "value2", "value3"))
	e3 := errorsServer.Send().
		Message(message.Response(http.StatusInternalServerError))

	e4 := errorsClient.Receive().
		Message(message.Response(http.StatusInternalServerError))

	checkErrors(t, expectedError, e1, e2, e3, e4)
}
