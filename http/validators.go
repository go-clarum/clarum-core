package http

import (
	"errors"
	"fmt"
	"github.com/goclarum/clarum/core/arrays"
	"net/http"
	"net/url"
	"strings"
)

func validateHeaders(action *Action, headers http.Header) error {
	for header, expectedValue := range action.headers {
		if receivedValue := headers.Get(header); expectedValue != receivedValue {
			return errors.New(fmt.Sprintf("Validation error: header <%s> mismatch. Expected [%s] but received [%s]",
				header, expectedValue, receivedValue))
		}
	}

	return nil
}

// validate query parameters based on these rules
//
//	-> validate that the param exists
//	-> that the value matches
func validateQueryParams(action *Action, params url.Values) error {
	for param, expectedValue := range action.queryParams {
		if receivedValues, exists := params[param]; exists {
			if !arrays.Contains(receivedValues, expectedValue) {
				return errors.New(fmt.Sprintf("Validation error: query params mismatch. Expected [%v] but received [%s]",
					expectedValue, receivedValues))
			}
		} else {
			return errors.New(fmt.Sprintf("Validation error: query param <%s> missing", param))
		}
	}

	return nil
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
