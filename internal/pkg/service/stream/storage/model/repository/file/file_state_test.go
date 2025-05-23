package file_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/keboola/keboola-sdk-go/v2/pkg/keboola"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/client/v3/concurrency"

	commonDeps "github.com/keboola/keboola-as-code/internal/pkg/service/common/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop/iterator"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/utctime"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/definition/key"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/model"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/test"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/test/dummy"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/etcdhelper"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/etcdlogger"
)

func TestFileRepository_StateTransition(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	clk := clockwork.NewFakeClockAt(utctime.MustParse("2000-01-01T01:00:00.000Z").Time())
	by := test.ByUser()

	// Fixtures
	projectID := keboola.ProjectID(123)
	branchKey := key.BranchKey{ProjectID: projectID, BranchID: 456}
	sourceKey := key.SourceKey{BranchKey: branchKey, SourceID: "my-source"}
	sinkKey := key.SinkKey{SourceKey: sourceKey, SinkID: "my-sink"}

	// Get services
	d, mocked := dependencies.NewMockedStorageScope(t, ctx, commonDeps.WithClock(clk))
	client := mocked.TestEtcdClient()
	defRepo := d.DefinitionRepository()
	storageRepo := d.StorageRepository()
	fileRepo := storageRepo.File()
	sliceRepo := storageRepo.Slice()
	volumeRepo := storageRepo.Volume()

	// Log etcd operations
	var etcdLogs bytes.Buffer
	rawClient := d.EtcdClient()
	rawClient.KV = etcdlogger.KVLogWrapper(rawClient.KV, &etcdLogs, etcdlogger.WithMinimal())

	// Register active volumes
	// -----------------------------------------------------------------------------------------------------------------
	{
		session, err := concurrency.NewSession(client)
		require.NoError(t, err)
		defer func() { require.NoError(t, session.Close()) }()
		test.RegisterWriterVolumes(t, ctx, volumeRepo, session, 1)
	}

	// Create parent branch, source, sink, file, slice
	// -----------------------------------------------------------------------------------------------------------------
	var fileKey model.FileKey
	{
		branch := test.NewBranch(branchKey)
		require.NoError(t, defRepo.Branch().Create(&branch, clk.Now(), by).Do(ctx).Err())
		source := test.NewSource(sourceKey)
		require.NoError(t, defRepo.Source().Create(&source, clk.Now(), by, "Create source").Do(ctx).Err())
		sink := dummy.NewSinkWithLocalStorage(sinkKey)
		require.NoError(t, defRepo.Sink().Create(&sink, clk.Now(), by, "Create sink").Do(ctx).Err())
		fileKey = model.FileKey{SinkKey: sinkKey, FileID: model.FileID{OpenedAt: utctime.From(clk.Now())}}
	}

	// Create the second file, the first file is switched to the Closing state
	// -----------------------------------------------------------------------------------------------------------------
	{
		clk.Advance(time.Hour)
		require.NoError(t, fileRepo.Rotate(sinkKey, clk.Now()).Do(ctx).Err())
	}

	// Mark all slices as Uploading
	// -----------------------------------------------------------------------------------------------------------------
	{
		clk.Advance(time.Hour)
		require.NoError(t, sliceRepo.ListIn(fileKey).ForEach(func(s model.Slice, header *iterator.Header) error {
			return sliceRepo.SwitchToUploading(s.SliceKey, clk.Now(), false).Do(ctx).Err()
		}).Do(ctx).Err())
	}

	// Mark all slices as Uploaded
	// -----------------------------------------------------------------------------------------------------------------
	{
		clk.Advance(time.Hour)
		require.NoError(t, sliceRepo.ListIn(fileKey).ForEach(func(s model.Slice, header *iterator.Header) error {
			return sliceRepo.SwitchToUploaded(s.SliceKey, clk.Now()).Do(ctx).Err()
		}).Do(ctx).Err())
	}

	// Switch file to the Importing state
	// -----------------------------------------------------------------------------------------------------------------
	{
		clk.Advance(time.Hour)
		require.NoError(t, fileRepo.SwitchToImporting(fileKey, clk.Now(), false).Do(ctx).Err())
	}

	// Switch file to the Imported state
	// -----------------------------------------------------------------------------------------------------------------
	var stateSwitchEtcdLogs string
	{
		clk.Advance(time.Hour)
		etcdLogs.Reset()
		require.NoError(t, fileRepo.SwitchToImported(fileKey, clk.Now()).Do(ctx).Err())
		stateSwitchEtcdLogs = etcdLogs.String()
	}

	// Check etcd logs
	// -----------------------------------------------------------------------------------------------------------------

	etcdlogger.AssertFromFile(t, `fixtures/file_state_ops_001.txt`, stateSwitchEtcdLogs)

	// Check etcd state - there is no file
	// -----------------------------------------------------------------------------------------------------------------
	etcdhelper.AssertKVsFromFile(t, client, `fixtures/file_state_snapshot_001.txt`, etcdhelper.WithIgnoredKeyPattern("^definition/|storage/file/all/|storage/slice/|storage/secret/|storage/volume/|storage/stats/"))
}
