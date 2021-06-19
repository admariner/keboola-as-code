package diff

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"keboola-as-code/src/api"
	"keboola-as-code/src/model"
	"keboola-as-code/src/utils"
	"reflect"
)

type typeName string

type Differ struct {
	projectDir  string
	metadataDir string
	ctx         context.Context
	api         *api.StorageApi
	logger      *zap.SugaredLogger
	stateLoaded bool
	state       *model.State
	results     []*Result
	typeCache   map[typeName][]reflect.StructField
	error       *utils.Error
}

type ResultState int

const (
	ResultNotSet ResultState = iota
	ResultNotEqual
	ResultEqual
	ResultOnlyInRemote
	ResultOnlyInLocal
)

type Result struct {
	State   ResultState
	Changes map[string]string
}

type Results struct {
	Results []*Result
}

func NewDiffer(projectDir, metadataDir string, ctx context.Context, a *api.StorageApi, logger *zap.SugaredLogger) *Differ {
	d := &Differ{
		projectDir:  projectDir,
		metadataDir: metadataDir,
		ctx:         ctx,
		api:         a,
		logger:      logger,
		state:       model.NewState(projectDir),
	}
	return d
}

func (d *Differ) LoadState() error {
	grp, ctx := errgroup.WithContext(d.ctx)
	grp.Go(d.loadRemoteState(ctx))
	grp.Go(d.loadLocalState())
	err := grp.Wait()
	if err == nil {
		d.stateLoaded = true
	}
	return err
}

func (d *Differ) Diff() (*Results, error) {
	if !d.stateLoaded {
		panic("LoadState() must be called before Diff()")
	}

	// Diff all states
	d.results = []*Result{}
	d.error = &utils.Error{}
	for _, objectState := range d.state.All() {
		result, err := d.diffObjectState(objectState)
		if err != nil {
			d.error.Add(err)
		} else {
			d.results = append(d.results, result)
		}
	}
	// Check errors
	var err error
	if d.error.Len() > 0 {
		err = fmt.Errorf("%s", d.error)
	}

	return &Results{d.results}, err
}

func (d *Differ) diffObjectState(state model.ObjectState) (*Result, error) {
	// Validate
	result := &Result{}
	remote := state.RemoteState()
	local := state.LocalState()
	if remote == nil && local == nil {
		panic(fmt.Errorf("both local and remote state are not set"))
	}
	if remote == nil {
		result.State = ResultOnlyInLocal
		return result, nil
	}
	if local == nil {
		result.State = ResultOnlyInLocal
		return result, nil
	}

	// Types must be same
	remoteType := reflect.TypeOf(remote)
	localType := reflect.TypeOf(local)
	if remoteType.String() != localType.String() {
		panic(fmt.Errorf("local(%s) and remote(%s) states must have same data type", remoteType, localType))
	}

	// All available fields for diff
	diffFields := d.diffFieldsForType(remoteType)
	if len(diffFields) == 0 {
		return nil, fmt.Errorf(`no field with tag "diff:true" in struct "%s"`, remoteType.String())
	}

	// Diff all diffFields
	result.Changes = make(map[string]string)
	remoteValues := reflect.ValueOf(remote).Elem()
	localValues := reflect.ValueOf(remote).Elem()
	for _, field := range diffFields {
		difference := cmp.Diff(
			remoteValues.FieldByName(field.Name).Interface(),
			localValues.FieldByName(field.Name).Interface(),
		)
		if len(difference) > 0 {
			result.Changes[field.Name] = difference
		}
	}
	if len(result.Changes) > 0 {
		result.State = ResultEqual
	} else {
		result.State = ResultNotEqual
	}

	return result, nil
}

func (d *Differ) diffFieldsForType(t reflect.Type) []reflect.StructField {
	if v, ok := d.typeCache[typeName(t.Name())]; ok {
		return v
	} else {
		var diffFields []reflect.StructField
		numFields := t.NumField()
		for i := 0; i < numFields; i++ {
			fieldType := t.Field(i)
			tag := fieldType.Tag.Get("diff")
			if tag == "true" {
				diffFields = append(diffFields, fieldType)
			}
		}
		d.typeCache[typeName(t.Name())] = diffFields
		return diffFields
	}
}

func (d *Differ) loadRemoteState(ctx context.Context) func() error {
	return func() error {
		d.logger.Debugf("Loading project remote state.")
		remoteErrors := d.api.LoadRemoteState(d.state, ctx)
		if remoteErrors.Len() > 0 {
			d.logger.Debugf("Project remote state load failed: %s", remoteErrors)
			return fmt.Errorf("cannot load project remote state: %s", remoteErrors)
		} else {
			d.logger.Debugf("Project remote state successfully loaded.")
		}
		return nil
	}
}

func (d *Differ) loadLocalState() func() error {
	return func() error {
		d.logger.Debugf("Loading project local state.")
		localErrors := model.LoadLocalState(d.state, d.projectDir, d.metadataDir)
		if localErrors.Len() > 0 {
			d.logger.Debugf("Project local state load failed: %s", localErrors)
			return fmt.Errorf("cannot load project local state: %s", localErrors)
		} else {
			d.logger.Debugf("Project local state successfully loaded.")
		}
		return nil
	}
}
