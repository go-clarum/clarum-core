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

	expected := Action{
		Method:         http.MethodPost,
		Url:            "http//localhost:8080",
		Path:           "my/api/v0",
		MessagePayload: "batman!",
	}

	if actionsEqual(actual, &expected) {
		t.Errorf("Action is not as expected.")
	}
}

func TestClone(t *testing.T) {
	action := Get("my-url").
		BaseUrl("http//localhost:8080").
		ContentType("text/plain").
		Payload("my payload")

	clonedAction := action.Clone()

	if clonedAction == action {
		t.Errorf("Action has not been cloned.")
	}

	if !actionsEqual(clonedAction, action) {
		t.Errorf("Actions are not equal.")
	}
}

func TestOverwriteWith(t *testing.T) {
	baseGet := Get("base-path").
		BaseUrl("http//localhost:8080").
		ContentType("text/plain").
		Payload("my initial payload")
	action1 := Post("post-path").
		BaseUrl("https//localhost:443").
		ContentType("application/json").
		Payload("my new payload")

	if actionsEqual(baseGet, action1) {
		t.Errorf("Actions should not be equal")
	}

	baseGet.OverwriteWith(action1)

	if !actionsEqual(baseGet, action1) {
		t.Errorf("Not all fields have been overwritten.")
	}
}

func actionsEqual(a1 *Action, a2 *Action) bool {

	if a1.Method != a2.Method {
		return false
	} else if a1.Url != a2.Url {
		return false
	} else if a1.Path != a2.Path {
		return false
	} else if !maps.Equal(a1.Headers, a2.Headers) {
		return false
	} else if !maps.Equal(a1.QueryParams, a2.QueryParams) {
		return false
	} else if a1.MessagePayload != a2.MessagePayload {
		return false
	}
	return true
}
