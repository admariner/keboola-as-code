package validaterow

import (
	"bytes"
	"context"

	"github.com/keboola/go-client/pkg/keboola"
	"github.com/keboola/go-utils/pkg/orderedmap"
	"go.opentelemetry.io/otel/trace"

	"github.com/keboola/keboola-as-code/internal/pkg/encoding/json/schema"
	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/telemetry"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

type Options struct {
	ComponentID keboola.ComponentID
	RowPath     string
}

type dependencies interface {
	Logger() log.Logger
	Tracer() trace.Tracer
	Fs() filesystem.Fs
	Components() *model.ComponentsMap
}

func Run(ctx context.Context, o Options, d dependencies) (err error) {
	ctx, span := d.Tracer().Start(ctx, "kac.lib.operation.project.local.validate.row")
	defer telemetry.EndSpan(span, &err)
	logger := d.Logger()

	// Get component
	component, err := d.Components().GetOrErr(o.ComponentID)
	if err != nil {
		return err
	}

	// Read file
	fs := d.Fs()
	f, err := fs.FileLoader().ReadJSONFile(filesystem.NewFileDef(filesystem.Join(fs.WorkingDir(), o.RowPath)))
	if err != nil {
		return err
	}

	// File cannot be empty
	if v, ok := f.Content.GetOrNil("parameters").(*orderedmap.OrderedMap); !ok || len(v.Keys()) == 0 {
		return errors.Errorf("configuration row is empty")
	}

	// Validate
	if len(component.SchemaRow) == 0 || bytes.Equal(component.SchemaRow, []byte("{}")) {
		logger.Warnf(`Component "%s" has no configuration row JSON schema.`, component.ID)
	} else if err := schema.ValidateContent(component.SchemaRow, f.Content); err != nil {
		return err
	}

	logger.Info("Validation done.")
	return nil
}
