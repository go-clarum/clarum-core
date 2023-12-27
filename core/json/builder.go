package json

import "log/slog"

type ComparatorBuilder struct {
	options
}

func Builder() *ComparatorBuilder {
	return &ComparatorBuilder{
		options{
			strictSizeCheck: true,
			pathsToIgnore:   []string{},
			logger:          slog.Default(),
		},
	}
}

func (builder *ComparatorBuilder) StrictSizeCheck(check bool) *ComparatorBuilder {
	builder.strictSizeCheck = check
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

func (builder *ComparatorBuilder) Comparator() *Comparator {
	return &Comparator{builder.options}
}
