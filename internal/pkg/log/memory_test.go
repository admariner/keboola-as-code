// nolint:forbidigo // allow usage of the "zap" package
package log

import (
	"context"
	"testing"
)

func TestMemoryLogger(t *testing.T) {
	t.Parallel()

	mem := NewMemoryLogger()
	mem.DebugCtx(context.Background(), `Debug message.`)
	mem.InfoCtx(context.Background(), `Info message.`)
	memWithCtx := mem.With("key1", "value1", "key2", "value2")
	memWithCtx.DebugCtx(context.Background(), `Debug message.`)
	memWithCtx.InfoCtx(context.Background(), `Info message.`)

	target := NewDebugLogger()
	mem.CopyLogsTo(target)

	expected := `
{"level":"debug","message":"Debug message."}
{"level":"info","message":"Info message."}
{"level":"debug","message":"Debug message.","key1": "value1", "key2": "value2"}
{"level":"info","message":"Info message.","key1": "value1", "key2": "value2"}
`
	AssertJSONMessages(t, expected, target.AllMessages())
}
