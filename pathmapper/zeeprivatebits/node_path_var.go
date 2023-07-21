package zeeprivatebits

import "utils/urlpath"

func newPathVariableNode(nextNode pathNode) pathNode {
	node := &pathVariableNode{nextNode: nextNode}
	node.nodeType = PathVarNode
	return node
}

type pathVariableNode struct {
	commonNonTerminalNode
	nextNode pathNode
}

func (n *pathVariableNode) mapToTerminalNode(entries []string, idx int) *terminalNode {
	if _, ok := n.entryFor(entries, idx); ok {
		return n.nextNode.mapToTerminalNode(entries, idx+1)
	}
	return noopTerminalNode // indicate this path no work!
}

func (n *pathVariableNode) String() string {
	return "(" + n.nodeType.String() + " -> /" + PathVariable.String() + ")"
}

func (n *pathVariableNode) populate(builder *urlpath.LinearBuilder) {
	n.nextNode.populate(builder.AddPathEntry(PathVariable.String()))
}

func (n *pathVariableNode) getNextNode() pathNode {
	return n.nextNode
}

func (n *pathVariableNode) setNextNode(nextNode pathNode) {
	n.nextNode = nextNode
}

func toVarNode(n pathNode) *pathVariableNode {
	return n.(*pathVariableNode)
}
