package zeeprivatebits

import "fast_url_path/processors"

type commonNonTerminalNode struct {
	commonNode
}

func (n *commonNonTerminalNode) isTerminal() bool {
	return false
}

func (n *commonNonTerminalNode) getProcessor() (processors.Processor, bool) {
	return nil, false
}
