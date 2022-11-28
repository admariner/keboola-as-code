package load

import (
	"context"

	"github.com/keboola/go-client/pkg/client"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/keboola/keboola-as-code/internal/pkg/encoding/json"
	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/state"
	"github.com/keboola/keboola-as-code/internal/pkg/telemetry"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

type InvalidRemoteStateError struct {
	error
}

type InvalidLocalStateError struct {
	error
}

func (e InvalidRemoteStateError) Unwrap() error {
	return e.error
}

func (e InvalidRemoteStateError) WriteError(w errors.Writer, level int, trace errors.StackTrace) {
	// Write underlying error
	w.WriteErrorLevel(level, e.error, trace)
}

func (e InvalidLocalStateError) Unwrap() error {
	return e.error
}

func (e InvalidLocalStateError) WriteError(w errors.Writer, level int, trace errors.StackTrace) {
	// Write underlying error
	w.WriteErrorLevel(level, e.error, trace)
}

type Options struct {
	LoadLocalState          bool
	LoadRemoteState         bool
	IgnoreNotFoundErr       bool
	IgnoreInvalidLocalState bool
}

type OptionsWithFilter struct {
	Options
	LocalFilter  *model.ObjectsFilter
	RemoteFilter *model.ObjectsFilter
}

func InitOptions(pull bool) Options {
	return Options{
		LoadLocalState:          true,
		LoadRemoteState:         pull,
		IgnoreNotFoundErr:       false,
		IgnoreInvalidLocalState: false,
	}
}

func DiffOptions() Options {
	return Options{
		LoadLocalState:          true,
		LoadRemoteState:         true,
		IgnoreNotFoundErr:       false,
		IgnoreInvalidLocalState: false,
	}
}

func PullOptions(force bool) Options {
	return Options{
		LoadLocalState:          true,
		LoadRemoteState:         true,
		IgnoreNotFoundErr:       false,
		IgnoreInvalidLocalState: force,
	}
}

func PushOptions() Options {
	return Options{
		LoadLocalState:          true,
		LoadRemoteState:         true,
		IgnoreNotFoundErr:       false,
		IgnoreInvalidLocalState: false,
	}
}

func PersistOptions() Options {
	return Options{
		LoadLocalState:          true,
		LoadRemoteState:         false,
		IgnoreNotFoundErr:       true,
		IgnoreInvalidLocalState: false,
	}
}

func LocalOperationOptions() Options {
	return Options{
		LoadLocalState:          true,
		LoadRemoteState:         false,
		IgnoreNotFoundErr:       false,
		IgnoreInvalidLocalState: false,
	}
}

type dependencies interface {
	Tracer() trace.Tracer
	Logger() log.Logger
	Components() *model.ComponentsMap
	StorageAPIClient() client.Sender
}

func Run(ctx context.Context, container state.ObjectsContainer, o OptionsWithFilter, d dependencies) (s *state.State, err error) {
	ctx, span := d.Tracer().Start(ctx, "kac.lib.operation.state.load")
	span.SetAttributes(attribute.Bool("remote.load", o.LoadRemoteState))
	span.SetAttributes(attribute.String("remote.filter", json.MustEncodeString(o.RemoteFilter, false)))
	span.SetAttributes(attribute.Bool("local.load", o.LoadLocalState))
	span.SetAttributes(attribute.String("local.filter", json.MustEncodeString(o.LocalFilter, false)))
	defer telemetry.EndSpan(span, &err)

	logger := d.Logger()
	loadOptions := state.LoadOptions{
		LoadLocalState:    o.LoadLocalState,
		LoadRemoteState:   o.LoadRemoteState,
		IgnoreNotFoundErr: o.IgnoreNotFoundErr,
		LocalFilter:       o.LocalFilter,
		RemoteFilter:      o.RemoteFilter,
	}

	// Create state
	projectState, err := state.New(ctx, container, d)
	if err != nil {
		return nil, err
	}

	// Load objects
	ok, localErr, remoteErr := projectState.Load(ctx, loadOptions)
	if ok {
		logger.Debugf("Project state has been successfully loaded.")
	} else {
		if remoteErr != nil {
			return nil, InvalidRemoteStateError{errors.PrefixError(remoteErr, "cannot load project remote state")}
		}
		if localErr != nil {
			if o.IgnoreInvalidLocalState {
				logger.Info(`Ignoring invalid local state.`)
			} else {
				return nil, InvalidLocalStateError{errors.PrefixError(localErr, "project local state is invalid")}
			}
		}
	}

	return projectState, nil
}
