package command

import "testing"

type ActionBuilder struct {
	endpoint *Endpoint
}

type TestActionBuilder struct {
	test *testing.T
	ActionBuilder
}

func (endpoint *Endpoint) In(t *testing.T) *TestActionBuilder {
	return &TestActionBuilder{
		test: t,
		ActionBuilder: ActionBuilder{
			endpoint: endpoint,
		},
	}
}

func (endpoint *Endpoint) Run() error {
	return endpoint.start()
}

func (endpoint *Endpoint) Stop() error {
	return endpoint.stop()
}

func (builder *TestActionBuilder) Run() {
	if err := builder.endpoint.start(); err != nil {
		builder.test.Error(err)
	}
}

func (builder *TestActionBuilder) Stop() {
	if err := builder.endpoint.stop(); err != nil {
		builder.test.Error(err)
	}
}
