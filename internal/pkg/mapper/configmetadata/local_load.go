package configmetadata

import (
	"github.com/keboola/keboola-as-code/internal/pkg/model"
)

// MapAfterLocalLoad - load metadata from manifest to config.
func (m *configMetadataMapper) MapAfterLocalLoad(recipe *model.LocalLoadRecipe) error {
	manifest, ok := recipe.ObjectManifest.(*model.ConfigManifest)
	if !ok {
		return nil
	}

	config, ok := recipe.Object.(*model.Config)
	if !ok {
		return nil
	}

	if manifest.Metadata != nil {
		config.Metadata = manifest.MetadataMap()
	}
	return nil
}
