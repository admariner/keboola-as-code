package links_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/fixtures"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
)

func TestRemoteSaveTranWithSharedCode(t *testing.T) {
	t.Parallel()
	mapperInst, context, logs := createMapper(t)

	// Shared code config with rows
	sharedCodeKey, sharedCodeRowsKeys := fixtures.CreateSharedCode(t, context.State, context.Naming)

	// Create transformation with shared code
	transformation := createInternalTranWithSharedCode(t, sharedCodeKey, sharedCodeRowsKeys, context)

	// Invoke
	internalObject := transformation.Local
	apiObject := internalObject.Clone().(*model.Config)
	recipe := &model.RemoteSaveRecipe{
		InternalObject: internalObject,
		ApiObject:      apiObject,
		ObjectManifest: transformation.Manifest(),
	}
	assert.NoError(t, mapperInst.MapBeforeRemoteSave(recipe))
	assert.Empty(t, logs.String())

	// Config ID and rows ID are set in Content
	id, found := apiObject.Content.Get(model.SharedCodeIdContentKey)
	assert.True(t, found)
	assert.Equal(t, sharedCodeKey.Id.String(), id)
	rows, found := apiObject.Content.Get(model.SharedCodeRowsIdContentKey)
	assert.True(t, found)
	assert.Equal(t, []interface{}{sharedCodeRowsKeys[0].ObjectId(), sharedCodeRowsKeys[1].ObjectId()}, rows)
}
