package zeeprivatebits

import (
	"utils/urlpath"
)

func newPathStringsMapNode(path string, nextNode pathNode) *pathStringsMapNode {
	node := &pathStringsMapNode{pathMap: make(map[string]pathNode)}
	node.nodeType = PathStringsMapNode
	node.add(path, nextNode)
	return node
}

func (n *pathStringsMapNode) add(path string, nextNode pathNode) *pathStringsMapNode {
	if _, ok := n.pathMap[path]; ok {
		panic("attempt to add duplicate path '" + path + "' to MapNode")
	}
	n.paths = append(n.paths, path)
	return n.update(path, nextNode)
}

func (n *pathStringsMapNode) update(path string, nextNode pathNode) *pathStringsMapNode {
	n.pathMap[path] = nextNode
	return n
}

func (n *pathStringsMapNode) get(path string) (pathNode, bool) {
	node, ok := n.pathMap[path]
	return node, ok
}

func (n *pathStringsMapNode) getAll() (entries []*pathEntryNextNodePair) {
	for _, path := range n.paths {
		node := n.pathMap[path]
		pair := &pathEntryNextNodePair{pathEntry: path, nextNode: node}
		entries = append(entries, pair)
	}
	return
}

type pathStringsMapNode struct {
	commonNonTerminalNode
	paths   []string
	pathMap map[string]pathNode
}

func (n *pathStringsMapNode) mapToTerminalNode(entries []string, idx int) *terminalNode {
	if entry, ok := n.entryFor(entries, idx); ok {
		if node, ok := n.pathMap[entry]; ok {
			return node.mapToTerminalNode(entries, idx+1)
		}
	}
	return noopTerminalNode // indicate this path no work!
}

func (n *pathStringsMapNode) String() string {
	return "(" + n.nodeType.String() + " -> /???)"
}

func (n *pathStringsMapNode) populate(builder *urlpath.LinearBuilder) {
	mapBuilder := builder.AddMap()
	for _, entry := range n.paths {
		node := n.pathMap[entry]
		node.populate(mapBuilder.AddOption(entry))
	}
}

func toStringsMapNode(n pathNode) *pathStringsMapNode {
	return n.(*pathStringsMapNode)
}
