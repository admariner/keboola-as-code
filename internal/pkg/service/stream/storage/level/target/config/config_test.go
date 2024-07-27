package config_test

import (
	"testing"
	"time"

	"github.com/c2h5oh/datasize"

	"github.com/keboola/keboola-as-code/internal/pkg/service/common/duration"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/level/target/config"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/test/testvalidation"
)

func TestConfig_Validation(t *testing.T) {
	t.Parallel()

	overMaximumCfg := config.NewConfig()
	overMaximumCfg.Import.Trigger = config.ImportTrigger{
		Count:    10000000 + 1,
		Size:     datasize.MustParseString("500MB") + 1,
		Interval: duration.From(24*time.Hour + 1),
	}

	// Test cases
	cases := testvalidation.TestCases[config.Config]{
		{
			Name: "empty",
			ExpectedError: `
- "operator.checkInterval" is a required field
- "operator.fileRotationTimeout" is a required field
- "operator.fileCloseTimeout" is a required field
- "operator.fileImportTimeout" is a required field
- "import.trigger.count" is a required field
- "import.trigger.size" is a required field
- "import.trigger.interval" is a required field
`,
			Value: config.Config{},
		},
		{
			Name: "over maximum",
			ExpectedError: `
- "import.trigger.count" must be 10,000,000 or less
- "import.trigger.size" must be 500MB or less
- "import.trigger.interval" must be 24h0m0s or less
`,
			Value: overMaximumCfg,
		},
		{
			Name:  "default",
			Value: config.NewConfig(),
		},
	}

	// Run test cases
	cases.Run(t)
}
