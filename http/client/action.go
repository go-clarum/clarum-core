package client

import (
	"github.com/goclarum/clarum/http/message"
	"testing"
)

// TODO
type ActionContext struct {
	test *testing.T
	ActionBuilder
}

// TODO
type ActionBuilder struct {
	endpoint *Endpoint
	send     bool
	message  message.Message
}

func (endpoint *Endpoint) In(t *testing.T) *ActionContext {
	return &ActionContext{
		test: t,
		ActionBuilder: ActionBuilder{
			endpoint: endpoint,
		},
	}
}

func (builder *ActionContext) Send() *ActionContext {
	builder.send = true
	return builder
}

func (endpoint *Endpoint) Send() *ActionBuilder {
	return &ActionBuilder{
		endpoint: endpoint,
		send:     true,
	}
}

func (builder *ActionContext) Receive() *ActionContext {
	builder.send = false
	return builder
}

func (endpoint *Endpoint) Receive() *ActionBuilder {
	return &ActionBuilder{
		endpoint: endpoint,
		send:     false,
	}
}

func (context *ActionContext) Message(message *message.Message) {

}

// finalise action execution; return error if test property is undefined
func (context *ActionBuilder) Message(message *message.Message) error {
	return nil
}
