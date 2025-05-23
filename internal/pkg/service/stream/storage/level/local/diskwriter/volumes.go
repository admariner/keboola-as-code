package diskwriter

import (
	"context"

	"github.com/jonboulle/clockwork"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/servicectx"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/level/local/events"
	volume "github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/level/local/volume/model"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/level/local/volume/opener"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

type Volumes struct {
	clock  clockwork.Clock
	logger log.Logger
	// events instance is passed to each volume and then to each writer
	events     *events.Events[Writer]
	collection *volume.Collection[*Volume]
}

type dependencies interface {
	Logger() log.Logger
	Clock() clockwork.Clock
	Process() *servicectx.Process
}

// OpenVolumes function detects and opens all volumes in the path.
func OpenVolumes(ctx context.Context, d dependencies, volumesPath string, config Config) (v *Volumes, err error) {
	v = &Volumes{
		clock:  d.Clock(),
		logger: d.Logger().WithComponent("storage.node.writer.volumes"),
		events: events.New[Writer](),
	}

	v.collection, err = opener.OpenVolumes(ctx, v.logger, volumesPath, func(spec volume.Spec) (*Volume, error) {
		return OpenVolume(ctx, v.logger, v.clock, config, spec, v.events)
	})
	if err != nil {
		return nil, err
	}

	// Graceful shutdown
	d.Process().OnShutdown(func(ctx context.Context) {
		v.logger.Info(ctx, "closing volumes")
		if err := v.collection.Close(ctx); err != nil {
			err := errors.PrefixError(err, "cannot close volumes")
			v.logger.Error(ctx, err.Error())
		}
		v.logger.Info(ctx, "closed volumes")
	})

	return v, nil
}

func (v *Volumes) Collection() *volume.Collection[*Volume] {
	return v.collection
}

func (v *Volumes) Events() *events.Events[Writer] {
	return v.events
}
