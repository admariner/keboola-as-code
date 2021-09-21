package version

import (
	"context"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/build"
	"github.com/keboola/keboola-as-code/internal/pkg/json"
	"github.com/keboola/keboola-as-code/internal/pkg/utils"
)

func TestCheckIfLatestVersionDev(t *testing.T) {
	c, _ := createMockedChecker(t)
	err := c.CheckIfLatest(build.DevVersionValue)
	assert.NotNil(t, err)
	assert.Equal(t, `skipped, found dev build`, err.Error())
}

func TestCheckIfLatestVersionEqual(t *testing.T) {
	c, logs := createMockedChecker(t)
	err := c.CheckIfLatest(`v1.2.3`)
	assert.Nil(t, err)
	assert.NotContains(t, logs.String(), `WARN`)
}

func TestCheckIfLatestVersionGreater(t *testing.T) {
	c, logs := createMockedChecker(t)
	err := c.CheckIfLatest(`v1.2.5`)
	assert.Nil(t, err)
	assert.NotContains(t, logs.String(), `WARN`)
}

func TestCheckIfLatestVersionLess(t *testing.T) {
	c, logs := createMockedChecker(t)
	err := c.CheckIfLatest(`v1.2.2`)
	assert.Nil(t, err)
	assert.Contains(t, logs.String(), `WARN  WARNING: A new version "v1.2.3" is available.`)
}

func createMockedChecker(t *testing.T) (*checker, *utils.Writer) {
	t.Helper()

	logger, logs := utils.NewDebugLogger()
	c := NewChecker(context.Background(), logger)
	resty := c.api.GetRestyClient()

	// Version check are disabled in tests by default
	utils.MustSetEnv(EnvVersionCheck, "")

	// Set short retry delay in tests
	resty.RetryWaitTime = 1 * time.Millisecond
	resty.RetryMaxWaitTime = 1 * time.Millisecond

	// Mocked resty transport
	httpmock.Activate()
	httpmock.ActivateNonDefault(resty.GetClient())
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Mocked body
	body := `
[
  {
    "tag_name": "v1.2.4",
    "assets": []
  },
  {
    "tag_name": "v1.2.3",
    "assets": [
      {
         "id": 123
      }
    ]
  }
]
`
	// Mocked response
	bodyJson := make([]interface{}, 0)
	json.MustDecodeString(body, &bodyJson)
	responder, err := httpmock.NewJsonResponder(200, bodyJson)
	assert.NoError(t, err)
	httpmock.RegisterResponder("GET", `=~.+repos/keboola/keboola-as-code/releases.+`, responder)

	return c, logs
}