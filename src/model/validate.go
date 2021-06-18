package model

import (
	"fmt"
	"keboola-as-code/src/utils"
	"keboola-as-code/src/validator"
)

type stateValidator struct {
	error *utils.Error
}

func validateState(state *State) *utils.Error {
	v := &stateValidator{}

	for _, b := range state.Branches() {
		v.validate("branch", b.Remote)
		v.validate("branch", b.Local)
		v.validate("branch manifest record", b.Manifest)
	}

	for _, c := range state.Components() {
		v.validate("component", c.Remote)
	}

	for _, c := range state.Configs() {
		v.validate("config", c.Remote)
		v.validate("config", c.Local)
		v.validate("config manifest record", c.Manifest)
	}

	for _, r := range state.ConfigRows() {
		v.validate("config row", r.Remote)
		v.validate("config row", r.Local)
		v.validate("config row manifest record", r.Manifest)
	}

	return v.error
}

func (s *stateValidator) AddError(err error) {
	s.error.Add(err)
}

func (s *stateValidator) validate(kind string, v interface{}) {
	if err := validator.Validate(v); err != nil {
		s.AddError(fmt.Errorf("%s is not valid: %s", kind, err))
	}
}
