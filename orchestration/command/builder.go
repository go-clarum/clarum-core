package command

import (
	"time"
)

// TODO
// documentation
// error messages
// unit tests

type Builder struct {
	name       string
	components []string
	warmup     time.Duration
}

// export this into the orchestration package to keep API clean
func Command() *Builder {
	return &Builder{}
}

func (ib *Builder) Components(components ...string) *Builder {
	ib.components = components
	return ib
}

func (ib *Builder) Warmup(warmup time.Duration) *Builder {
	ib.warmup = warmup
	return ib
}

func (ib *Builder) Build() *Endpoint {
	return newCommandEndpoint(ib.name, ib.components, ib.warmup)
}
