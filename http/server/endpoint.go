package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/goclarum/clarum/core/control"
	"github.com/goclarum/clarum/http/constants"
	"github.com/goclarum/clarum/http/internal/validators"
	"github.com/goclarum/clarum/http/message"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"
)

const contextNameKey = "endpointContext"

type Endpoint struct {
	name           string
	port           uint
	contentType    string
	server         *http.Server
	context        *context.Context
	requestChannel chan *http.Request
	sendChannel    chan message.Message
}

type endpointContext struct {
	endpointName   string
	requestChannel chan *http.Request
	sendChannel    chan message.Message
}

func NewServerEndpoint(name string, port uint, contentType string, timeout time.Duration) *Endpoint {
	ctx, cancelCtx := context.WithCancel(context.Background())
	sendChannel := make(chan message.Message)
	requestChannel := make(chan *http.Request)

	se := &Endpoint{
		name:           name,
		port:           port,
		contentType:    contentType,
		context:        &ctx,
		sendChannel:    sendChannel,
		requestChannel: requestChannel,
	}

	// feature: start automatically = true/false; to simulate connection errors
	se.start(ctx, cancelCtx, timeout)

	return se
}

// this Method is blocking, until a request is received
func (endpoint *Endpoint) receive(message *message.Message) error {
	logPrefix := serverLogPrefix(endpoint.name)
	slog.Debug(fmt.Sprintf("%s: message to receive %s", logPrefix, message.ToString()))
	messageToReceive := endpoint.getFinalMessage(message)

	request := <-endpoint.requestChannel
	slog.Debug(fmt.Sprintf("%s: validation message %s", logPrefix, messageToReceive.ToString()))

	return errors.Join(
		validators.ValidateHttpHeaders(logPrefix, messageToReceive, request.Header),
		validators.ValidateHttpQueryParams(logPrefix, messageToReceive, request.URL),
		validators.ValidateHttpBody(logPrefix, messageToReceive, request.Body))
}

func (endpoint *Endpoint) send(message *message.Message) {
	messageToSend := endpoint.getFinalMessage(message)
	endpoint.sendChannel <- *messageToSend
}

func (endpoint *Endpoint) getFinalMessage(message *message.Message) *message.Message {
	finalMessage := message.Clone()

	if len(finalMessage.Headers) == 0 || len(finalMessage.Headers[constants.ContentTypeHeaderName]) == 0 {
		finalMessage.ContentType(endpoint.contentType)
	}

	return finalMessage
}

func (endpoint *Endpoint) start(ctx context.Context, cancelCtx context.CancelFunc, timeout time.Duration) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", requestHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", endpoint.port),
		Handler:      mux,
		WriteTimeout: timeout,
		BaseContext: func(l net.Listener) context.Context {
			endpointContext := &endpointContext{
				endpointName:   endpoint.name,
				requestChannel: endpoint.requestChannel,
				sendChannel:    endpoint.sendChannel,
			}
			ctx = context.WithValue(ctx, contextNameKey, endpointContext)
			return ctx
		},
	}

	go func() {
		logPrefix := serverLogPrefix(endpoint.name)
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				fmt.Println(fmt.Sprintf("%s: closed", logPrefix))
			} else {
				fmt.Println(fmt.Sprintf("%s: error - %s", logPrefix, err))
			}
		} else {
			fmt.Println(fmt.Sprintf("%s: closed - %s", logPrefix, err))
		}

		cancelCtx()
	}()

	endpoint.server = server
}

// The requestHandler is started when the server receives a request.
// The request is sent to the requestChannel to be picked up by a test action (validation).
// After sending the request to the channel, the handler is blocked until the send() test action
// provides a response message. This way we can control, inside the test, when a response will be sent.
// The handler blocks until a timeout is triggered
// TODO: check how timeouts are handled
func requestHandler(resWriter http.ResponseWriter, request *http.Request) {
	control.RunningActions.Add(1)
	defer control.RunningActions.Done()

	ctx := request.Context().Value(contextNameKey).(*endpointContext)

	logPrefix := serverLogPrefix(ctx.endpointName)
	logIncomingRequest(logPrefix, request)
	ctx.requestChannel <- request
	messageToSend := <-ctx.sendChannel

	for header, value := range messageToSend.Headers {
		resWriter.Header().Set(header, value)
	}

	resWriter.WriteHeader(messageToSend.StatusCode)

	_, err := io.WriteString(resWriter, messageToSend.MessagePayload)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: could not write response body - %s", logPrefix, err))
	}
	logOutgoingResponse(logPrefix, messageToSend.StatusCode, messageToSend.MessagePayload, resWriter)
}

// we read the body 'as is' for logging, after which we put it back into the request
// with an open reader so that it can be read downstream again
func logIncomingRequest(logPrefix string, request *http.Request) {
	bodyBytes, _ := io.ReadAll(request.Body)
	bodyString := ""

	err := request.Body.Close()
	if err != nil {
		slog.Error(fmt.Sprintf("%s: could not read request body - %s", logPrefix, err))
	} else {
		request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		bodyString = string(bodyBytes)
	}

	slog.Info(fmt.Sprintf("%s: received request ["+
		"method: %s, "+
		"url: %s, "+
		"headers: %s, "+
		"payload: %s"+
		"]",
		logPrefix, request.Method, request.URL.String(), request.Header, bodyString))
}

func logOutgoingResponse(prefix string, statusCode int, payload string, res http.ResponseWriter) {
	slog.Info(fmt.Sprintf("%s: sending response ["+
		"status: %d, "+
		"headers: %s, "+
		"payload: %s"+
		"]",
		prefix, statusCode, res.Header(), payload))
}

func serverLogPrefix(endpointName string) string {
	return fmt.Sprintf("HTTP server %s", endpointName)
}
