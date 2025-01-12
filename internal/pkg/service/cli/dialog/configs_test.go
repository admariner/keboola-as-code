package dialog_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/model"
)

func TestSelectConfigInteractive(t *testing.T) {
	t.Parallel()

	// Dependencies
	dialog, _, console := createDialogs(t, true)

	// All configs
	config1 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "1"}, Name: `Config 1`}}
	config2 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "2"}, Name: `Config 2`}}
	config3 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "3"}, Name: `Config 3`}}
	allConfigs := []*model.ConfigWithRows{config1, config2, config3}

	// Interaction
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		assert.NoError(t, console.ExpectString("LABEL:"))

		assert.NoError(t, console.ExpectString("Config 1 (foo.bar:1)"))

		assert.NoError(t, console.ExpectString("Config 2 (foo.bar:2)"))

		assert.NoError(t, console.ExpectString("Config 3 (foo.bar:3)"))

		// down arrow -> select Config 2
		assert.NoError(t, console.SendDownArrow())
		assert.NoError(t, console.SendEnter())

		assert.NoError(t, console.ExpectEOF())
	}()

	// Run
	out, err := dialog.SelectConfig(allConfigs, `LABEL`)
	assert.Same(t, config2, out)
	assert.NoError(t, err)

	// Close terminal
	assert.NoError(t, console.Tty().Close())
	wg.Wait()
	assert.NoError(t, console.Close())
}

func TestSelectConfigByFlag(t *testing.T) {
	t.Parallel()

	// Dependencies
	dialog, o, _ := createDialogs(t, false)
	o.Set(`config`, `2`)

	// All configs
	config1 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "1"}, Name: `Config 1`}}
	config2 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "2"}, Name: `Config 2`}}
	config3 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "3"}, Name: `Config 3`}}
	allConfigs := []*model.ConfigWithRows{config1, config2, config3}

	// Run
	out, err := dialog.SelectConfig(allConfigs, `LABEL`)
	assert.Same(t, config2, out)
	assert.NoError(t, err)
}

func TestSelectConfigNonInteractive(t *testing.T) {
	t.Parallel()

	// Dependencies
	dialog, o, _ := createDialogs(t, false)
	o.Set(`non-interactive`, `true`)

	// All configs
	config1 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "1"}, Name: `Config 1`}}
	config2 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "2"}, Name: `Config 2`}}
	config3 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "3"}, Name: `Config 3`}}
	allConfigs := []*model.ConfigWithRows{config1, config2, config3}

	// Run
	_, err := dialog.SelectConfig(allConfigs, `LABEL`)
	assert.ErrorContains(t, err, "please specify config")
}

func TestSelectConfigMissing(t *testing.T) {
	t.Parallel()

	// Dependencies
	dialog, _, _ := createDialogs(t, false)

	// All configs
	config1 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "1"}, Name: `Config 1`}}
	config2 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "2"}, Name: `Config 2`}}
	config3 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "3"}, Name: `Config 3`}}
	allConfigs := []*model.ConfigWithRows{config1, config2, config3}

	// Run
	out, err := dialog.SelectConfig(allConfigs, `LABEL`)
	assert.Nil(t, out)
	assert.Error(t, err)
	assert.Equal(t, `please specify config`, err.Error())
}

func TestSelectConfigsInteractive(t *testing.T) {
	t.Parallel()

	// Dependencies
	dialog, _, console := createDialogs(t, true)

	// All configs
	config1 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "1"}, Name: `Config 1`}}
	config2 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "2"}, Name: `Config 2`}}
	config3 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "3"}, Name: `Config 3`}}
	config4 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "4"}, Name: `Config 4`}}
	config5 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "5"}, Name: `Config 5`}}
	allConfigs := []*model.ConfigWithRows{config1, config2, config3, config4, config5}

	// Interaction
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		assert.NoError(t, console.ExpectString("LABEL:"))

		assert.NoError(t, console.ExpectString("Config 1 (foo.bar:1)"))

		assert.NoError(t, console.ExpectString("Config 2 (foo.bar:2)"))

		assert.NoError(t, console.ExpectString("Config 3 (foo.bar:3)"))

		assert.NoError(t, console.ExpectString("Config 4 (foo.bar:4)"))

		assert.NoError(t, console.ExpectString("Config 5 (foo.bar:5)"))

		assert.NoError(t, console.SendDownArrow()) // -> Config 2

		assert.NoError(t, console.SendSpace()) // -> select

		assert.NoError(t, console.SendDownArrow()) // -> Config 3

		assert.NoError(t, console.SendDownArrow()) // -> Config 4

		assert.NoError(t, console.SendSpace()) // -> select

		assert.NoError(t, console.SendEnter()) // -> confirm

		assert.NoError(t, console.ExpectEOF())
	}()

	// Run
	out, err := dialog.SelectConfigs(allConfigs, `LABEL`)
	assert.Equal(t, []*model.ConfigWithRows{config2, config4}, out)
	assert.NoError(t, err)

	// Close terminal
	assert.NoError(t, console.Tty().Close())
	wg.Wait()
	assert.NoError(t, console.Close())
}

func TestSelectConfigsByFlag(t *testing.T) {
	t.Parallel()

	// Dependencies
	dialog, o, _ := createDialogs(t, false)
	o.Set(`configs`, `foo.bar:2, foo.bar:4`)

	// All configs
	config1 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "1"}, Name: `Config 1`}}
	config2 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "2"}, Name: `Config 2`}}
	config3 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "3"}, Name: `Config 3`}}
	config4 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "4"}, Name: `Config 4`}}
	config5 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "5"}, Name: `Config 5`}}
	allConfigs := []*model.ConfigWithRows{config1, config2, config3, config4, config5}

	// Run
	out, err := dialog.SelectConfigs(allConfigs, `LABEL`)
	assert.Equal(t, []*model.ConfigWithRows{config2, config4}, out)
	assert.NoError(t, err)
}

func TestSelectConfigsMissing(t *testing.T) {
	t.Parallel()

	// Dependencies
	dialog, _, _ := createDialogs(t, false)

	// All configs
	config1 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "1"}, Name: `Config 1`}}
	config2 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "2"}, Name: `Config 2`}}
	config3 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "3"}, Name: `Config 3`}}
	config4 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "4"}, Name: `Config 4`}}
	config5 := &model.ConfigWithRows{Config: &model.Config{ConfigKey: model.ConfigKey{BranchID: 1, ComponentID: `foo.bar`, ID: "5"}, Name: `Config 5`}}
	allConfigs := []*model.ConfigWithRows{config1, config2, config3, config4, config5}

	// Run
	out, err := dialog.SelectConfigs(allConfigs, `LABEL`)
	assert.Nil(t, out)
	assert.Error(t, err)
	assert.Equal(t, `please specify at least one config`, err.Error())
}
