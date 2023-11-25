package client

import (
	"github.com/goclarum/clarum/http/message"
	"testing"
)

// ReceiveActionBuilder used to configure a receive action on a client endpoint without the context of a test
// the method chain will end with the .Message() method which will return an error.
// The error will be a problem encountered during receiving or a validation error.
type ReceiveActionBuilder struct {
	endpoint *Endpoint
}

// TestReceiveActionBuilder used to configure a receive action on a client endpoint with the context of a test
// the method chain will end with the .Message() method which will not return anything.
// Any error encountered during receiving or validating will fail the test by calling t.Error().
type TestReceiveActionBuilder struct {
	test *testing.T
	ReceiveActionBuilder
}

func (testBuilder *TestReceiveActionBuilder) Message(message *message.ResponseMessage) {
	if err := testBuilder.endpoint.receive(message); err != nil {
		testBuilder.test.Error(err)
	}
}

func (builder *ReceiveActionBuilder) Message(message *message.ResponseMessage) error {
	return builder.endpoint.receive(message)
}
