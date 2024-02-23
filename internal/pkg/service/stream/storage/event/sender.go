// Package event provides the dispatch of events for platform telemetry purposes.
// Events contain slice/file statistics, for example, for billing purposes.
package event

import (
	"context"
	"fmt"
	"time"

	"github.com/keboola/go-client/pkg/client"
	"github.com/keboola/go-client/pkg/keboola"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/definition/key"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/model"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/statistics"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

const componentID = keboola.ComponentID("keboola.keboola-stream")

type Sender struct {
	logger log.Logger
}

type dependencies interface {
	Logger() log.Logger
}

func NewSender(d dependencies) *Sender {
	return &Sender{logger: d.Logger()}
}

type Params struct {
	ProjectID keboola.ProjectID
	SourceID  key.SourceID
	SinkID    key.SinkID
	Stats     statistics.Value
}

func (s *Sender) SendSliceUploadEvent(ctx context.Context, api *keboola.AuthorizedAPI, duration time.Duration, errPtr *error, slice model.Slice, stats statistics.Value) {
	// Get error
	var err error
	if errPtr != nil {
		err = *errPtr
	}

	// Catch panic
	panicErr := recover()
	if panicErr != nil {
		err = errors.Errorf(`%s`, panicErr)
	}

	formatMsg := func(err error) string {
		if err != nil {
			return "Slice upload failed."
		} else {
			return "Slice upload done."
		}
	}

	s.sendEvent(ctx, api, duration, "slice-upload", err, formatMsg, Params{
		ProjectID: slice.ProjectID,
		SourceID:  slice.SourceID,
		SinkID:    slice.SinkID,
		Stats:     stats,
	})

	// Throw panic
	if panicErr != nil {
		panic(panicErr)
	}
}

func (s *Sender) SendFileImportEvent(ctx context.Context, api *keboola.AuthorizedAPI, duration time.Duration, errPtr *error, file model.File, stats statistics.Value) {
	// Get error
	var err error
	if errPtr != nil {
		err = *errPtr
	}

	// Catch panic
	panicErr := recover()
	if panicErr != nil {
		err = errors.Errorf(`%s`, panicErr)
	}

	formatMsg := func(err error) string {
		if err != nil {
			return "File import failed."
		} else {
			return "File import done."
		}
	}

	s.sendEvent(ctx, api, duration, "file-import", err, formatMsg, Params{
		ProjectID: file.ProjectID,
		SourceID:  file.SourceID,
		SinkID:    file.SinkID,
		Stats:     stats,
	})

	// Throw panic
	if panicErr != nil {
		panic(panicErr)
	}
}

func (s *Sender) sendEvent(ctx context.Context, api *keboola.AuthorizedAPI, duration time.Duration, eventName string, err error, msg func(error) string, params Params) {
	event := &keboola.Event{
		ComponentID: componentID,
		Message:     msg(err),
		Type:        "info",
		Duration:    client.DurationSeconds(duration),
		Params: map[string]any{
			"eventName": eventName,
		},
		Results: map[string]any{
			"projectId": params.ProjectID,
			"sourceId":  params.SourceID,
			"sinkId":    params.SinkID,
		},
	}
	if err != nil {
		event.Type = "error"
		event.Results["error"] = fmt.Sprintf("%s", err)
	} else {
		event.Results["statistics"] = map[string]any{
			"firstRecordAt":    params.Stats.FirstRecordAt.String(),
			"lastRecordAt":     params.Stats.LastRecordAt.String(),
			"recordsCount":     params.Stats.RecordsCount,
			"slicesCount":      params.Stats.SlicesCount,
			"uncompressedSize": params.Stats.UncompressedSize.Bytes(),
			"compressedSize":   params.Stats.CompressedSize.Bytes(),
			"stagingSize":      params.Stats.StagingSize.Bytes(),
		}
	}

	event, err = api.CreateEventRequest(event).Send(ctx)
	if err == nil {
		s.logger.Debugf(ctx, "Sent \"%s\" event id: \"%s\"", eventName, event.ID)
	} else {
		s.logger.Warnf(ctx, "Cannot send \"%s\" event: %s", eventName, err)
	}
}
