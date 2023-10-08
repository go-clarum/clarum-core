package http

import (
	"fmt"
	"maps"
	"net/http"
	"strconv"
)

func Get(pathElements ...string) *Action {
	return &Action{
		method: http.MethodGet,
		path:   BuildPath("", pathElements...),
	}
}

func Head(pathElements ...string) *Action {
	return &Action{
		method: http.MethodHead,
		path:   BuildPath("", pathElements...),
	}
}

func Post(pathElements ...string) *Action {
	return &Action{
		method: http.MethodPost,
		path:   BuildPath("", pathElements...),
	}
}

func Put(pathElements ...string) *Action {
	return &Action{
		method: http.MethodPut,
		path:   BuildPath("", pathElements...),
	}
}

func Delete(pathElements ...string) *Action {
	return &Action{
		method: http.MethodDelete,
		path:   BuildPath("", pathElements...),
	}
}

func Options(pathElements ...string) *Action {
	return &Action{
		method: http.MethodOptions,
		path:   BuildPath("", pathElements...),
	}
}

func Trace(pathElements ...string) *Action {
	return &Action{
		method: http.MethodTrace,
		path:   BuildPath("", pathElements...),
	}
}

func Patch(pathElements ...string) *Action {
	return &Action{
		method: http.MethodPatch,
		path:   BuildPath("", pathElements...),
	}
}

func Response(statusCode int) *Action {
	return &Action{
		statusCode: statusCode,
	}
}

type Action struct {
	method      string
	statusCode  int
	baseUrl     string
	path        string
	headers     map[string]string
	queryParams map[string]string
	payload     string
}

// BaseUrl - While this should normally be configured only on the HTTP client, this is also allowed on the action so that
// a client can send a request to different targets.
// When used on an action passed to an HTTP server, it will do nothing.
func (action *Action) BaseUrl(baseUrl string) *Action {
	action.baseUrl = baseUrl
	return action
}

func (action *Action) Header(key string, value string) *Action {
	if action.headers == nil {
		action.headers = make(map[string]string)
	}

	action.headers[key] = value
	return action
}

func (action *Action) ContentType(value string) *Action {
	return action.Header(ContentTypeHeaderName, value)
}

func (action *Action) Authorization(value string) *Action {
	return action.Header(AuthorizationHeaderName, value)
}

func (action *Action) ETag(value string) *Action {
	return action.Header(ETagHeaderName, value)
}

func (action *Action) QueryParam(key string, value string) *Action {
	if action.queryParams == nil {
		action.queryParams = make(map[string]string)
	}

	action.queryParams[key] = value
	return action
}

func (action *Action) Payload(payload string) *Action {
	action.payload = payload
	return action
}

func (action *Action) Clone() *Action {
	return &Action{
		method:      action.method,
		statusCode:  action.statusCode,
		baseUrl:     action.baseUrl,
		path:        action.path,
		headers:     maps.Clone(action.headers),
		queryParams: maps.Clone(action.queryParams),
		payload:     action.payload,
	}
}

func (action *Action) OverwriteWith(overwriting *Action) *Action {

	if len(overwriting.method) > 0 {
		action.method = overwriting.method
	}
	if overwriting.statusCode > 0 {
		action.statusCode = overwriting.statusCode
	}
	if len(overwriting.baseUrl) > 0 {
		action.baseUrl = overwriting.baseUrl
	}
	if len(overwriting.path) > 0 {
		action.path = overwriting.path
	}
	if len(overwriting.headers) > 0 {
		action.headers = overwriting.headers
	}
	if len(overwriting.queryParams) > 0 {
		action.queryParams = overwriting.queryParams
	}
	if len(overwriting.payload) > 0 {
		action.payload = overwriting.payload
	}

	return action
}

func (action *Action) ToString() string {
	statusCodeText := "none"
	if action.statusCode > 0 {
		statusCodeText = strconv.Itoa(action.statusCode)
	}

	return fmt.Sprintf(
		"["+
			"method: %s, "+
			"statusCode: %s, "+
			"baseUrl: %s, "+
			"path: '%s', "+
			"headers: %s, "+
			"queryParams: %s, "+
			"payload: %s"+
			"]",
		action.method, statusCodeText, action.baseUrl, action.path,
		action.headers, action.queryParams, action.payload)
}
