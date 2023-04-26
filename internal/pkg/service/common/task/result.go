package task

import (
	"strings"

	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

type Result struct {
	result  string
	error   error
	outputs map[string]any
}

func OkResult(msg string) Result {
	if strings.TrimSpace(msg) == "" {
		panic(errors.New("message cannot be empty"))
	}
	return Result{result: msg}
}

func ErrResult(err error) Result {
	if err == nil {
		panic(errors.New("error cannot be nil"))
	}
	return Result{error: err}
}

func (r Result) Result() string {
	return r.result
}

func (r Result) Err() error {
	return r.error
}

func (r Result) IsErr() bool {
	return r.error != nil
}

// WithResult can be used to edit the result message later.
func (r Result) WithResult(result string) Result {
	if r.error == nil && r.result == "" {
		panic(errors.New(`result struct is empty, use task.OkResult(msg) or task.ErrResult(err) function instead`))
	}
	if r.error != nil {
		panic(errors.New(`result type is "error", not "ok", it cannot be modified`))
	}
	r.result = result
	return r
}

// WithErr can be used to edit the error message later.
func (r Result) WithErr(err error) Result {
	if r.error == nil && r.result == "" {
		panic(errors.New(`result struct is empty, use task.OkResult(msg) or task.ErrResult(err) function instead`))
	}
	if r.error == nil {
		panic(errors.New(`result type is "ok", not "error", it cannot be modified`))
	}
	r.error = err
	return r
}

// WithOutput adds some task operation output.
func (r Result) WithOutput(k string, v any) Result {
	if r.error == nil && r.result == "" {
		panic(errors.New(`result struct is empty, use task.OkResult(msg) or task.ErrResult(err) function first`))
	}

	// Clone map
	original := r.outputs
	r.outputs = make(map[string]any)
	for key, value := range original {
		r.outputs[key] = value
	}

	// Add new key
	r.outputs[k] = v
	return r
}