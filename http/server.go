package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/goclarum/clarum/core/control"
	"io"
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
	actionToExecute := se.getActionToExecute(action)

	request := <-se.requestChannel
	// debug logging - log entire request as is
	fmt.Println(fmt.Sprintf("HTTP server <%s> received request: %s", se.name, request.Method))

	if err := validateHeaders(actionToExecute, request.Header); err != nil {
		t.Errorf("HTTP server <%s>: %s", se.name, err)
	} else {
		// debug logging
		fmt.Println(fmt.Sprintf("HTTP server <%s> header validation successful", se.name))
	}

	if err := validateQueryParams(actionToExecute, request.URL.Query()); err != nil {
		t.Errorf("HTTP server <%s>: %s", se.name, err)
	} else {
		// debug logging
		fmt.Println(fmt.Sprintf("HTTP server <%s> query params validation successful", se.name))
	}
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
				fmt.Println(fmt.Sprintf("HTTP server <%s> closed", se.name))
			} else {
				fmt.Println(fmt.Sprintf("HTTP server <%s> error: %s", se.name, err))
			}
		} else {
			fmt.Println(fmt.Sprintf("HTTP server <%s> closed: %s", se.name, err))
		}

		cancelCtx()
	}()

	se.server = server
}

// The requestHandler is started on request in, reports the request so that it can be validated
// after which is blocked until the send method prepares a response. This way we can tell it
// inside the test, when to send the response.
// The handler blocks until a timeout is triggered // TODO: check how timeouts are handled
func requestHandler(resWriter http.ResponseWriter, req *http.Request) {
	control.RunningActions.Add(1)
	defer control.RunningActions.Done()

	ctx := req.Context().Value(contextNameKey).(*endpointContext)
	ctx.requestChannel <- req
	sendAction := <-ctx.sendChannel

	for header, value := range sendAction.headers {
		resWriter.Header().Set(header, value)
	}

	resWriter.WriteHeader(sendAction.statusCode)
	io.WriteString(resWriter, fmt.Sprintf("Hello, from server"))
}
