package zeeprivatebits

import (
	"utils/enums"
)

type PathSpecial int

const (
	PathVariable PathSpecial = iota + 1
	PathAny
	geInvalidPathSpecial
)

var enumPathSpecial = enums.Init("PathSpecial", geInvalidPathSpecial).
	Add(PathVariable, "{variable}").
	Add(PathAny, "{any}").
	Build()

func (enum PathSpecial) String() string {
	return enumPathSpecial.ToString(enum)
}

func isValidPathSpecial(it any) (val int, ok bool) {
	return enumPathSpecial.IsValid(it)
}
