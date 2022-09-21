package dialog_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/cli/prompt/interactive"
	"github.com/keboola/keboola-as-code/internal/pkg/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/env"
	fixtures "github.com/keboola/keboola-as-code/internal/pkg/fixtures/local"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/template"
	"github.com/keboola/keboola-as-code/internal/pkg/template/input"
	"github.com/keboola/keboola-as-code/internal/pkg/template/repository"
	createOp "github.com/keboola/keboola-as-code/pkg/lib/operation/template/local/test/create"
)

func TestAskCreateTemplateTestInteractive(t *testing.T) {
	t.Parallel()

	// Test dependencies
	dialog, console := createDialogs(t, true)
	d := dependencies.NewMockedDeps()
	addMockedObjectsResponses(d.MockedHttpTransport())

	// Set fake file editor
	dialog.Prompt.(*interactive.Prompt).SetEditor(`true`)

	// Prepare the template
	fs, err := fixtures.LoadFS("template-simple", env.Empty())
	assert.NoError(t, err)
	version, err := model.NewSemVersion("v0.0.1")
	assert.NoError(t, err)
	tmplRef := model.NewTemplateRef(model.TemplateRepository{Name: "my-repository"}, "tmpl1", version.String())
	versionRec := repository.VersionRecord{
		Version:     version,
		Description: "",
		Stable:      true,
		Components:  []string{},
		AbsPath:     model.NewAbsPath("", "v1"),
	}
	tmplRec := repository.TemplateRecord{
		Id:          tmplRef.TemplateId(),
		Name:        "Template 1",
		Description: "",
		AbsPath:     model.NewAbsPath("", "tmpl1"),
		Versions:    []repository.VersionRecord{versionRec},
	}

	tmpl, err := template.New(tmplRef, tmplRec, versionRec, fs, fs)
	assert.NoError(t, err)

	// Interaction
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		_, err := console.ExpectString("Default Group")
		assert.NoError(t, err)

		_, err = console.ExpectString("Default Step")
		assert.NoError(t, err)

		_, err = console.ExpectString("Default Step Description")
		assert.NoError(t, err)

		_, err = console.ExpectString("url description")
		assert.NoError(t, err)

		_, err = console.ExpectString("API URL:")
		assert.NoError(t, err)

		time.Sleep(20 * time.Millisecond)
		_, err = console.SendLine(`foo.bar.com`)
		assert.NoError(t, err)

		_, err = console.ExpectString(`Enter the name of the environment variable that will fill input "API Token". Note that it will get prefix KBC_SECRET_.`)
		assert.NoError(t, err)

		_, err = console.ExpectString("API Token:")
		assert.NoError(t, err)

		time.Sleep(20 * time.Millisecond)
		_, err = console.SendLine(`VAR1`)
		assert.NoError(t, err)

		_, err = console.ExpectEOF()
		assert.NoError(t, err)
	}()

	// Run
	d.Options().Set(`test-name`, `one`)
	opts, warnings, err := dialog.AskCreateTemplateTestOptions(tmpl, d.Options())
	assert.NoError(t, err)
	assert.NoError(t, console.Tty().Close())
	wg.Wait()
	assert.NoError(t, console.Close())

	// Assert
	assert.Equal(t, createOp.Options{
		TestName: "one",
		Inputs: input.Values{
			{
				Id:      "generic-url",
				Value:   "foo.bar.com",
				Skipped: false,
			},
			{
				Id:      "generic-token",
				Value:   "##KBC_SECRET_VAR1##",
				Skipped: false,
			},
		},
	}, opts)
	assert.Equal(t, []string{`Input "API Token" expects setting of environment variable "KBC_SECRET_VAR1".`}, warnings)
}