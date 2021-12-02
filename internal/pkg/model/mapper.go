package model

import (
	"go.uber.org/zap"

	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
)

type MapperContext struct {
	Logger *zap.SugaredLogger
	Fs     filesystem.Fs
	Naming *Naming
	State  *State
}

// LocalLoadRecipe - all items related to the object, when loading from local fs.
type LocalLoadRecipe struct {
	ObjectManifest                      // manifest record, eg *ConfigManifest
	Object         Object               // object, eg. Config
	Metadata       *filesystem.JsonFile // meta.json
	Configuration  *filesystem.JsonFile // config.json
	Description    *filesystem.File     // description.md
	ExtraFiles     []*filesystem.File   // extra files
}

// LocalSaveRecipe - all items related to the object, when saving to local fs.
type LocalSaveRecipe struct {
	ChangedFields  ChangedFields
	ObjectManifest                      // manifest record, eg *ConfigManifest
	Object         Object               // object, eg. Config
	Metadata       *filesystem.JsonFile // meta.json
	Configuration  *filesystem.JsonFile // config.json
	Description    *filesystem.File     // description.md
	ExtraFiles     []*filesystem.File   // extra files
	ToDelete       []string             // paths to delete, on save
}

// RemoteLoadRecipe - all items related to the object, when loading from Storage API.
type RemoteLoadRecipe struct {
	Manifest       ObjectManifest
	ApiObject      Object // eg. Config, original version, API representation
	InternalObject Object // eg. Config, modified version, internal representation
}

// RemoteSaveRecipe - all items related to the object, when saving to Storage API.
type RemoteSaveRecipe struct {
	ChangedFields  ChangedFields
	Manifest       ObjectManifest
	InternalObject Object // eg. Config, original version, internal representation
	ApiObject      Object // eg. Config, modified version, API representation
}

// PersistRecipe contains object to persist.
type PersistRecipe struct {
	ParentKey Key
	Manifest  ObjectManifest
}

type PathsGenerator interface {
	AddRenamed(path RenamedPath)
	RenameEnabled() bool // if true, existing paths will be renamed
}

// OnObjectPathUpdateEvent contains object with updated path.
type OnObjectPathUpdateEvent struct {
	PathsGenerator PathsGenerator
	ObjectState    ObjectState
	Renamed        bool
	OldPath        string
	NewPath        string
}
