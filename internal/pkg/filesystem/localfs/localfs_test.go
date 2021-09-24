package localfs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLocalFs(t *testing.T) {
	projectDir := t.TempDir()
	fs := New(projectDir)
	assert.Equal(t, projectDir, fs.BasePath())
}
