package message

import (
	"fmt"
	"github.com/goclarum/clarum/http/constants"
	"github.com/goclarum/clarum/http/internal/utils"
	"maps"
	"net/http"
	"strconv"
)

type Message struct {
	Method         string
	StatusCode     int
	Url            string
	Path           string
	Headers        map[string]string
	QueryParams    map[string]string
	MessagePayload string
}

func Get(pathElements ...string) *Message {
	return &Message{
		Method: http.MethodGet,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Head(pathElements ...string) *Message {
	return &Message{
		Method: http.MethodHead,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Post(pathElements ...string) *Message {
	return &Message{
		Method: http.MethodPost,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Put(pathElements ...string) *Message {
	return &Message{
		Method: http.MethodPut,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Delete(pathElements ...string) *Message {
	return &Message{
		Method: http.MethodDelete,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Options(pathElements ...string) *Message {
	return &Message{
		Method: http.MethodOptions,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Trace(pathElements ...string) *Message {
	return &Message{
		Method: http.MethodTrace,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Patch(pathElements ...string) *Message {
	return &Message{
		Method: http.MethodPatch,
		Path:   utils.BuildPath("", pathElements...),
	}
}

func Response(statusCode int) *Message {
	return &Message{
		StatusCode: statusCode,
	}
}

// BaseUrl - While this should normally be configured only on the HTTP client,
// this is also allowed on the message so that a client can send a request to different targets.
// When used on a message passed to an HTTP server, it will do nothing.
func (message *Message) BaseUrl(baseUrl string) *Message {
	message.Url = baseUrl
	return message
}

func (message *Message) Header(key string, value string) *Message {
	if message.Headers == nil {
		message.Headers = make(map[string]string)
	}

	message.Headers[key] = value
	return message
}

func (message *Message) ContentType(value string) *Message {
	return message.Header(constants.ContentTypeHeaderName, value)
}

func (message *Message) Authorization(value string) *Message {
	return message.Header(constants.AuthorizationHeaderName, value)
}

func (message *Message) ETag(value string) *Message {
	return message.Header(constants.ETagHeaderName, value)
}

func (message *Message) QueryParam(key string, value string) *Message {
	if message.QueryParams == nil {
		message.QueryParams = make(map[string]string)
	}

	message.QueryParams[key] = value
	return message
}

func (message *Message) Payload(payload string) *Message {
	message.MessagePayload = payload
	return message
}

func (message *Message) Clone() *Message {
	return &Message{
		Method:         message.Method,
		StatusCode:     message.StatusCode,
		Url:            message.Url,
		Path:           message.Path,
		Headers:        maps.Clone(message.Headers),
		QueryParams:    maps.Clone(message.QueryParams),
		MessagePayload: message.MessagePayload,
	}
}

func (message *Message) OverwriteWith(overwriting *Message) *Message {

	if len(overwriting.Method) > 0 {
		message.Method = overwriting.Method
	}
	if overwriting.StatusCode > 0 {
		message.StatusCode = overwriting.StatusCode
	}
	if len(overwriting.Url) > 0 {
		message.Url = overwriting.Url
	}
	if len(overwriting.Path) > 0 {
		message.Path = overwriting.Path
	}
	if len(overwriting.Headers) > 0 {
		message.Headers = overwriting.Headers
	}
	if len(overwriting.QueryParams) > 0 {
		message.QueryParams = overwriting.QueryParams
	}
	if len(overwriting.MessagePayload) > 0 {
		message.MessagePayload = overwriting.MessagePayload
	}

	return message
}

func (message *Message) ToString() string {
	statusCodeText := "none"
	if message.StatusCode > 0 {
		statusCodeText = strconv.Itoa(message.StatusCode)
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
		message.Method, statusCodeText, message.Url, message.Path,
		message.Headers, message.QueryParams, message.MessagePayload)
}
