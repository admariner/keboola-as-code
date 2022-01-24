package configmetadata

import (
	"github.com/keboola/keboola-as-code/internal/pkg/model"
)

// MapBeforeLocalSave - store config metadata to manifest.
func (m *configMetadataMapper) MapBeforeLocalSave(recipe *model.LocalSaveRecipe) error {
	manifest, ok := recipe.ObjectManifest.(*model.ConfigManifest)
	if !ok {
		return nil
	}

	config, ok := recipe.Object.(*model.Config)
	if !ok {
		return nil
	}

	manifest.Metadata = config.MetadataOrderedMap()
	recipe.ChangedFields.Remove(`metadata`)
	return nil
}
