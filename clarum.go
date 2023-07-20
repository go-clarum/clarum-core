package clarum

import (
	"fmt"
	"github.com/goclarum/clarum/core/control"
	"github.com/goclarum/clarum/http"
)

func Http() HttpBuilder {
	return &httpEndpointBuilder{}
}

type HttpBuilder interface {
	Client() *http.ClientEndpointBuilder
	Server() *http.ServerEndpointBuilder
}

type httpEndpointBuilder struct {
}

func (heb *httpEndpointBuilder) Client() *http.ClientEndpointBuilder {
	return http.Client()
}

func (heb *httpEndpointBuilder) Server() *http.ServerEndpointBuilder {
	return http.Server()
}

func Finish() {
	fmt.Println(fmt.Sprintf("Waiting for all actions to finish."))

	control.RunningActions.Wait()

	fmt.Println(fmt.Sprintf("All actions finished."))
}
