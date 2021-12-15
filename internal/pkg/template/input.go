package template

import (
	"fmt"
	"reflect"

	goValidator "github.com/go-playground/validator/v10"
	"github.com/umisama/go-regexpcache"

	"github.com/keboola/keboola-as-code/internal/pkg/validator"
)

type Inputs []*Input

func (i Inputs) Validate() error {
	validations := []validator.Validation{
		{
			Tag:  "template-input-id",
			Func: validateId,
		},
		{
			Tag:  "template-input-default",
			Func: validateDefault,
		},
		{
			Tag:  "template-input-options",
			Func: validateOptions,
		},
		{
			Tag:  "template-input-type",
			Func: validateType,
		},
	}
	return validator.Validate(i, validations...)
}

func validateId(fl goValidator.FieldLevel) bool {
	return regexpcache.MustCompile(`^[a-zA-Z0-9\.\_]+$`).MatchString(fl.Field().String())
}

// Default value must be of the same type as the Type or Options.
func validateDefault(fl goValidator.FieldLevel) bool {
	if fl.Field().Kind() == reflect.Ptr && fl.Field().IsNil() {
		return true
	}
	if fl.Field().IsZero() {
		return true
	}
	// Check if Default is present in Options
	if fl.Parent().FieldByName("Kind").String() == "select" || fl.Parent().FieldByName("Kind").String() == "multiselect" {
		for _, x := range fl.Parent().FieldByName("Options").Interface().([]Option) {
			if string(x) == fl.Field().String() {
				return true
			}
		}
		return false
	}
	err := checkTypeAgainstKind(fl.Field(), fl.Field().Kind().String())
	return err == nil
}

// Options must be filled only for select or multiselect Kind.
func validateOptions(fl goValidator.FieldLevel) bool {
	if fl.Parent().FieldByName("Kind").String() == "select" || fl.Parent().FieldByName("Kind").String() == "multiselect" {
		return fl.Field().Len() > 0
	}
	return fl.Field().Len() == 0
}

// Valid only for input Kind.
func validateType(fl goValidator.FieldLevel) bool {
	return fl.Parent().FieldByName("Kind").String() == "input"
}

type Input struct {
	Id          string      `json:"id" validate:"required,template-input-id"`
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description" validate:"required"`
	Default     interface{} `json:"default,omitempty" validate:"omitempty,template-input-default"`
	Kind        string      `json:"kind" validate:"required,oneof=input password textarea confirm select multiselect"`
	Type        string      `json:"type,omitempty" validate:"required_if=Kind input,omitempty,oneof=string int float64,template-input-type"`
	Options     []Option    `json:"options,omitempty" validate:"required_if=Type select Type multiselect,template-input-options"`
	Rules       string      `json:"rules,omitempty"`
	If          string      `json:"if,omitempty"`
}

func checkTypeAgainstKind(value interface{}, kind string) error {
	inputType := reflect.TypeOf(value).String()
	switch kind {
	case "password", "textarea":
		if inputType != reflect.String.String() {
			return fmt.Errorf("the input is of %s kind and should be a string, got %s instead", kind, inputType)
		}
	case "confirm":
		if inputType != reflect.Bool.String() {
			return fmt.Errorf("the input is of confirm kind and should be a bool, got %s instead", inputType)
		}
	}
	return nil
}

type Option string
