package api

import (
	"context"
	"github.com/stretchr/testify/assert"
	"keboola-as-code/src/model"
	"keboola-as-code/src/tests"
	"keboola-as-code/src/utils"
	"testing"
)

func TestApiWithToken(t *testing.T) {
	logger, _ := utils.NewDebugLogger()
	token := &model.Token{Id: "123", Token: "mytoken", Owner: model.TokenOwner{Id: 456, Name: "name"}}
	orgApi := NewStorageApi("foo.bar.com", context.Background(), logger, false)
	tokenApi := orgApi.WithToken(token)

	// Must be cloned, not modified
	assert.NotSame(t, orgApi, tokenApi)
	assert.Same(t, token, tokenApi.token)
	assert.Equal(t, "mytoken", tokenApi.http.resty.Header.Get("X-StorageApi-Token"))
}

func TestGetToken(t *testing.T) {
	tokenValue := tests.TestToken()
	api, logs := newStorageApi(t)
	token, err := api.GetToken(tokenValue)
	assert.NoError(t, err)
	assert.Regexp(t, `DEBUG  HTTP      GET https://.*/v2/storage/tokens/verify | 200 | .*`, logs.String())
	assert.Equal(t, tokenValue, token.Token)
	assert.Equal(t, tests.TestProjectId(), token.ProjectId())
	assert.Equal(t, tests.TestProjectName(), token.ProjectName())
}

func TestGetTokenEmpty(t *testing.T) {
	tokenValue := ""
	api, _ := newStorageApi(t)
	token, err := api.GetToken(tokenValue)
	assert.Error(t, err)
	apiErr := err.(*Error)
	assert.Equal(t, "Access token must be set", apiErr.Message)
	assert.Equal(t, "", apiErr.ErrCode)
	assert.Equal(t, 401, apiErr.HttpStatus())
	assert.Nil(t, token)
}

func TestGetTokenInvalid(t *testing.T) {
	tokenValue := "mytoken"
	api, _ := newStorageApi(t)
	token, err := api.GetToken(tokenValue)
	assert.Error(t, err)
	apiErr := err.(*Error)
	assert.Equal(t, "Invalid access token", apiErr.Message)
	assert.Equal(t, "storage.tokenInvalid", apiErr.ErrCode)
	assert.Equal(t, 401, apiErr.HttpStatus())
	assert.Nil(t, token)
}

func TestGetTokenExpired(t *testing.T) {
	tokenValue := tests.TestTokenExpired()
	api, _ := newStorageApi(t)
	token, err := api.GetToken(tokenValue)
	assert.Error(t, err)
	apiErr := err.(*Error)
	assert.Equal(t, "Invalid access token", apiErr.Message)
	assert.Equal(t, "storage.tokenInvalid", apiErr.ErrCode)
	assert.Equal(t, 401, apiErr.HttpStatus())
	assert.Nil(t, token)
}
