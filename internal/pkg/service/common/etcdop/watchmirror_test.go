package etcdop

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop/op"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop/serde"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/etcdhelper"
)

type testUser struct {
	FirstName string
	LastName  string
	Age       int
}

func TestMirror(t *testing.T) {
	t.Parallel()

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Create a typed prefix with some keys
	client := etcdhelper.ClientForTest(t, etcdhelper.TmpNamespace(t))
	pfx := NewTypedPrefix[testUser]("my/prefix", serde.NewJSON(serde.NoValidation))
	require.NoError(t, pfx.Key("key1").Put(testUser{FirstName: "John", LastName: "Brown", Age: 10}).DoOrErr(ctx, client))
	require.NoError(t, pfx.Key("key2").Put(testUser{FirstName: "Paul", LastName: "Green", Age: 20}).DoOrErr(ctx, client))

	// Setup mirroring of a prefix tree to the memory, with custom key and value mapping.
	// The result are in-memory KV pairs "<first name> <last name>" => <age>.
	logger := log.NewDebugLogger()
	mirror, errCh := SetupMirror(
		logger,
		pfx.GetAllAndWatch(ctx, client, clientv3.WithPrevKV()),
		func(kv *op.KeyValue, v testUser) string { return v.FirstName + " " + v.LastName },
		func(kv *op.KeyValue, v testUser) int { return v.Age },
	).
		WithFilter(func(event WatchEventT[testUser]) bool {
			return !strings.Contains(event.Kv.String(), "/ignore")
		}).
		StartMirroring(wg)

	// waitForSync:  it waits until the memory mirror is synchronized with the revision of the last change
	var header *op.Header
	waitForSync := func() {
		assert.Eventually(t, func() bool { return mirror.Revision() >= header.Revision }, time.Second, 100*time.Millisecond)
	}

	// Wait for initialization
	require.NoError(t, <-errCh)

	// Test state after initialization
	mirror.ToMap()

	// Insert
	header, err := pfx.Key("key3").Put(testUser{FirstName: "Luke", LastName: "Blue", Age: 30}).DoWithHeader(ctx, client)
	require.NoError(t, err)
	waitForSync()
	assert.Equal(t, map[string]int{
		"John Brown": 10,
		"Paul Green": 20,
		"Luke Blue":  30,
	}, mirror.ToMap())

	// Update
	header, err = pfx.Key("key1").Put(testUser{FirstName: "Jacob", LastName: "Brown", Age: 15}).DoWithHeader(ctx, client)
	require.NoError(t, err)
	waitForSync()
	assert.Equal(t, map[string]int{
		"Jacob Brown": 15,
		"Paul Green":  20,
		"Luke Blue":   30,
	}, mirror.ToMap())

	// Delete
	header, err = pfx.Key("key2").Delete().DoWithHeader(ctx, client)
	require.NoError(t, err)
	waitForSync()
	assert.Equal(t, map[string]int{
		"Jacob Brown": 15,
		"Luke Blue":   30,
	}, mirror.ToMap())

	// Filter
	header, err = pfx.Key("ignore").Put(testUser{FirstName: "Ignored", LastName: "User", Age: 50}).DoWithHeader(ctx, client)
	require.NoError(t, err)
	waitForSync()
	assert.Equal(t, map[string]int{
		"Jacob Brown": 15,
		"Luke Blue":   30,
	}, mirror.ToMap())
}
