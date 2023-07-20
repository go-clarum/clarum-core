package http

import (
	"context"
	"net/http"
	"time"
)

func Client() *ClientEndpointBuilder {
	return &ClientEndpointBuilder{}
}

type ClientEndpointBuilder struct {
	baseUrl     string
	contentType string
	name        string
	timeout     time.Duration
}

func (ceb *ClientEndpointBuilder) Name(name string) *ClientEndpointBuilder {
	ceb.name = name
	return ceb
}

func (ceb *ClientEndpointBuilder) BaseUrl(baseUrl string) *ClientEndpointBuilder {
	ceb.baseUrl = baseUrl
	return ceb
}

func (ceb *ClientEndpointBuilder) ContentType(contentType string) *ClientEndpointBuilder {
	ceb.contentType = contentType
	return ceb
}

func (ceb *ClientEndpointBuilder) Build() *ClientEndpoint {
	client := http.Client{
		Timeout: ceb.timeout,
	}

	return &ClientEndpoint{
		name:            ceb.name,
		baseUrl:         ceb.baseUrl,
		contentType:     ceb.contentType,
		client:          &client,
		responseChannel: make(chan *http.Response),
	}
}

func Server() *ServerEndpointBuilder {
	return &ServerEndpointBuilder{}
}

type ServerEndpointBuilder struct {
	contentType string
	port        uint
	name        string
	timeout     time.Duration
}

func (seb *ServerEndpointBuilder) Timeout(timeout time.Duration) *ServerEndpointBuilder {
	seb.timeout = timeout
	return seb
}

func (seb *ServerEndpointBuilder) Name(name string) *ServerEndpointBuilder {
	seb.name = name
	return seb
}

func (seb *ServerEndpointBuilder) Port(port uint) *ServerEndpointBuilder {
	seb.port = port
	return seb
}

func (seb *ServerEndpointBuilder) ContentType(contentType string) *ServerEndpointBuilder {
	seb.contentType = contentType
	return seb
}

func (seb *ServerEndpointBuilder) Build() *ServerEndpoint {
	ctx, cancelCtx := context.WithCancel(context.Background())
	sendChannel := make(chan Action)
	requestChannel := make(chan *http.Request)

	se := &ServerEndpoint{
		port:           seb.port,
		name:           seb.name,
		contentType:    seb.contentType,
		context:        &ctx,
		sendChannel:    sendChannel,
		requestChannel: requestChannel,
	}

	// feature: start automatically = true/false; to simulate connection errors
	se.start(ctx, cancelCtx, seb.timeout)

	return se
}
