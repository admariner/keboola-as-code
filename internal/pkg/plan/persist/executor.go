package persist

import (
	"context"

	"github.com/keboola/keboola-sdk-go/v2/pkg/keboola"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/state"
	"github.com/keboola/keboola-as-code/internal/pkg/state/local"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

type executor struct {
	*Plan
	*state.State
	logger  log.Logger
	tickets *keboola.TicketProvider
	uow     *local.UnitOfWork
	errors  errors.MultiError
}

func newExecutor(ctx context.Context, logger log.Logger, keboolaProjectAPI *keboola.AuthorizedAPI, projectState *state.State, plan *Plan) *executor {
	return &executor{
		Plan:    plan,
		State:   projectState,
		logger:  logger,
		tickets: keboola.NewTicketProvider(ctx, keboolaProjectAPI),
		uow:     projectState.LocalManager().NewUnitOfWork(ctx),
		errors:  errors.NewMultiError(),
	}
}

func (e *executor) invoke() error {
	for _, action := range e.actions {
		switch a := action.(type) {
		case *newObjectAction:
			e.persistNewObject(a)
		case *deleteManifestRecordAction:
			objectState, _ := e.Get(a.Key())
			e.uow.DeleteObject(objectState, a.ObjectManifest)
		default:
			panic(errors.Errorf(`unexpected type "%T"`, action))
		}
	}

	// Let's wait until all new IDs are generated
	if err := e.tickets.Resolve(); err != nil {
		e.errors.Append(err)
	}

	// Wait for all local operations
	if err := e.uow.Invoke(); err != nil {
		e.errors.Append(err)
	}

	return e.errors.ErrorOrNil()
}

func (e *executor) persistNewObject(action *newObjectAction) {
	// Generate unique ID
	e.tickets.Request(func(ticket *keboola.Ticket) {
		key := action.Key

		// Set new id to the key
		switch k := key.(type) {
		case model.ConfigKey:
			k.ID = keboola.ConfigID(ticket.ID)
			key = k
		case model.ConfigRowKey:
			k.ID = keboola.RowID(ticket.ID)
			key = k
		default:
			panic(errors.Errorf(`unexpected type "%s" of the persisted object "%s"`, key.Kind(), key.Desc()))
		}

		// The parent was not persisted for some error -> skip
		if action.ParentKey != nil && action.ParentKey.ObjectID() == `` {
			return
		}

		// Create manifest record
		record, found, err := e.Manifest().CreateOrGetRecord(key)
		if err != nil {
			e.errors.Append(err)
			return
		} else if found {
			panic(errors.Errorf(`unexpected state: manifest record "%s" exists, but it should not`, record))
		}

		// Invoke mapper
		err = e.Mapper().MapBeforePersist(e.Ctx(), &model.PersistRecipe{
			ParentKey: action.ParentKey,
			Manifest:  record,
		})
		if err != nil {
			e.errors.Append(err)
			return
		}

		// Update parent path - may be affected by relations
		if err := e.Manifest().ResolveParentPath(record); err != nil {
			e.errors.Append(errors.Errorf(`cannot resolve path: %w`, err))
			return
		}

		// Set local path
		record.SetRelativePath(action.GetRelativePath())

		// Load model
		e.uow.LoadObject(record, model.NoFilter())

		// Save to manifest.json
		if err := e.Manifest().PersistRecord(record); err != nil {
			e.errors.Append(err)
			return
		}

		// Setup related objects
		action.InvokeOnPersist(key)
	})
}
