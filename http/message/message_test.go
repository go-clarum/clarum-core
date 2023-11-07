package message

import (
	"maps"
	"net/http"
	"testing"
)

func TestHTTPVerbs(t *testing.T) {
	if Get().Method != http.MethodGet {
		t.Errorf("Expected %s.", http.MethodGet)
	}
	if Head().Method != http.MethodHead {
		t.Errorf("Expected %s.", http.MethodHead)
	}
	if Post().Method != http.MethodPost {
		t.Errorf("Expected %s.", http.MethodPost)
	}
	if Put().Method != http.MethodPut {
		t.Errorf("Expected %s.", http.MethodPut)
	}
	if Delete().Method != http.MethodDelete {
		t.Errorf("Expected %s.", http.MethodDelete)
	}
	if Options().Method != http.MethodOptions {
		t.Errorf("Expected %s.", http.MethodOptions)
	}
	if Trace().Method != http.MethodTrace {
		t.Errorf("Expected %s.", http.MethodTrace)
	}
	if Patch().Method != http.MethodPatch {
		t.Errorf("Expected %s.", http.MethodPatch)
	}
}

func TestBuilder(t *testing.T) {
	actual := Post("my", "api/v0").
		BaseUrl("http//localhost:8080").
		ContentType("text/plain").
		Payload("batman!")

	expected := Message{
		Method:         http.MethodPost,
		Url:            "http//localhost:8080",
		Path:           "my/api/v0",
		MessagePayload: "batman!",
	}

	if messagesEqual(actual, &expected) {
		t.Errorf("Message is not as expected.")
	}
}

func TestClone(t *testing.T) {
	message := Get("my-url").
		BaseUrl("http//localhost:8080").
		ContentType("text/plain").
		Payload("my payload")

	clonedMessage := message.Clone()

	if clonedMessage == message {
		t.Errorf("Message has not been cloned.")
	}

	if !messagesEqual(clonedMessage, message) {
		t.Errorf("Messages are not equal.")
	}
}

func TestOverwriteWith(t *testing.T) {
	baseGet := Get("base-path").
		BaseUrl("http//localhost:8080").
		ContentType("text/plain").
		Payload("my initial payload")
	postMessage := Post("post-path").
		BaseUrl("https//localhost:443").
		ContentType("application/json").
		Payload("my new payload")

	if messagesEqual(baseGet, postMessage) {
		t.Errorf("Messages should not be equal")
	}

	baseGet.OverwriteWith(postMessage)

	if !messagesEqual(baseGet, postMessage) {
		t.Errorf("Not all fields have been overwritten.")
	}
}

func messagesEqual(m1 *Message, m2 *Message) bool {

	if m1.Method != m2.Method {
		return false
	} else if m1.Url != m2.Url {
		return false
	} else if m1.Path != m2.Path {
		return false
	} else if !maps.Equal(m1.Headers, m2.Headers) {
		return false
	} else if !maps.Equal(m1.QueryParams, m2.QueryParams) {
		return false
	} else if m1.MessagePayload != m2.MessagePayload {
		return false
	}
	return true
}
