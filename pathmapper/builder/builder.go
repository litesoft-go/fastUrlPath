package builder

import (
	"fast_url_path/pathmapper/zeeprivatebits"
	"fast_url_path/processors"
)

type Builder interface {
	Register(processor processors.Processor, pathEntries ...any) Builder
	Build(noPathMatchedProcessor processors.Processor) (processors.ProcessorMapper, error)
}

//goland:noinspection GoUnusedFunction
func newBuilder() Builder {
	return zeeprivatebits.InnerBuilder()
}
