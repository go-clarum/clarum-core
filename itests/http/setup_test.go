package http

import (
	"github.com/goclarum/clarum"
	"os"
	"testing"
	"time"
)

var testClient = clarum.Http().Client().
	Name("testClient").
	BaseUrl("http://localhost:8083/myApp").
	Timeout(2000 * time.Millisecond).
	Build()

var testServer = clarum.Http().Server().
	Name("testServer").
	Port(8083).
	Build()

func TestMain(m *testing.M) {
	clarum.Setup()

	result := m.Run()

	clarum.Finish()

	os.Exit(result)
}
