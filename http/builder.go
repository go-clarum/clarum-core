package http

import (
	"github.com/goclarum/clarum/http/client"
	"github.com/goclarum/clarum/http/server"
)

type EndpointBuilder struct {
}

func (heb *EndpointBuilder) Client() *client.EndpointBuilder {
	return client.NewEndpointBuilder()
}

func (heb *EndpointBuilder) Server() *server.EndpointBuilder {
	return server.NewEndpointBuilder()
}
