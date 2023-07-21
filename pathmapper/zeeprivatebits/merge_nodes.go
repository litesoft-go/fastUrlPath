package zeeprivatebits

import (
	"errors"
)

// Notes:
// "PathLessOrNode(s)" (OrNodes) should be managed to make the following node type the last in
//  any sequence of checks
//
// "TerminalNode(s)" should only exist in the following locations:
// -- 1 - end of a node list,
// -- 2 - 2nd entry in the "right"-most node of an OR tree!
//
// "PathAnyNode(s)" should only exist in the following locations:
// -- 1 - just before end "TerminalNode" of a node list,
// -- 2 - 1st OR 2nd entry in the "right"-most node of an OR tree!
//
// "PathVarNode(s)" should be encountered BEFORE any "TerminalNode(s)" or "PathAnyNode(s)"!

type mergeFunc func(m *nodeMerger, path string, n1, n2 pathNode) (node pathNode, err error)

type nodeMerger struct {
	pairsMap map[string]mergeFunc
}

func newNodeMerger() *nodeMerger {
	pairsMap := make(map[string]mergeFunc)
	pairsMap["TerminalNode-TerminalNode"] = mergeIncompatible
	pairsMap["PathAnyNode-PathAnyNode"] = mergeIncompatible
	pairsMap["PathLessOrNode-PathLessOrNode"] = mergeIncompatible

	pairsMap["PathLessOrNode-PathAnyNode"] = mergeOrNodeWith
	pairsMap["PathLessOrNode-PathVarNode"] = mergeOrNodeWith
	pairsMap["PathLessOrNode-PathStringNode"] = mergeOrNodeWith
	pairsMap["PathLessOrNode-PathStringsMapNode"] = mergeOrNodeWith
	pairsMap["PathLessOrNode-TerminalNode"] = mergeOrNodeWith

	pairsMap["TerminalNode-PathAnyNode"] = mergeNonOrNodesToOrNode
	pairsMap["TerminalNode-PathVarNode"] = mergeNonOrNodesToOrNode
	pairsMap["TerminalNode-PathStringNode"] = mergeNonOrNodesToOrNode
	pairsMap["TerminalNode-PathStringsMapNode"] = mergeNonOrNodesToOrNode

	pairsMap["PathAnyNode-PathVarNode"] = mergeNonOrNodesToOrNode
	pairsMap["PathAnyNode-PathStringNode"] = mergeNonOrNodesToOrNode
	pairsMap["PathAnyNode-PathStringsMapNode"] = mergeNonOrNodesToOrNode

	pairsMap["PathVarNode-PathVarNode"] = identicalVarNodes
	pairsMap["PathVarNode-PathStringNode"] = mergeNonOrNodesToOrNode
	pairsMap["PathVarNode-PathStringsMapNode"] = mergeNonOrNodesToOrNode

	pairsMap["PathStringNode-PathStringNode"] = twoStringNodes
	pairsMap["PathStringNode-PathStringsMapNode"] = mergeStringWithStringsMapNode

	pairsMap["PathStringsMapNode-PathStringsMapNode"] = mergeStringsMapWithStringsMapNode

	return &nodeMerger{pairsMap}
}

func (m *nodeMerger) merge(path string, n1, n2 pathNode) (node pathNode, err error) {
	ns1 := n1.getType().String()
	ns2 := n2.getType().String()
	if mFunc, ok := m.pairsMap[ns1+"-"+ns2]; ok {
		return mFunc(m, path, n1, n2)
	}
	if mFunc, ok := m.pairsMap[ns2+"-"+ns1]; ok {
		return mFunc(m, path, n2, n1)
	}
	panic("no merge function registered for: " + ns1 + " x " + ns2)
}

func mergeIncompatible(m *nodeMerger, path string, n1, n2 pathNode) (node pathNode, err error) {
	return m.mergeError(path, n1, n2, "are mutually exclusive")
}

func mergeNonOrNodesToOrNode(m *nodeMerger, path string, node1, node2 pathNode) (node pathNode, err error) {
	if node, err = mergeOrNodeWith(m, path, newPathLessOrNode(), node1); err == nil {
		node, err = mergeOrNodeWith(m, path, node, node2)
	}
	return
}

func mergeOrNodeWith(m *nodeMerger, path string, orNode, newNode pathNode) (node pathNode, err error) {
	nOR := toOrNode(orNode)
	return nOR, nOR.add(m, path, newNode)
}

func identicalVarNodes(m *nodeMerger, path string, varNode1, varNode2 pathNode) (node pathNode, err error) {
	n1 := toVarNode(varNode1)
	n2 := toVarNode(varNode2)
	return mergeNextNodes(m, path, n1.getType().String(), n1, n2)
}

func twoStringNodes(m *nodeMerger, path string, stringNode1, stringNode2 pathNode) (node pathNode, err error) {
	n1 := toStringNode(stringNode1)
	n2 := toStringNode(stringNode2)
	sn1 := n1.pathEntry
	sn2 := n2.pathEntry
	if sn1 != sn2 { // need map - happy case!
		return newPathStringsMapNode(sn1, n1.nextNode).add(sn2, n2.nextNode), nil
	}
	return mergeNextNodes(m, path, n1.pathEntry, n1, n2)
}

func mergeStringWithStringsMapNode(m *nodeMerger, path string, stringNode, stringsMapNode pathNode) (node pathNode, err error) {
	return stringsMapNode, addSourcePairToMap(m, path, toStringNode(stringNode).toPair(), toStringsMapNode(stringsMapNode))
}

func mergeStringsMapWithStringsMapNode(m *nodeMerger, path string, stringsMapNode1, stringsMapNode2 pathNode) (node pathNode, err error) {
	entries := toStringsMapNode(stringsMapNode1).getAll()
	target := toStringsMapNode(stringsMapNode2)
	for _, entry := range entries {
		if err = addSourcePairToMap(m, path, entry, target); err != nil {
			return
		}
	}
	return target, nil
}

func addSourcePairToMap(m *nodeMerger, path string, source *pathEntryNextNodePair, target *pathStringsMapNode) error {
	strEntry := source.pathEntry
	strNextNode := source.nextNode
	mapNextNode, ok := target.get(strEntry)
	if !ok { // string NOT in map - happy case
		target.add(strEntry, strNextNode)
	} else if newNextNode, err := m.merge(mergePathAppend(path, strEntry), mapNextNode, strNextNode); err != nil {
		return err
	} else {
		target.update(strEntry, newNextNode)
	}
	return nil
}

func mergeNextNodes(m *nodeMerger, path, pathAppendWith string, n1, n2 hasNextNode) (pathNode, error) {
	newChildNode, err := m.merge(mergePathAppend(path, pathAppendWith), n1.getNextNode(), n2.getNextNode())
	if err == nil {
		n1.setNextNode(newChildNode)
	}
	return n1, err
}

func mergePathAppend(path, entry string) string {
	return path + entry + "/"
}

//goland:noinspection GoUnusedParameter
func (m *nodeMerger) mergeError(path string, n1, n2 pathNode, problem string) (node pathNode, err error) {
	ns1 := n1.String()
	ns2 := n2.String()

	prefix := "'" + ns1
	if ns1 == ns2 {
		prefix = "two " + prefix
	} else {
		prefix += "' and '" + ns2
	}
	return nil, errors.New(prefix + "', at '" + path + "', " + problem)
}
