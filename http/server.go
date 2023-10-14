package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/goclarum/clarum/core/control"
	"io"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"
)

const contextNameKey = "endpointContext"

type ServerEndpoint struct {
	port           uint
	name           string
	contentType    string
	server         *http.Server
	context        *context.Context
	requestChannel chan *http.Request
	sendChannel    chan Action
}

type endpointContext struct {
	endpointName   string
	requestChannel chan *http.Request
	sendChannel    chan Action
}

func (se *ServerEndpoint) Receive(t *testing.T, action *Action) {
	logPrefix := serverLogPrefix(se.name)
	slog.Debug(fmt.Sprintf("%s: action to receive: %s", logPrefix, action.ToString()))
	actionToExecute := se.getActionToExecute(action)

	request := <-se.requestChannel
	slog.Debug(fmt.Sprintf("%s: executing validation action: %s", logPrefix, actionToExecute.ToString()))

	validateHttpHeaders(t, logPrefix, actionToExecute, request.Header)
	validateHttpQueryParams(t, logPrefix, actionToExecute, request.URL)
	validateHttpBody(t, logPrefix, actionToExecute, request.Body)
}

func (se *ServerEndpoint) Send(action *Action) {
	actionToExecute := se.getActionToExecute(action)
	se.sendChannel <- *actionToExecute
}

func (ce *ServerEndpoint) getActionToExecute(action *Action) *Action {
	actionToExecute := action.Clone()

	if len(actionToExecute.headers) == 0 || len(actionToExecute.headers[ContentTypeHeaderName]) == 0 {
		actionToExecute.ContentType(ce.contentType)
	}

	return actionToExecute
}

func (se *ServerEndpoint) start(ctx context.Context, cancelCtx context.CancelFunc, timeout time.Duration) {
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
	sendAction := <-ctx.sendChannel

	for header, value := range sendAction.headers {
		resWriter.Header().Set(header, value)
	}

	resWriter.WriteHeader(sendAction.statusCode)

	io.WriteString(resWriter, sendAction.payload)
	logOutgoingResponse(logPrefix, sendAction.statusCode, sendAction.payload, resWriter)
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
