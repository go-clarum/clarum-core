package json

import "log/slog"

type ComparatorBuilder struct {
	options
}

func Builder() *ComparatorBuilder {
	return &ComparatorBuilder{
		options{
			strictObjectSizeCheck: true,
			strictArrayCheck:      true,
			pathsToIgnore:         []string{},
			logger:                slog.Default(),
		},
	}
}

func (builder *ComparatorBuilder) StrictObjectSizeCheck(check bool) *ComparatorBuilder {
	builder.strictObjectSizeCheck = check
	return builder
}

func (builder *ComparatorBuilder) StrictArrayCheck(check bool) *ComparatorBuilder {
	builder.strictArrayCheck = check
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
