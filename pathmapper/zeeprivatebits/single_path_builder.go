package zeeprivatebits

import (
	"errors"
	"fast_url_path/processors"
	"fmt"
	"utils"
)

func singlePathBuilder(processor processors.Processor, path ...any) (rootNode pathNode, err error) {
	err = checkParams(path, processor)
	if err == nil {
		return checkedParamsSinglePathBuilder(processor, path...)
	}
	return
}

func checkedParamsSinglePathBuilder(processor processors.Processor, path ...any) (rootNode pathNode, err error) {
	lastPathEntryAt := len(path) - 1
	rootNode, err = processLast(path, lastPathEntryAt, processor)
	for (err == nil) && (0 != lastPathEntryAt) {
		path = path[:lastPathEntryAt] // shorten!
		lastPathEntryAt--             // point to new lastPathEntry
		rootNode, err = processNonLast(path, lastPathEntryAt, rootNode)
	}
	return
}

func processLast(path []any, lastPathEntryAt int, processor processors.Processor) (pathNode, error) {
	if utils.IsNil(processor) {
		return nil, errors.New("no Processor provided")
	}
	terminalNode := newTerminalNode(processor)
	strPath, pathSpecial, err := checkPathEntry(path[lastPathEntryAt])
	if err == nil {
		switch pathSpecial {
		case PathAny:
			return newPathAnyNode(terminalNode), nil
		case PathVariable:
			return newPathVariableNode(terminalNode), nil
		default:
			if strPath != "" {
				return newPathStringNode(strPath, terminalNode), nil
			}
		}
		err = errors.New("unexpected pathSpecial: " + pathSpecial.String())
	}
	return nil, pathError(path, err)
}

func processNonLast(path []any, lastPathEntryAt int, lastRoot pathNode) (pathNode, error) {
	strPath, pathSpecial, err := checkPathEntry(path[lastPathEntryAt])
	if err == nil {
		switch pathSpecial {
		case PathAny:
			err = errors.New(pathSpecial.String() + " (PathAny) only supported on final urlpath entry")
		case PathVariable:
			return newPathVariableNode(lastRoot), nil
		default:
			if strPath != "" {
				return newPathStringNode(strPath, lastRoot), nil
			}
		}
		err = errors.New("unexpected pathSpecial: " + pathSpecial.String())
	}
	return nil, pathError(path, err)
}

func checkParams(path []any, processor processors.Processor) error {
	if utils.IsNil(processor) {
		return errors.New("no processors")
	}
	if utils.IsNil(path) || (len(path) == 0) {
		return errors.New("no urlpath")
	}
	return nil
}

func pathError(path []any, err error) error {
	collector := ""
	for _, s := range path {
		collector = fmt.Sprintf("%v/%v", collector, s)
	}
	return fmt.Errorf("%v %w", collector, err)
}

func checkPathEntry(entry any) (strPath string, pathSpecial PathSpecial, err error) {
	if utils.IsNil(entry) {
		err = errors.New("was null")
	} else {
		str, ok := entry.(string)
		if ok {
			strPath, err = cvtPathString(str)
		} else {
			pathSpecial, err = checkPathSpecial(entry)
		}
	}
	return
}

func cvtPathString(str string) (strPath string, err error) {
	err = utils.StringVisibleAsciiOrNonAsciiUTF8(str)
	if err == nil {
		strPath = str
	}
	return
}

func checkPathSpecial(entry any) (pathSpecial PathSpecial, err error) {
	val, ok := isValidPathSpecial(entry)
	if ok {
		pathSpecial = PathSpecial(val)
	} else {
		err = fmt.Errorf("expected 'PathSpecial', but was: %v", entry)
	}
	return
}
