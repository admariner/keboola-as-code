package iterator_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/keboola/go-utils/pkg/wildcards"
	"github.com/stretchr/testify/assert"
	etcd "go.etcd.io/etcd/client/v3"

	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop/iterator"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop/op"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/etcdhelper"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/etcdlogger"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/ioutil"
)

type testCase struct {
	name         string
	kvCount      int
	pageSize     int
	options      []iterator.Option
	expected     []result
	expectedLogs string
}

type result struct {
	key   string
	value string
}

func TestIterator(t *testing.T) {
	t.Parallel()

	cases := []testCase{
		{
			name:     "empty",
			kvCount:  0,
			pageSize: 3,
			expected: []result{},
			expectedLogs: `
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/", "some/prefix0")
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/", "some/prefix0") | rev: %d | count: 0 | %s
`,
		},
		{
			name:     "count 1, under page size",
			kvCount:  1,
			pageSize: 3,
			expected: []result{
				{key: "some/prefix/foo001", value: "bar001"},
			},
			expectedLogs: `
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/", "some/prefix0")
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/", "some/prefix0") | rev: %d | count: 1 | %s
`,
		},
		{
			name:     "count 1, equal to page size",
			kvCount:  1,
			pageSize: 1,
			expected: []result{
				{key: "some/prefix/foo001", value: "bar001"},
			},
			expectedLogs: `
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/", "some/prefix0")
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/", "some/prefix0") | rev: %d | count: 1 | %s
`,
		},
		{
			name:     "count 2, under page size",
			kvCount:  2,
			pageSize: 3,
			expected: []result{
				{key: "some/prefix/foo001", value: "bar001"},
				{key: "some/prefix/foo002", value: "bar002"},
			},
			expectedLogs: `
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/", "some/prefix0")
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/", "some/prefix0") | rev: %d | count: 2 | %s
`,
		},
		{
			name:     "count 3, equal to page size",
			kvCount:  3,
			pageSize: 3,
			expected: []result{
				{key: "some/prefix/foo001", value: "bar001"},
				{key: "some/prefix/foo002", value: "bar002"},
				{key: "some/prefix/foo003", value: "bar003"},
			},
			expectedLogs: `
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/", "some/prefix0")
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/", "some/prefix0") | rev: %d | count: 3 | %s
`,
		},
		{
			name:     "one on the second page",
			kvCount:  4,
			pageSize: 3,
			expected: []result{
				{key: "some/prefix/foo001", value: "bar001"},
				{key: "some/prefix/foo002", value: "bar002"},
				{key: "some/prefix/foo003", value: "bar003"},
				{key: "some/prefix/foo004", value: "bar004"},
			},
			expectedLogs: `
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/", "some/prefix0")
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/", "some/prefix0") | rev: %d | count: 4 | %s
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/foo004", "some/prefix0") | rev: %d
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/foo004", "some/prefix0") | rev: %d | count: 1 | %s
`,
		},
		{
			name:     "two on the second page",
			kvCount:  5,
			pageSize: 3,
			expected: []result{
				{key: "some/prefix/foo001", value: "bar001"},
				{key: "some/prefix/foo002", value: "bar002"},
				{key: "some/prefix/foo003", value: "bar003"},
				{key: "some/prefix/foo004", value: "bar004"},
				{key: "some/prefix/foo005", value: "bar005"},
			},
			expectedLogs: `
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/", "some/prefix0")
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/", "some/prefix0") | rev: %d | count: 5 | %s
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/foo004", "some/prefix0") | rev: %d
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/foo004", "some/prefix0") | rev: %d | count: 2 | %s
`,
		},
		{
			name:     "page size = 1",
			kvCount:  5,
			pageSize: 1,
			expected: []result{
				{key: "some/prefix/foo001", value: "bar001"},
				{key: "some/prefix/foo002", value: "bar002"},
				{key: "some/prefix/foo003", value: "bar003"},
				{key: "some/prefix/foo004", value: "bar004"},
				{key: "some/prefix/foo005", value: "bar005"},
			},
			expectedLogs: `
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/", "some/prefix0")
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/", "some/prefix0") | rev: %d | count: 5 | %s
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/foo002", "some/prefix0") | rev: %d
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/foo002", "some/prefix0") | rev: %d | count: 4 | %s
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/foo003", "some/prefix0") | rev: %d
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/foo003", "some/prefix0") | rev: %d | count: 3 | %s
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/foo004", "some/prefix0") | rev: %d
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/foo004", "some/prefix0") | rev: %d | count: 2 | %s
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/foo005", "some/prefix0") | rev: %d
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/foo005", "some/prefix0") | rev: %d | count: 1 | %s
`,
		},
		{
			name:     "WithFromSameRev = false",
			kvCount:  5,
			pageSize: 1,
			options:  []iterator.Option{iterator.WithFromSameRev(false)},
			expected: []result{
				{key: "some/prefix/foo001", value: "bar001"},
				{key: "some/prefix/foo002", value: "bar002"},
				{key: "some/prefix/foo003", value: "bar003"},
				{key: "some/prefix/foo004", value: "bar004"},
				{key: "some/prefix/foo005", value: "bar005"},
			},
			expectedLogs: `
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/", "some/prefix0")
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/", "some/prefix0") | rev: %d | count: 5 | %s
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/foo002", "some/prefix0")
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/foo002", "some/prefix0") | rev: %d | count: 4 | %s
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/foo003", "some/prefix0")
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/foo003", "some/prefix0") | rev: %d | count: 3 | %s
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/foo004", "some/prefix0")
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/foo004", "some/prefix0") | rev: %d | count: 2 | %s
ETCD_REQUEST[%d] ➡️  GET ["some/prefix/foo005", "some/prefix0")
ETCD_REQUEST[%d] ✔️️  GET ["some/prefix/foo005", "some/prefix0") | rev: %d | count: 1 | %s
`,
		},
	}

	for _, tc := range cases {
		var logs strings.Builder
		ctx := context.Background()
		client := etcdhelper.ClientForTest(t, etcdhelper.TmpNamespace(t))
		client.KV = etcdlogger.KVLogWrapper(client.KV, &logs)
		prefix := generateKVs(t, tc.kvCount, ctx, client)
		ops := append([]iterator.Option{iterator.WithPageSize(tc.pageSize)}, tc.options...)

		// Test iteration methods
		logs.Reset()
		actual := iterateAll(t, prefix.GetAll(ops...), ctx, client)
		assert.Equal(t, tc.expected, actual, tc.name)
		wildcards.Assert(t, tc.expectedLogs, logs.String(), tc.name)

		// Test All method
		logs.Reset()
		actualKvs, err := prefix.GetAll(ops...).Do(ctx, client).All()
		assert.NoError(t, err)
		actual = make([]result, 0)
		for _, kv := range actualKvs {
			actual = append(actual, result{key: string(kv.Key), value: string(kv.Value)})
		}
		assert.Equal(t, tc.expected, actual, tc.name)
		wildcards.Assert(t, tc.expectedLogs, logs.String(), tc.name)

		// Test ForEach method
		logs.Reset()
		itr := prefix.GetAll(ops...).Do(ctx, client)
		actual = make([]result, 0)
		assert.NoError(t, itr.ForEach(func(kv *op.KeyValue, header *iterator.Header) error {
			assert.NotNil(t, header)
			actual = append(actual, result{key: string(kv.Key), value: string(kv.Value)})
			return nil
		}))
		assert.Equal(t, tc.expected, actual, tc.name)
		wildcards.Assert(t, tc.expectedLogs, logs.String(), tc.name)
	}
}

func TestIterator_AllKeys(t *testing.T) {
	t.Parallel()

	var logs strings.Builder
	ctx := context.Background()
	client := etcdhelper.ClientForTest(t, etcdhelper.TmpNamespace(t))
	client.KV = etcdlogger.KVLogWrapper(client.KV, &logs)
	ops := []iterator.Option{iterator.WithPageSize(3)}

	// Generate keys
	prefix := etcdop.Prefix("") // <<<<< all!
	for i := 1; i <= 5; i++ {
		key := prefix.Key(fmt.Sprintf("foo/bar%03d", i))
		val := fmt.Sprintf("bar%03d", i)
		assert.NoError(t, key.Put(val).Do(ctx, client))
	}

	// Get all keys from the etcd
	logs.Reset()
	actualKvs, err := prefix.GetAll(ops...).Do(ctx, client).All()
	assert.NoError(t, err)
	actual := make([]result, 0)
	for _, kv := range actualKvs {
		actual = append(actual, result{key: string(kv.Key), value: string(kv.Value)})
	}
	assert.Equal(t, []result{
		{key: "foo/bar001", value: "bar001"},
		{key: "foo/bar002", value: "bar002"},
		{key: "foo/bar003", value: "bar003"},
		{key: "foo/bar004", value: "bar004"},
		{key: "foo/bar005", value: "bar005"},
	}, actual)
	wildcards.Assert(t, `
ETCD_REQUEST[%d] ➡️  GET ["<NUL>", "<NUL>")
ETCD_REQUEST[%d] ✔️️  GET ["<NUL>", "<NUL>") | rev: %d | count: 5 | %s
ETCD_REQUEST[%d] ➡️  GET ["foo/bar004", "<NUL>") | rev: %d
ETCD_REQUEST[%d] ✔️️  GET ["foo/bar004", "<NUL>") | rev: %d | count: 2 | %s
`, logs.String())
}

func TestIterator_Revision(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client := etcdhelper.ClientForTest(t, etcdhelper.TmpNamespace(t))

	prefix := etcdop.NewPrefix("some/prefix")

	// There are 3 keys
	assert.NoError(t, prefix.Key("foo001").Put("bar001").Do(ctx, client))
	assert.NoError(t, prefix.Key("foo002").Put("bar002").Do(ctx, client))
	assert.NoError(t, prefix.Key("foo003").Put("bar003").Do(ctx, client))

	// Get current revision
	r, err := prefix.Key("foo003").Get().Do(ctx, client)
	assert.NoError(t, err)
	revision := r.ModRevision

	// Add more keys
	assert.NoError(t, prefix.Key("foo004").Put("bar004").Do(ctx, client))
	assert.NoError(t, prefix.Key("foo005").Put("bar005").Do(ctx, client))

	// Get all WithRev
	var actual []result
	assert.NoError(
		t,
		prefix.
			GetAll(iterator.WithRev(revision)).Do(ctx, client).
			ForEach(func(kv *op.KeyValue, _ *iterator.Header) error {
				actual = append(actual, result{key: string(kv.Key), value: string(kv.Value)})
				return nil
			}),
	)

	// The iterator only sees the values in the revision
	assert.Equal(t, []result{
		{key: "some/prefix/foo001", value: "bar001"},
		{key: "some/prefix/foo002", value: "bar002"},
		{key: "some/prefix/foo003", value: "bar003"},
	}, actual)
}

func TestIterator_End(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client := etcdhelper.ClientForTest(t, etcdhelper.TmpNamespace(t))

	prefix := etcdop.NewPrefix("some/prefix")

	// There are 5 keys
	assert.NoError(t, prefix.Key("foo001").Put("bar001").Do(ctx, client))
	assert.NoError(t, prefix.Key("foo002").Put("bar002").Do(ctx, client))
	assert.NoError(t, prefix.Key("foo003").Put("bar003").Do(ctx, client))
	assert.NoError(t, prefix.Key("foo004").Put("bar004").Do(ctx, client))
	assert.NoError(t, prefix.Key("foo005").Put("bar005").Do(ctx, client))

	// Get all WithEnd, so only the first 3 keys are loaded
	var actual []result
	assert.NoError(
		t,
		prefix.
			GetAll(iterator.WithEnd("foo004")).Do(ctx, client).
			ForEach(func(kv *op.KeyValue, _ *iterator.Header) error {
				actual = append(actual, result{key: string(kv.Key), value: string(kv.Value)})
				return nil
			}),
	)

	// The iterator only sees the values in the revision
	assert.Equal(t, []result{
		{key: "some/prefix/foo001", value: "bar001"},
		{key: "some/prefix/foo002", value: "bar002"},
		{key: "some/prefix/foo003", value: "bar003"},
	}, actual)
}

func TestIterator_Value_UsedIncorrectly(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client := etcdhelper.ClientForTest(t, etcdhelper.TmpNamespace(t))
	prefix := generateKVs(t, 3, ctx, client)

	it := prefix.GetAll().Do(ctx, client)
	assert.PanicsWithError(t, "unexpected Value() call: Next() must be called first", func() {
		it.Value()
	})
}

func TestIterator_ForEachOp(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client := etcdhelper.ClientForTest(t, etcdhelper.TmpNamespace(t))
	out := ioutil.NewAtomicWriter()
	prefix := generateKVs(t, 5, ctx, client)

	// Define op
	getAllOp := prefix.GetAll(iterator.WithPageSize(2)).ForEachOp(func(value *op.KeyValue, header *iterator.Header) error {
		_, _ = out.WriteString(fmt.Sprintf("%s\n", string(value.Value)))
		return nil
	})

	// Run op
	tracker := op.NewTracker(client)
	assert.NoError(t, getAllOp.DoOrErr(ctx, tracker))

	// All requests can be tracked by the TrackerKV
	assert.Equal(t, []op.TrackedOp{
		{Type: op.GetOp, Key: []byte("some/prefix/"), RangeEnd: []byte("some/prefix0"), Count: 5},
		{Type: op.GetOp, Key: []byte("some/prefix/foo003"), RangeEnd: []byte("some/prefix0"), Count: 3},
		{Type: op.GetOp, Key: []byte("some/prefix/foo005"), RangeEnd: []byte("some/prefix0"), Count: 1},
	}, tracker.Operations())

	// All values have been received
	assert.Equal(t, strings.TrimSpace(`
bar001
bar002
bar003
bar004
bar005
`), strings.TrimSpace(out.String()))
}

func iterateAll(t *testing.T, def iterator.Definition, ctx context.Context, client *etcd.Client) []result {
	t.Helper()
	it := def.Do(ctx, client)
	actual := make([]result, 0)
	for it.Next() {
		kv := it.Value()
		actual = append(actual, result{key: string(kv.Key), value: string(kv.Value)})
	}
	assert.NoError(t, it.Err())
	return actual
}

func generateKVs(t *testing.T, count int, ctx context.Context, client *etcd.Client) etcdop.Prefix {
	t.Helper()

	// There are some keys before the prefix
	assert.NoError(t, etcdop.Key("some/abc").Put("foo").Do(ctx, client))

	// Create keys in the iterated prefix
	prefix := etcdop.NewPrefix("some/prefix")
	for i := 1; i <= count; i++ {
		key := prefix.Key(fmt.Sprintf("foo%03d", i))
		val := fmt.Sprintf("bar%03d", i)
		assert.NoError(t, key.Put(val).Do(ctx, client))
	}

	// There are some keys after the prefix
	assert.NoError(t, etcdop.Key("some/xyz").Put("foo").Do(ctx, client))

	return prefix
}
