package interactive_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/service/cli"
	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/prompt"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/testhelper/terminal"
)

func TestPrompt_Select(t *testing.T) {
	t.Parallel()

	// Create virtual console
	console, err := terminal.New(t)
	assert.NoError(t, err)
	p := cli.NewPrompt(console.Tty(), console.Tty(), console.Tty())

	// Interaction
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		assert.NoError(t, console.ExpectString("My Select"))

		assert.NoError(t, console.SendEnter()) // enter - default value

		assert.NoError(t, console.ExpectEOF())
	}()

	// Show select
	result, ok := p.Select(&prompt.Select{
		Label:      "My Select",
		Options:    []string{"value1", "value2", "value3"},
		UseDefault: true,
		Default:    "value2",
	})
	assert.NoError(t, console.Tty().Close())
	wg.Wait()
	assert.NoError(t, console.Close())

	// Assert
	assert.True(t, ok)
	assert.Equal(t, "value2", result)
}

func TestPrompt_SelectIndex(t *testing.T) {
	t.Parallel()

	// Create virtual console
	console, err := terminal.New(t)
	assert.NoError(t, err)
	p := cli.NewPrompt(console.Tty(), console.Tty(), console.Tty())

	// Interaction
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		assert.NoError(t, console.ExpectString("My Select"))

		assert.NoError(t, console.SendEnter()) // enter - default value

		assert.NoError(t, console.ExpectEOF())
	}()

	// Show select
	result, ok := p.SelectIndex(&prompt.SelectIndex{
		Label:      "My Select",
		Options:    []string{"value1", "value2", "value3"},
		UseDefault: true,
		Default:    1,
	})
	assert.NoError(t, console.Tty().Close())
	wg.Wait()
	assert.NoError(t, console.Close())

	// Assert
	assert.True(t, ok)
	assert.Equal(t, 1, result)
}

func TestPrompt_MultiSelect(t *testing.T) {
	t.Parallel()

	// Create virtual console
	console, err := terminal.New(t)
	assert.NoError(t, err)
	p := cli.NewPrompt(console.Tty(), console.Tty(), console.Tty())

	// Interaction
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		assert.NoError(t, console.ExpectString("My Select"))

		assert.NoError(t, console.SendEnter()) // enter - default value

		assert.NoError(t, console.ExpectEOF())
	}()

	// Show select
	result, ok := p.MultiSelect(&prompt.MultiSelect{
		Label:   "My Select",
		Options: []string{"value1", "value2", "value3"},
		Default: []string{"value1", "value3"},
	})
	assert.NoError(t, console.Tty().Close())
	wg.Wait()
	assert.NoError(t, console.Close())

	// Assert
	assert.True(t, ok)
	assert.Equal(t, []string{"value1", "value3"}, result)
}

func TestPrompt_MultiSelectIndex(t *testing.T) {
	t.Parallel()

	// Create virtual console
	console, err := terminal.New(t)
	assert.NoError(t, err)
	p := cli.NewPrompt(console.Tty(), console.Tty(), console.Tty())

	// Interaction
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		assert.NoError(t, console.ExpectString("My Select"))

		assert.NoError(t, console.SendEnter()) // enter - default value

		assert.NoError(t, console.ExpectEOF())
	}()

	// Show select
	result, ok := p.MultiSelectIndex(&prompt.MultiSelectIndex{
		Label:   "My Select",
		Options: []string{"value1", "value2", "value3"},
		Default: []int{0, 2},
	})
	assert.NoError(t, console.Tty().Close())
	wg.Wait()
	assert.NoError(t, console.Close())

	// Assert
	assert.True(t, ok)
	assert.Equal(t, []int{0, 2}, result)
}

func TestPrompt_ShowLeaveBlank(t *testing.T) {
	t.Parallel()

	// Create virtual console
	console, err := terminal.New(t)
	assert.NoError(t, err)
	p := cli.NewPrompt(console.Tty(), console.Tty(), console.Tty())

	// Interaction
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		assert.NoError(t, console.ExpectString("My input"))

		assert.NoError(t, console.ExpectString("Leave blank for default value."))

		assert.NoError(t, console.SendEnter()) // enter - default value

		assert.NoError(t, console.ExpectEOF())
	}()

	// Show select
	result, ok := p.Ask(&prompt.Question{
		Label:       "Default",
		Description: "My input",
		Help:        "help",
		Hidden:      true,
		Default:     "default",
	})
	assert.NoError(t, console.Tty().Close())
	wg.Wait()
	assert.NoError(t, console.Close())

	// Assert
	assert.True(t, ok)
	assert.Equal(t, "default", result)
}

func TestPrompt_HideLeaveBlank(t *testing.T) {
	t.Parallel()

	// Create virtual console
	console, err := terminal.New(t)
	assert.NoError(t, err)
	p := cli.NewPrompt(console.Tty(), console.Tty(), console.Tty())

	// Interaction
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		assert.NoError(t, console.ExpectString("My input"))

		assert.NoError(t, console.SendEnter()) // enter - default value

		assert.NoError(t, console.ExpectEOF())

		assert.NotContains(t, console.String(), "Leave blank for default value.")
	}()

	// Show select
	result, ok := p.Ask(&prompt.Question{
		Label:       "Default",
		Description: "My input",
		Help:        "help",
		Hidden:      true,
	})
	assert.NoError(t, console.Tty().Close())
	wg.Wait()
	assert.NoError(t, console.Close())

	// Assert
	assert.True(t, ok)
	assert.Equal(t, "", result)
}
