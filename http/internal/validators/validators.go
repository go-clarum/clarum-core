package validators

import (
	"errors"
	"fmt"
	"github.com/goclarum/clarum/core/arrays"
	clarumstrings "github.com/goclarum/clarum/core/validators/strings"
	"github.com/goclarum/clarum/http/message"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strings"
)

func ValidatePath(logPrefix string, expectedMessage *message.RequestMessage, actualUrl *url.URL) error {
	cleanedExpected := cleanPath(expectedMessage.Path)
	cleanedActual := cleanPath(actualUrl.Path)

	if cleanedExpected != cleanedActual {
		return handleError("%s: validation error - HTTP path mismatch - expected [%s] but received [%s]",
			logPrefix, cleanedExpected, cleanedActual)
	} else {
		slog.Info(fmt.Sprintf("%s: HTTP path validation successful", logPrefix))
	}

	return nil
}

func ValidateHttpMethod(logPrefix string, expectedMessage *message.RequestMessage, actualMethod string) error {
	if expectedMessage.Method != actualMethod {
		return handleError("%s: validation error - HTTP method mismatch - expected [%s] but received [%s]",
			logPrefix, expectedMessage.Method, actualMethod)
	} else {
		slog.Info(fmt.Sprintf("%s: HTTP method validation successful", logPrefix))
	}

	return nil
}

func ValidateHttpHeaders(logPrefix string, expectedMessage *message.Message, actualHeaders http.Header) error {
	if err := validateHeaders(expectedMessage, actualHeaders); err != nil {
		return handleError("%s: %s", logPrefix, err)
	} else {
		slog.Info(fmt.Sprintf("%s: header validation successful", logPrefix))
	}

	return nil
}

func validateHeaders(message *message.Message, headers http.Header) error {
	for header, expectedValue := range message.Headers {
		if receivedValue := headers.Get(header); expectedValue != receivedValue {
			return errors.New(fmt.Sprintf("validation error - header <%s> mismatch - expected [%s] but received [%s]",
				header, expectedValue, receivedValue))
		}
	}

	return nil
}

func ValidateHttpQueryParams(logPrefix string, expectedMessage *message.RequestMessage, actualUrl *url.URL) error {
	if err := validateQueryParams(expectedMessage, actualUrl.Query()); err != nil {
		return handleError("%s: %s", logPrefix, err)
	} else {
		slog.Info(fmt.Sprintf("%s: query params validation successful", logPrefix))
	}

	return nil
}

// validate query parameters based on these rules
//
//	-> validate that the param exists
//	-> that the value matches
func validateQueryParams(message *message.RequestMessage, params url.Values) error {
	for param, expectedValue := range message.QueryParams {
		if receivedValues, exists := params[param]; exists {
			if !arrays.Contains(receivedValues, expectedValue) {
				return errors.New(fmt.Sprintf("validation error - query param <%s> values mismatch - expected [%v] but received [%s]",
					param, expectedValue, receivedValues))
			}
		} else {
			return errors.New(fmt.Sprintf("validation error - query param <%s> missing", param))
		}
	}

	return nil
}

func ValidateHttpStatusCode(logPrefix string, expectedMessage *message.ResponseMessage, actualStatusCode int) error {
	if actualStatusCode != expectedMessage.StatusCode {
		return handleError("%s: validation error - HTTP status mismatch - expected [%d] but received [%d]", logPrefix, expectedMessage.StatusCode, actualStatusCode)
	} else {
		slog.Info(fmt.Sprintf("%s: HTTP status validation successful", logPrefix))
	}

	return nil
}

func ValidateHttpPayload(logPrefix string, expectedMessage *message.Message, actualPayload io.ReadCloser) error {
	defer closeBody(logPrefix, actualPayload)

	if clarumstrings.IsBlank(expectedMessage.MessagePayload) {
		slog.Info(fmt.Sprintf("%s: message payload is empty - no body validation will be done", logPrefix))
		return nil
	}

	bodyBytes, err := io.ReadAll(actualPayload)
	if err != nil {
		return handleError("%s: could not read response body - %s", logPrefix, err)
	}

	if err := validatePayload(expectedMessage, bodyBytes); err != nil {
		return handleError("%s: %s", logPrefix, err)
	} else {
		slog.Info(fmt.Sprintf("%s: payload validation successful", logPrefix))
	}

	return nil
}

func closeBody(logPrefix string, body io.ReadCloser) {
	if err := body.Close(); err != nil {
		slog.Error(fmt.Sprintf("%s: unable to close body - %s", logPrefix, err))
	}
}

func validatePayload(message *message.Message, payload []byte) error {
	receivedPayload := string(payload)

	if message.MessagePayload != receivedPayload { // plain text validation
		return errors.New(fmt.Sprintf("validation error - payload missmatch - expected [%s] but received [%s]",
			message.MessagePayload, receivedPayload))
	}

	return nil
}

func handleError(format string, a ...any) error {
	errorMessage := fmt.Sprintf(format, a...)
	slog.Error(errorMessage)
	return errors.New(errorMessage)
}

// path.Clean() does not remove leading "/", so we do that ourselves
func cleanPath(pathToClean string) string {
	return strings.TrimPrefix(path.Clean(pathToClean), "/")
}
