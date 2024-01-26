package target

import (
	"time"

	"github.com/c2h5oh/datasize"
)

// Config configures the target storage.
type Config struct {
	Import ImportConfig `configKey:"import"`
}

// ConfigPatch is same as the Config, but with optional/nullable fields.
// It is part of the definition.TableSink structure to allow modification of the default configuration.
type ConfigPatch struct {
	Import *ImportConfigPatch `json:"import,omitempty"`
}

func NewConfig() Config {
	return Config{
		Import: ImportConfig{
			MinInterval: 1 * time.Minute,
			Trigger: ImportTrigger{
				Count:    50000,
				Size:     5 * datasize.MB,
				Interval: 5 * time.Minute,
			},
		},
	}
}

// With copies values from the ConfigPatch, if any.
func (c Config) With(v ConfigPatch) Config {
	if v.Import != nil {
		c.Import = c.Import.With(*v.Import)
	}
	return c
}

// ---------------------------------------------------------------------------------------------------------------------

// ImportConfig configures the file import.
type ImportConfig struct {
	MinInterval time.Duration `configKey:"minInterval" configUsage:"Minimal interval between imports." validate:"required,minDuration=30s,maxDuration=30m"`
	Trigger     ImportTrigger `configKey:"trigger"`
}

// ImportConfigPatch is same as the ImportConfig, but with optional/nullable fields.
// It is part of the definition.TableSink structure to allow modification of the default configuration.
type ImportConfigPatch struct {
	Trigger *ImportTriggerPatch `json:"trigger,omitempty"`
}

// With copies values from the ConfigPatch, if any.
func (c ImportConfig) With(v ImportConfigPatch) ImportConfig {
	if v.Trigger != nil {
		c.Trigger = c.Trigger.With(*v.Trigger)
	}
	return c
}

// ---------------------------------------------------------------------------------------------------------------------

// ImportTrigger configures file import conditions, at least one must be met.
type ImportTrigger struct {
	Count    uint64            `json:"count" configKey:"count" configUsage:"Records count." validate:"required,min=1,max=10000000"`
	Size     datasize.ByteSize `json:"size" configKey:"size" configUsage:"Records size." validate:"required,minBytes=100B,maxBytes=500MB"`
	Interval time.Duration     `json:"interval" configKey:"interval" configUsage:"Duration from the last import." validate:"required,minDuration=60s,maxDuration=24h"`
}

// ImportTriggerPatch is same as the ImportTrigger, but with optional/nullable fields.
// It is part of the definition.TableSink structure to allow modification of the default configuration.
type ImportTriggerPatch struct {
	Count    *uint64            `json:"count,omitempty" configKey:"count"`
	Size     *datasize.ByteSize `json:"size,omitempty" configKey:"size"`
	Interval *time.Duration     `json:"interval,omitempty" configKey:"interval"`
}

// With copies values from the ImportTriggerPatch, if any.
func (c ImportTrigger) With(v ImportTriggerPatch) ImportTrigger {
	if v.Count != nil {
		c.Count = *v.Count
	}
	if v.Size != nil {
		c.Size = *v.Size
	}
	if v.Interval != nil {
		c.Interval = *v.Interval
	}
	return c
}
