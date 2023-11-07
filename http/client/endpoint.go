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

func (ce *Endpoint) send(t *testing.T, message *message.Message) {
	logPrefix := clientLogPrefix(ce.name)
	slog.Debug(fmt.Sprintf("%s: message to send: %s", logPrefix, message.ToString()))
	control.RunningActions.Add(1)

	go func() {
		defer control.RunningActions.Done()

		// from here
		messageToSend := ce.getMessageToSend(message)
		slog.Debug(fmt.Sprintf("%s: sending message: %s", logPrefix, messageToSend.ToString()))

		req, err := buildRequest(ce.name, messageToSend)
		if err != nil {
			t.Errorf("%s: canceled message. Error: %s", logPrefix, err)
		}
		//to here, move outside the goroutine and return error if t is nil

		logOutgoingRequest(logPrefix, message.MessagePayload, req)
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

func (ce *Endpoint) receive(t *testing.T, message *message.Message) {
	logPrefix := clientLogPrefix(ce.name)
	slog.Debug(fmt.Sprintf("%s: message to receive: %s", logPrefix, message.ToString()))

	response := <-ce.responseChannel

	messageToReceive := ce.getMessageToReceive(message)
	slog.Debug(fmt.Sprintf("%s: validating message: %s", logPrefix, messageToReceive.ToString()))

	validators.ValidateHttpStatusCode(t, logPrefix, messageToReceive, response.StatusCode)
	validators.ValidateHttpHeaders(t, logPrefix, messageToReceive, response.Header)
	validators.ValidateHttpBody(t, logPrefix, messageToReceive, response.Body)
}

// put missing data into a message to send
func (ce *Endpoint) getMessageToSend(message *message.Message) *message.Message {
	messageToSend := message.Clone()

	if len(messageToSend.Url) == 0 {
		messageToSend.Url = ce.baseUrl
	}
	if len(messageToSend.Headers) == 0 || len(messageToSend.Headers[constants.ContentTypeHeaderName]) == 0 {
		messageToSend.ContentType(ce.contentType)
	}

	return messageToSend
}

// put missing data into message to receive
func (ce *Endpoint) getMessageToReceive(message *message.Message) *message.Message {
	messageToReceive := message.Clone()

	if len(messageToReceive.Headers) == 0 || len(messageToReceive.Headers[constants.ContentTypeHeaderName]) == 0 {
		messageToReceive.ContentType(ce.contentType)
	}

	return messageToReceive
}

func buildRequest(prefix string, message *message.Message) (*http.Request, error) {
	url := utils.BuildPath(message.Url, message.Path)

	req, err := http.NewRequest(message.Method, url, bytes.NewBufferString(message.MessagePayload))
	if err != nil {
		slog.Error(fmt.Sprintf("%s: error: %s", prefix, err))
		return nil, err
	}

	for header, value := range message.Headers {
		req.Header.Set(header, value)
	}

	qParams := req.URL.Query()
	for key, value := range message.QueryParams {
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
