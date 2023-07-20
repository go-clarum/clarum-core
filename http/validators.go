package http

import (
	"errors"
	"fmt"
	"github.com/goclarum/clarum/core/arrays"
	"net/http"
	"net/url"
)

func validateHeaders(action *Action, headers http.Header) error {
	for header, value := range action.headers {
		if value != headers.Get(header) {
			return errors.New(fmt.Sprintf("validation error: header <%s> mismatch", header))
		}
	}

	return nil
}

// validate query parameters based on these rules
//
//	-> validate that the param exists
//	-> that the value matches
func validateQueryParams(action *Action, params url.Values) error {
	for param, value := range action.queryParams {
		if paramValues, exists := params[param]; exists {
			if !arrays.Contains(paramValues, value) {
				return errors.New(fmt.Sprintf("validation error: query params mismatch: expected %v, actual %s", paramValues, value))
			}
		} else {
			return errors.New(fmt.Sprintf("validation error: query param <%s> missing", param))
		}
	}

	return nil
}
