package message

import (
	"fmt"
	"strconv"
)

type ResponseMessage struct {
	Message
	StatusCode int
}

func Response(statusCode int) *ResponseMessage {
	return &ResponseMessage{
		StatusCode: statusCode,
	}
}

func (response *ResponseMessage) Header(key string, value string) *ResponseMessage {
	response.Message.header(key, value)
	return response
}

func (response *ResponseMessage) ContentType(value string) *ResponseMessage {
	response.Message.contentType(value)
	return response
}

func (response *ResponseMessage) ETag(value string) *ResponseMessage {
	response.Message.eTag(value)
	return response
}

func (response *ResponseMessage) Payload(payload string) *ResponseMessage {
	response.Message.MessagePayload = payload
	return response
}

func (response *ResponseMessage) Clone() *ResponseMessage {
	return &ResponseMessage{
		StatusCode: response.StatusCode,
		Message:    response.Message.clone(),
	}
}

func (response *ResponseMessage) ToString() string {
	statusCodeText := "none"
	if response.StatusCode > 0 {
		statusCodeText = strconv.Itoa(response.StatusCode)
	}

	return fmt.Sprintf(
		"["+
			"StatusCode: %s, "+
			"Headers: %s, "+
			"MessagePayload: %s"+
			"]",
		statusCodeText, response.Headers, response.MessagePayload)
}
