package zeeprivatebits

import (
	"fast_url_path/processors"
	"utils/urlpath"
)

type pathNode interface {
	getType() NodeType

	isType(nodeType NodeType) bool

	mapToTerminalNode(entries []string, idx int) *terminalNode

	String() string

	isTerminal() bool

	getProcessor() (processors.Processor, bool)

	populate(builder *urlpath.LinearBuilder)
}

type commonNode struct {
	nodeType NodeType
}

func (n *commonNode) getType() NodeType {
	return n.nodeType
}

func (n *commonNode) isType(nodeType NodeType) bool {
	return n.nodeType == nodeType
}

func (n *commonNode) entryFor(entries []string, idx int) (entry string, ok bool) {
	if idx < len(entries) {
		return entries[idx], true
	}
	return
}
