package mapper

import (
	bufferDesign "github.com/keboola/keboola-as-code/api/buffer"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/api/gen/buffer"
	taskModel "github.com/keboola/keboola-as-code/internal/pkg/service/common/task"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

func (m Mapper) TaskPayload(model *taskModel.Task) (r *buffer.Task) {
	out := &buffer.Task{
		ID:        model.TaskID,
		URL:       formatTaskURL(m.bufferAPIHost, model.Key),
		CreatedAt: model.CreatedAt.String(),
	}

	if model.FinishedAt != nil {
		v := model.FinishedAt.String()
		out.FinishedAt = &v
	}

	if model.Duration != nil {
		v := model.Duration.Milliseconds()
		out.Duration = &v
	}

	switch {
	case model.IsProcessing():
		out.Status = bufferDesign.TaskStatusProcessing
	case model.IsSuccessful():
		out.Status = bufferDesign.TaskStatusSuccess
		out.IsFinished = true
		out.Result = &model.Result
	case model.IsFailed():
		out.Status = bufferDesign.TaskStatusError
		out.IsFinished = true
		out.Error = &model.Error
	default:
		panic(errors.New("unexpected task status"))
	}

	return out
}
