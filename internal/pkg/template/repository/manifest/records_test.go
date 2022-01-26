package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplateRecord_LatestVersion_NotFound(t *testing.T) {
	t.Parallel()
	template := TemplateRecord{}
	_, found := template.LatestVersion()
	assert.False(t, found)
}

func TestTemplateRecord_LatestVersion_Found(t *testing.T) {
	t.Parallel()
	template := TemplateRecord{
		Versions: []VersionRecord{
			{
				Version:     version(`0.0.1`),
				Stable:      false,
				Description: `Version 0`,
			},
			{
				Version:     version(`1.2.3`),
				Stable:      true,
				Description: `Version 1`,
			},
			{
				Version:     version(`0.1.0`),
				Stable:      false,
				Description: `Version 0.1`,
			},
		},
	}

	v, found := template.LatestVersion()
	assert.True(t, found)
	assert.Equal(t, VersionRecord{
		Version:     version(`1.2.3`),
		Stable:      true,
		Description: `Version 1`,
	}, v)
}

func TestTemplateRecord_GetByPath_NotFound(t *testing.T) {
	t.Parallel()
	r := &TemplateRecord{}
	value, found := r.GetByPath(`v1`)
	assert.Empty(t, value)
	assert.False(t, found)
}

func TestTemplateRecord_GetByPath_Found(t *testing.T) {
	t.Parallel()
	r := &TemplateRecord{}
	version1 := version(`v1.2.3`)
	r.AddVersion(version1)
	version2 := version(`v2.0.0`)
	r.AddVersion(version2)
	value, found := r.GetByPath(`v1`)
	assert.Equal(t, version1, value.Version)
	assert.True(t, found)
}