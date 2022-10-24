package local

import (
	"github.com/keboola/go-utils/pkg/orderedmap"

	"github.com/keboola/keboola-as-code/internal/pkg/encoding/json/schema"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

func (m *Manager) createObject(key model.Key, name string) (model.Object, error) {
	switch k := key.(type) {
	case model.ConfigKey:
		component, err := m.state.Components().GetOrErr(k.ComponentId)
		if err != nil {
			return nil, err
		}
		content, err := generateContent(component.Schema, component.EmptyConfig)
		if err != nil {
			return nil, err
		}
		config := &model.Config{
			ConfigKey: k,
			Name:      name,
			Content:   content,
		}
		if component.IsTransformation() {
			config.Transformation = &model.Transformation{}
		}
		if component.IsOrchestrator() {
			config.Orchestration = &model.Orchestration{}
		}
		return config, nil
	case model.ConfigRowKey:
		component, err := m.state.Components().GetOrErr(k.ComponentId)
		if err != nil {
			return nil, err
		}
		content, err := generateContent(component.SchemaRow, component.EmptyConfigRow)
		if err != nil {
			return nil, err
		}
		return &model.ConfigRow{
			ConfigRowKey: k,
			Name:         name,
			Content:      content,
		}, nil
	default:
		panic(errors.Errorf(`unexpected type "%T"`, key))
	}
}

func generateContent(schemaDef []byte, defaultConfig *orderedmap.OrderedMap) (*orderedmap.OrderedMap, error) {
	// Use default configuration if defined in the component's metadata
	if len(defaultConfig.Keys()) > 0 {
		return defaultConfig, nil
	}

	// Otherwise, generate configuration from the JSON schema
	return schema.GenerateDocument(schemaDef)
}
