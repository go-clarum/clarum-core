package client

import (
	"bytes"
	"fmt"
	"github.com/goclarum/clarum/core/control"
	"github.com/goclarum/clarum/http/constants"
	"github.com/goclarum/clarum/http/internal/utils"
	"github.com/goclarum/clarum/http/internal/validators"
	"github.com/goclarum/clarum/http/message"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"
)

type Endpoint struct {
	name            string
	baseUrl         string
	contentType     string
	client          *http.Client
	responseChannel chan *http.Response
}

func NewEndpoint(name string, baseUrl string, contentType string, timeout time.Duration) *Endpoint {
	client := http.Client{
		Timeout: timeout,
	}

	return &Endpoint{
		name:            name,
		baseUrl:         baseUrl,
		contentType:     contentType,
		client:          &client,
		responseChannel: make(chan *http.Response),
	}
}

func (ce *Endpoint) send(t *testing.T, action *message.Action) {
	logPrefix := clientLogPrefix(ce.name)
	slog.Debug(fmt.Sprintf("%s: action to send: %s", logPrefix, action.ToString()))
	control.RunningActions.Add(1)

	go func() {
		defer control.RunningActions.Done()

		// from here
		actionToExecute := ce.getSendActionToExecute(action)
		slog.Debug(fmt.Sprintf("%s: executing action: %s", logPrefix, actionToExecute.ToString()))

		req, err := buildRequest(ce.name, actionToExecute)
		if err != nil {
			t.Errorf("%s: canceled action. Error: %s", logPrefix, err)
		}
		//to here, move outside the goroutine and return error if t is nil

		logOutgoingRequest(logPrefix, action.MessagePayload, req)
		res, err := ce.client.Do(req)
		logIncomingResponse(logPrefix, res)

		// TODO: handle technical errors
		//  check: socket connection, connection refused, connection timeout
		if err != nil {
			t.Errorf("%s: error on response: %s", logPrefix, err)
		}

		ce.responseChannel <- res
	}()
}

func (ce *Endpoint) receive(t *testing.T, action *message.Action) {
	logPrefix := clientLogPrefix(ce.name)
	slog.Debug(fmt.Sprintf("%s: action to receive: %s", logPrefix, action.ToString()))

	response := <-ce.responseChannel

	actionToExecute := ce.getReceiveActionToExecute(action)
	slog.Debug(fmt.Sprintf("%s: executing validation action: %s", logPrefix, actionToExecute.ToString()))

	validators.ValidateHttpStatusCode(t, logPrefix, actionToExecute, response.StatusCode)
	validators.ValidateHttpHeaders(t, logPrefix, actionToExecute, response.Header)
	validators.ValidateHttpBody(t, logPrefix, actionToExecute, response.Body)
}

// put missing data into a send action
func (ce *Endpoint) getSendActionToExecute(action *message.Action) *message.Action {
	actionToExecute := action.Clone()

	if len(actionToExecute.Url) == 0 {
		actionToExecute.Url = ce.baseUrl
	}
	if len(actionToExecute.Headers) == 0 || len(actionToExecute.Headers[constants.ContentTypeHeaderName]) == 0 {
		actionToExecute.ContentType(ce.contentType)
	}

	return actionToExecute
}

// put missing data into a receive action
func (ce *Endpoint) getReceiveActionToExecute(action *message.Action) *message.Action {
	actionToExecute := action.Clone()

	if len(actionToExecute.Headers) == 0 || len(actionToExecute.Headers[constants.ContentTypeHeaderName]) == 0 {
		actionToExecute.ContentType(ce.contentType)
	}

	return actionToExecute
}

func buildRequest(prefix string, action *message.Action) (*http.Request, error) {
	url := utils.BuildPath(action.Url, action.Path)

	req, err := http.NewRequest(action.Method, url, bytes.NewBufferString(action.MessagePayload))
	if err != nil {
		slog.Error(fmt.Sprintf("%s: error: %s", prefix, err))
		return nil, err
	}

	for header, value := range action.Headers {
		req.Header.Set(header, value)
	}

	qParams := req.URL.Query()
	for key, value := range action.QueryParams {
		qParams.Add(key, value)
	}
	req.URL.RawQuery = qParams.Encode()

	return req, nil
}

func logOutgoingRequest(prefix string, payload string, req *http.Request) {
	slog.Info(fmt.Sprintf("%s: sending request: ["+
		"method: %s, "+
		"url: %s, "+
		"headers: %s, "+
		"payload: %s"+
		"]",
		prefix, req.Method, req.URL, req.Header, payload))
}

// we read the body 'as is' for logging, after which we put it back into the response
// with an open reader so that it can be read downstream again
func logIncomingResponse(prefix string, res *http.Response) {
	bodyBytes, _ := io.ReadAll(res.Body)
	bodyString := ""

	err := res.Body.Close()
	if err != nil {
		slog.Error(fmt.Sprintf("%s: could not read response body: %s", prefix, err))
	} else {
		res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		bodyString = string(bodyBytes)
	}

	slog.Info(fmt.Sprintf("%s: received response: ["+
		"status: %s, "+
		"headers: %s, "+
		"payload: %s"+
		"]",
		prefix, res.Status, res.Header, bodyString))
}

func clientLogPrefix(endpointName string) string {
	return fmt.Sprintf("HTTP client %s", endpointName)
}
