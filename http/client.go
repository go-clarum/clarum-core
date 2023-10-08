package http

import (
	"bytes"
	"fmt"
	"github.com/goclarum/clarum/core/control"
	"github.com/goclarum/clarum/core/validators/strings"
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
	slog.Debug(fmt.Sprintf("HTTP client %s - action to send: %s", ce.name, action.ToString()))
	control.RunningActions.Add(1)

	go func() {
		defer control.RunningActions.Done()

		actionToExecute := ce.getSendActionToExecute(action)
		slog.Debug(fmt.Sprintf("HTTP client %s executing action: %s", ce.name, actionToExecute.ToString()))

		req, err := buildRequest(ce.name, actionToExecute)
		if err != nil {
			t.Errorf("HTTP client %s canceled action. Error: %s", ce.name, err)
		}

		logRequest(ce.name, action.payload, req)
		res, err := ce.client.Do(req)
		logResponse(ce.name, res)

		// TODO: handle technical errors
		//  check: socket connection, connection refused, connection timeout
		if err != nil {
			t.Errorf("HTTP client %s error on response: %s", ce.name, err)
		}

		ce.responseChannel <- res
	}()
}

func (ce *ClientEndpoint) Receive(t *testing.T, action *Action) {
	slog.Debug(fmt.Sprintf("HTTP client %s - action to receive: %s", ce.name, action.ToString()))

	response := <-ce.responseChannel

	actionToExecute := ce.getReceiveActionToExecute(action)
	slog.Debug(fmt.Sprintf("HTTP client %s executing validation action: %s", ce.name, actionToExecute.ToString()))

	validateHttpStatusCode(t, ce.name, response, actionToExecute)
	validateHttpHeaders(t, ce.name, actionToExecute, response)
	validateHttpBody(t, ce.name, actionToExecute, response)
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

func buildRequest(endpointName string, action *Action) (*http.Request, error) {
	url := BuildPath(action.baseUrl, action.path)

	req, err := http.NewRequest(action.method, url, bytes.NewBufferString(action.payload))
	if err != nil {
		slog.Error(fmt.Sprintf("HTTP client %s error: %s", endpointName, err))
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

func validateHttpStatusCode(t *testing.T, endpointName string, response *http.Response, actionToExecute *Action) {
	if response.StatusCode != actionToExecute.statusCode {
		t.Errorf("HTTP client %s validation error: HTTP status mismatch", endpointName)
	} else {
		slog.Debug(fmt.Sprintf("HTTP client %s HTTP status validation successful", endpointName))
	}
}

func validateHttpHeaders(t *testing.T, endpointName string, actionToExecute *Action, response *http.Response) {
	if err := validateHeaders(actionToExecute, response.Header); err != nil {
		t.Errorf("HTTP client %s: %s", endpointName, err)
	} else {
		slog.Debug(fmt.Sprintf("HTTP client %s header validation successful", endpointName))
	}
}

func validateHttpBody(t *testing.T, endpointName string, actionToExecute *Action, response *http.Response) {
	defer closeBody(endpointName, response.Body)

	if strings.IsBlank(actionToExecute.payload) {
		slog.Debug(fmt.Sprintf("HTTP client %s action payload is empty. No body validation will be done", endpointName))
		return
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("HTTP client %s: could not read response body: %s", endpointName, err)
	}

	if err := validatePayload(actionToExecute, bodyBytes); err != nil {
		t.Errorf("HTTP client %s: %s", endpointName, err)
	} else {
		slog.Debug(fmt.Sprintf("HTTP client %s payload validation successful", endpointName))
	}
}

func closeBody(endpointName string, body io.ReadCloser) {
	if err := body.Close(); err != nil {
		slog.Error(fmt.Sprintf("HTTP client %s unable to close body: %s", endpointName, err))
	}
}

func logRequest(endpointName string, payload string, req *http.Request) {
	slog.Info(fmt.Sprintf("HTTP client %s sending request: ["+
		"method: %s, "+
		"url: %s, "+
		"headers: %s, "+
		"payload: %s"+
		"]",
		endpointName, req.Method, req.URL, req.Header, payload))
}

// we read the body 'as is' for logging, after which we put it back into the response
// with an open reader so that it can be read downstream again
func logResponse(endpointName string, res *http.Response) {
	bodyBytes, _ := io.ReadAll(res.Body)
	bodyString := ""

	err := res.Body.Close()
	if err != nil {
		slog.Error(fmt.Sprintf("HTTP client %s: could not read response body: %s", endpointName, err))
	} else {
		res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		bodyString = string(bodyBytes)
	}

	slog.Info(fmt.Sprintf("HTTP client %s received response: ["+
		"status: %s, "+
		"headers: %s, "+
		"payload: %s"+
		"]",
		endpointName, res.Status, res.Header, bodyString))
}
