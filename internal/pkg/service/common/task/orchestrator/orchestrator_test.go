package orchestrator_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/keboola/go-client/pkg/keboola"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop/serde"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/task"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/task/orchestrator"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/etcdhelper"
	"github.com/keboola/keboola-as-code/internal/pkg/validator"
)

type testResource struct {
	ProjectID       keboola.ProjectID
	DistributionKey string
	ID              string
}

func TestOrchestrator(t *testing.T) {
	t.Parallel()

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	etcdCfg := etcdhelper.TmpNamespace(t)
	client := etcdhelper.ClientForTest(t, etcdCfg)

	d1 := dependencies.NewMocked(t,
		dependencies.WithCtx(ctx),
		dependencies.WithEtcdConfig(etcdCfg),
		dependencies.WithEnabledOrchestrator(),
		dependencies.WithNodeID("node1"),
	)
	grp1, err := d1.DistributionNode().Group("my-group")
	require.NoError(t, err)
	node1 := orchestrator.NewNode(d1)

	d2 := dependencies.NewMocked(t,
		dependencies.WithCtx(ctx),
		dependencies.WithEtcdConfig(etcdCfg),
		dependencies.WithEnabledOrchestrator(),
		dependencies.WithNodeID("node2"),
	)
	grp2, err := d2.DistributionNode().Group("my-group")
	require.NoError(t, err)
	node2 := orchestrator.NewNode(d2)

	pfx := etcdop.NewTypedPrefix[testResource]("my/prefix", serde.NewJSON(validator.New().Validate))

	// Orchestrator config
	config := orchestrator.Config[testResource]{
		Name: "some.task",
		Source: orchestrator.Source[testResource]{
			WatchPrefix:     pfx,
			RestartInterval: time.Minute,
		},
		DistributionKey: func(event etcdop.WatchEventT[testResource]) string {
			return event.Value.DistributionKey
		},
		Lock: func(event etcdop.WatchEventT[testResource]) string {
			// Define a custom lock name
			return "custom-lock"
		},
		TaskKey: func(event etcdop.WatchEventT[testResource]) task.Key {
			resource := event.Value
			return task.Key{
				ProjectID: resource.ProjectID,
				TaskID:    task.ID("my-prefix/some.task/" + resource.ID),
			}
		},
		TaskCtx: func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(ctx, time.Minute)
		},
		TaskFactory: func(event etcdop.WatchEventT[testResource]) task.Fn {
			return func(ctx context.Context, logger log.Logger) task.Result {
				logger.Info(ctx, "message from the task")
				return task.OkResult(event.Value.ID)
			}
		},
	}

	// Create orchestrator per each node
	assert.NoError(t, <-node1.Start(grp1, config))
	assert.NoError(t, <-node2.Start(grp2, config))

	// Put some key to trigger the task
	v := testResource{ProjectID: 1000, DistributionKey: "foo", ID: "ResourceID"}
	assert.NoError(t, pfx.Key("key1").Put(client, v).Do(ctx).Err())

	// Wait for task on the node 2
	assert.Eventually(t, func() bool {
		return d2.DebugLogger().CompareJSONMessages(`{"level":"debug","message":"lock released%s"}`) == nil
	}, 5*time.Second, 10*time.Millisecond, "timeout")

	// Wait for "not assigned" message form the node 1
	assert.Eventually(t, func() bool {
		return d1.DebugLogger().CompareJSONMessages(`{"level":"debug","message":"not assigned%s"}`) == nil
	}, 5*time.Second, 10*time.Millisecond, "timeout")

	cancel()
	wg.Wait()
	d1.Process().Shutdown(ctx, errors.New("bye bye 1"))
	d1.Process().WaitForShutdown()
	d2.Process().Shutdown(ctx, errors.New("bye bye 2"))
	d2.Process().WaitForShutdown()

	expected := `
{"level":"info","message":"ready","component":"orchestrator","task":"some.task"}
{"level":"info","message":"assigned \"1000/my-prefix/some.task/ResourceID\"","component":"orchestrator","task":"some.task"}
{"level":"info","message":"started task","component":"task","task":"1000/my-prefix/some.task/ResourceID/%s"}
{"level":"debug","message":"lock acquired \"runtime/lock/task/custom-lock\"","component":"task","task":"1000/my-prefix/some.task/ResourceID/%s"}
{"level":"info","message":"message from the task","component":"task","task":"1000/my-prefix/some.task/ResourceID/%s"}
{"level":"info","message":"task succeeded (%s): ResourceID","component":"task","task":"1000/my-prefix/some.task/ResourceID/%s"}
{"level":"debug","message":"lock released \"runtime/lock/task/custom-lock\"","component":"task","task":"1000/my-prefix/some.task/ResourceID/%s"}
`
	d2.DebugLogger().AssertJSONMessages(t, expected)

	expected = `
{"level":"info","message":"ready","component":"orchestrator","task":"some.task"}
{"level":"debug","message":"not assigned \"1000/my-prefix/some.task/ResourceID\", distribution key \"foo\"","component":"orchestrator","task":"some.task"}
`
	d1.DebugLogger().AssertJSONMessages(t, expected)
}

func TestOrchestrator_StartTaskIf(t *testing.T) {
	t.Parallel()

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	etcdCfg := etcdhelper.TmpNamespace(t)
	client := etcdhelper.ClientForTest(t, etcdCfg)

	d := dependencies.NewMocked(t,
		dependencies.WithCtx(ctx),
		dependencies.WithEtcdConfig(etcdCfg),
		dependencies.WithNodeID("node1"),
		dependencies.WithEnabledOrchestrator(),
	)

	dist, err := d.DistributionNode().Group("my-group")
	require.NoError(t, err)

	node := orchestrator.NewNode(d)

	pfx := etcdop.NewTypedPrefix[testResource]("my/prefix", serde.NewJSON(validator.New().Validate))

	// Orchestrator config
	config := orchestrator.Config[testResource]{
		Name: "some.task",
		Source: orchestrator.Source[testResource]{
			WatchPrefix:     pfx,
			RestartInterval: time.Minute,
		},
		DistributionKey: func(event etcdop.WatchEventT[testResource]) string {
			return event.Value.DistributionKey
		},
		TaskKey: func(event etcdop.WatchEventT[testResource]) task.Key {
			resource := event.Value
			return task.Key{
				ProjectID: resource.ProjectID,
				TaskID:    task.ID("my-prefix/some.task/" + resource.ID),
			}
		},
		TaskCtx: func() (context.Context, context.CancelFunc) {
			return context.WithTimeout(ctx, time.Minute)
		},
		StartTaskIf: func(event etcdop.WatchEventT[testResource]) (string, bool) {
			if event.Value.ID == "GoodID" { // <<<<<<<<<<<<<<<<<<<<
				return "", true
			}
			return "StartTaskIf condition evaluated as false", false
		},
		TaskFactory: func(event etcdop.WatchEventT[testResource]) task.Fn {
			return func(ctx context.Context, logger log.Logger) task.Result {
				logger.Info(ctx, "message from the task")
				return task.OkResult(event.Value.ID)
			}
		},
	}

	assert.NoError(t, <-node.Start(dist, config))
	v1 := testResource{ProjectID: 1000, DistributionKey: "foo", ID: "BadID"}
	v2 := testResource{ProjectID: 1000, DistributionKey: "foo", ID: "GoodID"}
	assert.NoError(t, pfx.Key("key1").Put(client, v1).Do(ctx).Err())
	assert.NoError(t, pfx.Key("key2").Put(client, v2).Do(ctx).Err())
	assert.Eventually(t, func() bool {
		return d.DebugLogger().CompareJSONMessages(`{"level":"debug","message":"lock released%s"}`) == nil
	}, 5*time.Second, 10*time.Millisecond, "timeout")

	cancel()
	wg.Wait()
	d.Process().Shutdown(ctx, errors.New("bye bye 1"))
	d.Process().WaitForShutdown()

	expected := `
{"level":"info","message":"ready","component":"orchestrator","task":"some.task"}
{"level":"debug","message":"skipped \"1000/my-prefix/some.task/BadID\", StartTaskIf condition evaluated as false","component":"orchestrator","task":"some.task"}
{"level":"info","message":"assigned \"1000/my-prefix/some.task/GoodID\"","component":"orchestrator","task":"some.task"}
{"level":"info","message":"started task","component":"task","task":"1000/my-prefix/some.task/GoodID/%s"}
{"level":"debug","message":"lock acquired \"runtime/lock/task/1000/my-prefix/some.task/GoodID\"","component":"task","task":"1000/my-prefix/some.task/GoodID/%s"}
{"level":"info","message":"message from the task","component":"task","task":"1000/my-prefix/some.task/GoodID/%s"}
{"level":"info","message":"task succeeded (%s): GoodID","component":"task","task":"1000/my-prefix/some.task/GoodID/%s"}
{"level":"debug","message":"lock released \"runtime/lock/task/1000/my-prefix/some.task/GoodID\"","component":"task","task":"1000/my-prefix/some.task/GoodID/%s"}
`
	d.DebugLogger().AssertJSONMessages(t, expected)
}

func TestOrchestrator_RestartInterval(t *testing.T) {
	t.Parallel()

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	etcdCfg := etcdhelper.TmpNamespace(t)
	client := etcdhelper.ClientForTest(t, etcdCfg)

	restartInterval := time.Millisecond
	clk := clock.NewMock()
	d := dependencies.NewMocked(t,
		dependencies.WithCtx(ctx),
		dependencies.WithClock(clk),
		dependencies.WithEtcdConfig(etcdCfg),
		dependencies.WithNodeID("node1"),
		dependencies.WithEnabledOrchestrator(),
	)

	dist, err := d.DistributionNode().Group("my-group")
	require.NoError(t, err)

	node := orchestrator.NewNode(d)

	pfx := etcdop.NewTypedPrefix[testResource]("my/prefix", serde.NewJSON(validator.New().Validate))

	// Orchestrator config
	config := orchestrator.Config[testResource]{
		Name: "some.task",
		Source: orchestrator.Source[testResource]{
			WatchPrefix:     pfx,
			RestartInterval: restartInterval,
		},
		DistributionKey: func(event etcdop.WatchEventT[testResource]) string {
			return event.Value.DistributionKey
		},
		TaskKey: func(event etcdop.WatchEventT[testResource]) task.Key {
			resource := event.Value
			return task.Key{
				ProjectID: resource.ProjectID,
				TaskID:    task.ID("my-prefix/some.task/" + resource.ID),
			}
		},
		TaskCtx: func() (context.Context, context.CancelFunc) {
			// Each orchestrator task must have a deadline.
			return context.WithTimeout(ctx, time.Minute)
		},
		TaskFactory: func(event etcdop.WatchEventT[testResource]) task.Fn {
			return func(ctx context.Context, logger log.Logger) task.Result {
				logger.Info(ctx, "message from the task")
				return task.OkResult(event.Value.ID)
			}
		},
	}

	// Create orchestrator per each node
	assert.NoError(t, <-node.Start(dist, config))

	// Put some key to trigger the task
	v := testResource{ProjectID: 1000, DistributionKey: "foo", ID: "ResourceID"}
	assert.NoError(t, pfx.Key("key1").Put(client, v).Do(ctx).Err())
	assert.Eventually(t, func() bool {
		return d.DebugLogger().CompareJSONMessages(`{"level":"debug","message":"lock released%s"}`) == nil
	}, 5*time.Second, 10*time.Millisecond, "timeout")
	d.DebugLogger().Truncate()

	// 3x restart interval
	clk.Add(restartInterval)
	assert.Eventually(t, func() bool {
		return d.DebugLogger().CompareJSONMessages(`{"level":"debug","message":"restart"}`) == nil
	}, 5*time.Second, 10*time.Millisecond, "timeout")

	// Put some key to trigger the task
	assert.Eventually(t, func() bool {
		return d.DebugLogger().CompareJSONMessages(`{"level":"debug","message":"lock released%s"}`) == nil
	}, 5*time.Second, 10*time.Millisecond, "timeout")

	cancel()
	wg.Wait()
	d.Process().Shutdown(ctx, errors.New("bye bye"))
	d.Process().WaitForShutdown()

	expected := `
{"level":"debug","message":"restart","component":"orchestrator","task":"some.task"}
{"level":"info","message":"assigned \"1000/my-prefix/some.task/ResourceID\"","component":"orchestrator","task":"some.task"}
{"level":"info","message":"started task","component":"task","task":"1000/my-prefix/some.task/ResourceID/%s"}
{"level":"debug","message":"lock acquired \"runtime/lock/task/1000/my-prefix/some.task/ResourceID\"","component":"task","task":"1000/my-prefix/some.task/ResourceID/%s"}
{"level":"info","message":"message from the task","component":"task","task":"1000/my-prefix/some.task/ResourceID/%s"}
{"level":"info","message":"task succeeded (0s): ResourceID","component":"task","task":"1000/my-prefix/some.task/ResourceID/%s"}
{"level":"debug","message":"lock released \"runtime/lock/task/1000/my-prefix/some.task/ResourceID\"","component":"task","task":"1000/my-prefix/some.task/ResourceID/%s"}
`
	d.DebugLogger().AssertJSONMessages(t, expected)
}
