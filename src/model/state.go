package model

type ObjectState interface {
	Level() int       // hierarchical level, "1" for branch, "2" for config, ...
	Key() Key         // unique key for all objects
	ObjectId() string // eg. config id
	Kind() Kind       // branch, config, ...
	RelativePath() string
	Manifest() Record
	HasLocalState() bool
	LocalState() Object
	HasRemoteState() bool
	RemoteState() Object
}

type RecordProvider interface {
}

type BranchState struct {
	*BranchManifest
	Remote *Branch `validate:"omitempty,dive"`
	Local  *Branch `validate:"omitempty,dive"`
}

type ConfigState struct {
	*ConfigManifest
	Component *Component `validate:"dive"`
	Remote    *Config    `validate:"omitempty,dive"`
	Local     *Config    `validate:"omitempty,dive"`
}

type ConfigRowState struct {
	*ConfigRowManifest
	Remote *ConfigRow `validate:"omitempty,dive"`
	Local  *ConfigRow `validate:"omitempty,dive"`
}

func (b *BranchState) HasLocalState() bool {
	return b.Local != nil
}

func (c *ConfigState) HasLocalState() bool {
	return c.Local != nil
}

func (r *ConfigRowState) HasLocalState() bool {
	return r.Local != nil
}

func (b *BranchState) LocalState() Object {
	return b.Local
}

func (c *ConfigState) LocalState() Object {
	return c.Local
}

func (r *ConfigRowState) LocalState() Object {
	return r.Local
}

func (b *BranchState) HasRemoteState() bool {
	return b.Remote != nil
}

func (c *ConfigState) HasRemoteState() bool {
	return c.Remote != nil
}

func (r *ConfigRowState) HasRemoteState() bool {
	return r.Remote != nil
}

func (b *BranchState) RemoteState() Object {
	return b.Remote
}

func (c *ConfigState) RemoteState() Object {
	return c.Remote
}

func (r *ConfigRowState) RemoteState() Object {
	return r.Remote
}

func (b *BranchState) Manifest() Record {
	return b.BranchManifest
}

func (c *ConfigState) Manifest() Record {
	return c.ConfigManifest
}

func (r *ConfigRowState) Manifest() Record {
	return r.ConfigRowManifest
}

func (b *BranchState) GetName() string {
	if b.Remote != nil {
		return b.Remote.Name
	}
	if b.Local != nil {
		return b.Local.Name
	}
	return ""
}

func (c *ConfigState) GetName() string {
	if c.Remote != nil {
		return c.Remote.Name
	}
	if c.Local != nil {
		return c.Local.Name
	}
	return ""
}

func (r *ConfigRowState) GetName() string {
	if r.Remote != nil {
		return r.Remote.Name
	}
	if r.Local != nil {
		return r.Local.Name
	}
	return ""
}