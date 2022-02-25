package dialog

import (
	"context"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cast"

	"github.com/keboola/keboola-as-code/internal/pkg/cli/options"
	"github.com/keboola/keboola-as-code/internal/pkg/cli/prompt"
	"github.com/keboola/keboola-as-code/internal/pkg/json"
	"github.com/keboola/keboola-as-code/internal/pkg/project"
	"github.com/keboola/keboola-as-code/internal/pkg/template"
	"github.com/keboola/keboola-as-code/internal/pkg/template/input"
	useTemplate "github.com/keboola/keboola-as-code/pkg/lib/operation/project/local/template/use"
	loadState "github.com/keboola/keboola-as-code/pkg/lib/operation/state/load"
)

const inputsFileFlag = "inputs-file"

type contextKey string

type useTmplDialogDeps interface {
	Options() *options.Options
	ProjectState(loadOptions loadState.Options) (*project.State, error)
}

type useTmplDialog struct {
	*Dialogs
	deps             useTmplDialogDeps
	loadStateOptions loadState.Options
	inputsFile       map[string]interface{} // inputs values loaded from a file specified by inputsFileFlag
	out              useTemplate.Options
	context          context.Context        // for input.ValidateUserInput
	inputsValues     map[string]interface{} // for input.Available
}

// AskUseTemplateOptions - dialog for using the template in the project.
func (p *Dialogs) AskUseTemplateOptions(inputs *template.Inputs, d useTmplDialogDeps, loadStateOptions loadState.Options) (useTemplate.Options, error) {
	dialog := &useTmplDialog{
		Dialogs:          p,
		deps:             d,
		loadStateOptions: loadStateOptions,
		context:          context.Background(),
		inputsValues:     make(map[string]interface{}),
	}
	return dialog.ask(inputs)
}

func (d *useTmplDialog) ask(inputs *input.Inputs) (useTemplate.Options, error) {
	// Load inputs file
	opts := d.deps.Options()
	if opts.IsSet(inputsFileFlag) {
		path := opts.GetString(inputsFileFlag)
		content, err := os.ReadFile(path) // nolint:forbidigo // file may be outside the project, so the OS package is used
		if err != nil {
			return d.out, fmt.Errorf(`cannot read inputs file "%s": %w`, path, err)
		}
		if err := json.Decode(content, &d.inputsFile); err != nil {
			return d.out, fmt.Errorf(`cannot decode inputs file "%s": %w`, path, err)
		}
	}

	// Load state
	projectState, err := d.deps.ProjectState(d.loadStateOptions)
	if err != nil {
		return d.out, err
	}

	// Target branch
	targetBranch, err := d.SelectBranch(opts, projectState.LocalObjects().Branches(), `Select the target branch`)
	if err != nil {
		return d.out, err
	}
	d.out.TargetBranch = targetBranch.BranchKey

	// User inputs
	if err := d.askInputs(inputs); err != nil {
		return d.out, err
	}

	return d.out, nil
}

// addInputValue from CLI dialog or inputs file.
func (d *useTmplDialog) addInputValue(value interface{}, inputDef input.Input) error {
	// Convert
	value, err := inputDef.Type.ParseValue(value)
	if err != nil {
		return fmt.Errorf("invalid template input: %w", err)
	}

	// Validate
	if err := inputDef.ValidateUserInput(value, d.context); err != nil {
		return fmt.Errorf("invalid template input: %w", err)
	}

	// Add
	d.context = context.WithValue(d.context, contextKey(inputDef.Id), value)
	d.inputsValues[inputDef.Id] = value
	d.out.Inputs = append(d.out.Inputs, template.InputValue{Id: inputDef.Id, Value: value})
	return nil
}

func (d *useTmplDialog) askInputs(inputs *input.Inputs) error {
	for _, inputDef := range inputs.All() {
		if result, err := inputDef.Available(d.inputsValues); err != nil {
			return err
		} else if !result {
			continue
		}
		if err := d.askInput(inputDef); err != nil {
			return err
		}
	}
	return nil
}

func (d *useTmplDialog) askInput(inputDef input.Input) error {
	// Use value from the inputs file, if it is present
	if v, found := d.inputsFile[inputDef.Id]; found {
		// Validate and save
		return d.addInputValue(v, inputDef)
	}

	// Ask for input
	switch inputDef.Kind {
	case input.KindInput, input.KindHidden, input.KindTextarea:
		question := &prompt.Question{
			Label:       inputDef.Name,
			Description: inputDef.Description,
			Validator: func(raw interface{}) error {
				value, err := inputDef.Type.ParseValue(raw)
				if err != nil {
					return err
				}
				return inputDef.ValidateUserInput(value, d.context)
			},
			Default: cast.ToString(inputDef.Default),
			Hidden:  inputDef.Kind == input.KindHidden,
		}

		var value string
		if inputDef.Kind == input.KindTextarea {
			value, _ = d.Editor("txt", question)
		} else {
			value, _ = d.Ask(question)
		}

		// Save value
		if err := d.addInputValue(value, inputDef); err != nil {
			return err
		}
	case input.KindConfirm:
		confirm := &prompt.Confirm{
			Label:       inputDef.Name,
			Description: inputDef.Description,
		}
		confirm.Default, _ = inputDef.Default.(bool)
		return d.addInputValue(d.Confirm(confirm), inputDef)
	case input.KindSelect:
		selectPrompt := &prompt.SelectIndex{
			Label:       inputDef.Name,
			Description: inputDef.Description,
			Options:     inputDef.Options.Names(),
			UseDefault:  true,
			Validator: func(answerRaw interface{}) error {
				answer := answerRaw.(survey.OptionAnswer)
				return inputDef.ValidateUserInput(answer.Value, d.context)
			},
		}
		if inputDef.Default != nil {
			if _, index, found := inputDef.Options.GetById(inputDef.Default.(string)); found {
				selectPrompt.Default = index
			}
		}
		selectedIndex, _ := d.SelectIndex(selectPrompt)
		return d.addInputValue(inputDef.Options[selectedIndex].Id, inputDef)
	case input.KindMultiSelect:
		multiSelect := &prompt.MultiSelectIndex{
			Label:       inputDef.Name,
			Description: inputDef.Description,
			Options:     inputDef.Options.Names(),
			Validator: func(answersRaw interface{}) error {
				answers := answersRaw.([]survey.OptionAnswer)
				values := make([]string, len(answers))
				for i, v := range answers {
					values[i] = v.Value
				}
				return inputDef.ValidateUserInput(values, d.context)
			},
		}
		// Default indices
		if inputDef.Default != nil {
			defaultIndices := make([]int, 0)
			for _, id := range inputDef.Default.([]interface{}) {
				if _, index, found := inputDef.Options.GetById(id.(string)); found {
					defaultIndices = append(defaultIndices, index)
				}
			}
			multiSelect.Default = defaultIndices
		}
		// Selected values
		selectedIndices, _ := d.MultiSelectIndex(multiSelect)
		selectedValues := make([]string, 0)
		for _, selectedIndex := range selectedIndices {
			selectedValues = append(selectedValues, inputDef.Options[selectedIndex].Id)
		}
		// Save value
		return d.addInputValue(selectedValues, inputDef)
	}

	return nil
}
