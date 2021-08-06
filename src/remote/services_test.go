package remote

import (
	"context"
	"fmt"
	"testing"

	"keboola-as-code/src/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetEncryptionApiUrl(t *testing.T) {
	logger, _ := utils.NewDebugLogger()
	api := NewStorageApi("connection.keboola.com", context.Background(), logger, false)
	encryptionApiUrl, _ := api.GetEncryptionApiUrl()

	assert.NotEmpty(t, encryptionApiUrl)
	assert.Equal(t, encryptionApiUrl, "https://encryption.keboola.com")
}

func TestErrorGetEncryptionApiUrl(t *testing.T) {
	logger, _ := utils.NewDebugLogger()
	api := NewStorageApi("connection.foobar.keboola.com", context.Background(), logger, false)
	_, err := api.GetEncryptionApiUrl()
	fmt.Printf("%v", err)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "no such host")
}