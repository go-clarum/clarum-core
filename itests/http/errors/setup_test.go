package errors

import (
	"errors"
	"github.com/goclarum/clarum"
	"os"
	"testing"
	"time"
)

var errorsClient = clarum.Http().Client().
	Name("errorsClient").
	Timeout(2000 * time.Millisecond).
	Build()

var errorsServer = clarum.Http().Server().
	Name("errorsServer").
	Port(8083).
	Build()

func TestMain(m *testing.M) {
	clarum.Setup()

	result := m.Run()

	clarum.Finish()

	os.Exit(result)
}

func checkErrors(t *testing.T, expectedError string, actionErrors ...error) {
	allErrors := errors.Join(actionErrors...)

	if allErrors == nil {
		t.Error("One error expected, but there was none.")
	} else {
		if allErrors.Error() != expectedError {
			t.Errorf("Unexpected errors: %s", allErrors)
		}
	}
}
