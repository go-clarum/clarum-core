package client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/goclarum/clarum/core/control"
	"github.com/goclarum/clarum/http/constants"
	"github.com/goclarum/clarum/http/internal/utils"
	"github.com/goclarum/clarum/http/internal/validators"
	"github.com/goclarum/clarum/http/message"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Endpoint struct {
	name            string
	baseUrl         string
	contentType     string
	client          *http.Client
	responseChannel chan *responsePair
}

type responsePair struct {
	response *http.Response
	error    error
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
		responseChannel: make(chan *responsePair),
	}
}

func (endpoint *Endpoint) send(message *message.Message) error {
	logPrefix := clientLogPrefix(endpoint.name)
	slog.Debug(fmt.Sprintf("%s: message to send: %s", logPrefix, message.ToString()))

	messageToSend := endpoint.getMessageToSend(message)
	slog.Debug(fmt.Sprintf("%s: sending message: %s", logPrefix, messageToSend.ToString()))

	req, err := buildRequest(endpoint.name, messageToSend)
	// we return error here directly and not in the goroutine below
	// this way we can signal to the test synchronously that there was an error
	if err != nil {
		errorMessage := fmt.Sprintf("%s: canceled message. Error: %s", logPrefix, err)
		slog.Error(errorMessage)
		return errors.New(errorMessage)
	}

	control.RunningActions.Add(1)

	go func() {
		defer control.RunningActions.Done()

		logOutgoingRequest(logPrefix, message.MessagePayload, req)
		res, err := endpoint.client.Do(req)

		// we log the error here directly, but will do error handling downstream
		if err != nil {
			slog.Error(fmt.Sprintf("%s: error on response: %s", logPrefix, err))
		} else {
			logIncomingResponse(logPrefix, res)
		}

		endpoint.responseChannel <- &responsePair{
			response: res,
			error:    err,
		}
	}()

	return nil
}

func (endpoint *Endpoint) receive(message *message.Message) error {
	logPrefix := clientLogPrefix(endpoint.name)
	slog.Debug(fmt.Sprintf("%s: message to receive: %s", logPrefix, message.ToString()))

	responsePair := <-endpoint.responseChannel

	// TODO: handle technical errors
	//  check: socket connection, connection refused, connection timeout
	if responsePair.error != nil {
		errorMessage := fmt.Sprintf("%s: error while receiving response: %s", logPrefix, responsePair.error)
		slog.Error(errorMessage)
		return errors.New(errorMessage)
	}

	messageToReceive := endpoint.getMessageToReceive(message)
	slog.Debug(fmt.Sprintf("%s: validating message: %s", logPrefix, messageToReceive.ToString()))

	// TODO: check errors.Join();
	if err := validators.ValidateHttpStatusCode(logPrefix, messageToReceive, responsePair.response.StatusCode); err != nil {
		return err
	}
	if err := validators.ValidateHttpHeaders(logPrefix, messageToReceive, responsePair.response.Header); err != nil {
		return err
	}
	if err := validators.ValidateHttpBody(logPrefix, messageToReceive, responsePair.response.Body); err != nil {
		return err
	}

	return nil
}

// Put missing data into a message to send: baseUrl & ContentType Header
func (endpoint *Endpoint) getMessageToSend(message *message.Message) *message.Message {
	messageToSend := message.Clone()

	if len(messageToSend.Url) == 0 {
		messageToSend.Url = endpoint.baseUrl
	}
	if len(messageToSend.Headers) == 0 || len(messageToSend.Headers[constants.ContentTypeHeaderName]) == 0 {
		messageToSend.ContentType(endpoint.contentType)
	}

	return messageToSend
}

// Put missing data into message to receive: ContentType Header
func (endpoint *Endpoint) getMessageToReceive(message *message.Message) *message.Message {
	messageToReceive := message.Clone()

	if len(messageToReceive.Headers) == 0 || len(messageToReceive.Headers[constants.ContentTypeHeaderName]) == 0 {
		messageToReceive.ContentType(endpoint.contentType)
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
