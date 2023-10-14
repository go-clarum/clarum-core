package http

import (
	"bytes"
	"fmt"
	"github.com/goclarum/clarum/core/control"
	"io"
	"log/slog"
	"net/http"
	"testing"
)

type ClientEndpoint struct {
	name            string
	baseUrl         string
	contentType     string
	client          *http.Client
	responseChannel chan *http.Response
}

func (ce *ClientEndpoint) Send(t *testing.T, action *Action) {
	logPrefix := clientLogPrefix(ce.name)
	slog.Debug(fmt.Sprintf("%s: action to send: %s", logPrefix, action.ToString()))
	control.RunningActions.Add(1)

	go func() {
		defer control.RunningActions.Done()

		actionToExecute := ce.getSendActionToExecute(action)
		slog.Debug(fmt.Sprintf("%s: executing action: %s", logPrefix, actionToExecute.ToString()))

		req, err := buildRequest(ce.name, actionToExecute)
		if err != nil {
			t.Errorf("%s: canceled action. Error: %s", logPrefix, err)
		}

		logOutgoingRequest(logPrefix, action.payload, req)
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

func (ce *ClientEndpoint) Receive(t *testing.T, action *Action) {
	logPrefix := clientLogPrefix(ce.name)
	slog.Debug(fmt.Sprintf("%s: action to receive: %s", logPrefix, action.ToString()))

	response := <-ce.responseChannel

	actionToExecute := ce.getReceiveActionToExecute(action)
	slog.Debug(fmt.Sprintf("%s: executing validation action: %s", logPrefix, actionToExecute.ToString()))

	validateHttpStatusCode(t, logPrefix, actionToExecute, response.StatusCode)
	validateHttpHeaders(t, logPrefix, actionToExecute, response.Header)
	validateHttpBody(t, logPrefix, actionToExecute, response.Body)
}

// put missing data into a send action
func (ce *ClientEndpoint) getSendActionToExecute(action *Action) *Action {
	actionToExecute := action.Clone()

	if len(actionToExecute.baseUrl) == 0 {
		actionToExecute.baseUrl = ce.baseUrl
	}
	if len(actionToExecute.headers) == 0 || len(actionToExecute.headers[ContentTypeHeaderName]) == 0 {
		actionToExecute.ContentType(ce.contentType)
	}

	return actionToExecute
}

// put missing data into a receive action
func (ce *ClientEndpoint) getReceiveActionToExecute(action *Action) *Action {
	actionToExecute := action.Clone()

	if len(actionToExecute.headers) == 0 || len(actionToExecute.headers[ContentTypeHeaderName]) == 0 {
		actionToExecute.ContentType(ce.contentType)
	}

	return actionToExecute
}

func buildRequest(prefix string, action *Action) (*http.Request, error) {
	url := BuildPath(action.baseUrl, action.path)

	req, err := http.NewRequest(action.method, url, bytes.NewBufferString(action.payload))
	if err != nil {
		slog.Error(fmt.Sprintf("%s: error: %s", prefix, err))
		return nil, err
	}

	for header, value := range action.headers {
		req.Header.Set(header, value)
	}

	qParams := req.URL.Query()
	for key, value := range action.queryParams {
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
