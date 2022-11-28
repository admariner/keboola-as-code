package dependencies

import (
	"net/http"
	"testing"

	"github.com/keboola/go-client/pkg/storageapi"
	"github.com/stretchr/testify/assert"
	goaHttp "goa.design/goa/v3/http"
)

func TestNewProjectDeps_MasterTokenRequiredError(t *testing.T) {
	t.Parallel()
	d := NewMockedDeps()
	token := storageapi.Token{IsMaster: false}
	_, err := newProjectDeps(d, d, token)
	assert.Error(t, err)
	assert.Equal(t, "a master token of a project administrator is required", err.Error())
	assert.Equal(t, http.StatusUnauthorized, err.(goaHttp.Statuser).StatusCode())
}