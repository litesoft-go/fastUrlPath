package zeeprivatebits

import (
	"fast_url_path/processors"
	"utils"
	"utils/urlpath"
)

var noopTerminalNode = terminalNodeCommonNew("NoOp")

func newTerminalNode(processor processors.Processor) *terminalNode {
	if utils.IsNil(processor) {
		panic("processors required")
	}
	node := terminalNodeCommonNew("Processor")
	node.processor = processor
	node.hasProcessor = true
	return node
}

func terminalNodeCommonNew(text string) *terminalNode {
	node := &terminalNode{text: text}
	node.nodeType = TerminalNode
	return node
}

type terminalNode struct {
	commonNode
	processor    processors.Processor
	text         string
	hasProcessor bool
}

// if this is called then we must have fallen OFF the path
func (n *terminalNode) mapToTerminalNode(entries []string, idx int) *terminalNode {
	if _, ok := n.entryFor(entries, idx); !ok {
		return n
	}
	return noopTerminalNode // indicate this path no work!
}

func (n *terminalNode) String() string {
	return "{" + n.text + "}"
}

func (n *terminalNode) isTerminal() bool {
	return true
}

func (n *terminalNode) getProcessor() (processors.Processor, bool) {
	return n.processor, n.hasProcessor
}

func (n *terminalNode) populate(builder *urlpath.LinearBuilder) {
	builder.AddRawText(" " + n.String())
}

//goland:noinspection GoUnusedFunction
func toTerminalNode(n pathNode) *terminalNode {
	return n.(*terminalNode)
}
