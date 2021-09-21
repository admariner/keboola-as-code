package state

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"go.uber.org/zap"

	"github.com/keboola/keboola-as-code/internal/pkg/local"
	"github.com/keboola/keboola-as-code/internal/pkg/manifest"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/remote"
	"github.com/keboola/keboola-as-code/internal/pkg/utils"
	"github.com/keboola/keboola-as-code/internal/pkg/validator"
)

// State - Local and Remote state of the project.
type State struct {
	*Options
	*model.State
	*model.PathsState
	mutex        *sync.Mutex
	localManager *local.Manager
	remoteErrors *utils.Error
	localErrors  *utils.Error
}

type Options struct {
	manifest             *manifest.Manifest
	api                  *remote.StorageApi
	context              context.Context
	logger               *zap.SugaredLogger
	LoadLocalState       bool
	LoadRemoteState      bool
	IgnoreMarkedToDelete bool // dev config/row marked with [TO DELETE] in the name will be ignored
	SkipNotFoundErr      bool // not found error will be ignored
}

func NewOptions(m *manifest.Manifest, api *remote.StorageApi, ctx context.Context, logger *zap.SugaredLogger) *Options {
	return &Options{
		manifest:             m,
		api:                  api,
		context:              ctx,
		logger:               logger,
		IgnoreMarkedToDelete: true,
	}
}

// LoadState - remote and local.
func LoadState(options *Options) (state *State, ok bool) {
	state = newState(options)

	// Token and manifest project ID must be same
	if state.manifest.Project.Id != state.api.ProjectId() {
		state.AddLocalError(fmt.Errorf("used token is from the project \"%d\", but it must be from the project \"%d\"", state.api.ProjectId(), state.manifest.Project.Id))
		return state, false
	}

	// Log allowed branches
	state.logger.Debugf(`Allowed branches: %s`, state.manifest.Content.AllowedBranches)

	if state.LoadRemoteState {
		state.logger.Debugf("Loading project remote state.")
		state.doLoadRemoteState()
	}

	if state.LoadLocalState {
		state.logger.Debugf("Loading local state.")
		state.doLoadLocalState()
	}

	state.validate()

	ok = state.LocalErrors().Len() == 0 && state.RemoteErrors().Len() == 0
	return state, ok
}

func newState(options *Options) *State {
	s := &State{
		Options:      options,
		mutex:        &sync.Mutex{},
		remoteErrors: utils.NewMultiError(),
		localErrors:  utils.NewMultiError(),
	}
	s.State = model.NewState(options.api.Components())
	s.PathsState = model.NewPathsState(s.manifest.ProjectDir, s.localErrors)
	s.localManager = local.NewManager(options.logger, options.manifest, s.api.Components())
	return s
}

func (s *State) Manifest() *manifest.Manifest {
	return s.manifest
}

func (s *State) ProjectDir() string {
	return s.manifest.ProjectDir
}

func (s *State) Naming() model.Naming {
	return s.manifest.Naming
}

func (s *State) LocalManager() *local.Manager {
	return s.localManager
}

func (s *State) UntrackedDirs() (dirs []string) {
	for _, path := range s.UntrackedPaths() {
		if !utils.IsDir(filepath.Join(s.manifest.ProjectDir, path)) {
			continue
		}
		dirs = append(dirs, path)
	}
	return dirs
}

func (s *State) LogUntrackedPaths(logger *zap.SugaredLogger) {
	untracked := s.UntrackedPaths()
	if len(untracked) > 0 {
		logger.Warn("Unknown paths found:")
		for _, path := range untracked {
			logger.Warn("\t- ", path)
		}
	}
}

func (s *State) RemoteErrors() *utils.Error {
	return s.remoteErrors
}

func (s *State) LocalErrors() *utils.Error {
	return s.localErrors
}

func (s *State) AddRemoteError(err error) {
	s.remoteErrors.Append(err)
}

func (s *State) AddLocalError(err error) {
	s.localErrors.Append(err)
}

func (s *State) All() []model.ObjectState {
	return s.State.All(s.manifest.SortBy)
}

func (s *State) Branches() (branches []*model.BranchState) {
	return s.State.Branches(s.manifest.SortBy)
}

func (s *State) Configs() (configs []*model.ConfigState) {
	return s.State.Configs(s.manifest.SortBy)
}

func (s *State) ConfigRows() (rows []*model.ConfigRowState) {
	return s.State.ConfigRows(s.manifest.SortBy)
}

func (s *State) SetRemoteState(remote model.Object) (model.ObjectState, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Skip ignored objects
	if s.manifest.IsObjectIgnored(remote) {
		return nil, nil
	}

	// Get or create state
	state, err := s.GetOrCreate(remote.Key())
	if err != nil {
		s.AddRemoteError(err)
		return nil, nil
	}

	state.SetRemoteState(remote)
	if !state.HasManifest() {
		// Generate manifest record
		m, _, err := s.manifest.CreateOrGetRecord(remote.Key())
		if err != nil {
			return nil, err
		}
		state.SetManifest(m)

		// Generate local path
		if err := s.localManager.UpdatePaths(state, false); err != nil {
			return nil, err
		}
	}
	return state, nil
}

func (s *State) SetLocalState(local model.Object, record model.Record) model.ObjectState {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	state, err := s.GetOrCreate(local.Key())
	if err != nil {
		s.AddLocalError(err)
		return nil
	}

	state.SetLocalState(local)
	state.SetManifest(record)
	for _, path := range record.GetRelatedPaths() {
		s.MarkTracked(path)
	}
	return state
}

func (s *State) validate() {
	for _, component := range s.Components().AllLoaded() {
		if err := validator.Validate(component); err != nil {
			s.AddLocalError(utils.PrefixError(fmt.Sprintf(`component \"%s\" is not valid`, component.Key()), err))
		}
	}
	for _, objectState := range s.All() {
		if objectState.HasRemoteState() {
			if err := validator.Validate(objectState.RemoteState()); err != nil {
				s.AddRemoteError(utils.PrefixError(fmt.Sprintf(`%s is not valid`, objectState.Desc()), err))
			}
		}

		if objectState.HasLocalState() {
			if err := validator.Validate(objectState.LocalState()); err != nil {
				s.AddLocalError(utils.PrefixError(fmt.Sprintf(`%s is not valid`, objectState.Desc()), err))
			}
		}
	}
}