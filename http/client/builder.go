package client

import (
	"time"
)

type EndpointBuilder struct {
	baseUrl     string
	contentType string
	name        string
	timeout     time.Duration
}

func NewEndpointBuilder() *EndpointBuilder {
	return &EndpointBuilder{}
}

func (ceb *EndpointBuilder) Name(name string) *EndpointBuilder {
	ceb.name = name
	return ceb
}

func (ceb *EndpointBuilder) BaseUrl(baseUrl string) *EndpointBuilder {
	ceb.baseUrl = baseUrl
	return ceb
}

func (ceb *EndpointBuilder) ContentType(contentType string) *EndpointBuilder {
	ceb.contentType = contentType
	return ceb
}

func (ceb *EndpointBuilder) Build() *Endpoint {
	return NewEndpoint(ceb.name, ceb.baseUrl, ceb.contentType, ceb.timeout)
}
