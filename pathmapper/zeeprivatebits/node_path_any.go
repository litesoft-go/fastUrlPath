package zeeprivatebits

import "utils/urlpath"

func newPathAnyNode(nextNode *terminalNode) pathNode {
	node := &pathAnyNode{nextNode: nextNode}
	node.nodeType = PathVarNode
	return node
}

type pathAnyNode struct {
	commonNonTerminalNode
	nextNode *terminalNode
}

func (n *pathAnyNode) mapToTerminalNode(entries []string, idx int) *terminalNode {
	if _, ok := n.entryFor(entries, idx); ok {
		return n.nextNode
	}
	return noopTerminalNode // indicate this path no work!
}

func (n *pathAnyNode) String() string {
	return "(" + n.nodeType.String() + " -> /" + PathAny.String() + ")"
}

func (n *pathAnyNode) populate(builder *urlpath.LinearBuilder) {
	n.nextNode.populate(builder.AddPathEntry(PathAny.String()))
}

//goland:noinspection GoUnusedFunction
func toAnyNode(n pathNode) *pathAnyNode {
	return n.(*pathAnyNode)
}
