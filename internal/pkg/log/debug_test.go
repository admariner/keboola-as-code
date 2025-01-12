package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDebugLogger_All(t *testing.T) {
	t.Parallel()
	logger := NewDebugLogger()
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")
	assert.Equal(t, "DEBUG  debug\nINFO  info\nWARN  warn\nERROR  error\n", logger.AllMessages())
}

func TestNewDebugLogger_Debug(t *testing.T) {
	t.Parallel()
	logger := NewDebugLogger()
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Errorf("error")
	assert.Equal(t, "DEBUG  debug\n", logger.DebugMessages())
}

func TestNewDebugLogger_Info(t *testing.T) {
	t.Parallel()
	logger := NewDebugLogger()
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Errorf("error")
	assert.Equal(t, "INFO  info\n", logger.InfoMessages())
}

func TestNewDebugLogger_Warn(t *testing.T) {
	t.Parallel()
	logger := NewDebugLogger()
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Errorf("error")
	assert.Equal(t, "WARN  warn\n", logger.WarnMessages())
}

func TestNewDebugLogger_WarnOrError(t *testing.T) {
	t.Parallel()
	logger := NewDebugLogger()
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")
	assert.Equal(t, "WARN  warn\nERROR  error\n", logger.WarnAndErrorMessages())
}

func TestNewDebugLogger_Error(t *testing.T) {
	t.Parallel()
	logger := NewDebugLogger()
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Errorf("error")
	assert.Equal(t, "ERROR  error\n", logger.ErrorMessages())
}
