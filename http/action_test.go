package http

import (
	"github.com/goclarum/clarum/core/maps"
	"net/http"
	"testing"
)

func TestHTTPVerbs(t *testing.T) {
	if Get().method != http.MethodGet {
		t.Errorf("Expected %s.", http.MethodGet)
	}
	if Head().method != http.MethodHead {
		t.Errorf("Expected %s.", http.MethodHead)
	}
	if Post().method != http.MethodPost {
		t.Errorf("Expected %s.", http.MethodPost)
	}
	if Put().method != http.MethodPut {
		t.Errorf("Expected %s.", http.MethodPut)
	}
	if Delete().method != http.MethodDelete {
		t.Errorf("Expected %s.", http.MethodDelete)
	}
	if Options().method != http.MethodOptions {
		t.Errorf("Expected %s.", http.MethodOptions)
	}
	if Trace().method != http.MethodTrace {
		t.Errorf("Expected %s.", http.MethodTrace)
	}
	if Patch().method != http.MethodPatch {
		t.Errorf("Expected %s.", http.MethodPatch)
	}
}

func TestBuilder(t *testing.T) {
	actual := Post("my", "api/v0").
		BaseUrl("http//localhost:8080").
		ContentType("text/plain")

	expected := Action{
		method:  http.MethodPost,
		baseUrl: "http//localhost:8080",
		path:    "my/api/v0",
	}

	if actionsEqual(actual, &expected) {
		t.Errorf("Action is not as expected.")
	}
}

func TestClone(t *testing.T) {
	action := Get("my-url").
		BaseUrl("http//localhost:8080").
		ContentType("text/plain")

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
		ContentType("text/plain")
	action1 := Post("post-path").
		BaseUrl("https//localhost:443").
		ContentType("application/json")

	if actionsEqual(baseGet, action1) {
		t.Errorf("Actions should not be equal")
	}

	baseGet.OverwriteWith(action1)

	if !actionsEqual(baseGet, action1) {
		t.Errorf("Not all fields have been overwritten.")
	}
}

func actionsEqual(a1 *Action, a2 *Action) bool {

	if a1.method != a2.method {
		return false
	} else if a1.baseUrl != a2.baseUrl {
		return false
	} else if a1.path != a2.path {
		return false
	} else if !maps.EqualString(a1.headers, a2.headers) {
		return false
	} else if !maps.EqualString(a1.queryParams, a2.queryParams) {
		return false
	}

	return true
}
