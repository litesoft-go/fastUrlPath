package zeeprivatebits

import "utils/urlpath"

func newPathStringNode(pathEntry string, nextNode pathNode) pathNode {
	node := &pathStringNode{}
	node.nodeType = PathStringNode
	node.pathEntry = pathEntry
	node.nextNode = nextNode
	return node
}

func (n *pathStringNode) toPair() *pathEntryNextNodePair {
	return &(n.pathEntryNextNodePair)
}

type pathEntryNextNodePair struct {
	pathEntry string
	nextNode  pathNode
}

type pathStringNode struct {
	commonNonTerminalNode
	pathEntryNextNodePair
}

func (n *pathStringNode) mapToTerminalNode(entries []string, idx int) *terminalNode {
	if entry, ok := n.entryFor(entries, idx); ok && (entry == n.pathEntry) {
		return n.nextNode.mapToTerminalNode(entries, idx+1)
	}
	return noopTerminalNode // indicate this path no work!
}

func (n *pathStringNode) String() string {
	return "(" + n.nodeType.String() + " -> /" + n.pathEntry + ")"
}

func (n *pathStringNode) populate(builder *urlpath.LinearBuilder) {
	n.nextNode.populate(builder.AddPathEntry(n.pathEntry))
}

func (n *pathStringNode) getNextNode() pathNode {
	return n.nextNode
}

func (n *pathStringNode) setNextNode(nextNode pathNode) {
	n.nextNode = nextNode
}

func toStringNode(n pathNode) *pathStringNode {
	return n.(*pathStringNode)
}
