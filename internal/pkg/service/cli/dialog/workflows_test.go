package dialog_test

import (
	"sync"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/service/cli/cmd/ci"
	genWorkflows "github.com/keboola/keboola-as-code/pkg/lib/operation/project/local/workflows/generate"
)

func TestAskWorkflowsOptionsInteractive(t *testing.T) {
	t.Parallel()

	dialog, o, console := createDialogs(t, true)

	// Default values are defined by options
	flags := pflag.NewFlagSet(``, pflag.ExitOnError)
	ci.WorkflowsCmdFlags(flags)
	assert.NoError(t, o.BindPFlags(flags))

	// Interaction
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		assert.NoError(t, console.ExpectString(`Generate "validate" workflow?`))

		assert.NoError(t, console.SendLine(`n`)) // no

		assert.NoError(t, console.ExpectString(`Generate "push" workflow?`))

		assert.NoError(t, console.SendEnter()) // enter - yes

		assert.NoError(t, console.ExpectString(`Generate "pull" workflow?`))

		assert.NoError(t, console.SendLine(`n`)) // no

		assert.NoError(t, console.ExpectString(`Please select the main GitHub branch name:`))

		assert.NoError(t, console.SendEnter()) // enter - main

		assert.NoError(t, console.ExpectEOF())
	}()

	// Run
	out := dialog.AskWorkflowsOptions()
	assert.Equal(t, genWorkflows.Options{
		Validate:   false,
		Push:       true,
		Pull:       false,
		MainBranch: `main`,
	}, out)

	// Close terminal
	assert.NoError(t, console.Tty().Close())
	wg.Wait()
	assert.NoError(t, console.Close())
}

func TestAskWorkflowsOptionsByFlag(t *testing.T) {
	t.Parallel()

	dialog, o, _ := createDialogs(t, true)
	o.Set(`ci-validate`, `false`)
	o.Set(`ci-push`, `true`)
	o.Set(`ci-pull`, `false`)
	o.Set(`ci-main-branch`, `main`)

	// Default values are defined by options
	flags := pflag.NewFlagSet(``, pflag.ExitOnError)
	ci.WorkflowsCmdFlags(flags)
	assert.NoError(t, o.BindPFlags(flags))

	// Run
	out := dialog.AskWorkflowsOptions()
	assert.Equal(t, genWorkflows.Options{
		Validate:   false,
		Push:       true,
		Pull:       false,
		MainBranch: `main`,
	}, out)
}
