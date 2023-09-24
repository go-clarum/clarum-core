package http

import (
	"github.com/goclarum/clarum"
	"os"
	"testing"
)

var Client1 = clarum.Http().Client().
	Name("client1").
	BaseUrl("http://localhost:8083/myApp").
	ContentType("application/json").
	Build()

var Server1 = clarum.Http().Server().
	Name("server1").
	Port(8083).
	ContentType("application/json").
	Build()

func TestMain(m *testing.M) {
	clarum.Setup()

	result := m.Run()

	clarum.Finish()

	os.Exit(result)
}
