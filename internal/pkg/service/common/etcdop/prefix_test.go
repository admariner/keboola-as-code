package etcdop

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	etcd "go.etcd.io/etcd/client/v3"

	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop/serde"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/etcdhelper"
)

func TestPrefix(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := etcdhelper.ClientForTest(t, etcdhelper.TmpNamespace(t))

	root := Prefix("my")
	pfx := root.Add("prefix")
	key0 := Key("key0")
	key1 := pfx.Key("key1")
	key2 := pfx.Key("key2")

	err := key0.Put("out of the prefix").Do(ctx, client)
	assert.NoError(t, err)

	// AtLeastOneExists - not found
	found, err := pfx.AtLeastOneExists().Do(ctx, client)
	assert.NoError(t, err)
	assert.False(t, found)

	// Count - 0
	count, err := pfx.Count().Do(ctx, client)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// GetAll - empty
	kvs, err := pfx.GetAll().Do(ctx, client).All()
	assert.NoError(t, err)
	assert.Empty(t, kvs)

	// GetOne - empty
	kv, err := pfx.GetOne(etcd.WithSort(etcd.SortByKey, etcd.SortDescend)).Do(ctx, client)
	assert.NoError(t, err)
	assert.Nil(t, kv)

	// DeleteAll - empty
	deleted, err := pfx.DeleteAll().Do(ctx, client)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), deleted)

	// ---
	err = key1.Put("foo").Do(ctx, client)
	assert.NoError(t, err)

	// AtLeastOneExists - found 1
	found, err = pfx.AtLeastOneExists().Do(ctx, client)
	assert.NoError(t, err)
	assert.True(t, found)

	// Count - 1
	count, err = pfx.Count().Do(ctx, client)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// GetAll - found 1
	kvs, err = pfx.GetAll().Do(ctx, client).All()
	assert.NoError(t, err)
	assert.NotEmpty(t, kvs)
	assert.Len(t, kvs, 1)
	assert.Equal(t, []byte("foo"), kvs[0].Value)

	// DeleteAll - deleted 1
	deleted, err = pfx.DeleteAll().Do(ctx, client)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), deleted)

	// ---
	err = key1.Put("foo").Do(ctx, client)
	assert.NoError(t, err)
	err = key2.Put("bar").Do(ctx, client)
	assert.NoError(t, err)

	// AtLeastOneExists - found 2
	found, err = pfx.AtLeastOneExists().Do(ctx, client)
	assert.NoError(t, err)
	assert.True(t, found)

	// Count - 2
	count, err = pfx.Count().Do(ctx, client)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)

	// GetAll - found 2
	kvs, err = pfx.GetAll().Do(ctx, client).All()
	assert.NoError(t, err)
	assert.NotEmpty(t, kvs)
	assert.Len(t, kvs, 2)
	assert.Equal(t, []byte("foo"), kvs[0].Value)
	assert.Equal(t, []byte("bar"), kvs[1].Value)

	// GetOne - 2 exists
	kv, err = pfx.GetOne(etcd.WithSort(etcd.SortByKey, etcd.SortDescend)).Do(ctx, client)
	assert.NoError(t, err)
	assert.Equal(t, []byte("my/prefix/key2"), kv.Key)
	assert.Equal(t, []byte("bar"), kv.Value)

	// DeleteAll - deleted 2
	deleted, err = pfx.DeleteAll().Do(ctx, client)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), deleted)
}

func TestTypedPrefix(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := etcdhelper.ClientForTest(t, etcdhelper.TmpNamespace(t))

	pfx := typedPrefixForTest()
	key0 := Key("key0")
	key1 := pfx.Key("key1")
	key2 := pfx.Key("key2")

	err := key0.Put("out of the prefix").Do(ctx, client)
	assert.NoError(t, err)

	// AtLeastOneExists - not found
	found, err := pfx.AtLeastOneExists().Do(ctx, client)
	assert.NoError(t, err)
	assert.False(t, found)

	// Count - 0
	count, err := pfx.Count().Do(ctx, client)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// GetAll - empty
	kvs, err := pfx.GetAll().Do(ctx, client).All()
	assert.NoError(t, err)
	assert.Empty(t, kvs)

	// GetOne - empty
	kv, err := pfx.GetOne(etcd.WithSort(etcd.SortByKey, etcd.SortDescend)).Do(ctx, client)
	assert.NoError(t, err)
	assert.Nil(t, kv)

	// DeleteAll - empty
	deleted, err := pfx.DeleteAll().Do(ctx, client)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), deleted)

	// ---
	err = key1.Put("foo").Do(ctx, client)
	assert.NoError(t, err)

	// AtLeastOneExists - found 1
	found, err = pfx.AtLeastOneExists().Do(ctx, client)
	assert.NoError(t, err)
	assert.True(t, found)

	// Count - 1
	count, err = pfx.Count().Do(ctx, client)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// GetAll - found 1
	kvs, err = pfx.GetAll().Do(ctx, client).All()
	assert.NoError(t, err)
	assert.NotEmpty(t, kvs)
	assert.Len(t, kvs, 1)
	assert.Equal(t, fooType("foo"), kvs[0].Value)

	// DeleteAll - deleted 1
	deleted, err = pfx.DeleteAll().Do(ctx, client)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), deleted)

	// ---
	err = key1.Put("foo").Do(ctx, client)
	assert.NoError(t, err)
	err = key2.Put("bar").Do(ctx, client)
	assert.NoError(t, err)

	// AtLeastOneExists - found 2
	found, err = pfx.AtLeastOneExists().Do(ctx, client)
	assert.NoError(t, err)
	assert.True(t, found)

	// Count - 2
	count, err = pfx.Count().Do(ctx, client)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)

	// GetAll - found 2
	kvs, err = pfx.GetAll().Do(ctx, client).All()
	assert.NoError(t, err)
	assert.NotEmpty(t, kvs)
	assert.Len(t, kvs, 2)
	assert.Equal(t, fooType("foo"), kvs[0].Value)
	assert.Equal(t, fooType("bar"), kvs[1].Value)

	// GetOne - 2 exists
	kv, err = pfx.GetOne(etcd.WithSort(etcd.SortByKey, etcd.SortDescend)).Do(ctx, client)
	assert.NoError(t, err)
	assert.Equal(t, "my/prefix/key2", kv.Key())
	assert.Equal(t, fooType("bar"), kv.Value)

	// DeleteAll - deleted 2
	deleted, err = pfx.DeleteAll().Do(ctx, client)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), deleted)
}

func BenchmarkPrefix_AtLestOneExists(b *testing.B) {
	ctx := context.Background()
	client := etcdhelper.ClientForTest(b, etcdhelper.TmpNamespace(b))

	pfx := Prefix("my/prefix/")
	if err := Key("key0").Put("foo").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}
	if err := pfx.Key("key1").Put("bar").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}
	if err := pfx.Key("key2").Put("baz").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		found, err := pfx.AtLeastOneExists().Do(ctx, client)
		if err != nil || !found {
			b.Fatalf("unexpected result")
		}
	}
}

func BenchmarkPrefix_Count(b *testing.B) {
	ctx := context.Background()
	client := etcdhelper.ClientForTest(b, etcdhelper.TmpNamespace(b))

	pfx := Prefix("my/prefix/")
	if err := Key("key0").Put("foo").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}
	if err := pfx.Key("key1").Put("bar").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}
	if err := pfx.Key("key2").Put("baz").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		count, err := pfx.Count().Do(ctx, client)
		if err != nil || count != int64(2) {
			b.Fatalf("unexpected result")
		}
	}
}

func BenchmarkPrefix_GetAll(b *testing.B) {
	ctx := context.Background()
	client := etcdhelper.ClientForTest(b, etcdhelper.TmpNamespace(b))

	pfx := Prefix("my/prefix/")
	if err := Key("key0").Put("foo").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key%04d", i)
		if err := pfx.Key(key).Put("bar").Do(ctx, client); err != nil {
			b.Fatalf(`cannot create etcd key "%s": %s`, key, err)
		}
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		kvs, err := pfx.GetAll().Do(ctx, client).All()
		if err != nil || len(kvs) != 100 {
			b.Fatalf("unexpected result")
		}
	}
}

func BenchmarkPrefix_DeleteAll(b *testing.B) {
	ctx := context.Background()
	client := etcdhelper.ClientForTest(b, etcdhelper.TmpNamespace(b))

	pfx := Prefix("my/prefix/")
	if err := Key("key0").Put("foo").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}
	if err := pfx.Key("key1").Put("bar").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}
	if err := pfx.Key("key2").Put("baz").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		deleted, err := pfx.DeleteAll().Do(ctx, client)
		if err != nil || (i == 0 != (deleted == 2)) { // xor
			b.Fatalf("unexpected result")
		}
	}
}

func BenchmarkPrefixT_AtLestOneExists(b *testing.B) {
	ctx := context.Background()
	client := etcdhelper.ClientForTest(b, etcdhelper.TmpNamespace(b))

	pfx := typedPrefixForTest()
	if err := Key("key0").Put("foo").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}
	if err := pfx.Key("key1").Put("bar").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}
	if err := pfx.Key("key2").Put("baz").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		found, err := pfx.AtLeastOneExists().Do(ctx, client)
		if err != nil || !found {
			b.Fatalf("unexpected result")
		}
	}
}

func BenchmarkPrefixT_Count(b *testing.B) {
	ctx := context.Background()
	client := etcdhelper.ClientForTest(b, etcdhelper.TmpNamespace(b))

	pfx := typedPrefixForTest()
	if err := Key("key0").Put("foo").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}
	if err := pfx.Key("key1").Put("bar").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}
	if err := pfx.Key("key2").Put("baz").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		count, err := pfx.Count().Do(ctx, client)
		if err != nil || count != int64(2) {
			b.Fatalf("unexpected result")
		}
	}
}

func BenchmarkPrefixT_GetAll(b *testing.B) {
	ctx := context.Background()
	client := etcdhelper.ClientForTest(b, etcdhelper.TmpNamespace(b))

	pfx := typedPrefixForTest()
	if err := Key("key0").Put("foo").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key%04d", i)
		if err := pfx.Key(key).Put("bar").Do(ctx, client); err != nil {
			b.Fatalf(`cannot create etcd key "%s": %s`, key, err)
		}
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		kvs, err := pfx.GetAll().Do(ctx, client).All()
		if err != nil || len(kvs) != 100 {
			b.Fatalf("unexpected result")
		}
	}
}

func BenchmarkPrefixT_DeleteAll(b *testing.B) {
	ctx := context.Background()
	client := etcdhelper.ClientForTest(b, etcdhelper.TmpNamespace(b))

	pfx := typedPrefixForTest()
	if err := Key("key0").Put("foo").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}
	if err := pfx.Key("key1").Put("bar").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}
	if err := pfx.Key("key2").Put("baz").Do(ctx, client); err != nil {
		b.Fatalf("cannot create etcd key: %s", err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		deleted, err := pfx.DeleteAll().Do(ctx, client)
		if err != nil || (i == 0 != (deleted == 2)) { // xor
			b.Fatalf("unexpected result")
		}
	}
}

func prefixForTest() Prefix {
	return Prefix("my").Add("prefix")
}

func typedPrefixForTest() PrefixT[fooType] {
	return PrefixT[fooType]{
		prefix: prefixForTest(),
		serde:  serde.NewJSON(serde.NoValidation),
	}
}
