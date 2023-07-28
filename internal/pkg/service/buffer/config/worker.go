package config

import (
	"time"

	"github.com/c2h5oh/datasize"

	"github.com/keboola/keboola-as-code/internal/pkg/env"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/store/model"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/cliconfig"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
	"github.com/keboola/keboola-as-code/internal/pkg/validator"
)

const (
	WorkerEnvPrefix = "BUFFER_WORKER_"
	// DefaultCheckConditionsInterval defines how often it will be checked upload and import conditions.
	DefaultCheckConditionsInterval = 5 * time.Second
	// DefaultMinimalUploadInterval defines minimal interval between two slice upload operations.
	DefaultMinimalUploadInterval = 5 * time.Second
	// DefaultMinimalImportInterval defines minimal interval between two import file operations.
	DefaultMinimalImportInterval = 30 * time.Second
	// DefaultCleanupInterval defines how often old tasks and files will be checked and deleted.
	DefaultCleanupInterval = 15 * time.Minute
)

// WorkerConfig of the Buffer Worker.
type WorkerConfig struct {
	ServiceConfig           `mapstructure:",squash"`
	UniqueID                string           `mapstructure:"unique-id" usage:"Unique process ID, auto-generated by default."`
	CheckConditionsInterval time.Duration    `mapstructure:"check-conditions-interval" usage:"How often will upload and import conditions be checked."`
	MinimalUploadInterval   time.Duration    `mapstructure:"min-upload-interval" usage:"Minimal interval between two slice upload operations."`
	MinimalImportInterval   time.Duration    `mapstructure:"min-import-interval" usage:"Minimal interval between two file import operations."`
	UploadConditions        model.Conditions `mapstructure:"upload-conditions"`
	MetricsListenAddress    string           `mapstructure:"metrics-listen-address" usage:"Prometheus /metrics HTTP endpoint listen address."`
	ConditionsCheck         bool             `mapstructure:"enable-conditions-check" usage:"Enable conditions check functionality."`
	CloseSlices             bool             `mapstructure:"enable-close-slices" usage:"Enable close slices functionality."`
	UploadSlices            bool             `mapstructure:"enable-upload-slices" usage:"Enable upload slices functionality."`
	RetryFailedSlices       bool             `mapstructure:"enable-retry-failed-slices" usage:"Enable retry for failed slices."`
	CloseFiles              bool             `mapstructure:"enable-close-files" usage:"Enable close files functionality."`
	ImportFiles             bool             `mapstructure:"enable-import-files" usage:"Enable import files functionality."`
	RetryFailedFiles        bool             `mapstructure:"enable-retry-failed-files" usage:"Enable retry for failed files."`
	TasksCleanup            bool             `mapstructure:"tasks-cleanup-enabled" usage:"Enable periodical tasks cleanup functionality."`
	TasksCleanupInterval    time.Duration    `mapstructure:"tasks-cleanup-interval" usage:"How often will old tasks be deleted."`
}

func NewWorkerConfig() WorkerConfig {
	return WorkerConfig{
		ServiceConfig:           NewServiceConfig(),
		UniqueID:                "",
		CheckConditionsInterval: DefaultCheckConditionsInterval,
		MinimalUploadInterval:   DefaultMinimalUploadInterval,
		MinimalImportInterval:   DefaultMinimalImportInterval,
		UploadConditions: model.Conditions{
			Count: 1000,
			Size:  1 * datasize.MB,
			Time:  1 * time.Minute,
		},
		MetricsListenAddress: "0.0.0.0:9000",
		ConditionsCheck:      true,
		CloseSlices:          true,
		UploadSlices:         true,
		RetryFailedSlices:    true,
		CloseFiles:           true,
		ImportFiles:          true,
		RetryFailedFiles:     true,
		TasksCleanup:         true,
		TasksCleanupInterval: DefaultCleanupInterval,
	}
}

type WorkerOption func(c *WorkerConfig)

func BindWorkerConfig(args []string, envs env.Provider) (WorkerConfig, error) {
	cfg := NewWorkerConfig()
	err := cfg.LoadFrom(args, envs)
	return cfg, err
}

func (c *WorkerConfig) LoadFrom(args []string, envs env.Provider) error {
	return cliconfig.LoadTo(c, args, envs, WorkerEnvPrefix)
}

func (c *WorkerConfig) Dump() string {
	if kvs, err := cliconfig.Dump(c); err != nil {
		panic(err)
	} else {
		return kvs.String()
	}
}

func (c *WorkerConfig) Normalize() {
	c.ServiceConfig.Normalize()
}

func (c *WorkerConfig) Validate() error {
	v := validator.New()
	errs := errors.NewMultiError()
	if err := c.ServiceConfig.Validate(); err != nil {
		errs.Append(err)
	}
	if c.CheckConditionsInterval <= 0 {
		errs.Append(errors.Errorf(`check conditions interval must be positive time.Duration, found "%v"`, c.CheckConditionsInterval))
	}
	if c.MinimalUploadInterval < c.CheckConditionsInterval {
		errs.Append(errors.Errorf(`minimal upload interval "%v" must be >= than check conditions interval "%v"`, c.MinimalUploadInterval, c.CheckConditionsInterval))
	}
	if c.MinimalImportInterval < c.CheckConditionsInterval {
		errs.Append(errors.Errorf(`minimal import interval "%v" must be >= than check conditions interval "%v"`, c.MinimalImportInterval, c.CheckConditionsInterval))
	}
	if c.CheckConditionsInterval <= 0 {
		errs.Append(errors.Errorf(`check conditions interval must be positive time.Duration, found "%v"`, c.CheckConditionsInterval))
	}
	if c.UploadConditions.Count <= 0 {
		errs.Append(errors.Errorf(`upload conditions count must be positive number, found "%v"`, c.UploadConditions.Count))
	}
	if c.UploadConditions.Time <= 0 {
		errs.Append(errors.Errorf(`upload conditions time must be positive time.Duration, found "%v"`, c.UploadConditions.Time.String()))
	}
	if c.UploadConditions.Size <= 0 {
		errs.Append(errors.Errorf(`upload conditions size must be positive number, found "%v"`, c.UploadConditions.Size.String()))
	}
	if c.MetricsListenAddress == "" {
		errs.Append(errors.New("metrics listen address is not set"))
	} else if err := v.ValidateValue(c.MetricsListenAddress, "hostname_port"); err != nil {
		errs.Append(errors.Errorf(`metrics address "%s" is not valid`, c.MetricsListenAddress))
	}
	if c.TasksCleanupInterval <= 0 {
		errs.Append(errors.Errorf(`tasks cleanup interval must be positive time.Duration, found "%v"`, c.TasksCleanupInterval))
	}
	return errs.ErrorOrNil()
}

func (c WorkerConfig) Apply(ops ...WorkerOption) WorkerConfig {
	for _, o := range ops {
		o(&c)
	}
	return c
}

func WithCleanupInterval(v time.Duration) WorkerOption {
	return func(c *WorkerConfig) {
		c.TasksCleanupInterval = v
	}
}

func WithCheckConditionsInterval(v time.Duration) WorkerOption {
	return func(c *WorkerConfig) {
		c.CheckConditionsInterval = v
	}
}

func WithUploadConditions(v model.Conditions) WorkerOption {
	return func(c *WorkerConfig) {
		c.UploadConditions = v
	}
}

// WithConditionsCheck enables/disables the conditions checker.
func WithConditionsCheck(v bool) WorkerOption {
	return func(c *WorkerConfig) {
		c.ConditionsCheck = v
	}
}

// WithCleanup enables/disables etcd cleanup task.
func WithCleanup(v bool) WorkerOption {
	return func(c *WorkerConfig) {
		c.TasksCleanup = v
	}
}

// WithCloseSlices enables/disables the "close slices" task.
func WithCloseSlices(v bool) WorkerOption {
	return func(c *WorkerConfig) {
		c.CloseSlices = v
	}
}

// WithUploadSlices enables/disables the "upload slices" task.
func WithUploadSlices(v bool) WorkerOption {
	return func(c *WorkerConfig) {
		c.UploadSlices = v
	}
}

// WithRetryFailedSlices enables/disables the "retry failed uploads" task.
func WithRetryFailedSlices(v bool) WorkerOption {
	return func(c *WorkerConfig) {
		c.RetryFailedSlices = v
	}
}

// WithCloseFiles enables/disables the "close files" task.
func WithCloseFiles(v bool) WorkerOption {
	return func(c *WorkerConfig) {
		c.CloseFiles = v
	}
}

// WithImportFiles enables/disables the "upload file" task.
func WithImportFiles(v bool) WorkerOption {
	return func(c *WorkerConfig) {
		c.ImportFiles = v
	}
}

// WithRetryFailedFiles enables/disables the "retry failed imports" task.
func WithRetryFailedFiles(v bool) WorkerOption {
	return func(c *WorkerConfig) {
		c.RetryFailedFiles = v
	}
}
