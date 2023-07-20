package http

import (
	"net/http"
	"testing"
)

func TestValidateHeadersOK(t *testing.T) {
	action := createTestActionWithHeaders()
	req := createRealRequest()

	if err := validateHeaders(action, req.Header); err != nil {
		t.Errorf("No header validation error expected, but got %s", err)
	}
}

func TestValidateHeadersError(t *testing.T) {
	action := createTestActionWithHeaders()
	action.ETag("1234")

	req := createRealRequest()

	err := validateHeaders(action, req.Header)

	if err == nil {
		t.Errorf("Header validation error expected, but got none")
	}

	if err.Error() != "validation error: header <ETag> mismatch" {
		t.Errorf("Header validation error mismatch")
	}
}

func TestValidateQueryParamsOK(t *testing.T) {
	action := Get("myPath").
		QueryParam("param1", "value1").
		QueryParam("param2", "value2")

	req := createRealRequest()
	qParams := req.URL.Query()
	qParams.Set("param1", "value1")
	qParams.Set("param2", "value2")
	req.URL.RawQuery = qParams.Encode()

	if err := validateQueryParams(action, req.URL.Query()); err != nil {
		t.Errorf("No query param validation error expected, but got %s", err)
	}
}

func TestValidateQueryParamsParamMismatch(t *testing.T) {
	action := Get("myPath").
		QueryParam("param1", "value1").
		QueryParam("param2", "value2")

	req := createRealRequest()
	qParams := req.URL.Query()
	qParams.Set("param1", "value1")
	qParams.Set("param3", "value2")
	req.URL.RawQuery = qParams.Encode()

	err := validateQueryParams(action, req.URL.Query())
	if err == nil {
		t.Errorf("Query param validation error expected, but got none")
	}

	if err.Error() != "validation error: query param <param2> missing" {
		t.Errorf("Query param validation error mismatch")
	}
}

func TestValidateQueryParamsValueMismatch(t *testing.T) {
	action := Get("myPath").
		QueryParam("param1", "value1").
		QueryParam("param2", "value2")

	req := createRealRequest()
	qParams := req.URL.Query()
	qParams.Set("param1", "value1")
	qParams.Set("param2", "value22")
	req.URL.RawQuery = qParams.Encode()

	err := validateQueryParams(action, req.URL.Query())
	if err == nil {
		t.Errorf("Query param validation error expected, but got none")
	}

	if err.Error() != "validation error: query params mismatch: expected [value22], actual value2" {
		t.Errorf("Query param validation error mismatch")
	}
}

func createTestActionWithHeaders() *Action {
	return Post("myPath").
		Header("Connection", "keep-alive").
		ContentType("application/json").
		Authorization("Bearer 0b79bab50daca910b000d4f1a2b675d604257e42").
		ETag("33a64df551425fcc55e4d42a148795d9f25f89d4")
}

func createRealRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "myPath", nil)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set(ContentTypeHeaderName, "application/json")
	req.Header.Set(AuthorizationHeaderName, "Bearer 0b79bab50daca910b000d4f1a2b675d604257e42")
	req.Header.Set(ETagHeaderName, "33a64df551425fcc55e4d42a148795d9f25f89d4")

	return req
}
