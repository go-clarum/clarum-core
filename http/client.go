package http

import (
	"fmt"
	"github.com/goclarum/clarum/core/control"
	"io"
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

func (ce *ClientEndpoint) Send(action *Action) {
	control.RunningActions.Add(1)

	go func() {
		defer control.RunningActions.Done()

		actionToExecute := ce.getActionToExecute(action)
		url := BuildPath(actionToExecute.baseUrl, actionToExecute.path)

		req, err := http.NewRequest(action.method, url, nil)
		if err != nil {
			fmt.Println(fmt.Sprintf("HTTP client <%s> error: %s", ce.name, err))
		}

		for header, value := range actionToExecute.headers {
			req.Header.Set(header, value)
		}

		qParams := req.URL.Query()
		for key, value := range actionToExecute.queryParams {
			qParams.Add(key, value)
		}
		req.URL.RawQuery = qParams.Encode()

		res, err := ce.client.Do(req)

		// TODO: handle technical errors here
		//   debug logging - log entire response as is
		if err != nil {
			fmt.Println(fmt.Sprintf("HTTP client <%s> error: %s", ce.name, err))
			return
		}

		ce.responseChannel <- res
	}()
}

func (ce *ClientEndpoint) Receive(t *testing.T, action *Action) {
	actionToExecute := ce.getActionToExecute(action)
	response := <-ce.responseChannel

	// debug logging - log validating message

	if response.StatusCode != actionToExecute.statusCode {
		t.Errorf("HTTP client <%s> validation error: HTTP status mismatch", ce.name)
	}

	if err := validateHeaders(actionToExecute, response.Header); err != nil {
		t.Errorf("HTTP client <%s>: %s", ce.name, err)
	} else {
		// debug logging
		fmt.Println(fmt.Sprintf("HTTP client <%s> header validation successful", ce.name))
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(fmt.Sprintf("HTTP client <%s>: could not read response body: %s", ce.name, err))
	}
	fmt.Println(fmt.Sprintf("HTTP Client <%s> response payload: %s", ce.name, body))
}

func (ce *ClientEndpoint) getActionToExecute(action *Action) *Action {
	actionToExecute := action.Clone()

	if len(actionToExecute.baseUrl) == 0 {
		actionToExecute.baseUrl = ce.baseUrl
	}
	if len(actionToExecute.headers) == 0 || len(actionToExecute.headers[ContentTypeHeaderName]) == 0 {
		actionToExecute.ContentType(ce.contentType)
	}

	return actionToExecute
}
