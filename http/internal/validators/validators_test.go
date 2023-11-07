package validators

import (
	"github.com/goclarum/clarum/http/constants"
	"github.com/goclarum/clarum/http/message"
	"net/http"
	"testing"
)

func TestValidateHeadersOK(t *testing.T) {
	expectedMessage := createTestMessageWithHeaders()
	req := createRealRequest()

	if err := validateHeaders(expectedMessage, req.Header); err != nil {
		t.Errorf("No header validation error expected, but got %s", err)
	}
}

func TestValidateHeadersError(t *testing.T) {
	expectedMessage := createTestMessageWithHeaders()
	expectedMessage.ETag("1234")

	req := createRealRequest()

	err := validateHeaders(expectedMessage, req.Header)

	if err == nil {
		t.Errorf("Header validation error expected, but got none")
	}

	if err.Error() != "Validation error: header <ETag> mismatch. Expected [1234] but received [33a64df551425fcc55e4d42a148795d9f25f89d4]" {
		t.Errorf("Header validation error message is unexpected")
	}
}

func TestValidateQueryParamsOK(t *testing.T) {
	expectedMessage := message.Get("myPath").
		QueryParam("param1", "value1").
		QueryParam("param2", "value2")

	req := createRealRequest()
	qParams := req.URL.Query()
	qParams.Set("param1", "value1")
	qParams.Set("param2", "value2")
	req.URL.RawQuery = qParams.Encode()

	if err := validateQueryParams(expectedMessage, req.URL.Query()); err != nil {
		t.Errorf("No query param validation error expected, but got %s", err)
	}
}

func TestValidateQueryParamsParamMismatch(t *testing.T) {
	expectedMessage := message.Get("myPath").
		QueryParam("param1", "value1").
		QueryParam("param2", "value2")

	req := createRealRequest()
	qParams := req.URL.Query()
	qParams.Set("param1", "value1")
	qParams.Set("param3", "value2")
	req.URL.RawQuery = qParams.Encode()

	err := validateQueryParams(expectedMessage, req.URL.Query())
	if err == nil {
		t.Errorf("Query param validation error expected, but got none")
	}

	if err.Error() != "Validation error: query param <param2> missing" {
		t.Errorf("Query param validation error message is unexpected")
	}
}

func TestValidateQueryParamsValueMismatch(t *testing.T) {
	expectedMessage := message.Get("myPath").
		QueryParam("param1", "value1").
		QueryParam("param2", "value2")

	req := createRealRequest()
	qParams := req.URL.Query()
	qParams.Set("param1", "value1")
	qParams.Set("param2", "value22")
	req.URL.RawQuery = qParams.Encode()

	err := validateQueryParams(expectedMessage, req.URL.Query())
	if err == nil {
		t.Errorf("Query param validation error expected, but got none")
	}

	if err.Error() != "Validation error: query param <param2> values mismatch. Expected [value2] but received [[value22]]" {
		t.Errorf("Query param validation error message is unexpected")
	}
}

func createTestMessageWithHeaders() *message.Message {
	return message.Post("myPath").
		Header("Connection", "keep-alive").
		ContentType("application/json").
		Authorization("Bearer 0b79bab50daca910b000d4f1a2b675d604257e42").
		ETag("33a64df551425fcc55e4d42a148795d9f25f89d4")
}

func createRealRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "myPath", nil)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set(constants.ContentTypeHeaderName, "application/json")
	req.Header.Set(constants.AuthorizationHeaderName, "Bearer 0b79bab50daca910b000d4f1a2b675d604257e42")
	req.Header.Set(constants.ETagHeaderName, "33a64df551425fcc55e4d42a148795d9f25f89d4")

	return req
}
