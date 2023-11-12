package server

import (
	"testing"
)

// TestSendActionBuilder used to initiate a send or receive action on a client endpoint
// with the context of a test
type TestActionBuilder struct {
	test     *testing.T
	endpoint *Endpoint
}

func (endpoint *Endpoint) In(t *testing.T) *TestActionBuilder {
	return &TestActionBuilder{
		test:     t,
		endpoint: endpoint,
	}
}

func (endpoint *Endpoint) Send() *SendActionBuilder {
	return &SendActionBuilder{
		endpoint: endpoint,
	}
}

func (endpoint *Endpoint) Receive() *ReceiveActionBuilder {
	return &ReceiveActionBuilder{
		endpoint: endpoint,
	}
}

func (testBuilder *TestActionBuilder) Receive() *TestReceiveActionBuilder {
	return &TestReceiveActionBuilder{
		test: testBuilder.test,
		ReceiveActionBuilder: ReceiveActionBuilder{
			endpoint: testBuilder.endpoint,
		},
	}
}
