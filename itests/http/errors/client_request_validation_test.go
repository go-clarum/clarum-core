package errors

import (
	"github.com/goclarum/clarum/http/message"
	"testing"
)

// The following tests check client send request validation errors.

func TestClientSendNilMessage(t *testing.T) {
	expectedError := "HTTP client errorsClient: message to send is nil"

	e1 := errorsClient.Send().Message(nil)

	checkErrors(t, expectedError, e1)
}

func TestClientSendInvalidMessageUrl(t *testing.T) {
	expectedError := "HTTP client errorsClient: message to send is invalid - missing url"

	e1 := errorsClient.Send().Message(message.Get())

	checkErrors(t, expectedError, e1)
}

func TestClientSendInvalidMessageMethod(t *testing.T) {
	expectedError := "HTTP client errorsClient: message to send is invalid - missing HTTP method"

	request := &message.Message{
		Url: "something",
	}
	e1 := errorsClient.Send().Message(request)

	checkErrors(t, expectedError, e1)
}
