package manifest

import (
	"fmt"
	"github.com/iancoleman/orderedmap"
	"keboola-as-code/src/json"
	"keboola-as-code/src/model"
	"keboola-as-code/src/utils"
	"keboola-as-code/src/validator"
	"path/filepath"
	"sync"
)

const (
	MetadataDir = ".keboola"
	FileName    = "manifest.json"
	MetaFile    = "meta.json"
	ConfigFile  = "config.json"
	RowsDir     = "rows"
	SortById    = "id"
	SortByPath  = "path"
)

type Manifest struct {
	*Content    `validate:"required,dive"` // content of the file, updated only on LoadManifest() and Save()
	ProjectDir  string                     `validate:"required"` // project root
	MetadataDir string                     `validate:"required"` // inside ProjectDir
	changed     bool                       // is content changed?
	records     orderedmap.OrderedMap      // common map for all: branches, configs and rows manifests
	lock        *sync.Mutex
}

type Content struct {
	Version  int                       `json:"version" validate:"required,min=1,max=1"`
	Project  *Project                  `json:"project" validate:"required"`
	SortBy   string                    `json:"sortBy" validate:"oneof=id path"`
	Naming   *LocalNaming              `json:"naming" validate:"required"`
	Branches []*BranchManifest         `json:"branches" validate:"dive"`
	Configs  []*ConfigManifestWithRows `json:"configurations" validate:"dive"`
}

type Record interface {
	Kind() model.Kind           // eg. branch, config, config row -> used in logs
	Key() model.Key             // unique key for map
	SortKey(sort string) string // unique key for sorting
	IsInvalid() bool            // invalid records are skipped on save
	GetPaths() Paths            // define the location of the files
	MetaFilePath() string       // path to the meta.json file
	ConfigFilePath() string     // path to the config.json file
}

type state struct {
	invalid bool
}

type Paths struct {
	Path       string `json:"path" validate:"required"`
	ParentPath string `json:"-"`
}

type Project struct {
	Id      int    `json:"id" validate:"required"`
	ApiHost string `json:"apiHost" validate:"required,hostname"`
}

type BranchManifest struct {
	state
	model.BranchKey
	Paths
}

type ConfigManifest struct {
	state
	model.ConfigKey
	Paths
}

type ConfigManifestWithRows struct {
	*ConfigManifest
	Rows []*ConfigRowManifest `json:"rows"`
}

type ConfigRowManifest struct {
	state
	model.ConfigRowKey
	Paths
}

func NewManifest(projectId int, apiHost string, projectDir, metadataDir string) (*Manifest, error) {
	m := newManifest(projectId, apiHost, projectDir, metadataDir)
	err := m.validate()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func newManifest(projectId int, apiHost string, projectDir, metadataDir string) *Manifest {
	return &Manifest{
		ProjectDir:  projectDir,
		MetadataDir: metadataDir,
		records:     *utils.NewOrderedMap(),
		Content: &Content{
			Version:  1,
			Project:  &Project{Id: projectId, ApiHost: apiHost},
			SortBy:   SortById,
			Naming:   DefaultNaming(),
			Branches: make([]*BranchManifest, 0),
			Configs:  make([]*ConfigManifestWithRows, 0),
		},
		lock: &sync.Mutex{},
	}
}

func LoadManifest(projectDir string, metadataDir string) (*Manifest, error) {
	// Exists?
	path := filepath.Join(metadataDir, FileName)
	if !utils.IsFile(path) {
		return nil, fmt.Errorf("manifest \"%s\" not found", utils.RelPath(projectDir, path))
	}

	// Read JSON file
	m := newManifest(0, "", projectDir, metadataDir)
	if err := json.ReadFile(metadataDir, FileName, &m.Content, "manifest"); err != nil {
		return nil, err
	}

	// Resolve parent path, set parent IDs, store in records
	branchById := make(map[int]*BranchManifest)
	for _, branch := range m.Content.Branches {
		branch.ResolveParentPath()
		branchById[branch.Id] = branch
		m.records.Set(branch.Key().String(), branch)
	}
	for _, configWithRows := range m.Content.Configs {
		config := configWithRows.ConfigManifest
		branch, found := branchById[config.BranchId]
		if !found {
			return nil, fmt.Errorf("branch \"%d\" not found in the manifest - referenced from the config \"%s:%s\" in \"%s\"", config.BranchId, config.ComponentId, config.Id, path)
		}
		config.ResolveParentPath(branch)
		m.records.Set(config.Key().String(), config)
		for _, row := range configWithRows.Rows {
			row.BranchId = config.BranchId
			row.ComponentId = config.ComponentId
			row.ConfigId = config.Id
			row.ResolveParentPath(config)
			m.records.Set(row.Key().String(), row)
		}
	}

	// Validate
	if err := m.validate(); err != nil {
		return nil, err
	}

	// Return
	return m, nil
}

func (m *Manifest) Save() error {
	// Sort record in manifest + ensure order of processing: branch, config, configRow
	m.records.Sort(func(a *orderedmap.Pair, b *orderedmap.Pair) bool {
		return a.Value().(Record).SortKey(m.SortBy) < b.Value().(Record).SortKey(m.SortBy)
	})

	// Convert records map to slices
	configsMap := make(map[string]*ConfigManifestWithRows)
	m.Content.Branches = make([]*BranchManifest, 0)
	m.Content.Configs = make([]*ConfigManifestWithRows, 0)
	for _, key := range m.records.Keys() {
		r, _ := m.records.Get(key)
		record := r.(Record)

		// Skip invalid
		if record.IsInvalid() {
			continue
		}

		switch v := record.(type) {
		case *BranchManifest:
			m.Content.Branches = append(m.Content.Branches, v)
		case *ConfigManifest:
			config := &ConfigManifestWithRows{v, make([]*ConfigRowManifest, 0)}
			configsMap[config.String()] = config
			m.Content.Configs = append(m.Content.Configs, config)
		case *ConfigRowManifest:
			config, found := configsMap[v.ConfigKey().String()]
			if !found {
				panic(fmt.Errorf(`config with key "%s" not found"`, v.ConfigKey().String()))
			}
			config.Rows = append(config.Rows, v)
		default:
			panic(fmt.Errorf(`unexpected type "%T"`, record))
		}
	}

	// Validate
	err := m.validate()
	if err != nil {
		return err
	}

	// Write JSON file
	return json.WriteFile(m.MetadataDir, FileName, m.Content, "manifest")
}

func (m *Manifest) IsChanged() bool {
	return m.changed
}

func (m *Manifest) Path() string {
	return filepath.Join(m.MetadataDir, FileName)
}

func (m *Manifest) validate() error {
	if err := validator.Validate(m); err != nil {
		return fmt.Errorf("manifest is not valid: %s", err)
	}
	return nil
}

func (m *Manifest) GetRecords() orderedmap.OrderedMap {
	return m.records
}

func (m *Manifest) MustGetRecord(key model.Key) Record {
	record, found := m.GetRecord(key)
	if !found {
		panic(fmt.Errorf("manifest record with key \"%s\"", key.String()))
	}
	return record
}

func (m *Manifest) GetRecord(key model.Key) (Record, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if r, found := m.records.Get(key.String()); found {
		return r.(Record), found
	}
	return nil, false
}

func (m *Manifest) CreateRecordFor(key model.Key) (record Record) {
	switch v := key.(type) {
	case model.BranchKey:
		record = &BranchManifest{BranchKey: v}
	case model.ConfigKey:
		record = &ConfigManifest{ConfigKey: v}
	case model.ConfigRowKey:
		record = &ConfigRowManifest{ConfigRowKey: v}
	default:
		panic(fmt.Errorf("unexpected type \"%s\"", key))
	}
	m.SetRecord(record)
	return record
}

func (m *Manifest) SetRecord(record Record) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.changed = true
	m.records.Set(record.Key().String(), record)
}

func (m *Manifest) DeleteRecord(value model.ValueWithKey) {
	m.DeleteRecordByKey(value.Key())
}

func (m *Manifest) DeleteRecordByKey(key model.Key) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.changed = true
	m.records.Delete(key.String())
}

func (o Paths) GetPaths() Paths {
	return o
}

func (o Paths) RelativePath() string {
	return filepath.Join(o.ParentPath, o.Path)
}

func (o Paths) AbsolutePath(projectDir string) string {
	return filepath.Join(projectDir, o.RelativePath())
}

func (o Paths) MetaFilePath() string {
	return filepath.Join(o.RelativePath(), MetaFile)
}

func (o Paths) ConfigFilePath() string {
	return filepath.Join(o.RelativePath(), ConfigFile)
}

func (b BranchManifest) Kind() model.Kind {
	return model.Kind{Name: "branch", Abbr: "B"}
}

func (c ConfigManifest) Kind() model.Kind {
	return model.Kind{Name: "config", Abbr: "C"}
}

func (r ConfigRowManifest) Kind() model.Kind {
	return model.Kind{Name: "config row", Abbr: "R"}
}

func (s *state) IsInvalid() bool {
	return s.invalid
}

func (s *state) SetInvalid() {
	s.invalid = true
}

func (b *BranchManifest) ResolveParentPath() {
	b.ParentPath = ""
}

func (c *ConfigManifest) ResolveParentPath(branchManifest *BranchManifest) {
	c.ParentPath = filepath.Join(branchManifest.ParentPath, branchManifest.Path)
}

func (r *ConfigRowManifest) ResolveParentPath(configManifest *ConfigManifest) {
	r.ParentPath = filepath.Join(configManifest.ParentPath, configManifest.Path, RowsDir)
}

func (b BranchManifest) SortKey(sort string) string {
	if sort == SortByPath {
		return fmt.Sprintf("01_branch_%s", b.RelativePath())
	} else {
		return b.BranchKey.String()
	}

}

func (c ConfigManifest) SortKey(sort string) string {
	if sort == SortByPath {
		return fmt.Sprintf("02_config_%s", c.RelativePath())
	} else {
		return c.ConfigKey.String()
	}

}

func (r ConfigRowManifest) SortKey(sort string) string {
	if sort == SortByPath {
		return fmt.Sprintf("03_row_%s", r.RelativePath())
	} else {
		return r.ConfigRowKey.String()
	}

}
