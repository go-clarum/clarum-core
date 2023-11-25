package server

import (
	"github.com/goclarum/clarum/http/message"
)

// SendActionBuilder used to configure a send action on a client endpoint without the context of a test
// the method chain will end with the .Message() method which will return an error.
// The error will be a problem encountered during sending.
type SendActionBuilder struct {
	endpoint *Endpoint
}

func (builder *SendActionBuilder) Message(message *message.ResponseMessage) {
	builder.endpoint.send(message)
}
