package validators

import (
	"errors"
	"fmt"
	"github.com/goclarum/clarum/core/arrays"
	clarumstrings "github.com/goclarum/clarum/core/validators/strings"
	"github.com/goclarum/clarum/http/constants"
	"github.com/goclarum/clarum/http/message"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func ValidateHttpHeaders(t *testing.T, logPrefix string, message *message.Message, headers http.Header) {
	if err := validateHeaders(message, headers); err != nil {
		t.Errorf("%s: %s", logPrefix, err)
	} else {
		slog.Debug(fmt.Sprintf("%s: header validation successful", logPrefix))
	}
}

func validateHeaders(message *message.Message, headers http.Header) error {
	for header, expectedValue := range message.Headers {
		if receivedValue := headers.Get(header); expectedValue != receivedValue {
			return errors.New(fmt.Sprintf("Validation error: header <%s> mismatch. Expected [%s] but received [%s]",
				header, expectedValue, receivedValue))
		}
	}

	return nil
}

func ValidateHttpQueryParams(t *testing.T, logPrefix string, message *message.Message, url *url.URL) {
	if err := validateQueryParams(message, url.Query()); err != nil {
		t.Errorf("%s: %s", logPrefix, err)
	} else {
		slog.Debug(fmt.Sprintf("%s: query params validation successful", logPrefix))
	}
}

// validate query parameters based on these rules
//
//	-> validate that the param exists
//	-> that the value matches
func validateQueryParams(message *message.Message, params url.Values) error {
	for param, expectedValue := range message.QueryParams {
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

func ValidateHttpStatusCode(t *testing.T, logPrefix string, message *message.Message, statusCode int) {
	if statusCode != message.StatusCode {
		t.Errorf("%s: validation error: HTTP status mismatch. Expected [%d] but received [%d]", logPrefix, message.StatusCode, statusCode)
	} else {
		slog.Debug(fmt.Sprintf("%s: HTTP status validation successful", logPrefix))
	}
}

func ValidateHttpBody(t *testing.T, logPrefix string, message *message.Message, body io.ReadCloser) {
	defer closeBody(logPrefix, body)

	if clarumstrings.IsBlank(message.MessagePayload) {
		slog.Debug(fmt.Sprintf("%s: message payload is empty. No body validation will be done", logPrefix))
		return
	}

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		t.Errorf("%s: could not read response body: %s", logPrefix, err)
	}

	if err := validatePayload(message, bodyBytes); err != nil {
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

func validatePayload(message *message.Message, payload []byte) error {
	contentTypeHeader := message.Headers[constants.ContentTypeHeaderName]
	receivedPayload := string(payload)

	// we let the message decide what kind of validation we do
	if strings.Contains(contentTypeHeader, constants.ContentTypeJsonHeader) {
		// do json validation
	} else if message.MessagePayload != receivedPayload { // plain text validation
		return errors.New(fmt.Sprintf("Validation error: payload missmatch. Expected [%s] but received [%s]",
			message.MessagePayload, receivedPayload))
	}

	return nil
}
