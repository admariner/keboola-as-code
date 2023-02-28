package service_test

import (
	"compress/gzip"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/keboola/go-client/pkg/keboola"
	"github.com/keboola/go-utils/pkg/wildcards"
	"github.com/stretchr/testify/assert"
	etcd "go.etcd.io/etcd/client/v3"

	"github.com/keboola/keboola-as-code/internal/pkg/idgenerator"
	bufferDependencies "github.com/keboola/keboola-as-code/internal/pkg/service/buffer/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/store/model"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/worker/service"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/etcdhelper"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/strhelper"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/testhelper"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/testproject"
)

// TestSliceUploadTask - there are 2 slices, one is empty.
// Both sliced are closed, but only the non-empty slice is uploaded to file storage.
func TestSliceUploadTask(t *testing.T) {
	t.Parallel()

	project := testproject.GetTestProjectForTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test dependencies
	clk := clock.NewMock()
	clk.Set(time.Time{}.Add(time.Second))
	etcdNamespace := "unit-" + t.Name() + "-" + idgenerator.Random(8)
	client := etcdhelper.ClientForTestWithNamespace(t, etcdNamespace)
	opts := []dependencies.MockedOption{
		dependencies.WithClock(clk),
		dependencies.WithEtcdNamespace(etcdNamespace),
		dependencies.WithTestProject(project),
	}

	// Create file, download credentials are used to verify file content
	uploadCredentials, err := project.KeboolaProjectAPI().CreateFileResourceRequest("slice-upload-task-test", keboola.WithIsSliced(true)).Send(ctx)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	downloadCredentials, err := project.KeboolaProjectAPI().GetFileWithCredentialsRequest(uploadCredentials.ID).Send(ctx)
	assert.NoError(t, err)

	// Create receivers, exports and records
	apiDeps1 := bufferDependencies.NewMockedDeps(t, append(opts, dependencies.WithUniqueID("api-node-1"))...)
	apiDeps2 := bufferDependencies.NewMockedDeps(t, append(opts, dependencies.WithUniqueID("api-node-2"))...)
	str := apiDeps1.Store()
	emptySliceKey := createExport(t, "my-receiver-1", "my-export-1", ctx, clk, client, str, uploadCredentials)
	clk.Add(time.Minute)
	slice1Key := createExport(t, "my-receiver-2", "my-export-2", ctx, clk, client, str, uploadCredentials)
	clk.Add(time.Minute)
	createRecords(t, ctx, clk, apiDeps1, slice1Key.ReceiverKey, 1, 3)
	createRecords(t, ctx, clk, apiDeps2, slice1Key.ReceiverKey, 4, 4)
	assert.Eventually(t, func() bool {
		count, err := str.CountRecords(ctx, slice1Key)
		assert.NoError(t, err)
		return count == 7
	}, time.Second, 10*time.Millisecond)
	clk.Add(time.Minute) // trigger API nodes sync
	assertStateBeforeUpload(t, client)

	// Start worker node
	workerDeps := bufferDependencies.NewMockedDeps(t, append(opts, dependencies.WithUniqueID("my-worker"))...)
	workerDeps.DebugLogger().ConnectTo(testhelper.VerboseStdout())
	_, err = service.New(
		workerDeps,
		service.WithCheckConditions(false),
		service.WithCleanup(false),
		service.WithCloseSlices(true),
		service.WithUploadSlices(true),
		service.WithRetryFailedSlices(false),
		service.WithCloseFiles(false),
		service.WithImportFiles(false),
		service.WithRetryFailedFiles(false),
	)
	assert.NoError(t, err)

	// Get slices
	emptySlice, err := str.GetSlice(ctx, emptySliceKey)
	assert.NoError(t, err)
	slice1, err := str.GetSlice(ctx, slice1Key)
	assert.NoError(t, err)

	// Close the empty slice (a new slice is created)
	clk.Add(10 * time.Second)
	_, err = str.SwapSlice(ctx, &emptySlice)
	clk.Add(time.Minute) // trigger API nodes sync
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		count, err := apiDeps1.Schema().Slices().Uploaded().Count().Do(ctx, client)
		assert.NoError(t, err)
		return count == 1
	}, 10*time.Second, 100*time.Millisecond)
	clk.Add(10 * time.Second)

	// Close the slice1 (slice 2 is created)
	slice2, err := str.SwapSlice(ctx, &slice1)
	clk.Add(time.Minute) // trigger API nodes sync
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		count, err := apiDeps1.Schema().Slices().Uploaded().Count().Do(ctx, client)
		assert.NoError(t, err)
		return count == 2
	}, 10*time.Second, 100*time.Millisecond)

	// Check content of the uploaded slice
	AssertUploadedSlice(t, ctx, downloadCredentials, slice1, strings.TrimLeft(`
1,0001-01-01T00:02:02.000Z,1.2.3.4,"{""key"":""value001""}","{""Content-Type"":""application/json""}","""---value001---"""
2,0001-01-01T00:02:03.000Z,1.2.3.4,"{""key"":""value002""}","{""Content-Type"":""application/json""}","""---value002---"""
3,0001-01-01T00:02:04.000Z,1.2.3.4,"{""key"":""value003""}","{""Content-Type"":""application/json""}","""---value003---"""
4,0001-01-01T00:02:05.000Z,1.2.3.4,"{""key"":""value004""}","{""Content-Type"":""application/json""}","""---value004---"""
5,0001-01-01T00:02:06.000Z,1.2.3.4,"{""key"":""value005""}","{""Content-Type"":""application/json""}","""---value005---"""
6,0001-01-01T00:02:07.000Z,1.2.3.4,"{""key"":""value006""}","{""Content-Type"":""application/json""}","""---value006---"""
7,0001-01-01T00:02:08.000Z,1.2.3.4,"{""key"":""value007""}","{""Content-Type"":""application/json""}","""---value007---"""
`, "\n"))

	// Check content of the uploaded manifest
	AssertUploadedManifest(t, ctx, downloadCredentials, `
{"entries":[{"url":"%s00010101000101.gz"}]}
`)

	// Create some records also in the slice2
	clk.Add(time.Minute)
	createRecords(t, ctx, clk, apiDeps1, slice2.ReceiverKey, 8, 1)
	createRecords(t, ctx, clk, apiDeps2, slice2.ReceiverKey, 9, 2)

	// Close the slice2 (a new slice is created)
	_, err = str.SwapSlice(ctx, &slice2)
	clk.Add(time.Minute) // trigger API nodes sync
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		count, err := apiDeps1.Schema().Slices().Uploaded().Count().Do(ctx, client)
		assert.NoError(t, err)
		return count == 3
	}, 10*time.Second, 100*time.Millisecond)

	// Check content of the uploaded slice
	AssertUploadedSlice(t, ctx, downloadCredentials, slice2, strings.TrimLeft(`
8,0001-01-01T00:06:29.000Z,1.2.3.4,"{""key"":""value008""}","{""Content-Type"":""application/json""}","""---value008---"""
9,0001-01-01T00:06:30.000Z,1.2.3.4,"{""key"":""value009""}","{""Content-Type"":""application/json""}","""---value009---"""
10,0001-01-01T00:06:31.000Z,1.2.3.4,"{""key"":""value010""}","{""Content-Type"":""application/json""}","""---value010---"""
`, "\n"))

	// Check content of the uploaded manifest
	AssertUploadedManifest(t, ctx, downloadCredentials, `
{"entries":[{"url":"%s.gz"},{"url":"%s.gz"}]}
`)

	// Shutdown
	apiDeps1.Process().Shutdown(errors.New("bye bye API 1"))
	apiDeps1.Process().WaitForShutdown()
	apiDeps2.Process().Shutdown(errors.New("bye bye API 2"))
	apiDeps2.Process().WaitForShutdown()
	workerDeps.Process().Shutdown(errors.New("bye bye Worker"))
	workerDeps.Process().WaitForShutdown()

	// Check etcd state
	assertStateAfterUpload(t, client)

	// Check "close slice" logs
	wildcards.Assert(t, `
[task][slice.close/%s]INFO  started task "00000123/my-receiver-1/slice.close/%s"
[task][slice.close/%s]DEBUG  lock acquired "runtime/lock/task/slice.close/00000123/my-receiver-1/%s"
[task][slice.close/%s]INFO  waiting until all API nodes switch to a revision >= %s
[task][slice.close/%s]INFO  task succeeded (%s): slice closed
[task][slice.close/%s]DEBUG  lock released "runtime/lock/task/slice.close/00000123/my-receiver-1/%s"
`, strhelper.FilterLines(`^\[task\]\[slice.close\/0001-01-01T00:03:`, workerDeps.DebugLogger().AllMessages()))
	wildcards.Assert(t, `
[task][slice.close/%s]INFO  started task "00000123/my-receiver-2/slice.close/%s"
[task][slice.close/%s]DEBUG  lock acquired "runtime/lock/task/slice.close/00000123/my-receiver-2/%s"
[task][slice.close/%s]INFO  waiting until all API nodes switch to a revision >= %s
[task][slice.close/%s]INFO  task succeeded (%s): slice closed
[task][slice.close/%s]DEBUG  lock released "runtime/lock/task/slice.close/00000123/my-receiver-2/%s"
`, strhelper.FilterLines(`^\[task\]\[slice.close\/0001-01-01T00:04:`, workerDeps.DebugLogger().AllMessages()))
	wildcards.Assert(t, `
[task][slice.close/%s]INFO  started task "00000123/my-receiver-2/slice.close/%s"
[task][slice.close/%s]DEBUG  lock acquired "runtime/lock/task/slice.close/00000123/my-receiver-2/%s"
[task][slice.close/%s]INFO  waiting until all API nodes switch to a revision >= %s
[task][slice.close/%s]INFO  task succeeded (%s): slice closed
[task][slice.close/%s]DEBUG  lock released "runtime/lock/task/slice.close/00000123/my-receiver-2/%s"
`, strhelper.FilterLines(`^\[task\]\[slice.close\/0001-01-01T00:06:`, workerDeps.DebugLogger().AllMessages()))

	// Check "upload slice" logs
	wildcards.Assert(t, `
[task][slice.upload/%s]INFO  started task "00000123/my-receiver-1/slice.upload/%s"
[task][slice.upload/%s]DEBUG  lock acquired "runtime/lock/task/slice.upload/00000123/my-receiver-1/%s"
[task][slice.upload/%s]INFO  task succeeded (%s): skipped upload of the empty slice
[task][slice.upload/%s]DEBUG  lock released "runtime/lock/task/slice.upload/00000123/my-receiver-1/%s"
`, strhelper.FilterLines(`^\[task\]\[slice.upload\/0001-01-01T00:03:`, workerDeps.DebugLogger().AllMessages()))
	wildcards.Assert(t, `
[task][slice.upload/%s]INFO  started task "00000123/my-receiver-2/slice.upload/%s"
[task][slice.upload/%s]DEBUG  lock acquired "runtime/lock/task/slice.upload/00000123/my-receiver-2/%s"
[task][slice.upload/%s]INFO  task succeeded (%s): slice uploaded
[task][slice.upload/%s]DEBUG  lock released "runtime/lock/task/slice.upload/00000123/my-receiver-2/%s"
`, strhelper.FilterLines(`^\[task\]\[slice.upload\/0001-01-01T00:04:`, workerDeps.DebugLogger().AllMessages()))
	wildcards.Assert(t, `
[task][slice.upload/%s]INFO  started task "00000123/my-receiver-2/slice.upload/%s"
[task][slice.upload/%s]DEBUG  lock acquired "runtime/lock/task/slice.upload/00000123/my-receiver-2/%s"
[task][slice.upload/%s]INFO  task succeeded (%s): slice uploaded
[task][slice.upload/%s]DEBUG  lock released "runtime/lock/task/slice.upload/00000123/my-receiver-2/%s"
`, strhelper.FilterLines(`^\[task\]\[slice.upload\/0001-01-01T00:06:`, workerDeps.DebugLogger().AllMessages()))
}

func assertStateBeforeUpload(t *testing.T, client *etcd.Client) {
	t.Helper()
	etcdhelper.AssertKVsString(t, client, `
<<<<<
config/export/00000123/my-receiver-1/my-export-1
-----
%A
>>>>>

<<<<<
config/export/00000123/my-receiver-2/my-export-2
-----
%A
>>>>>

<<<<<
config/mapping/revision/00000123/my-receiver-1/my-export-1/00000001
-----
%A
>>>>>

<<<<<
config/mapping/revision/00000123/my-receiver-2/my-export-2/00000001
-----
%A
>>>>>

<<<<<
config/receiver/00000123/my-receiver-1
-----
%A
>>>>>

<<<<<
config/receiver/00000123/my-receiver-2
-----
%A
>>>>>

<<<<<
file/opened/00000123/my-receiver-1/my-export-1/0001-01-01T00:00:01.000Z
-----
%A
>>>>>

<<<<<
file/opened/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z
-----
%A
>>>>>

<<<<<
record/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:02:02.000Z_%s
-----
<<~~id~~>>,0001-01-01T00:02:02.000Z,1.2.3.4,"{""key"":""value001""}","{""Content-Type"":""application/json""}","""---value001---"""
>>>>>

<<<<<
record/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:02:03.000Z_%s
-----
<<~~id~~>>,0001-01-01T00:02:03.000Z,1.2.3.4,"{""key"":""value002""}","{""Content-Type"":""application/json""}","""---value002---"""
>>>>>

<<<<<
record/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:02:04.000Z_%s
-----
<<~~id~~>>,0001-01-01T00:02:04.000Z,1.2.3.4,"{""key"":""value003""}","{""Content-Type"":""application/json""}","""---value003---"""
>>>>>

<<<<<
record/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:02:05.000Z_%s
-----
<<~~id~~>>,0001-01-01T00:02:05.000Z,1.2.3.4,"{""key"":""value004""}","{""Content-Type"":""application/json""}","""---value004---"""
>>>>>

<<<<<
record/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:02:06.000Z_%s
-----
<<~~id~~>>,0001-01-01T00:02:06.000Z,1.2.3.4,"{""key"":""value005""}","{""Content-Type"":""application/json""}","""---value005---"""
>>>>>

<<<<<
record/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:02:07.000Z_%s
-----
<<~~id~~>>,0001-01-01T00:02:07.000Z,1.2.3.4,"{""key"":""value006""}","{""Content-Type"":""application/json""}","""---value006---"""
>>>>>

<<<<<
record/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:02:08.000Z_%s
-----
<<~~id~~>>,0001-01-01T00:02:08.000Z,1.2.3.4,"{""key"":""value007""}","{""Content-Type"":""application/json""}","""---value007---"""
>>>>>

<<<<<
runtime/api/node/watcher/cached/revision/api-node-1 (lease)
-----
%A
>>>>>

<<<<<
runtime/api/node/watcher/cached/revision/api-node-2 (lease)
-----
%A
>>>>>

<<<<<
secret/export/token/00000123/my-receiver-1/my-export-1
-----
%A
>>>>>

<<<<<
secret/export/token/00000123/my-receiver-2/my-export-2
-----
%A
>>>>>

<<<<<
slice/active/opened/writing/00000123/my-receiver-1/my-export-1/0001-01-01T00:00:01.000Z/0001-01-01T00:00:01.000Z
-----
{
  "projectId": 123,
  "receiverId": "my-receiver-1",
  "exportId": "my-export-1",
  "fileId": "0001-01-01T00:00:01.000Z",
  "sliceId": "0001-01-01T00:00:01.000Z",
  "state": "active/opened/writing",
  "mapping": {
%A
  },
  "storageResource": {
%A
  },
  "sliceNumber": 1
}
>>>>>

<<<<<
slice/active/opened/writing/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:01:01.000Z
-----
{
  "projectId": 123,
  "receiverId": "my-receiver-2",
  "exportId": "my-export-2",
  "fileId": "0001-01-01T00:01:01.000Z",
  "sliceId": "0001-01-01T00:01:01.000Z",
  "state": "active/opened/writing",
  "mapping": {
%A
  },
  "storageResource": {
%A
  },
  "sliceNumber": 1
}
>>>>>

<<<<<
stats/received/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:01:01.000Z/api-node-1
-----
{
  "projectId": 123,
  "receiverId": "my-receiver-2",
  "exportId": "my-export-2",
  "fileId": "0001-01-01T00:01:01.000Z",
  "sliceId": "0001-01-01T00:01:01.000Z",
  "nodeId": "api-node-1",
  "lastRecordAt": "0001-01-01T00:02:04.000Z",
  "recordsCount": 3,
  "recordsSize": "396B",
  "bodySize": "54B"
}
>>>>>

<<<<<
stats/received/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:01:01.000Z/api-node-2
-----
{
  "projectId": 123,
  "receiverId": "my-receiver-2",
  "exportId": "my-export-2",
  "fileId": "0001-01-01T00:01:01.000Z",
  "sliceId": "0001-01-01T00:01:01.000Z",
  "nodeId": "api-node-2",
  "lastRecordAt": "0001-01-01T00:02:08.000Z",
  "recordsCount": 4,
  "recordsSize": "528B",
  "bodySize": "72B"
}
>>>>>
`)
}

func assertStateAfterUpload(t *testing.T, client *etcd.Client) {
	t.Helper()
	etcdhelper.AssertKVsString(t, client, `
<<<<<
config/export/00000123/my-receiver-1/my-export-1
-----
%A
>>>>>

<<<<<
config/export/00000123/my-receiver-2/my-export-2
-----
%A
>>>>>

<<<<<
config/mapping/revision/00000123/my-receiver-1/my-export-1/00000001
-----
%A
>>>>>

<<<<<
config/mapping/revision/00000123/my-receiver-2/my-export-2/00000001
-----
%A
>>>>>

<<<<<
config/receiver/00000123/my-receiver-1
-----
%A
>>>>>

<<<<<
config/receiver/00000123/my-receiver-2
-----
%A
>>>>>

<<<<<
file/opened/00000123/my-receiver-1/my-export-1/0001-01-01T00:00:01.000Z
-----
%A
>>>>>

<<<<<
file/opened/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z
-----
%A
>>>>>

<<<<<
runtime/last/record/id/00000123/my-receiver-2/my-export-2
-----
10
>>>>>

<<<<<
secret/export/token/00000123/my-receiver-1/my-export-1
-----
%A
>>>>>

<<<<<
secret/export/token/00000123/my-receiver-2/my-export-2
-----
%A
>>>>>

<<<<<
slice/active/closed/uploaded/00000123/my-receiver-1/my-export-1/0001-01-01T00:00:01.000Z/0001-01-01T00:00:01.000Z
-----
{
  "projectId": 123,
  "receiverId": "my-receiver-1",
  "exportId": "my-export-1",
  "fileId": "0001-01-01T00:00:01.000Z",
  "sliceId": "0001-01-01T00:00:01.000Z",
  "state": "active/closed/uploaded",
  "mapping": {
%A
  },
  "storageResource": {
%A
  },
  "sliceNumber": 1,
  "closingAt": "%s",
  "uploadingAt": "%s",
  "uploadedAt": "%s",
  "isEmpty": true
}
>>>>>

<<<<<
slice/active/closed/uploaded/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:01:01.000Z
-----
{
  "projectId": 123,
  "receiverId": "my-receiver-2",
  "exportId": "my-export-2",
  "fileId": "0001-01-01T00:01:01.000Z",
  "sliceId": "0001-01-01T00:01:01.000Z",
  "state": "active/closed/uploaded",
  "mapping": {
%A
  },
  "storageResource": {
%A
  },
  "sliceNumber": 1,
  "closingAt": "%s",
  "uploadingAt": "%s",
  "uploadedAt": "%s",
  "statistics": {
    "lastRecordAt": "0001-01-01T00:02:08.000Z",
    "recordsCount": 7,
    "recordsSize": "924B",
    "bodySize": "126B",
    "fileSize": "861B",
    "fileGZipSize": "%s"
  },
  "idRange": {
    "start": 1,
    "count": 7
  }
}
>>>>>

<<<<<
slice/active/closed/uploaded/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:04:28.000Z
-----
{
  "projectId": 123,
  "receiverId": "my-receiver-2",
  "exportId": "my-export-2",
  "fileId": "0001-01-01T00:01:01.000Z",
  "sliceId": "0001-01-01T00:04:28.000Z",
  "state": "active/closed/uploaded",
  "mapping": {
%A
  },
  "storageResource": {
%A
  },
  "sliceNumber": 2,
  "closingAt": "%s",
  "uploadingAt": "%s",
  "uploadedAt": "%s",
  "statistics": {
    "lastRecordAt": "0001-01-01T00:06:31.000Z",
    "recordsCount": 3,
    "recordsSize": "396B",
    "bodySize": "54B",
    "fileSize": "370B",
    "fileGZipSize": "%s"
  },
  "idRange": {
    "start": 8,
    "count": 3
  }
}
>>>>>

<<<<<
slice/active/opened/writing/00000123/my-receiver-1/my-export-1/0001-01-01T00:00:01.000Z/0001-01-01T00:03:18.000Z
-----
{
  "projectId": 123,
  "receiverId": "my-receiver-1",
  "exportId": "my-export-1",
  "fileId": "0001-01-01T00:00:01.000Z",
  "sliceId": "0001-01-01T00:03:18.000Z",
  "state": "active/opened/writing",
  "mapping": {
%A
  },
  "storageResource": {
%A
  },
  "sliceNumber": 2
}
>>>>>

<<<<<
slice/active/opened/writing/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:06:31.000Z
-----
{
  "projectId": 123,
  "receiverId": "my-receiver-2",
  "exportId": "my-export-2",
  "fileId": "0001-01-01T00:01:01.000Z",
  "sliceId": "0001-01-01T00:06:31.000Z",
  "state": "active/opened/writing",
  "mapping": {
%A
  },
  "storageResource": {
%A
  },
  "sliceNumber": 3
}
>>>>>

<<<<<
task/00000123/my-receiver-1/slice.close/%s
-----
{
  "projectId": 123,
  "receiverId": "my-receiver-1",
  "type": "slice.close",
  "taskId": "%s",
  "createdAt": "%s",
  "finishedAt": "%s",
  "workerNode": "my-worker",
  "lock": "slice.close/00000123/my-receiver-1/my-export-1/0001-01-01T00:00:01.000Z/0001-01-01T00:00:01.000Z",
  "result": "slice closed",
  "duration": %d
}
>>>>>

<<<<<
task/00000123/my-receiver-1/slice.upload/%s
-----
{
  "projectId": 123,
  "receiverId": "my-receiver-1",
  "type": "slice.upload",
  "taskId": "%s",
  "createdAt": "%s",
  "finishedAt": "%s",
  "workerNode": "my-worker",
  "lock": "slice.upload/00000123/my-receiver-1/my-export-1/0001-01-01T00:00:01.000Z/0001-01-01T00:00:01.000Z",
  "result": "skipped upload of the empty slice",
  "duration": %d
}
>>>>>

<<<<<
task/00000123/my-receiver-2/slice.close/%s
-----
{
  "projectId": 123,
  "receiverId": "my-receiver-2",
  "type": "slice.close",
  "taskId": "%s",
  "createdAt": "%s",
  "finishedAt": "%s",
  "workerNode": "my-worker",
  "lock": "slice.close/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:01:01.000Z",
  "result": "slice closed",
  "duration": %d
}
>>>>>

<<<<<
task/00000123/my-receiver-2/slice.close/%s
-----
{
  "projectId": 123,
  "receiverId": "my-receiver-2",
  "type": "slice.close",
  "taskId": "%s",
  "createdAt": "%s",
  "finishedAt": "%s",
  "workerNode": "my-worker",
  "lock": "slice.close/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:04:28.000Z",
  "result": "slice closed",
  "duration": %d
}
>>>>>

<<<<<
task/00000123/my-receiver-2/slice.upload/%s
-----
{
  "projectId": 123,
  "receiverId": "my-receiver-2",
  "type": "slice.upload",
  "taskId": "%s",
  "createdAt": "%s",
  "finishedAt": "%s",
  "workerNode": "my-worker",
  "lock": "slice.upload/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:01:01.000Z",
  "result": "slice uploaded",
  "duration": %d
}
>>>>>

<<<<<
task/00000123/my-receiver-2/slice.upload/%s
-----
{
  "projectId": 123,
  "receiverId": "my-receiver-2",
  "type": "slice.upload",
  "taskId": "%s",
  "createdAt": "%s",
  "finishedAt": "%s",
  "workerNode": "my-worker",
  "lock": "slice.upload/00000123/my-receiver-2/my-export-2/0001-01-01T00:01:01.000Z/0001-01-01T00:04:28.000Z",
  "result": "slice uploaded",
  "duration": %d
}
>>>>>
`)
}

func AssertUploadedSlice(t *testing.T, ctx context.Context, file *keboola.FileDownloadCredentials, slice model.Slice, expected string) {
	t.Helper()

	// Read file content
	r, err := keboola.DownloadSliceReader(ctx, file, slice.Filename())
	assert.NoError(t, err)
	gz, err := gzip.NewReader(r)
	assert.NoError(t, err)
	data, err := io.ReadAll(gz)
	assert.NoError(t, err)
	assert.NoError(t, r.Close())
	assert.NoError(t, gz.Close())

	// Compare
	assert.Equal(t, expected, string(data))
}

func AssertUploadedManifest(t *testing.T, ctx context.Context, file *keboola.FileDownloadCredentials, expected string) {
	t.Helper()

	// Read file content
	r, err := keboola.DownloadSliceReader(ctx, file, keboola.ManifestFileName)
	assert.NoError(t, err)
	data, err := io.ReadAll(r)
	assert.NoError(t, err)
	assert.NoError(t, r.Close())

	// Compare
	wildcards.Assert(t, expected, string(data))
}
