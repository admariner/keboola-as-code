package task_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/keboola/go-utils/pkg/wildcards"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	etcd "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdclient"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/task"
	"github.com/keboola/keboola-as-code/internal/pkg/telemetry"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/etcdhelper"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/ioutil"
)

func TestSuccessfulTask(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	etcdCredentials := etcdhelper.TmpNamespace(t)
	client := etcdhelper.ClientForTest(t, etcdCredentials)

	lock := "my-lock"
	taskType := "some.task"
	tKey := task.Key{
		ProjectID: 123,
		TaskID:    task.ID("my-receiver/my-export/" + taskType),
	}

	logs := ioutil.NewAtomicWriter()
	tel := newTestTelemetryWithFilter(t)

	// Create nodes
	node1, _ := createNode(t, etcdCredentials, logs, tel, "node1")
	node2, _ := createNode(t, etcdCredentials, logs, tel, "node2")
	logs.Truncate()
	tel.Reset()

	// Start a task
	taskWork := make(chan struct{})
	taskDone := make(chan struct{})
	_, err := node1.StartTask(task.Config{
		Key:  tKey,
		Type: taskType,
		Lock: lock,
		Context: func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(ctx, time.Minute)
		},
		Operation: func(ctx context.Context, logger log.Logger) task.Result {
			defer close(taskDone)
			<-taskWork
			logger.Info("some message from the task (1)")
			return task.OkResult("some result (1)").WithOutput("key", "value")
		},
	})
	assert.NoError(t, err)
	_, err = node2.StartTask(task.Config{
		Key:  tKey,
		Type: taskType,
		Lock: lock,
		Context: func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(ctx, time.Minute)
		},
		Operation: func(ctx context.Context, logger log.Logger) task.Result {
			assert.Fail(t, "should not be called")
			return task.ErrResult(errors.New("the task should not be called"))
		},
	})
	assert.NoError(t, err)

	// Check etcd state during task
	etcdhelper.AssertKVsString(t, client, `
<<<<<
runtime/lock/task/my-lock (lease)
-----
node1
>>>>>

<<<<<
task/123/my-receiver/my-export/some.task/%s
-----
{
  "projectId": 123,
  "taskId": "my-receiver/my-export/some.task/%s",
  "type": "some.task",
  "createdAt": "%s",
  "node": "node1",
  "lock": "runtime/lock/task/my-lock"
}
>>>>>
`)

	// Wait for task to finish
	finishTaskAndWait(t, client, taskWork, taskDone)

	// Check etcd state after task
	etcdhelper.AssertKVsString(t, client, `
<<<<<
task/123/my-receiver/my-export/some.task/%s
-----
{
  "projectId": 123,
  "taskId": "my-receiver/my-export/some.task/%s",
  "type": "some.task",
  "createdAt": "%s",
  "finishedAt": "%s",
  "node": "node1",
  "lock": "runtime/lock/task/my-lock",
  "result": "some result (1)",
  "outputs": {
    "key": "value"
  },
  "duration": %d
}
>>>>>
`)

	// Start another task with the same lock (lock is free)
	taskWork = make(chan struct{})
	taskDone = make(chan struct{})
	_, err = node2.StartTask(task.Config{
		Key:  tKey,
		Type: taskType,
		Lock: lock,
		Context: func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(ctx, time.Minute)
		},
		Operation: func(ctx context.Context, logger log.Logger) task.Result {
			defer close(taskDone)
			<-taskWork
			logger.Info("some message from the task (2)")
			return task.OkResult("some result (2)")
		},
	})
	assert.NoError(t, err)

	// Wait for task to finish
	finishTaskAndWait(t, client, taskWork, taskDone)

	// Check etcd state after second task
	etcdhelper.AssertKVsString(t, client, `
<<<<<
task/123/my-receiver/my-export/some.task/%s
-----
{
  "projectId": 123,
  "taskId": "my-receiver/my-export/some.task/%s",
  "type": "some.task",
  "createdAt": "%s",
  "finishedAt": "%s",
  "node": "node1",
  "lock": "runtime/lock/task/my-lock",
  "result": "some result (1)",
  "outputs": {
    "key": "value"
  },
  "duration": %d
}
>>>>>

<<<<<
task/123/my-receiver/my-export/some.task/%s
-----
{
  "projectId": 123,
  "taskId": "my-receiver/my-export/some.task/%s",
  "type": "some.task",
  "createdAt": "%s",
  "finishedAt": "%s",
  "node": "node2",
  "lock": "runtime/lock/task/my-lock",
  "result": "some result (2)",
  "duration": %d
}
>>>>>
`)

	// Check logs
	wildcards.Assert(t, `
[node1][task][123/my-receiver/my-export/some.task/%s]INFO  started task
[node1][task][123/my-receiver/my-export/some.task/%s]DEBUG  lock acquired "runtime/lock/task/my-lock"
[node2][task][123/my-receiver/my-export/some.task/%s]INFO  task ignored, the lock "runtime/lock/task/my-lock" is in use
[node1][task][123/my-receiver/my-export/some.task/%s]INFO  some message from the task (1)
[node1][task][123/my-receiver/my-export/some.task/%s]INFO  task succeeded (%s): some result (1) outputs: {"key":"value"}
[node1][task][123/my-receiver/my-export/some.task/%s]DEBUG  lock released "runtime/lock/task/my-lock"
[node2][task][123/my-receiver/my-export/some.task/%s]INFO  started task
[node2][task][123/my-receiver/my-export/some.task/%s]DEBUG  lock acquired "runtime/lock/task/my-lock"
[node2][task][123/my-receiver/my-export/some.task/%s]INFO  some message from the task (2)
[node2][task][123/my-receiver/my-export/some.task/%s]INFO  task succeeded (%s): some result (2)
[node2][task][123/my-receiver/my-export/some.task/%s]DEBUG  lock released "runtime/lock/task/my-lock"
`, logs.String())

	// Check spans
	tel.AssertSpans(t,
		tracetest.SpanStubs{
			{
				Name:     "keboola.go.task",
				SpanKind: trace.SpanKindInternal,
				SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    tel.TraceID(1),
					SpanID:     tel.SpanID(1),
					TraceFlags: trace.FlagsSampled,
				}),
				Status: tracesdk.Status{Code: codes.Ok},
				Attributes: []attribute.KeyValue{
					attribute.String("resource.name", "some.task"),
					attribute.String("task_id", "<dynamic>"),
					attribute.String("task_type", "some.task"),
					attribute.String("lock", "<dynamic>"),
					attribute.String("node", "node1"),
					attribute.String("created_at", "<dynamic>"),
					attribute.String("project_id", "123"),
					attribute.String("duration_sec", "<dynamic>"),
					attribute.String("finished_at", "<dynamic>"),
					attribute.Bool("is_success", true),
					attribute.String("result_outputs.key", "value"),
				},
			},
			{
				Name:     "keboola.go.task",
				SpanKind: trace.SpanKindInternal,
				SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    tel.TraceID(2),
					SpanID:     tel.SpanID(2),
					TraceFlags: trace.FlagsSampled,
				}),
				Status: tracesdk.Status{Code: codes.Ok},
				Attributes: []attribute.KeyValue{
					attribute.String("resource.name", "some.task"),
					attribute.String("task_id", "<dynamic>"),
					attribute.String("task_type", "some.task"),
					attribute.String("lock", "<dynamic>"),
					attribute.String("node", "node2"),
					attribute.String("created_at", "<dynamic>"),
					attribute.String("project_id", "123"),
					attribute.String("duration_sec", "<dynamic>"),
					attribute.String("finished_at", "<dynamic>"),
					attribute.Bool("is_success", true),
				},
			},
		},
		telemetry.WithSpanAttributeMapper(func(attr attribute.KeyValue) attribute.KeyValue {
			switch attr.Key {
			case "task_id", "lock", "created_at", "duration_sec", "finished_at":
				return attribute.String(string(attr.Key), "<dynamic>")
			}
			return attr
		}),
	)

	// Check metrics
	histBounds := []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000, 2500, 5000, 7500, 10000} // ms
	tel.AssertMetrics(t, []metricdata.Metrics{
		{
			Name:        "keboola.go.task.running",
			Description: "Background running tasks count.",
			Data: metricdata.Sum[int64]{
				Temporality: 1,
				DataPoints: []metricdata.DataPoint[int64]{
					{
						Value: 0,
						Attributes: attribute.NewSet(
							attribute.String("task_type", "some.task"),
						),
					},
				},
			},
		},
		{
			Name:        "keboola.go.task.duration",
			Description: "Background task duration.",
			Unit:        "ms",
			Data: metricdata.Histogram[float64]{
				Temporality: 1,
				DataPoints: []metricdata.HistogramDataPoint[float64]{
					{
						Count:  2,
						Bounds: histBounds,
						Attributes: attribute.NewSet(
							attribute.String("task_type", "some.task"),
							attribute.Bool("is_success", true),
						),
					},
				},
			},
		},
	})
}

func TestFailedTask(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	etcdCredentials := etcdhelper.TmpNamespace(t)
	client := etcdhelper.ClientForTest(t, etcdCredentials)

	lock := "my-lock"
	taskType := "some.task"
	tKey := task.Key{
		ProjectID: 123,
		TaskID:    task.ID("my-receiver/my-export/" + taskType),
	}

	logs := ioutil.NewAtomicWriter()
	tel := newTestTelemetryWithFilter(t)

	// Create nodes
	node1, _ := createNode(t, etcdCredentials, logs, tel, "node1")
	node2, _ := createNode(t, etcdCredentials, logs, tel, "node2")
	logs.Truncate()
	tel.Reset()

	// Start a task
	taskWork := make(chan struct{})
	taskDone := make(chan struct{})
	_, err := node1.StartTask(task.Config{
		Key:  tKey,
		Type: taskType,
		Lock: lock,
		Context: func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(ctx, time.Minute)
		},
		Operation: func(ctx context.Context, logger log.Logger) task.Result {
			defer close(taskDone)
			<-taskWork
			logger.Info("some message from the task (1)")
			return task.
				ErrResult(task.WrapUserError(errors.New("some error (1) - expected"))).
				WithOutput("key", "value")
		},
	})
	assert.NoError(t, err)
	_, err = node2.StartTask(task.Config{
		Key:  tKey,
		Type: taskType,
		Lock: lock,
		Context: func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(ctx, time.Minute)
		},
		Operation: func(ctx context.Context, logger log.Logger) task.Result {
			assert.Fail(t, "should not be called")
			return task.ErrResult(errors.New("the task should not be called"))
		},
	})
	assert.NoError(t, err)

	// Check etcd state during task
	etcdhelper.AssertKVsString(t, client, `
<<<<<
runtime/lock/task/my-lock (lease)
-----
node1
>>>>>

<<<<<
task/123/my-receiver/my-export/some.task/%s
-----
{
  "projectId": 123,
  "taskId": "my-receiver/my-export/some.task/%s",
  "type": "some.task",
  "createdAt": "%s",
  "node": "node1",
  "lock": "runtime/lock/task/my-lock"
}
>>>>>
`)

	// Wait for task to finish
	finishTaskAndWait(t, client, taskWork, taskDone)

	// Check etcd state after task
	etcdhelper.AssertKVsString(t, client, `
<<<<<
task/123/my-receiver/my-export/some.task/%s
-----
{
  "projectId": 123,
  "taskId": "my-receiver/my-export/some.task/%s",
  "type": "some.task",
  "createdAt": "%s",
  "finishedAt": "%s",
  "node": "node1",
  "lock": "runtime/lock/task/my-lock",
  "error": "some error (1) - expected",
  "outputs": {
    "key": "value"
  },
  "duration": %d
}
>>>>>
`)

	// Start another task with the same lock (lock is free)
	taskWork = make(chan struct{})
	taskDone = make(chan struct{})
	_, err = node2.StartTask(task.Config{
		Key:  tKey,
		Type: taskType,
		Lock: lock,
		Context: func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(ctx, time.Minute)
		},
		Operation: func(ctx context.Context, logger log.Logger) task.Result {
			defer close(taskDone)
			<-taskWork
			logger.Info("some message from the task (2)")
			return task.ErrResult(errors.New("some error (2) - unexpected"))
		},
	})
	assert.NoError(t, err)

	// Wait for task to finish
	finishTaskAndWait(t, client, taskWork, taskDone)

	// Check etcd state after second task
	etcdhelper.AssertKVsString(t, client, `
<<<<<
task/123/my-receiver/my-export/some.task/%s
-----
{
  "projectId": 123,
  "taskId": "my-receiver/my-export/some.task/%s",
  "type": "some.task",
  "createdAt": "%s",
  "finishedAt": "%s",
  "node": "node1",
  "lock": "runtime/lock/task/my-lock",
  "error": "some error (1) - expected",
  "outputs": {
    "key": "value"
  },
  "duration": %d
}
>>>>>

<<<<<
task/123/my-receiver/my-export/some.task/%s
-----
{
  "projectId": 123,
  "taskId": "my-receiver/my-export/some.task/%s",
  "type": "some.task",
  "createdAt": "%s",
  "finishedAt": "%s",
  "node": "node2",
  "lock": "runtime/lock/task/my-lock",
  "error": "some error (2) - unexpected",
  "duration": %d
}
>>>>>
`)

	// Check logs
	wildcards.Assert(t, `
[node1][task][123/my-receiver/my-export/some.task/%s]INFO  started task
[node1][task][123/my-receiver/my-export/some.task/%s]DEBUG  lock acquired "runtime/lock/task/my-lock"
[node2][task][123/my-receiver/my-export/some.task/%s]INFO  task ignored, the lock "runtime/lock/task/my-lock" is in use
[node1][task][123/my-receiver/my-export/some.task/%s]INFO  some message from the task (1)
[node1][task][123/my-receiver/my-export/some.task/%s]WARN  task failed (%s): some error (1) - expected %A outputs: {"key":"value"}
[node1][task][123/my-receiver/my-export/some.task/%s]DEBUG  lock released "runtime/lock/task/my-lock"
[node2][task][123/my-receiver/my-export/some.task/%s]INFO  started task
[node2][task][123/my-receiver/my-export/some.task/%s]DEBUG  lock acquired "runtime/lock/task/my-lock"
[node2][task][123/my-receiver/my-export/some.task/%s]INFO  some message from the task (2)
[node2][task][123/my-receiver/my-export/some.task/%s]WARN  task failed (%s): some error (2) - unexpected [%s]
[node2][task][123/my-receiver/my-export/some.task/%s]DEBUG  lock released "runtime/lock/task/my-lock"
`, logs.String())

	// Check spans
	tel.AssertSpans(t,
		tracetest.SpanStubs{
			{
				Name:     "keboola.go.task",
				SpanKind: trace.SpanKindInternal,
				SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    tel.TraceID(1),
					SpanID:     tel.SpanID(1),
					TraceFlags: trace.FlagsSampled,
				}),
				Status: tracesdk.Status{Code: codes.Error, Description: "some error (1) - expected"},
				Attributes: []attribute.KeyValue{
					attribute.String("resource.name", "some.task"),
					attribute.String("task_id", "<dynamic>"),
					attribute.String("task_type", "some.task"),
					attribute.String("lock", "<dynamic>"),
					attribute.String("node", "node1"),
					attribute.String("created_at", "<dynamic>"),
					attribute.String("project_id", "123"),
					attribute.String("duration_sec", "<dynamic>"),
					attribute.String("finished_at", "<dynamic>"),
					attribute.Bool("is_success", false),
					attribute.Bool("is_application_error", false),
					attribute.String("error", "some error (1) - expected"),
					attribute.String("error_type", "other"),
					attribute.String("result_outputs.key", "value"),
				},
				Events: []tracesdk.Event{
					{
						Name: "exception",
						Attributes: []attribute.KeyValue{
							attribute.String("exception.type", "*task.UserError"),
							attribute.String("exception.message", "some error (1) - expected"),
						},
					},
				},
			},
			{
				Name:     "keboola.go.task",
				SpanKind: trace.SpanKindInternal,
				SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    tel.TraceID(2),
					SpanID:     tel.SpanID(2),
					TraceFlags: trace.FlagsSampled,
				}),
				Status: tracesdk.Status{Code: codes.Error, Description: "some error (2) - unexpected"},
				Attributes: []attribute.KeyValue{
					attribute.String("resource.name", "some.task"),
					attribute.String("task_id", "<dynamic>"),
					attribute.String("task_type", "some.task"),
					attribute.String("lock", "<dynamic>"),
					attribute.String("node", "node2"),
					attribute.String("created_at", "<dynamic>"),
					attribute.String("project_id", "123"),
					attribute.String("duration_sec", "<dynamic>"),
					attribute.String("finished_at", "<dynamic>"),
					attribute.Bool("is_success", false),
					attribute.Bool("is_application_error", true),
					attribute.String("error", "some error (2) - unexpected"),
					attribute.String("error_type", "other"),
				},
				Events: []tracesdk.Event{
					{
						Name: "exception",
						Attributes: []attribute.KeyValue{
							attribute.String("exception.type", "*errors.withStack"),
							attribute.String("exception.message", "some error (2) - unexpected"),
						},
					},
				},
			},
		},
		telemetry.WithSpanAttributeMapper(func(attr attribute.KeyValue) attribute.KeyValue {
			switch attr.Key {
			case "task_id", "lock", "created_at", "duration_sec", "finished_at":
				return attribute.String(string(attr.Key), "<dynamic>")
			}
			return attr
		}),
	)

	// Check metrics
	histBounds := []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000, 2500, 5000, 7500, 10000} // ms
	tel.AssertMetrics(t,
		[]metricdata.Metrics{
			{
				Name:        "keboola.go.task.running",
				Description: "Background running tasks count.",
				Data: metricdata.Sum[int64]{
					Temporality: 1,
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Value: 0,
							Attributes: attribute.NewSet(
								attribute.String("task_type", "some.task"),
							),
						},
					},
				},
			},
			{
				Name:        "keboola.go.task.duration",
				Description: "Background task duration.",
				Unit:        "ms",
				Data: metricdata.Histogram[float64]{
					Temporality: 1,
					DataPoints: []metricdata.HistogramDataPoint[float64]{
						// Expected error, so it will not be taken as an error in the metrics.
						{
							Count:  1,
							Bounds: histBounds,
							Attributes: attribute.NewSet(
								attribute.String("task_type", "some.task"),
								attribute.Bool("is_success", false),
								attribute.Bool("is_application_error", false),
								attribute.String("error_type", "other"),
							),
						},
						// Unexpected error
						{
							Count:  1,
							Bounds: histBounds,
							Attributes: attribute.NewSet(
								attribute.String("task_type", "some.task"),
								attribute.Bool("is_success", false),
								attribute.Bool("is_application_error", true),
								attribute.String("error_type", "other"),
							),
						},
					},
				},
			},
		},
		telemetry.WithDataPointSortKey(func(attrs attribute.Set) string {
			if v, _ := attrs.Value("is_application_error"); v.AsBool() {
				return "1"
			}
			return "0"
		}),
	)
}

func TestTaskTimeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	etcdCredentials := etcdhelper.TmpNamespace(t)
	client := etcdhelper.ClientForTest(t, etcdCredentials)

	lock := "my-lock"
	taskType := "some.task"
	tKey := task.Key{
		ProjectID: 123,
		TaskID:    task.ID("my-receiver/my-export/" + taskType),
	}

	tel := newTestTelemetryWithFilter(t)

	// Create node and
	node1, d := createNode(t, etcdCredentials, nil, tel, "node1")
	logger := d.DebugLogger()
	logger.Truncate()
	tel.Reset()

	// Start task
	_, err := node1.StartTask(task.Config{
		Key:  tKey,
		Type: taskType,
		Lock: lock,
		Context: func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(ctx, 10*time.Millisecond) // <<<<<<<<<<<<<<<<<
		},
		Operation: func(ctx context.Context, logger log.Logger) task.Result {
			// Sleep for 10 seconds
			select {
			case <-time.After(10 * time.Second):
			case <-ctx.Done():
				return task.ErrResult(ctx.Err())
			}
			return task.ErrResult(errors.New("invalid state, task should time out"))
		},
	})
	assert.NoError(t, err)

	// Wait for the task
	assert.Eventually(t, func() bool {
		return strings.Contains(logger.AllMessages(), "DEBUG  lock released")
	}, 5*time.Second, 100*time.Millisecond, logger.AllMessages())

	// Check etcd state after task
	etcdhelper.AssertKVsString(t, client, `
<<<<<
task/123/my-receiver/my-export/some.task/%s
-----
{
  "projectId": 123,
  "taskId": "my-receiver/my-export/some.task/%s",
  "type": "some.task",
  "createdAt": "%s",
  "finishedAt": "%s",
  "node": "node1",
  "lock": "runtime/lock/task/my-lock",
  "error": "context deadline exceeded",
  "duration": %d
}
>>>>>
`)

	// Check logs
	wildcards.Assert(t, `
[node1][task][123/my-receiver/my-export/some.task/%s]INFO  started task
[node1][task][123/my-receiver/my-export/some.task/%s]DEBUG  lock acquired "runtime/lock/task/my-lock"
[node1][task][123/my-receiver/my-export/some.task/%s]WARN  task failed (%s): context deadline exceeded
[node1][task][123/my-receiver/my-export/some.task/%s]DEBUG  lock released "runtime/lock/task/my-lock"
`, logger.AllMessages())

	// Check spans
	tel.AssertSpans(t,
		tracetest.SpanStubs{
			{
				Name:     "keboola.go.task",
				SpanKind: trace.SpanKindInternal,
				SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    tel.TraceID(1),
					SpanID:     tel.SpanID(1),
					TraceFlags: trace.FlagsSampled,
				}),
				Status: tracesdk.Status{Code: codes.Error, Description: "context deadline exceeded"},
				Attributes: []attribute.KeyValue{
					attribute.String("resource.name", "some.task"),
					attribute.String("task_id", "<dynamic>"),
					attribute.String("task_type", "some.task"),
					attribute.String("lock", "<dynamic>"),
					attribute.String("node", "node1"),
					attribute.String("created_at", "<dynamic>"),
					attribute.String("project_id", "123"),
					attribute.String("duration_sec", "<dynamic>"),
					attribute.String("finished_at", "<dynamic>"),
					attribute.Bool("is_success", false),
					attribute.Bool("is_application_error", true),
					attribute.String("error", "context deadline exceeded"),
					attribute.String("error_type", "deadline_exceeded"),
				},
				Events: []tracesdk.Event{
					{
						Name: "exception",
						Attributes: []attribute.KeyValue{
							attribute.String("exception.type", "context.deadlineExceededError"),
							attribute.String("exception.message", "context deadline exceeded"),
						},
					},
				},
			},
		},
		telemetry.WithSpanAttributeMapper(func(attr attribute.KeyValue) attribute.KeyValue {
			switch attr.Key {
			case "task_id", "lock", "created_at", "duration_sec", "finished_at":
				return attribute.String(string(attr.Key), "<dynamic>")
			}
			return attr
		}),
	)

	// Check metrics
	histBounds := []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000, 2500, 5000, 7500, 10000} // ms
	tel.AssertMetrics(t,
		[]metricdata.Metrics{
			{
				Name:        "keboola.go.task.running",
				Description: "Background running tasks count.",
				Data: metricdata.Sum[int64]{
					Temporality: 1,
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Value: 0,
							Attributes: attribute.NewSet(
								attribute.String("task_type", "some.task"),
							),
						},
					},
				},
			},
			{
				Name:        "keboola.go.task.duration",
				Description: "Background task duration.",
				Unit:        "ms",
				Data: metricdata.Histogram[float64]{
					Temporality: 1,
					DataPoints: []metricdata.HistogramDataPoint[float64]{
						{
							Count:  1,
							Bounds: histBounds,
							Attributes: attribute.NewSet(
								attribute.String("task_type", "some.task"),
								attribute.Bool("is_success", false),
								attribute.Bool("is_application_error", true),
								attribute.String("error_type", "deadline_exceeded"),
							),
						},
					},
				},
			},
		},
	)
}

func TestWorkerNodeShutdownDuringTask(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	etcdCredentials := etcdhelper.TmpNamespace(t)
	client := etcdhelper.ClientForTest(t, etcdCredentials)

	lock := "my-lock"
	taskType := "some.task"
	tKey := task.Key{
		ProjectID: 123,
		TaskID:    task.ID("my-receiver/my-export/" + taskType),
	}

	logs := ioutil.NewAtomicWriter()
	tel := newTestTelemetryWithFilter(t)

	// Create node
	node1, d := createNode(t, etcdCredentials, logs, tel, "node1")
	tel.Reset()
	logs.Truncate()

	// Start a task
	taskWork := make(chan struct{})
	taskDone := make(chan struct{})
	etcdhelper.ExpectModification(t, client, func() {
		_, err := node1.StartTask(task.Config{
			Key:  tKey,
			Type: taskType,
			Lock: lock,
			Context: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(ctx, time.Minute)
			},
			Operation: func(ctx context.Context, logger log.Logger) task.Result {
				defer close(taskDone)
				<-taskWork
				logger.Info("some message from the task")
				return task.OkResult("some result")
			},
		})
		assert.NoError(t, err)
	})

	// Shutdown node
	shutdownDone := make(chan struct{})
	d.Process().Shutdown(errors.New("some reason"))
	go func() {
		defer close(shutdownDone)
		d.Process().WaitForShutdown()
	}()

	// Wait for task to finish
	time.Sleep(100 * time.Millisecond)
	finishTaskAndWait(t, client, taskWork, taskDone)

	// Wait for shutdown
	select {
	case <-time.After(time.Second):
		assert.Fail(t, "timeout")
	case <-shutdownDone:
	}

	// Check etcd state
	etcdhelper.AssertKVsString(t, client, `
<<<<<
task/123/my-receiver/my-export/some.task/%s
-----
{
  "projectId": 123,
  "taskId": "my-receiver/my-export/some.task/%s",
  "type": "some.task",
  "createdAt": "%s",
  "finishedAt": "%s",
  "node": "node1",
  "lock": "runtime/lock/task/my-lock",
  "result": "some result",
  "duration": %d
}
>>>>>
`)

	// Check logs
	wildcards.Assert(t, `
[node1][task][%s]INFO  started task
[node1][task][%s]DEBUG  lock acquired "runtime/lock/task/my-lock"
[node1]INFO  exiting (some reason)
[node1][task]INFO  received shutdown request
[node1][task]INFO  waiting for "1" tasks to be finished
[node1][task][123/my-receiver/my-export/some.task/%s]INFO  some message from the task
[node1][task][123/my-receiver/my-export/some.task/%s]INFO  task succeeded (%s): some result
[node1][task][123/my-receiver/my-export/some.task/%s]DEBUG  lock released "runtime/lock/task/my-lock"
[node1][task][etcd-session]INFO  closing etcd session
[node1][task][etcd-session]INFO  closed etcd session | %s
[node1][task]INFO  shutdown done
[node1][etcd-client]INFO  closing etcd connection
[node1][etcd-client]INFO  closed etcd connection | %s
[node1]INFO  exited
`, logs.String())
}

func newTestTelemetryWithFilter(t *testing.T) telemetry.ForTest {
	t.Helper()
	return telemetry.
		NewForTest(t).
		SetSpanFilter(func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) bool {
			// Ignore etcd spans
			return !strings.HasPrefix(spanName, "etcd")
		})
}

func createNode(t *testing.T, etcdCredentials etcdclient.Credentials, logs io.Writer, tel telemetry.ForTest, nodeName string) (*task.Node, dependencies.Mocked) {
	t.Helper()
	d := createDeps(t, etcdCredentials, logs, tel, nodeName)
	node, err := task.NewNode(d)
	require.NoError(t, err)
	d.DebugLogger().Truncate()
	return node, d
}

func createDeps(t *testing.T, etcdCredentials etcdclient.Credentials, logs io.Writer, tel telemetry.ForTest, nodeName string) dependencies.Mocked {
	t.Helper()

	d := dependencies.NewMocked(
		t,
		dependencies.WithUniqueID(nodeName),
		dependencies.WithLoggerPrefix(fmt.Sprintf("[%s]", nodeName)),
		dependencies.WithTelemetry(tel),
		dependencies.WithEtcdCredentials(etcdCredentials),
	)

	// Connect logs output
	if logs != nil {
		d.DebugLogger().ConnectTo(logs)
	}

	return d
}

func finishTaskAndWait(t *testing.T, client *etcd.Client, taskWork, taskDone chan struct{}) {
	t.Helper()

	// Wait for update of the task in etcd
	etcdhelper.ExpectModification(t, client, func() {
		// Finish work in the task
		close(taskWork)

		// Wait for goroutine
		select {
		case <-time.After(time.Second):
			assert.Fail(t, "timeout")
		case <-taskDone:
		}
	})
}
