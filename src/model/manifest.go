package model

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"keboola-as-code/src/json"
	"keboola-as-code/src/utils"
	"keboola-as-code/src/validator"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

const (
	MetadataDir      = ".keboola"
	MetaFile         = "meta.json"
	ConfigFile       = "config.json"
	RowsDir          = "rows"
	ManifestFileName = "manifest.json"
)

type Manifest struct {
	Path        string            `json:"-"`
	ProjectDir  string            `json:"-" validate:"required"`
	MetadataDir string            `json:"-" validate:"required"`
	Version     int               `json:"version" validate:"required,min=1,max=1"`
	Project     *ProjectManifest  `json:"project" validate:"required"`
	Branches    []*BranchManifest `json:"branches"`
	Configs     []*ConfigManifest `json:"configurations"`
}

type ManifestX interface {
	Kind() string
	KindAbbr() string
}

type ManifestPath struct {
	Path       string `json:"path" validate:"required"`
	ParentPath string `validate:"required"` // generated, not in JSON
}

type ProjectManifest struct {
	Id      int    `json:"id" validate:"required,min=1"`
	ApiHost string `json:"apiHost" validate:"required,hostname"`
}

type BranchManifest struct {
	*ManifestPath
	Id int `json:"id" validate:"required,min=1"`
}

type ConfigManifest struct {
	*ManifestPath
	BranchId    int                  `json:"branchId" validate:"required"`
	ComponentId string               `json:"componentId" validate:"required"`
	Id          string               `json:"id" validate:"required,min=1"`
	Rows        []*ConfigRowManifest `json:"rows"`
}

type ConfigRowManifest struct {
	*ManifestPath
	Id          string `json:"id" validate:"required,min=1"`
	BranchId    int    `json:"-" validate:"required"` // generated, not in JSON
	ComponentId string `json:"-" validate:"required"` // generated, not in JSON
	ConfigId    string `json:"-" validate:"required"` // generated, not in JSON
}

func NewManifest(projectId int, apiHost string, projectDir, metadataDir string) (*Manifest, error) {
	m := &Manifest{
		ProjectDir:  projectDir,
		MetadataDir: metadataDir,
		Version:     1,
		Project:     &ProjectManifest{Id: projectId, ApiHost: apiHost},
		Branches:    make([]*BranchManifest, 0),
		Configs:     make([]*ConfigManifest, 0),
	}
	err := m.validate()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func LoadManifest(projectDir string, metadataDir string) (*Manifest, error) {
	// Exists?
	path := filepath.Join(metadataDir, ManifestFileName)
	if !utils.IsFile(path) {
		return nil, fmt.Errorf("manifest \"%s\" not found", utils.RelPath(projectDir, path))
	}

	// Load file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read from manifest \"%s\": %s", utils.RelPath(projectDir, path), err)
	}

	// Decode JSON
	m := &Manifest{ProjectDir: projectDir, MetadataDir: metadataDir}
	err = json.Decode(data, m)
	if err != nil {
		return nil, fmt.Errorf("manifest \"%s\" is not valid: %s", utils.RelPath(projectDir, path), err)
	}

	// Resolve paths and set parents
	branchById := make(map[int]*BranchManifest)
	for _, branch := range m.Branches {
		branch.ResolvePaths()
		branchById[branch.Id] = branch
	}
	for _, config := range m.Configs {
		branch, found := branchById[config.BranchId]
		if !found {
			return nil, fmt.Errorf("branch \"%d\" not found in manifest - referenced from the config \"%s:%s\" in \"%s\"", config.BranchId, config.ComponentId, config.Id, m.Path)
		}
		config.ResolvePaths(branch)
		for _, row := range config.Rows {
			row.BranchId = config.BranchId
			row.ComponentId = config.ComponentId
			row.ConfigId = config.Id
			row.ResolvePaths(config)
		}
	}

	// Set path
	m.Path = path

	// Validate
	err = m.validate()
	if err != nil {
		return nil, err
	}

	// Return
	return m, nil
}

func (m *Manifest) Save() error {
	// Set path
	m.Path = filepath.Join(m.MetadataDir, ManifestFileName)

	// Validate
	err := m.validate()
	if err != nil {
		return err
	}

	// Encode JSON
	data, err := json.Encode(m, true)
	if err != nil {
		return err
	}

	// Write file
	return os.WriteFile(m.Path, data, 0650)
}

func (m *Manifest) LoadModel(path *ManifestPath, model interface{}) error {
	modelType := reflect.TypeOf(model).Elem()
	modelValue := reflect.ValueOf(model).Elem()

	// Load meta file
	metaFields := utils.GetFieldsWithTag("metaFile", "true", modelType, model)
	if len(metaFields) > 0 {
		metadataContent := make(map[string]interface{})
		metadataFile := filepath.Join(path.AbsolutePath(m.ProjectDir), MetaFile)
		content, err := os.ReadFile(metadataFile)
		if err != nil {
			return err
		}
		if err := json.Decode(content, &metadataContent); err != nil {
			return err
		}
		for _, field := range metaFields {
			// Use JSON name if present
			name := field.Name
			jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
			if jsonName != "" {
				name = jsonName
			}

			// Check if present
			value, ok := metadataContent[name]
			if !ok {
				return fmt.Errorf("missing field \"%s\" in metadata file \"%s\"", name, utils.RelPath(m.ProjectDir, metadataFile))
			}

			// Set value
			modelValue.FieldByName(field.Name).Set(reflect.ValueOf(value))
		}
	}

	// Load config file
	configFields := utils.GetFieldsWithTag("configFile", "true", modelType, model)
	if len(configFields) > 1 {
		panic(fmt.Errorf("struct \"%s\" has multiple fields with tag `configFile:\"true\"`, but only one allowed", modelType))
	} else if len(configFields) == 1 {
		configFile := filepath.Join(path.RelativePath(), ConfigFile)
		if !utils.IsFile(configFile) {
			return fmt.Errorf("%s JSON file \"%s\" not found", "X", configFile)
		}

		v := make(map[string]interface{})
		if err := readJsonFile(m.ProjectDir, configFile, v); err != nil {
			return fmt.Errorf("%s JSON file \"%s\" is invalid: %s", "x", configFile, err)
		}
		modelValue.FieldByName(configFields[0].Name).Set(reflect.ValueOf(v))
	}

	return nil
}

func (m *Manifest) SaveModel(path *ManifestPath, model interface{}, logger *zap.SugaredLogger) error {
	// Write metadata file
	if metadata := m.toMetadataFile(model); metadata != nil {
		metadataFile := filepath.Join(path.AbsolutePath(m.ProjectDir), MetaFile)
		metadataJson, err := json.Encode(metadata, true)
		if err != nil {
			return err
		}
		if err := os.WriteFile(metadataFile, metadataJson, 0650); err != nil {
			return err
		}
		logger.Debugf("Saved \"%s\"", metadataFile)
	}

	// Write config file
	if configContent := m.toConfigFile(model); configContent != nil {
		configFile := filepath.Join(path.AbsolutePath(m.ProjectDir), ConfigFile)
		configJson, err := json.Encode(configContent, true)
		if err != nil {
			return err
		}
		if err := os.WriteFile(configFile, configJson, 0650); err != nil {
			return err
		}
		logger.Debugf("Saved \"%s\"", configFile)
	}

	return nil
}

func (m *Manifest) DeleteModel(path *ManifestPath, model interface{}, logger *zap.SugaredLogger) error {
	// Delete metadata file
	if metadata := m.toMetadataFile(model); metadata != nil {
		metadataFile := filepath.Join(path.AbsolutePath(m.ProjectDir), MetaFile)
		if err := os.Remove(metadataFile); err != nil {
			return err
		}
		logger.Debugf("Removed \"%s\"", metadataFile)
	}

	// Delete config file
	if configContent := m.toConfigFile(model); configContent != nil {
		configFile := filepath.Join(path.AbsolutePath(m.ProjectDir), ConfigFile)
		if err := os.Remove(configFile); err != nil {
			return err
		}
		logger.Debugf("Removed \"%s\"", configFile)
	}

	return nil
}

func (m *Manifest) validate() error {
	if err := validator.Validate(m); err != nil {
		return fmt.Errorf("manifest is not valid: %s", err)
	}
	return nil
}

func (m *Manifest) toMetadataFile(model interface{}) map[string]interface{} {
	target := make(map[string]interface{})
	modelType := reflect.TypeOf(model).Elem()
	modelValue := reflect.ValueOf(model).Elem()
	for _, field := range utils.GetFieldsWithTag("metaFile", "true", modelType, model) {
		// Use JSON name if present
		name := field.Name
		jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
		if jsonName != "" {
			name = jsonName
		}

		// Get field value
		target[name] = modelValue.FieldByName(field.Name).Interface()
	}
	return target
}

func (m *Manifest) toConfigFile(model interface{}) map[string]interface{} {
	modelType := reflect.TypeOf(model).Elem()
	modelValue := reflect.ValueOf(model).Elem()
	fields := utils.GetFieldsWithTag("configFile", "true", modelType, model)

	// Check number of fields
	if len(fields) > 1 {
		panic(fmt.Errorf("struct \"%s\" has multiple fields with tag `configFile:\"true\"`, but only one allowed", modelType))
	} else if len(fields) == 0 {
		return nil
	}

	// Ok, return map
	return modelValue.FieldByName(fields[0].Name).Interface().(map[string]interface{})
}

func (o *ManifestPath) RelativePath() string {
	return filepath.Join(o.ParentPath, o.Path)
}

func (o *ManifestPath) AbsolutePath(projectDir string) string {
	return filepath.Join(projectDir, o.RelativePath())
}

func (o *ManifestPath) MetadataFilePath() string {
	return filepath.Join(o.RelativePath(), MetaFile)
}

func (o *ManifestPath) ConfigFilePath() string {
	return filepath.Join(o.RelativePath(), ConfigFile)
}

func (b *BranchManifest) Kind() string {
	return "config"
}

func (c *ConfigManifest) Kind() string {
	return "config"
}

func (r *ConfigRowManifest) Kind() string {
	return "config row"
}

func (b *BranchManifest) KindAbbr() string {
	return "B"
}

func (c *ConfigManifest) KindAbbr() string {
	return "C"
}

func (r *ConfigRowManifest) KindAbbr() string {
	return "R"
}

func (b *BranchManifest) ResolvePaths() {
	b.ParentPath = ""
}

func (c *ConfigManifest) ResolvePaths(b *BranchManifest) {
	c.ParentPath = filepath.Join(b.ParentPath, b.Path)
}

func (r *ConfigRowManifest) ResolvePaths(c *ConfigManifest) {
	r.ParentPath = filepath.Join(c.ParentPath, c.Path, RowsDir)
}

func (b *BranchManifest) ToModel(m *Manifest) (*Branch, error) {
	branch := &Branch{Id: b.Id}
	if err := m.LoadModel(b.ManifestPath, branch); err != nil {
		return nil, err
	}
	return branch, nil
}

func (c *ConfigManifest) ToModel(m *Manifest) (*Config, error) {
	config := &Config{BranchId: c.BranchId, ComponentId: c.ComponentId, Id: c.Id, Rows: make([]*ConfigRow, 0)}
	if err := m.LoadModel(c.ManifestPath, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (r *ConfigRowManifest) ToModel(m *Manifest) (*ConfigRow, error) {
	row := &ConfigRow{BranchId: r.BranchId, ComponentId: r.ComponentId, ConfigId: r.ConfigId, Id: r.Id}
	if err := m.LoadModel(r.ManifestPath, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (b *Branch) GenerateManifest() *BranchManifest {
	manifest := &BranchManifest{Id: b.Id}
	if b.IsDefault {
		manifest.Path = "main"
	} else {
		manifest.Path = generatePath(cast.ToString(b.Id), b.Name)
	}
	manifest.ResolvePaths()
	return manifest
}

func (c *Config) GenerateManifest(b *BranchManifest, component *Component) *ConfigManifest {
	manifest := &ConfigManifest{BranchId: c.BranchId, ComponentId: c.ComponentId, Id: c.Id}
	manifest.Path = filepath.Join(component.Type, c.ComponentId, generatePath(c.Id, c.Name))
	manifest.ResolvePaths(b)
	return manifest
}

func (r *ConfigRow) GenerateManifest(c *ConfigManifest) *ConfigRowManifest {
	manifest := &ConfigRowManifest{BranchId: r.BranchId, ComponentId: r.ComponentId, ConfigId: r.ConfigId, Id: r.Id}
	manifest.Path = generatePath(r.Id, r.Name)
	manifest.ResolvePaths(c)
	return manifest
}

//func readConfigFile(kind, projectDir, relPath string, v interface{}) error {
//	path := filepath.Join(projectDir, relPath)
//	if !utils.IsFile(path) {
//		return fmt.Errorf("%s JSON file \"%s\" not found", kind, relPath)
//	}
//
//	if err := readJsonFile(projectDir, path, v); err != nil {
//		return fmt.Errorf("%s JSON file \"%s\" is invalid: %s", kind, relPath, err)
//	}
//	return nil
//}
//
//func readMetadataFile(kind, projectDir, relPath string, v interface{}) error {
//	path := filepath.Join(projectDir, relPath)
//	if !utils.IsFile(path) {
//		return fmt.Errorf("%s metadata JSON file \"%s\" not found", kind, relPath)
//	}
//
//	if err := readJsonFile(projectDir, path, v); err != nil {
//		return fmt.Errorf("%s metadata JSON file \"%s\" is invalid: %s", kind, relPath, err)
//	}
//	return nil
//}

func readJsonFile(projectDir string, path string, v interface{}) error {
	// Read meta file
	if !utils.IsFile(path) {
		return fmt.Errorf("file not found \"%s\"", utils.RelPath(projectDir, path))
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read file \"%s\"", utils.RelPath(projectDir, path))
	}

	// Decode meta file
	err = json.Decode(content, v)
	if err != nil {
		return err
	}
	return nil
}

func generatePath(id string, name string) string {
	name = regexp.
		MustCompile(`[^a-zA-Z0-9-]]`).
		ReplaceAllString(strcase.ToDelimited(name, '-'), "-")
	return id + "-" + name
}
