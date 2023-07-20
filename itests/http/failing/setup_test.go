package http

import (
	"github.com/goclarum/clarum"
	"log"
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
	log.Println(log.Ldate|log.Ltime|log.Lshortfile, "My main test setup")

	result := m.Run()

	clarum.Finish()

	os.Exit(result)
}
