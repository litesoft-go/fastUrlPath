package zeeprivatebits

import (
	"fmt"
	"strconv"
	"testing"
	"utils"
	"utils/urlpath"
)

func paths(entries ...[]any) [][]any {
	return entries
}

func TestBuilder_build(t *testing.T) {
	tests := []struct {
		paths   [][]any
		wantErr bool
		tree    string
	}{
		{paths(sliceOf("p1", "p2")),
			false, "/p1/p2 {Processor}"},
		{paths(sliceOf("p1", PathVariable, "p3", PathAny)),
			false, "/p1/{variable}/p3/{any} {Processor}"},
		{paths(sliceOf("p1", "p2"),
			sliceOf("p1", PathVariable, "p3", PathAny),
			sliceOf("p1", "p2b", "p3"),
			sliceOf("p1", "p2b", "p3", PathAny)),
			false, utils.LinesToString("",
			"/p1/┤¿p2 {Processor}",
			"     ¿p2b/p3/┤{any} {Processor}",
			"             ┤ {Processor}",
			"    ┤{variable}/p3/{any} {Processor}",
			"")},
	}
	for i, tt := range tests {
		name := fmt.Sprintf("builder(urlpath merging)[%d]", i)
		t.Run(name, func(t *testing.T) {
			mb := innerBuilder()
			for _, pathEntries := range tt.paths {
				mb.Register(processorNotNil, pathEntries...)
			}
			got, err := mb.innerBuild(nil)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("%v error = %v, wantErr %v", name, err, tt.wantErr)
				}
			} else {
				builder := urlpath.NewGraphBuilder()
				got.graph.populate(builder.GetLinearBuilder())
				actual := builder.String()
				expected := tt.tree
				success := actual == expected
				if !success {
					t.Error(name +
						"\n--- expected (" + strconv.Itoa(len(expected)) + "):" +
						"\n" + expected +
						"\n--- actual (" + strconv.Itoa(len(actual)) + "):" +
						"\n" + actual + "\n")
				}
			}
		})
	}
}
