package source_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/keboola/keboola-sdk-go/v2/pkg/keboola"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	serviceErrors "github.com/keboola/keboola-as-code/internal/pkg/service/common/errors"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/utctime"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/definition"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/definition/key"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/test"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/etcdhelper"
)

func TestSourceRepository_Undelete(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	now := utctime.MustParse("2000-01-01T01:00:00.000Z").Time()
	by := test.ByUser()

	d, mocked := dependencies.NewMockedServiceScope(t, ctx)
	client := mocked.TestEtcdClient()
	repo := d.DefinitionRepository().Source()
	ignoredEtcdKeys := etcdhelper.WithIgnoredKeyPattern("^(definition/branch)")

	// Fixtures
	projectID := keboola.ProjectID(123)
	branchKey := key.BranchKey{ProjectID: projectID, BranchID: 567}
	sourceKey := key.SourceKey{BranchKey: branchKey, SourceID: "my-source"}

	// Create branch - ok
	// -----------------------------------------------------------------------------------------------------------------
	{
		branch := test.NewBranch(branchKey)
		require.NoError(t, d.DefinitionRepository().Branch().Create(&branch, now, by).Do(ctx).Err())
	}

	// Undelete - not found
	// -----------------------------------------------------------------------------------------------------------------
	if err := repo.Undelete(sourceKey, now, by).Do(ctx).Err(); assert.Error(t, err) {
		assert.Equal(t, `deleted source "my-source" not found in the branch`, err.Error())
		serviceErrors.AssertErrorStatusCode(t, http.StatusNotFound, err)
	}

	// Create - ok
	// -----------------------------------------------------------------------------------------------------------------
	{
		source := test.NewSource(sourceKey)
		require.NoError(t, repo.Create(&source, now, by, "Create source").Do(ctx).Err())
	}

	// Get - ok
	// -----------------------------------------------------------------------------------------------------------------
	{
		require.NoError(t, repo.Get(sourceKey).Do(ctx).Err())
	}

	// SoftDelete - ok
	// -----------------------------------------------------------------------------------------------------------------
	{
		now = now.Add(time.Hour)
		require.NoError(t, repo.SoftDelete(sourceKey, now, by).Do(ctx).Err())
		etcdhelper.AssertKVsFromFile(t, client, "fixtures/source_undelete_snapshot_001.txt", ignoredEtcdKeys)
	}

	// Get - not found
	// -----------------------------------------------------------------------------------------------------------------
	{
		err := repo.Get(sourceKey).Do(ctx).Err()
		if assert.Error(t, err) {
			assert.Equal(t, `source "my-source" not found in the branch`, err.Error())
		}
	}

	// Undelete - ok
	// -----------------------------------------------------------------------------------------------------------------
	{
		now = now.Add(time.Hour)
		require.NoError(t, repo.Undelete(sourceKey, now.Add(time.Hour), by).Do(ctx).Err())
		etcdhelper.AssertKVsFromFile(t, client, "fixtures/source_undelete_snapshot_002.txt", ignoredEtcdKeys)
	}

	// Get - ok
	// -----------------------------------------------------------------------------------------------------------------
	{
		require.NoError(t, repo.Get(sourceKey).Do(ctx).Err())
	}

	// GetDeleted - not found
	// -----------------------------------------------------------------------------------------------------------------
	{
		err := repo.GetDeleted(sourceKey).Do(ctx).Err()
		if assert.Error(t, err) {
			assert.Equal(t, `deleted source "my-source" not found in the branch`, err.Error())
		}
	}
}

func TestSourceRepository_UndeleteSourcesOnBranchUndelete(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	now := utctime.MustParse("2000-01-01T01:00:00.000Z").Time()
	by := test.ByUser()

	d, mocked := dependencies.NewMockedServiceScope(t, ctx)
	client := mocked.TestEtcdClient()
	repo := d.DefinitionRepository().Source()
	ignoredEtcdKeys := etcdhelper.WithIgnoredKeyPattern("^(definition/source/version)")

	// Fixtures
	projectID := keboola.ProjectID(123)
	branchKey := key.BranchKey{ProjectID: projectID, BranchID: 567}
	sourceKey1 := key.SourceKey{BranchKey: branchKey, SourceID: "my-source-1"}
	sourceKey2 := key.SourceKey{BranchKey: branchKey, SourceID: "my-source-2"}
	sourceKey3 := key.SourceKey{BranchKey: branchKey, SourceID: "my-source-3"}

	// Create branch
	// -----------------------------------------------------------------------------------------------------------------
	{
		branch := test.NewBranch(branchKey)
		require.NoError(t, d.DefinitionRepository().Branch().Create(&branch, now, by).Do(ctx).Err())
	}

	// Create sources
	// -----------------------------------------------------------------------------------------------------------------
	var source1, source2, source3 definition.Source
	{
		source1 = test.NewSource(sourceKey1)
		require.NoError(t, repo.Create(&source1, now, by, "Create source").Do(ctx).Err())
		source2 = test.NewSource(sourceKey2)
		require.NoError(t, repo.Create(&source2, now, by, "Create source").Do(ctx).Err())
		source3 = test.NewSource(sourceKey3)
		require.NoError(t, repo.Create(&source3, now, by, "Create source").Do(ctx).Err())
	}

	// Delete Source1
	// -----------------------------------------------------------------------------------------------------------------
	{
		now = now.Add(time.Hour)
		require.NoError(t, repo.SoftDelete(sourceKey1, now, by).Do(ctx).Err())
	}

	// Delete Branch
	// -----------------------------------------------------------------------------------------------------------------
	{
		now = now.Add(time.Hour)
		require.NoError(t, d.DefinitionRepository().Branch().SoftDelete(branchKey, now, by).Do(ctx).Err())
	}

	// Undelete Branch
	// -----------------------------------------------------------------------------------------------------------------
	{
		now = now.Add(time.Hour)
		require.NoError(t, d.DefinitionRepository().Branch().Undelete(branchKey, now, by).Do(ctx).Err())
		etcdhelper.AssertKVsFromFile(t, client, "fixtures/source_undelete_snapshot_003.txt", ignoredEtcdKeys)
	}
	{
		var err error

		// Source1 has been deleted before the Branch deletion, so it remains deleted.
		source1, err = repo.GetDeleted(sourceKey1).Do(ctx).ResultOrErr()
		require.NoError(t, err)
		assert.True(t, source1.IsDeleted())

		// Source2 and Source2
		source2, err = repo.Get(sourceKey2).Do(ctx).ResultOrErr()
		require.NoError(t, err)
		assert.False(t, source2.IsDeleted())
		source3, err = repo.Get(sourceKey3).Do(ctx).ResultOrErr()
		require.NoError(t, err)
		assert.False(t, source3.IsDeleted())
	}
}
