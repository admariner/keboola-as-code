package storageapi_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/keboola/keboola-as-code/internal/pkg/api/client/storageapi"
	"github.com/keboola/keboola-as-code/internal/pkg/env"
	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/testproject"
)

func TestApiWithToken(t *testing.T) {
	t.Parallel()
	logger := log.NewDebugLogger()
	token := model.Token{Id: "123", Token: "mytoken", Owner: model.TokenOwner{Id: 456, Name: "name"}}
	orgApi := New(context.Background(), logger, "foo.bar.com", false)
	tokenApi := orgApi.WithToken(token)

	// Must be cloned, not modified
	assert.NotEqual(t, orgApi, tokenApi)
	assert.Equal(t, token, tokenApi.Token())
	assert.Equal(t, "mytoken", tokenApi.RestyClient().Header.Get("X-StorageApi-Token"))
}

func TestGetToken(t *testing.T) {
	t.Parallel()
	project := testproject.GetTestProject(t, env.Empty())
	logger := log.NewDebugLogger()
	api := New(context.Background(), logger, project.StorageAPIHost(), false)

	tokenValue := project.StorageAPIToken()
	token, err := api.GetToken(tokenValue)
	assert.NoError(t, err)
	assert.Regexp(t, `DEBUG  HTTP      GET https://.*/v2/storage/tokens/verify | 200 | .*`, logger.AllMessages())
	assert.Equal(t, tokenValue, token.Token)
	assert.Equal(t, project.ID(), token.ProjectId())
	assert.NotEmpty(t, token.ProjectName())
}

func TestGetTokenEmpty(t *testing.T) {
	t.Parallel()
	project := testproject.GetTestProject(t, env.Empty())
	logger := log.NewDebugLogger()
	api := New(context.Background(), logger, project.StorageAPIHost(), false)

	tokenValue := ""
	token, err := api.GetToken(tokenValue)
	assert.Error(t, err)
	apiErr := err.(*Error)
	assert.Equal(t, "Access token must be set", apiErr.Message)
	assert.Equal(t, "", apiErr.ErrCode)
	assert.Equal(t, 401, apiErr.StatusCode())
	assert.Empty(t, token)
}

func TestGetTokenInvalid(t *testing.T) {
	t.Parallel()
	project := testproject.GetTestProject(t, env.Empty())
	logger := log.NewDebugLogger()
	api := New(context.Background(), logger, project.StorageAPIHost(), false)

	tokenValue := "mytoken"
	token, err := api.GetToken(tokenValue)
	assert.Error(t, err)
	apiErr := err.(*Error)
	assert.Equal(t, "Invalid access token", apiErr.Message)
	assert.Equal(t, "storage.tokenInvalid", apiErr.ErrCode)
	assert.Equal(t, 401, apiErr.StatusCode())
	assert.Empty(t, token)
}
