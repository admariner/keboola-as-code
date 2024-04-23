package jsonnet

import (
	"bytes"
	"encoding/json"
	"sync"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/google/go-jsonnet/parser"

	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

type VM struct {
	lock    sync.Mutex
	vm      *jsonnet.VM
	err     error
	payload any
}

// GetPayload returns the payload of the current evaluation.
// Can be used by custom functions for request specific values.
func (vm *VM) GetPayload() any {
	return vm.payload
}

func (vm *VM) Evaluate(code string, payload any) (jsonOut string, err error) {
	if vm.err != nil {
		return "", vm.err
	}

	vm.lock.Lock()
	defer vm.lock.Unlock()

	vm.payload = payload

	astNode, err := ToAst(code, "")
	if err != nil {
		return "", err
	}

	// Pre-process
	node := ast.Clone(astNode)
	if err := parser.PreprocessAst(&node, vm.vm.GlobalVars()...); err != nil {
		return "", err
	}

	// Evaluate
	jsonContent, err := vm.vm.Evaluate(node)
	if err != nil {
		return "", errors.Errorf(`jsonnet error: %w`, err)
	}

	// Format (go-jsonnet library use 3 space indent)
	var out bytes.Buffer

	err = json.Compact(&out, []byte(jsonContent))
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func (vm *VM) Validate(code string) error {
	_, err := parser.SnippetToAst(code, "", vm.vm.GlobalVars()...)
	return err
}
