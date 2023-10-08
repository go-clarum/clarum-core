package http

import (
	"github.com/goclarum/clarum"
	"os"
	"testing"
)

var testClient = clarum.Http().Client().
	Name("testClient").
	BaseUrl("http://localhost:8083/myApp").
	ContentType("application/json").
	Build()

var testServer = clarum.Http().Server().
	Name("testServer").
	Port(8083).
	ContentType("application/json").
	Build()

func TestMain(m *testing.M) {
	clarum.Setup()

	result := m.Run()

	clarum.Finish()

	os.Exit(result)
}
