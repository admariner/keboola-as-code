//nolint:forbidigo
package api

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/env"
	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/etcdhelper"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/testhelper"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/testhelper/runner"
)

const (
	instanceIDPlaceholder = "<<INSTANCE_ID>>"
)

// TestTemplatesApiE2E runs one Templates API functional test per each subdirectory.
func TestTemplatesApiE2E(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("Skipping API E2E tests on Windows")
	}

	binaryPath := testhelper.CompileBinary(t, "templates-api", "build-templates-api")

	runner.
		NewRunner(t).
		ForEachTest(func(test *runner.Test) {
			var repositories string
			if test.TestDirFS().Exists("repository") {
				repositories = fmt.Sprintf("keboola|file://%s", filepath.Join(test.TestDirFS().BasePath(), "repository"))
			} else {
				repositories = "keboola|https://github.com/keboola/keboola-as-code-templates.git|main"
			}
			addArgs := []string{fmt.Sprintf("--repositories=%s", repositories)}

			// Connect to the etcd
			etcdCredentials := etcdhelper.TmpNamespaceFromEnv(t, "TEMPLATES_API_ETCD_")
			etcdClient := etcdhelper.ClientForTest(t, etcdCredentials)

			addEnvs := env.FromMap(map[string]string{
				"TEMPLATES_API_DATADOG_ENABLED":  "false",
				"TEMPLATES_API_STORAGE_API_HOST": test.TestProject().StorageAPIHost(),
				"TEMPLATES_API_ETCD_NAMESPACE":   etcdCredentials.Namespace,
				"TEMPLATES_API_ETCD_ENDPOINT":    etcdCredentials.Endpoint,
				"TEMPLATES_API_ETCD_USERNAME":    etcdCredentials.Username,
				"TEMPLATES_API_ETCD_PASSWORD":    etcdCredentials.Password,
				"TEMPLATES_API_PUBLIC_ADDRESS":   "https://templates.keboola.local",
			})

			requestDecoratorFn := func(request *runner.APIRequestDef) {
				// Replace placeholder by instance ID.
				if strings.Contains(request.Path, instanceIDPlaceholder) {
					result := make(map[string]any)
					_, err := test.
						APIClient().
						R().
						SetResult(&result).
						SetHeader("X-StorageApi-Token", test.TestProject().StorageAPIToken().Token).
						Get("/v1/project/default/instances")

					instances := result["instances"].([]any)
					if assert.NoError(t, err) && assert.Equal(t, 1, len(instances)) {
						instanceId := instances[0].(map[string]any)["instanceId"].(string)
						request.Path = strings.ReplaceAll(request.Path, instanceIDPlaceholder, instanceId)
					}
				}
			}

			// Run the test
			test.Run(
				runner.WithInitProjectState(),
				runner.WithRunAPIServerAndRequests(
					binaryPath,
					addArgs,
					addEnvs,
					requestDecoratorFn,
				),
				runner.WithAssertProjectState(),
			)

			// Write current etcd KVs
			etcdDump, err := etcdhelper.DumpAllToString(context.Background(), etcdClient)
			assert.NoError(test.T(), err)
			assert.NoError(test.T(), test.WorkingDirFS().WriteFile(filesystem.NewRawFile("actual-etcd-kvs.txt", etcdDump)))

			// Assert current etcd state against expected state.
			expectedEtcdKVsPath := "expected-etcd-kvs.txt"
			if test.TestDirFS().IsFile(expectedEtcdKVsPath) {
				// Compare expected and actual kvs
				etcdhelper.AssertKVsString(test.T(), etcdClient, test.ReadFileFromTestDir(expectedEtcdKVsPath))
			}
		})
}
