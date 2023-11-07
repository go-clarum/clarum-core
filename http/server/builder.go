package server

import (
	"time"
)

func NewEndpointBuilder() *EndpointBuilder {
	return &EndpointBuilder{}
}

type EndpointBuilder struct {
	contentType string
	port        uint
	name        string
	timeout     time.Duration
}

func (seb *EndpointBuilder) Timeout(timeout time.Duration) *EndpointBuilder {
	seb.timeout = timeout
	return seb
}

func (seb *EndpointBuilder) Name(name string) *EndpointBuilder {
	seb.name = name
	return seb
}

func (seb *EndpointBuilder) Port(port uint) *EndpointBuilder {
	seb.port = port
	return seb
}

func (seb *EndpointBuilder) ContentType(contentType string) *EndpointBuilder {
	seb.contentType = contentType
	return seb
}

func (seb *EndpointBuilder) Build() *Endpoint {
	return NewServerEndpoint(seb.name, seb.port, seb.contentType, seb.timeout)
}
