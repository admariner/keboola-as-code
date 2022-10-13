package diffop

import (
	"context"
	"fmt"
	"sort"

	"github.com/keboola/keboola-as-code/internal/pkg/diff"
	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/state/local"
	"github.com/keboola/keboola-as-code/internal/pkg/state/remote"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

// Plan is based on the diff results.
type Plan struct {
	name                string
	actions             []*action
	allowedRemoteDelete bool
}

func NewPlan(name string) *Plan {
	return &Plan{name: name}
}

func (p *Plan) Empty() bool {
	return len(p.actions) == 0
}

func (p *Plan) Name() string {
	return p.name
}

func (p *Plan) AllowRemoteDelete() {
	p.allowedRemoteDelete = true
}

func (p *Plan) Invoke(logger log.Logger, ctx context.Context, localManager *local.Manager, remoteManager *remote.Manager, changeDescription string) error {
	executor := newExecutor(p, logger, ctx, localManager, remoteManager, changeDescription)
	return executor.invoke()
}

func (p *Plan) Log(logger log.Logger) {
	writer := logger.InfoWriter()
	writer.WriteString(fmt.Sprintf(`Plan for "%s" operation:`, p.Name()))
	actions := p.actions
	sort.SliceStable(actions, func(i, j int) bool {
		return actions[i].Path() < actions[j].Path()
	})

	if len(actions) == 0 {
		writer.WriteStringIndent(1, "no difference")
	} else {
		skippedDeleteCount := 0
		for _, action := range actions {
			msg := action.String()
			if !p.allowedRemoteDelete &&
				(action.action == ActionDeleteRemote) {
				msg += " - SKIPPED"
				skippedDeleteCount++
			}
			writer.WriteStringIndent(1, msg)
		}

		if skippedDeleteCount > 0 {
			writer.WriteString("Skipped remote objects deletion, use \"--force\" to delete them.")
		}
	}
}

func (p *Plan) Validate() error {
	errs := errors.NewMultiError()
	for _, action := range p.actions {
		if err := action.validate(); err != nil {
			errs.Append(err)
		}
	}

	if errs.Len() > 0 {
		return errors.PrefixErrorf(errs, `cannot perform the "%s" operation`, p.Name())
	}

	return nil
}

func (p *Plan) Add(result *diff.Result, actionType ActionType) {
	p.actions = append(p.actions, &action{Result: result, action: actionType})
}
