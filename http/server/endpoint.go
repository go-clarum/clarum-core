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
	"testing"
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
func (se *Endpoint) receive(t *testing.T, message *message.Message) {
	logPrefix := serverLogPrefix(se.name)
	slog.Debug(fmt.Sprintf("%s: message to receive: %s", logPrefix, message.ToString()))
	messageToReceive := se.getFinalMessage(message)

	request := <-se.requestChannel
	slog.Debug(fmt.Sprintf("%s: validating message: %s", logPrefix, messageToReceive.ToString()))

	validators.ValidateHttpHeaders(t, logPrefix, messageToReceive, request.Header)
	validators.ValidateHttpQueryParams(t, logPrefix, messageToReceive, request.URL)
	validators.ValidateHttpBody(t, logPrefix, messageToReceive, request.Body)
}

func (se *Endpoint) send(message *message.Message) {
	messageToSend := se.getFinalMessage(message)
	// can we refactor this to send the response instead of the message?
	se.sendChannel <- *messageToSend
}

func (ce *Endpoint) getFinalMessage(message *message.Message) *message.Message {
	finalMessage := message.Clone()

	if len(finalMessage.Headers) == 0 || len(finalMessage.Headers[constants.ContentTypeHeaderName]) == 0 {
		finalMessage.ContentType(ce.contentType)
	}

	return finalMessage
}

func (se *Endpoint) start(ctx context.Context, cancelCtx context.CancelFunc, timeout time.Duration) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", requestHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", se.port),
		Handler:      mux,
		WriteTimeout: timeout,
		BaseContext: func(l net.Listener) context.Context {
			endpointContext := &endpointContext{
				endpointName:   se.name,
				requestChannel: se.requestChannel,
				sendChannel:    se.sendChannel,
			}
			ctx = context.WithValue(ctx, contextNameKey, endpointContext)
			return ctx
		},
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				fmt.Println(fmt.Sprintf("%s: closed", serverLogPrefix(se.name)))
			} else {
				fmt.Println(fmt.Sprintf("%s: error: %s", serverLogPrefix(se.name), err))
			}
		} else {
			fmt.Println(fmt.Sprintf("%s: closed: %s", serverLogPrefix(se.name), err))
		}

		cancelCtx()
	}()

	se.server = server
}

// The requestHandler is started on request in, reports the request so that it can be validated
// after which is blocked until the send method prepares a response. This way we can tell it
// inside the test, when to send the response.
// The handler blocks until a timeout is triggered // TODO: check how timeouts are handled
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
		slog.Error(fmt.Sprintf("%s: could not write response body: %s", logPrefix, err))
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
		slog.Error(fmt.Sprintf("%s: could not read request body: %s", logPrefix, err))
	} else {
		request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		bodyString = string(bodyBytes)
	}

	slog.Info(fmt.Sprintf("%s: received request: ["+
		"method: %s, "+
		"url: %s, "+
		"headers: %s, "+
		"payload: %s"+
		"]",
		logPrefix, request.Method, request.URL.String(), request.Header, bodyString))
}

func logOutgoingResponse(prefix string, statusCode int, payload string, res http.ResponseWriter) {
	slog.Info(fmt.Sprintf("%s: sending response: ["+
		"status: %d, "+
		"headers: %s, "+
		"payload: %s"+
		"]",
		prefix, statusCode, res.Header(), payload))
}

func serverLogPrefix(endpointName string) string {
	return fmt.Sprintf("HTTP server %s", endpointName)
}
