package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

func TestProcessPanic(t *testing.T) {
	t.Parallel()
	logger := log.NewDebugLogger()
	logFilePath := "/foo/bar.log"
	exitCode := ProcessPanic(context.Background(), errors.New("test"), logger, logFilePath)
	assert.Equal(t, 1, exitCode)

	log.AssertJSONMessages(t, `
{"level":"debug","message":"Unexpected panic: test"}
{"level":"debug","message":"Trace:%A"}
{"level":"info","message":"%ATo help us diagnose the problem you can send us a crash report.%A"}
`, logger.AllMessages())
}
