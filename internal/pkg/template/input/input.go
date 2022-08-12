package input

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"
	"github.com/umisama/go-regexpcache"

	"github.com/keboola/keboola-as-code/internal/pkg/utils"
)

type Inputs []Input

// OauthAccountsSupportedComponents contains list of components supported by input kind KindOauthAccounts.
var OauthAccountsSupportedComponents = map[string]bool{ // nolint:gochecknoglobals
	"keboola.ex-google-analytics-v4": true,
	"keboola.ex-google-ads":          true,
	"keboola.ex-facebook-ads":        true,
	"keboola.ex-facebook":            true,
	"keboola.ex-instagram":           true,
}

// OauthAccountsEmptyValue contains empty values used by CLI dialog.
// KindOauthAccounts must be configured in UI,
// but at least empty keys must be generated in CLI,
// so values can be found during the upgrade operation.
var OauthAccountsEmptyValue = map[string]any{ // nolint:gochecknoglobals
	"keboola.ex-google-analytics-v4": map[string]any{
		"profiles":   []any{},
		"properties": []any{},
	},
	"keboola.ex-google-ads": map[string]any{
		"customerId":           []any{},
		"onlyEnabledCustomers": true,
	},
	"keboola.ex-facebook-ads": map[string]any{
		"accounts": map[string]any{},
	},
	"keboola.ex-facebook": map[string]any{
		"accounts": map[string]any{},
	},
	"keboola.ex-instagram": map[string]any{
		"accounts": map[string]any{},
	},
}

func NewInputs() *Inputs {
	inputs := make(Inputs, 0)
	return &inputs
}

func (i Inputs) ValidateDefinitions() error {
	errors := utils.NewMultiError()

	// Validate rules
	if err := validateDefinitions(i); err != nil {
		errors.Append(err)
	}

	// Enhance error messages
	for index, item := range errors.Errors {
		// Replace input index by input ID. Example:
		//   before: [123].default
		//   after:  input "my-input": default
		msg := regexpcache.
			MustCompile(`^\[(\d+)\]\.`).
			ReplaceAllStringFunc(item.Error(), func(s string) string {
				return fmt.Sprintf(`input "%s": `, i.GetIndex(cast.ToInt(strings.Trim(s, "[]."))).Id)
			})
		errors.Errors[index] = fmt.Errorf(msg)
	}

	return errors.ErrorOrNil()
}

func (i *Inputs) Add(input Input) {
	*i = append(*i, input)
}

func (i *Inputs) GetIndex(index int) Input {
	return (*i)[index]
}

func (i *Inputs) All() []Input {
	out := make([]Input, len(*i))
	copy(out, *i)
	return out
}

type Values []Value

func (v Values) ToMap() map[string]Value {
	res := map[string]Value{}
	for _, val := range v {
		res[val.Id] = val
	}
	return res
}

type Value struct {
	Id      string
	Value   any
	Skipped bool // the input was skipped, [showIf=false OR step has been skipped], the default value was used
}

type Input struct {
	Id           string  `json:"id" validate:"required,template-input-id"`
	Name         string  `json:"name" validate:"required,min=1,max=25"`
	Description  string  `json:"description" validate:"max=60"`
	Type         Type    `json:"type" validate:"required,template-input-type,template-input-type-for-kind"`
	Kind         Kind    `json:"kind" validate:"required,template-input-kind"`
	Default      any     `json:"default,omitempty" validate:"omitempty,template-input-default-value,template-input-default-options"`
	Rules        Rules   `json:"rules,omitempty" validate:"omitempty,template-input-rules"`
	If           If      `json:"showIf,omitempty" validate:"omitempty,template-input-if"`
	Options      Options `json:"options,omitempty" validate:"template-input-options"`
	ComponentId  string  `json:"componentId,omitempty" validate:"required_if=Kind oauth"`
	OauthInputId string  `json:"oauthInputId,omitempty" validate:"required_if=Kind oauthAccounts"`
}

// ValidateUserInput validates input from the template user using Input.Rules.
func (i Input) ValidateUserInput(userInput any) error {
	if err := i.Type.ValidateValue(reflect.ValueOf(userInput)); err != nil {
		return fmt.Errorf("%s %w", i.Name, err)
	}
	return i.Rules.ValidateValue(userInput, i.Id)
}

// Available decides if the input should be visible to user according to Input.If.
func (i Input) Available(params map[string]any) (bool, error) {
	result, err := i.If.Evaluate(params)
	if err != nil {
		return false, utils.PrefixError(fmt.Sprintf(`invalid input "%s"`, i.Id), err)
	}
	return result, nil
}

func (i Input) DefaultOrEmpty() any {
	if i.Default != nil {
		return i.Default
	}

	return i.Type.EmptyValue()
}

func (i Input) Empty() any {
	return i.Type.EmptyValue()
}
