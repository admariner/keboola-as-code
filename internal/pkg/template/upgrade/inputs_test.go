package upgrade

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/template/input"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/orderedmap"
)

type configInput struct {
	id         string
	inMetadata bool
	inContent  bool
	value      interface{}
}

type rowInput struct {
	id         string
	inMetadata bool
	inContent  bool
	value      interface{}
}

type testCase struct {
	name            string
	configInputs    []configInput
	rowInputs       []rowInput
	templateInputs  input.StepsGroups
	expected        input.StepsGroups
	configuredSteps []string
}

func TestExportInputsValues(t *testing.T) {
	t.Parallel()

	testCases := []testCase{
		{
			name:         "no-input",
			configInputs: []configInput{},
			rowInputs:    []rowInput{},
			templateInputs: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps:    input.Steps{{Inputs: input.Inputs{}}},
			}},
			expected: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps:    input.Steps{{Inputs: input.Inputs{}}},
			}},
		},
		{
			name:         "value-not-present",
			configInputs: []configInput{},
			rowInputs:    []rowInput{},
			templateInputs: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeString,
						Kind:    input.KindInput,
						Default: "default",
					},
				}}},
			}},
			expected: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeString,
						Kind:    input.KindInput,
						Default: "default",
					},
				}}},
			}},
		},
		{
			name: "input-only-in-metadata",
			configInputs: []configInput{
				{
					id:         "input1",
					inMetadata: true,
					inContent:  false,
					value:      "my value",
				},
			},
			rowInputs: []rowInput{},
			templateInputs: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeString,
						Kind:    input.KindInput,
						Default: "default",
					},
				}}},
			}},
			expected: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeString,
						Kind:    input.KindInput,
						Default: "default",
					},
				}}},
			}},
		},
		{
			name: "value-present-in-config",
			configInputs: []configInput{
				{
					id:         "input1",
					inMetadata: true,
					inContent:  true,
					value:      "my value",
				},
			},
			rowInputs: []rowInput{},
			templateInputs: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeString,
						Kind:    input.KindInput,
						Default: "default",
					},
				}}},
			}},
			expected: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeString,
						Kind:    input.KindInput,
						Default: "my value", // <<<<<<<<<
					},
				}}},
			}},
			configuredSteps: []string{"g01-s01"},
		},
		{
			name: "value-present-in-config-invalid-type-1",
			configInputs: []configInput{
				{
					id:         "input1",
					inMetadata: true,
					inContent:  true,
					value:      123, // <<<<<<
				},
			},
			rowInputs: []rowInput{},
			templateInputs: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeString,
						Kind:    input.KindInput,
						Default: "default",
					},
				}}},
			}},
			expected: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeString,
						Kind:    input.KindInput,
						Default: "default",
					},
				}}},
			}},
		},
		{
			name: "value-present-in-config-invalid-type-2",
			configInputs: []configInput{
				{
					id:         "input1",
					inMetadata: true,
					inContent:  true,
					value:      "foo", // <<<<<<
				},
			},
			rowInputs: []rowInput{},
			templateInputs: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeInt,
						Kind:    input.KindInput,
						Default: 123,
					},
				}}},
			}},
			expected: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeInt,
						Kind:    input.KindInput,
						Default: 123,
					},
				}}},
			}},
		},
		{
			name:         "value-present-in-row",
			configInputs: []configInput{},
			rowInputs: []rowInput{
				{
					id:         "input1",
					inMetadata: true,
					inContent:  true,
					value:      "my value",
				},
			},
			templateInputs: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeString,
						Kind:    input.KindInput,
						Default: "default",
					},
				}}},
			}},
			expected: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeString,
						Kind:    input.KindInput,
						Default: "my value",
					},
				}}},
			}},
			configuredSteps: []string{"g01-s01"},
		},
		{
			name:         "value-present-in-row-invalid-type",
			configInputs: []configInput{},
			rowInputs: []rowInput{
				{
					id:         "input1",
					inMetadata: true,
					inContent:  true,
					value:      123, // <<<<<<
				},
			},
			templateInputs: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeString,
						Kind:    input.KindInput,
						Default: "default",
					},
				}}},
			}},
			expected: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeString,
						Kind:    input.KindInput,
						Default: "default",
					},
				}}},
			}},
		},
		{
			name: "value-present-multiple-times",
			configInputs: []configInput{
				{
					id:         "input1",
					inMetadata: true,
					inContent:  true,
					value:      "value 1",
				},
				{
					id:         "input1",
					inMetadata: true,
					inContent:  true,
					value:      "value 2",
				},
				{
					id:         "input1",
					inMetadata: true,
					inContent:  true,
					value:      123,
				},
			},
			rowInputs: []rowInput{
				{
					id:         "input1",
					inMetadata: true,
					inContent:  true,
					value:      "value 3",
				},
				{
					id:         "input1",
					inMetadata: true,
					inContent:  true,
					value:      "value 4",
				},
			},
			templateInputs: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeString,
						Kind:    input.KindInput,
						Default: "default",
					},
				}}},
			}},
			expected: input.StepsGroups{{
				Required: input.RequiredOptional,
				Steps: input.Steps{{Inputs: input.Inputs{
					{
						Id:      "input1",
						Type:    input.TypeString,
						Kind:    input.KindInput,
						Default: "value 4",
					},
				}}},
			}},
			configuredSteps: []string{"g01-s01"},
		},
	}

	// Test all cases
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.run(t)
		})
	}
}

func (tc testCase) run(t *testing.T) {
	t.Helper()

	// Create objects
	d := dependencies.NewTestContainer()
	state := d.EmptyState()
	branchKey := model.BranchKey{Id: 123}
	configKey := model.ConfigKey{BranchId: 123, ComponentId: "foo.bar", Id: "111"}
	configRowKey := model.ConfigRowKey{BranchId: 123, ComponentId: "foo.bar", ConfigId: "111", Id: "222"}
	configMetadata := model.ConfigMetadata{}
	configContent := orderedmap.New()
	rowContent := orderedmap.New()

	// Set instance ID
	instanceId := "12345"
	configMetadata.SetTemplateInstance("repo", "template", instanceId)
	configMetadata.SetConfigTemplateId("configInTemplate")
	configMetadata.AddRowTemplateId(configRowKey.Id, "rowInTemplate")

	// Add config inputs
	for index, inputDef := range tc.configInputs {
		contentKey := fmt.Sprintf("foo.bar.item%d", index)
		if inputDef.inMetadata {
			configMetadata.AddInputUsage(inputDef.id, orderedmap.KeyFromStr(contentKey))
		}
		if inputDef.inContent {
			assert.NoError(t, configContent.SetNested(contentKey, inputDef.value))
		}
	}

	// Add row inputs
	for index, inputDef := range tc.rowInputs {
		contentKey := fmt.Sprintf("foo.bar.item%d", index)
		if inputDef.inMetadata {
			configMetadata.AddRowInputUsage(configRowKey.Id, inputDef.id, orderedmap.KeyFromStr(contentKey))
		}
		if inputDef.inContent {
			assert.NoError(t, rowContent.SetNested(contentKey, inputDef.value))
		}
	}

	// Set objects to state
	assert.NoError(t, state.Set(&model.ConfigState{
		ConfigManifest: &model.ConfigManifest{ConfigKey: configKey},
		Local:          &model.Config{ConfigKey: configKey, Metadata: configMetadata, Content: configContent},
	}))
	assert.NoError(t, state.Set(&model.ConfigRowState{
		ConfigRowManifest: &model.ConfigRowManifest{ConfigRowKey: configRowKey},
		Local:             &model.ConfigRow{ConfigRowKey: configRowKey, Content: rowContent},
	}))

	// Assert inputs
	actual := ExportInputsValues(log.NewNopLogger(), state, branchKey, instanceId, tc.templateInputs)
	assert.Equal(t, tc.expected, actual.ToValue())

	// Assert steps state
	var configuredSteps []string
	_ = actual.VisitSteps(func(group *input.StepsGroupExt, step *input.StepExt) error {
		if step.Show {
			configuredSteps = append(configuredSteps, step.Id)
		}
		return nil
	})
	assert.Equal(t, tc.configuredSteps, configuredSteps)
}
