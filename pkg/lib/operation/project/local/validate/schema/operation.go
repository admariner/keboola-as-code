package validateschema

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"github.com/keboola/keboola-as-code/internal/pkg/encoding/json/schema"
	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/telemetry"
)

type Options struct {
	SchemaPath string
	FilePath   string
}

type dependencies interface {
	Logger() log.Logger
	Tracer() trace.Tracer
	Fs() filesystem.Fs
}

func Run(ctx context.Context, o Options, d dependencies) (err error) {
	ctx, span := d.Tracer().Start(ctx, "kac.lib.operation.project.local.validate.schema")
	defer telemetry.EndSpan(span, &err)
	logger := d.Logger()

	// Read schema
	s, err := d.Fs().FileLoader().ReadRawFile(filesystem.NewFileDef(o.SchemaPath))
	if err != nil {
		return err
	}

	// Read file
	fs := d.Fs()
	f, err := d.Fs().FileLoader().ReadJSONFile(filesystem.NewFileDef(filesystem.Join(fs.WorkingDir(), o.FilePath)))
	if err != nil {
		return err
	}

	// Validate
	if err := schema.ValidateContent([]byte(s.Content), f.Content); err != nil {
		return err
	}

	logger.Info("Validation done.")
	return nil
}