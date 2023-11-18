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
)

func ValidateHttpMethod(logPrefix string, message *message.Message, receivedMethod string) error {
	if message.Method != receivedMethod {
		return handleError("%s: validation error - HTTP method mismatch - expected [%s] but received [%s]",
			logPrefix, message.Method, receivedMethod)
	} else {
		slog.Info(fmt.Sprintf("%s: HTTP method validation successful", logPrefix))
	}

	return nil
}

func ValidateHttpHeaders(logPrefix string, message *message.Message, headers http.Header) error {
	if err := validateHeaders(message, headers); err != nil {
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

func ValidateHttpQueryParams(logPrefix string, message *message.Message, url *url.URL) error {
	if err := validateQueryParams(message, url.Query()); err != nil {
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
func validateQueryParams(message *message.Message, params url.Values) error {
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

func ValidateHttpStatusCode(logPrefix string, message *message.Message, statusCode int) error {
	if statusCode != message.StatusCode {
		return handleError("%s: validation error - HTTP status mismatch - expected [%d] but received [%d]", logPrefix, message.StatusCode, statusCode)
	} else {
		slog.Info(fmt.Sprintf("%s: HTTP status validation successful", logPrefix))
	}

	return nil
}

func ValidateHttpBody(logPrefix string, message *message.Message, body io.ReadCloser) error {
	defer closeBody(logPrefix, body)

	if clarumstrings.IsBlank(message.MessagePayload) {
		slog.Info(fmt.Sprintf("%s: message payload is empty - no body validation will be done", logPrefix))
		return nil
	}

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return handleError("%s: could not read response body - %s", logPrefix, err)
	}

	if err := validatePayload(message, bodyBytes); err != nil {
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
