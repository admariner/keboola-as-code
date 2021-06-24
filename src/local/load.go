package local

import (
	"keboola-as-code/src/json"
	"keboola-as-code/src/manifest"
	"keboola-as-code/src/model"
	"keboola-as-code/src/utils"
	"reflect"
)

// LoadModel from manifest and disk
func LoadModel(projectDir string, record manifest.Record, target interface{}) error {
	errors := &utils.Error{}

	// Load values from meta file
	errPrefix := record.Kind().Name + " metadata"
	if err := utils.ReadTaggedFields(projectDir, record.MetaFilePath(), model.MetaFileTag, errPrefix, target); err != nil {
		errors.Add(err)
	}

	// Load config file content
	errPrefix = record.Kind().Name
	if configField := utils.GetOneFieldWithTag(model.ConfigFileTag, target); configField != nil {
		content := utils.NewOrderedMap()
		modelValue := reflect.ValueOf(target).Elem()
		if err := json.ReadFile(projectDir, record.ConfigFilePath(), &content, errPrefix); err == nil {
			modelValue.FieldByName(configField.Name).Set(reflect.ValueOf(content))
		} else {
			errors.Add(err)
		}
	}

	if errors.Len() > 0 {
		return errors
	}

	return nil
}

func LoadBranch(projectDir string, b *manifest.BranchManifest) (*model.Branch, error) {
	branch := &model.Branch{BranchKey: b.BranchKey}
	if err := LoadModel(projectDir, b, branch); err != nil {
		return nil, err
	}
	return branch, nil
}

func LoadConfig(projectDir string, c *manifest.ConfigManifest) (*model.Config, error) {
	config := &model.Config{ConfigKey: c.ConfigKey}
	if err := LoadModel(projectDir, c, config); err != nil {
		return nil, err
	}
	return config, nil
}

func LoadConfigRow(projectDir string, r *manifest.ConfigRowManifest) (*model.ConfigRow, error) {
	row := &model.ConfigRow{ConfigRowKey: r.ConfigRowKey}
	if err := LoadModel(projectDir, r, row); err != nil {
		return nil, err
	}
	return row, nil
}
