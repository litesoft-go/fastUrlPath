package zeeprivatebits

type nextNodeGetter interface {
	getNextNode() pathNode
}

type nextNodeSetter interface {
	setNextNode(nextNode pathNode)
}

type nextNodeAccessor interface {
	nextNodeGetter
	nextNodeSetter
}

type hasNextNode interface {
	pathNode
	nextNodeAccessor
}
