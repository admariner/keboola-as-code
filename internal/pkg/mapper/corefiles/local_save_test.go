package corefiles_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/fixtures"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/orderedmap"
)

func TestSaveCoreFiles(t *testing.T) {
	t.Parallel()
	state, _ := createStateWithMapper(t)

	// Recipe
	manifest := &fixtures.MockedManifest{}
	object := &fixtures.MockedObject{
		Foo1:  "1",
		Foo2:  "2",
		Meta1: "3",
		Meta2: "4",
		Config: orderedmap.FromPairs([]orderedmap.Pair{
			{
				Key:   "foo",
				Value: "bar",
			},
		}),
	}
	recipe := model.NewLocalSaveRecipe(manifest, object, model.NewChangedFields())

	// No files
	assert.Empty(t, recipe.Files.All())

	// Call mapper
	assert.NoError(t, state.Mapper().MapBeforeLocalSave(recipe))

	// Files are generated
	expectedFiles := model.NewFilesToSave()
	expectedFiles.
		Add(
			filesystem.NewJsonFile(state.NamingGenerator().MetaFilePath(manifest.Path()),
				orderedmap.FromPairs([]orderedmap.Pair{
					{Key: "myKey", Value: "3"},
					{Key: "Meta2", Value: "4"},
				}),
			),
		).
		AddTag(model.FileTypeJson).
		AddTag(model.FileKindObjectMeta)
	expectedFiles.
		Add(
			filesystem.NewJsonFile(state.NamingGenerator().ConfigFilePath(manifest.Path()),
				orderedmap.FromPairs([]orderedmap.Pair{
					{Key: "foo", Value: "bar"},
				}),
			),
		).
		AddTag(model.FileTypeJson).
		AddTag(model.FileKindObjectConfig)
	assert.Equal(t, expectedFiles, recipe.Files)
}