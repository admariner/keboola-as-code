package etcdhelper_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/utils/etcdhelper"
)

func TestDumpAll(t *testing.T) {
	t.Parallel()
	client := etcdhelper.ClientForTest(t, etcdhelper.TmpNamespace(t))

	// Put keys
	_, err := client.Put(context.Background(), "key1", "value1")
	assert.NoError(t, err)
	_, err = client.Put(context.Background(), "key2", "value2")
	assert.NoError(t, err)
	_, err = client.Put(context.Background(), "key3/key4", `{"foo1": "bar1", "foo2": ["bar2", "bar3"]}`)
	assert.NoError(t, err)

	// Dump
	dump, err := etcdhelper.DumpAllToString(context.Background(), client)
	assert.NoError(t, err)

	expected := `
<<<<<
key1
-----
value1
>>>>>

<<<<<
key2
-----
value2
>>>>>

<<<<<
key3/key4
-----
{
  "foo1": "bar1",
  "foo2": [
    "bar2",
    "bar3"
  ]
}
>>>>>
`
	assert.Equal(t, strings.TrimLeft(expected, "\n"), dump)
}
