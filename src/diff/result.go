package diff

import (
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"keboola-as-code/src/api"
	"keboola-as-code/src/client"
)

type ResultState int

const (
	ResultNotSet ResultState = iota
	ResultNotEqual
	ResultEqual
	ResultOnlyInRemote
	ResultOnlyInLocal
)

type Result interface {
	State() ResultState
	Changes() []string
	LocalPath() string
	SaveRemote(logger *zap.SugaredLogger, workers *errgroup.Group, pool *client.Pool, a *api.StorageApi) error
	DeleteRemote(logger *zap.SugaredLogger, workers *errgroup.Group, pool *client.Pool, a *api.StorageApi) error
	DeleteLocal(logger *zap.SugaredLogger, workers *errgroup.Group) error
	SaveLocal(logger *zap.SugaredLogger, workers *errgroup.Group) error
}

type Results struct {
	Results []Result
}

type resultData struct {
	state   ResultState
	changes []string
}

func (d resultData) State() ResultState {
	return d.state
}

func (d resultData) Changes() []string {
	return d.changes
}

type BranchDiff struct {
	resultData
	*BranchState
}
type ConfigDiff struct {
	resultData
	*ConfigState
}
type ConfigRowDiff struct {
	resultData
	*ConfigRowState
}

func (b *BranchDiff) LocalPath() string {
	return b.BranchState.Manifest.RelativePath
}

func (c *ConfigDiff) LocalPath() string {
	return c.ConfigState.Manifest.RelativePath
}

func (r *ConfigRowDiff) LocalPath() string {
	return r.ConfigRowState.Manifest.RelativePath
}
