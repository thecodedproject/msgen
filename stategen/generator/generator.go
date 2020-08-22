package generator

import(
	"github.com/thecodedproject/msgen/generator/files/common"
	"io"
)

func Generate(
	inputFile string,
	stateStructName string,
	writer io.Writer,
) error {

	t, err := common.BaseTemplate().Parse(backendsTmpl)
	if err != nil {
		return err
	}

	return t.Execute(writer, struct{}{})
}

var backendsTmpl = `package state

imports(
	"testing"
)

type Backends interface {
}

type backendsOption func(*backendsImpl)

func NewBackendsForTesting(
	_ testing.TB,
	opts ...backendsOption,
) *State {

	var s State
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}
`
