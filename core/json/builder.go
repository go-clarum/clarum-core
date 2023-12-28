package json

import "log/slog"

type ComparatorBuilder struct {
	options
}

func Builder() *ComparatorBuilder {
	return &ComparatorBuilder{
		options{
			strictObjectSizeCheck: true,
			pathsToIgnore:         []string{},
			logger:                slog.Default(),
			recorder:              &noopRecorder{},
		},
	}
}

func (builder *ComparatorBuilder) StrictObjectSizeCheck(check bool) *ComparatorBuilder {
	builder.strictObjectSizeCheck = check
	return builder
}

func (builder *ComparatorBuilder) PathsToIgnore(paths ...string) *ComparatorBuilder {
	builder.pathsToIgnore = append(builder.pathsToIgnore, paths...)
	return builder
}

func (builder *ComparatorBuilder) Logger(logger *slog.Logger) *ComparatorBuilder {
	builder.logger = logger
	return builder
}

func (builder *ComparatorBuilder) Recorder(recorder Recorder) *ComparatorBuilder {
	builder.recorder = recorder
	return builder
}

func (builder *ComparatorBuilder) Comparator() *Comparator {
	return &Comparator{builder.options}
}
