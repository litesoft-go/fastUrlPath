package zeeprivatebits

import (
	"utils/enums"
)

type NodeType int

const (
	PathAnyNode NodeType = iota + 1
	PathStringNode
	PathStringsMapNode
	PathVarNode
	PathLessOrNode // has NO urlpath!
	TerminalNode
	geInvalidNodeType
)

var enumNodeType = enums.Init("NodeType", geInvalidNodeType).
	Add(PathAnyNode, "PathAnyNode").
	Add(PathStringNode, "PathStringNode").
	Add(PathStringsMapNode, "PathStringsMapNode").
	Add(PathVarNode, "PathVarNode").
	Add(PathLessOrNode, "PathLessOrNode").
	Add(TerminalNode, "TerminalNode").
	Build()

func (enum NodeType) String() string {
	return enumNodeType.ToString(enum)
}
