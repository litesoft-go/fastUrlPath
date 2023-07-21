package zeeprivatebits

import (
	"fmt"
	"utils/urlpath"
)

func newPathLessOrNode() pathNode {
	node := &pathLessOrNode{tNode: noopTerminalNode}
	node.nodeType = PathLessOrNode
	return node
}

type pathLessOrNode struct {
	commonNonTerminalNode
	prioritizedNodes []pathNode // groups: Stringy, Var, Any
	tNode            *terminalNode
}

func (n *pathLessOrNode) mapToTerminalNode(entries []string, idx int) *terminalNode {
	if _, ok := n.entryFor(entries, idx); !ok {
		return n.tNode
	}
	for _, node := range n.prioritizedNodes {
		tNode := node.mapToTerminalNode(entries, idx)
		if _, ok := tNode.getProcessor(); ok {
			return tNode
		}
	}
	return noopTerminalNode // indicate this path no work!
}

func (n *pathLessOrNode) insertAbove(at int, node pathNode) {
	n.prioritizedNodes = append(append(append([]pathNode{}, n.prioritizedNodes[:at]...), node), n.prioritizedNodes[at:]...)
}

func (n *pathLessOrNode) add(m *nodeMerger, path string, node pathNode) (err error) {
	if node.isType(TerminalNode) {
		if _, ok := n.tNode.getProcessor(); ok {
			return fmt.Errorf("%v, at '%v' already has a 'real' TerminalNode", n.nodeType, path)
		}
		n.tNode = toTerminalNode(node)
		return
	}
	newPriority := priorityFor(node)
	if newPriority == -1 {
		return fmt.Errorf("%v, at '%v' does not support '%v' entries", n.nodeType, path, node.getType())
	}
	above, equal, at := n.shouldInsertAt(newPriority)
	if above {
		n.insertAbove(at, node)
	} else if equal {
		n.prioritizedNodes[at], err = m.merge(path, n.prioritizedNodes[at], node)
	} else {
		n.prioritizedNodes = append(n.prioritizedNodes, node)
	}
	return
}

func (n *pathLessOrNode) shouldInsertAt(newPriority int) (above, equal bool, at int) {
	for i, existingNode := range n.prioritizedNodes {
		above, equal = shouldInsert(newPriority, existingNode)
		if above || equal {
			at = i
			return
		}
	}
	return
}

func shouldInsert(newPriority int, existingNode pathNode) (above, equal bool) {
	existingPriority := priorityFor(existingNode) // won't return -1 -- Filtered above
	above = newPriority < existingPriority
	equal = newPriority == existingPriority
	return
}

func priorityFor(node pathNode) int {
	switch node.getType() {
	case PathStringNode, PathStringsMapNode:
		return 0
	case PathVarNode:
		return 1
	case PathAnyNode:
		return 2
	}
	return -1
}

func (n *pathLessOrNode) String() (str string) {
	str += "("
	for _, node := range n.prioritizedNodes[1:] {
		str += node.getType().String() + " | "
	}
	str += n.tNode.String() + ")"
	return
}

func (n *pathLessOrNode) populate(builder *urlpath.LinearBuilder) {
	orBuilder := builder.AddOr()
	for _, node := range n.prioritizedNodes {
		node.populate(orBuilder.NextOption())
	}
	if _, ok := n.tNode.getProcessor(); ok {
		n.tNode.populate(orBuilder.NextOption())
	}
}

func toOrNode(n pathNode) *pathLessOrNode {
	return n.(*pathLessOrNode)
}
