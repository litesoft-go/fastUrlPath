package zeeprivatebits

import (
	"fast_url_path/processors"
	"utils/urlpath"
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func sliceOf(entries ...any) []any {
	return entries
}

var processorNil processors.Processor = nil
var processorNotNil processors.Processor = func(http.ResponseWriter, *http.Request) (statusCode int, err error) {
	return
}

func Test_checkPathSpecial(t *testing.T) {
	tests := []struct {
		name            string
		arg             any
		wantPathSpecial PathSpecial
		wantErr         bool
	}{
		{"PathVariable", PathVariable, PathVariable, false},
		{"-1", -1, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPathSpecial, err := checkPathSpecial(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkPathSpecial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPathSpecial != tt.wantPathSpecial {
				t.Errorf("checkPathSpecial() gotPathSpecial = %v, want %v", gotPathSpecial, tt.wantPathSpecial)
			}
			if err != nil {
				fmt.Println("    error:", tt.name, "--", err)
			}
		})
	}
}

func Test_cvtPathString(t *testing.T) {
	tests := []struct {
		name        string
		arg         string
		wantStrPath string
		wantErr     bool
	}{
		{"empty", "", "", true},
		{"space", " ", "", true},
		{"a", "a", "a", false},
		{"a-b", "a-b", "a-b", false},
		{"a b", "a b", "", true},
		{"a{nl}b", "a\nb", "", true},
		{"a{31}b", string([]byte{'a', 31, 'b'}), "", true},
		{"a{127}b", string([]byte{'a', 127, 'b'}), "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStrPath, err := cvtPathString(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("cvtPathString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStrPath != tt.wantStrPath {
				t.Errorf("cvtPathString() gotStrPath = %v, want %v", gotStrPath, tt.wantStrPath)
			}
			if err != nil {
				fmt.Println("    error:", tt.name, "--", err)
			}
		})
	}
}

func Test_checkParams(t *testing.T) {
	tests := []struct {
		argProcessor processors.Processor
		argPath      []any
		wantErr      bool
	}{
		{processorNil, sliceOf("p"), true},
		{processorNotNil, sliceOf(), true},
		{processorNotNil, nil, true},
		{processorNotNil, sliceOf("p"), false},
	}
	for i, tt := range tests {
		name := fmt.Sprintf("checkParams[%d]", i)
		t.Run(name, func(t *testing.T) {
			if err := checkParams(tt.argPath, tt.argProcessor); (err != nil) != tt.wantErr {
				t.Errorf("%v error = %v, wantErr %v", name, err, tt.wantErr)
			}
		})
	}
}

func Test_singlePathBuilder(t *testing.T) {
	tests := []struct {
		resultPath string
		argPath    []any
		wantErr    bool
	}{
		{"", sliceOf(nil), true},
		{"/p {Processor}", sliceOf("p"), false},
		{"/p1/p2 {Processor}", sliceOf("p1", "p2"), false},
		{"/p1/{variable}/p3/{any} {Processor}", sliceOf("p1", PathVariable, "p3", PathAny), false},
	}
	for i, tt := range tests {
		name := fmt.Sprintf("checkedParamsSinglePathBuilder[%d]", i)
		t.Run(name, func(t *testing.T) {
			rootNode, err := checkedParamsSinglePathBuilder(processorNotNil, tt.argPath...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("%v, error = %v, wantErr %v", name, err, tt.wantErr)
				}
			} else {
				builder := urlpath.NewGraphBuilder()
				rootNode.populate(builder.GetLinearBuilder())
				actual := builder.String()
				expected := tt.resultPath
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
