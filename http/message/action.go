package message

import (
	"fmt"
	"github.com/goclarum/clarum/http/constants"
	"github.com/goclarum/clarum/http/internal/utils"
	"maps"
	"net/http"
	"strconv"
)

type Action struct {
	Method         string
	StatusCode     int
	Url            string
	Path           string
	Headers        map[string]string
	QueryParams    map[string]string
	MessagePayload string
}

func Get(pathElements ...string) *Action {
	return &Action{
		Method: http.MethodGet,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Head(pathElements ...string) *Action {
	return &Action{
		Method: http.MethodHead,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Post(pathElements ...string) *Action {
	return &Action{
		Method: http.MethodPost,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Put(pathElements ...string) *Action {
	return &Action{
		Method: http.MethodPut,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Delete(pathElements ...string) *Action {
	return &Action{
		Method: http.MethodDelete,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Options(pathElements ...string) *Action {
	return &Action{
		Method: http.MethodOptions,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Trace(pathElements ...string) *Action {
	return &Action{
		Method: http.MethodTrace,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Patch(pathElements ...string) *Action {
	return &Action{
		Method: http.MethodPatch,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Response(statusCode int) *Action {
	return &Action{
		StatusCode: statusCode,
	}
}

// BaseUrl - While this should normally be configured only on the HTTP client, this is also allowed on the action so that
// a client can send a request to different targets.
// When used on an action passed to an HTTP server, it will do nothing.
func (action *Action) BaseUrl(baseUrl string) *Action {
	action.Url = baseUrl
	return action
}

func (action *Action) Header(key string, value string) *Action {
	if action.Headers == nil {
		action.Headers = make(map[string]string)
	}

	action.Headers[key] = value
	return action
}

func (action *Action) ContentType(value string) *Action {
	return action.Header(constants.ContentTypeHeaderName, value)
}

func (action *Action) Authorization(value string) *Action {
	return action.Header(constants.AuthorizationHeaderName, value)
}

func (action *Action) ETag(value string) *Action {
	return action.Header(constants.ETagHeaderName, value)
}

func (action *Action) QueryParam(key string, value string) *Action {
	if action.QueryParams == nil {
		action.QueryParams = make(map[string]string)
	}

	action.QueryParams[key] = value
	return action
}

func (action *Action) Payload(payload string) *Action {
	action.MessagePayload = payload
	return action
}

func (action *Action) Clone() *Action {
	return &Action{
		Method:         action.Method,
		StatusCode:     action.StatusCode,
		Url:            action.Url,
		Path:           action.Path,
		Headers:        maps.Clone(action.Headers),
		QueryParams:    maps.Clone(action.QueryParams),
		MessagePayload: action.MessagePayload,
	}
}

func (action *Action) OverwriteWith(overwriting *Action) *Action {

	if len(overwriting.Method) > 0 {
		action.Method = overwriting.Method
	}
	if overwriting.StatusCode > 0 {
		action.StatusCode = overwriting.StatusCode
	}
	if len(overwriting.Url) > 0 {
		action.Url = overwriting.Url
	}
	if len(overwriting.Path) > 0 {
		action.Path = overwriting.Path
	}
	if len(overwriting.Headers) > 0 {
		action.Headers = overwriting.Headers
	}
	if len(overwriting.QueryParams) > 0 {
		action.QueryParams = overwriting.QueryParams
	}
	if len(overwriting.MessagePayload) > 0 {
		action.MessagePayload = overwriting.MessagePayload
	}

	return action
}

func (action *Action) ToString() string {
	statusCodeText := "none"
	if action.StatusCode > 0 {
		statusCodeText = strconv.Itoa(action.StatusCode)
	}

	return fmt.Sprintf(
		"["+
			"Method: %s, "+
			"StatusCode: %s, "+
			"Url: %s, "+
			"Path: '%s', "+
			"Headers: %s, "+
			"QueryParams: %s, "+
			"MessagePayload: %s"+
			"]",
		action.Method, statusCodeText, action.Url, action.Path,
		action.Headers, action.QueryParams, action.MessagePayload)
}
