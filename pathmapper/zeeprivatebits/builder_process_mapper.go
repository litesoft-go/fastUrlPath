package zeeprivatebits

import (
	"errors"
	"fast_url_path/processors"
	"net/http"
	"net/url"
	"strings"
	"utils"
)

type builder struct {
	merger   *nodeMerger
	graph    pathNode
	count    int
	firstErr error
}

func InnerBuilder() builder.Builder {
	return innerBuilder()
}

// Register - processors and path -- path entries are limited to "string" or "pathSpecials"
func (b *builder) Register(processor processors.Processor, pathEntries ...any) builder.Builder {
	node, err := singlePathBuilder(processor, pathEntries...)
	if err == nil {
		if b.count == 0 {
			b.graph = node
		} else {
			b.graph, err = b.merger.merge("/", node, b.graph)
		}
		b.count++
	}
	if (err != nil) && (b.firstErr == nil) {
		b.firstErr = err
	}
	return b
}

func (b *builder) Build(noPathMatchedProcessor processors.Processor) (processors.ProcessorMapper, error) {
	return b.innerBuild(noPathMatchedProcessor)
}

func innerBuilder() *builder {
	return &builder{merger: newNodeMerger()}
}

func (b *builder) innerBuild(noPathMatchedProcessor processors.Processor) (*processMapper, error) {
	if b.firstErr != nil {
		return nil, b.firstErr
	}
	if b.count == 0 {
		return nil, errors.New("no processors-urlpath pairs registered")
	}
	if utils.IsNil(noPathMatchedProcessor) {
		noPathMatchedProcessor = processors.Return404
	}
	m := &processMapper{graph: b.graph, noPathMatchedProcessor: noPathMatchedProcessor}
	b.graph = nil
	b.merger = nil
	return m, nil
}

type processMapper struct {
	graph                  pathNode
	noPathMatchedProcessor processors.Processor
}

//goland:noinspection GoUnusedParameter
func (m *processMapper) FromRequest(req *http.Request) (processor processors.Processor, err error) {
	if req == nil {
		return m.FromUrl(req.URL)
	}
	err = errors.New("no 'req' provided at ProcessorMapper.FromRequest")
	return
}

func (m *processMapper) FromUrl(url *url.URL) (processor processors.Processor, err error) {
	if url != nil {
		return m.FromPath(url.Path)
	}
	err = errors.New("no 'url' provided at ProcessorMapper.FromUrl")
	return
}

func (m *processMapper) FromPath(path string) (processors.Processor, error) {
	for strings.HasSuffix(path, "/") || strings.HasSuffix(path, " ") { // remove all trailing slashes and spaces!
		path = path[:len(path)-1]
	}
	return m.fromPathEntries(strings.Split(path, "/")...)
}

func (m *processMapper) fromPathEntries(pathEntries ...string) (processors.Processor, error) {
	node := m.graph.mapToTerminalNode(pathEntries, 0)
	processor, ok := node.getProcessor()
	if !ok {
		processor = m.noPathMatchedProcessor
	}
	return processor, nil
}
