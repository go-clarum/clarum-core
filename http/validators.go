package http

import (
	"errors"
	"fmt"
	"github.com/goclarum/clarum/core/arrays"
	clmStrings "github.com/goclarum/clarum/core/validators/strings"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func validateHttpHeaders(t *testing.T, logPrefix string, actionToExecute *Action, headers http.Header) {
	if err := validateHeaders(actionToExecute, headers); err != nil {
		t.Errorf("%s: %s", logPrefix, err)
	} else {
		slog.Debug(fmt.Sprintf("%s: header validation successful", logPrefix))
	}
}

func validateHeaders(action *Action, headers http.Header) error {
	for header, expectedValue := range action.headers {
		if receivedValue := headers.Get(header); expectedValue != receivedValue {
			return errors.New(fmt.Sprintf("Validation error: header <%s> mismatch. Expected [%s] but received [%s]",
				header, expectedValue, receivedValue))
		}
	}

	return nil
}

func validateHttpQueryParams(t *testing.T, logPrefix string, action *Action, url *url.URL) {
	if err := validateQueryParams(action, url.Query()); err != nil {
		t.Errorf("%s: %s", logPrefix, err)
	} else {
		slog.Debug(fmt.Sprintf("%s: query params validation successful", logPrefix))
	}
}

// validate query parameters based on these rules
//
//	-> validate that the param exists
//	-> that the value matches
func validateQueryParams(action *Action, params url.Values) error {
	for param, expectedValue := range action.queryParams {
		if receivedValues, exists := params[param]; exists {
			if !arrays.Contains(receivedValues, expectedValue) {
				return errors.New(fmt.Sprintf("Validation error: query param <%s> values mismatch. Expected [%v] but received [%s]",
					param, expectedValue, receivedValues))
			}
		} else {
			return errors.New(fmt.Sprintf("Validation error: query param <%s> missing", param))
		}
	}

	return nil
}

func validateHttpStatusCode(t *testing.T, logPrefix string, action *Action, statusCode int) {
	if statusCode != action.statusCode {
		t.Errorf("%s: validation error: HTTP status mismatch. Expected [%d] but received [%d]", logPrefix, action.statusCode, statusCode)
	} else {
		slog.Debug(fmt.Sprintf("%s: HTTP status validation successful", logPrefix))
	}
}

func validateHttpBody(t *testing.T, logPrefix string, action *Action, body io.ReadCloser) {
	defer closeBody(logPrefix, body)

	if clmStrings.IsBlank(action.payload) {
		slog.Debug(fmt.Sprintf("%s: action payload is empty. No body validation will be done", logPrefix))
		return
	}

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("%s: could not read response body: %s", logPrefix, err)
	}

	if err := validatePayload(action, bodyBytes); err != nil {
		t.Errorf("%s: %s", logPrefix, err)
	} else {
		slog.Debug(fmt.Sprintf("%s: payload validation successful", logPrefix))
	}
}

func closeBody(logPrefix string, body io.ReadCloser) {
	if err := body.Close(); err != nil {
		slog.Error(fmt.Sprintf("%s: unable to close body: %s", logPrefix, err))
	}
}

func validatePayload(action *Action, payload []byte) error {
	contentTypeHeader := action.headers[ContentTypeHeaderName]
	receivedPayload := string(payload)

	// we let the action decide what kind of validation we do
	if strings.Contains(contentTypeHeader, ContentTypeJsonHeader) {
		// do json validation
	} else if action.payload != receivedPayload { // plain text validation
		return errors.New(fmt.Sprintf("Validation error: payload missmatch. Expected [%s] but received [%s]",
			action.payload, receivedPayload))
	}

	return nil
}
