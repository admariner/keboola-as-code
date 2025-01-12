package condition

import (
	"testing"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/statistics"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/store/model"
)

func TestEvaluate_DefaultImportConditions(t *testing.T) {
	t.Parallel()

	now, _ := time.Parse(time.RFC3339, "2010-01-01T01:01:01Z")
	before01Min := now.Add(-1 * time.Minute)
	before20Min := now.Add(-20 * time.Minute)

	// Defaults
	cn := model.DefaultImportConditions()

	// Defaults not met
	res, desc := evaluate(cn, now, before01Min, statistics.Value{
		RecordsCount: 50,
		RecordsSize:  1 * datasize.KB,
	})
	assert.False(t, res)
	assert.Equal(t, "no condition met", desc)

	// Default count met
	res, desc = evaluate(cn, now, before01Min, statistics.Value{
		RecordsCount: 20000,
		RecordsSize:  1 * datasize.MB,
	})
	assert.True(t, res)
	assert.Equal(t, "count threshold met, received: 20000 rows, threshold: 10000 rows", desc)

	// Default size met
	res, desc = evaluate(cn, now, before01Min, statistics.Value{
		RecordsCount: 100,
		RecordsSize:  10 * datasize.MB,
	})
	assert.True(t, res)
	assert.Equal(t, "size threshold met, received: 10.0 MB, threshold: 5.0 MB", desc)

	// Default time met
	res, desc = evaluate(cn, now, before20Min, statistics.Value{
		RecordsCount: 100,
		RecordsSize:  1 * datasize.KB,
	})
	assert.True(t, res)
	assert.Equal(t, "time threshold met, opened at: 2010-01-01T00:41:01.000Z, passed: 20m0s threshold: 5m0s", desc)
}

func TestEvaluate(t *testing.T) {
	t.Parallel()

	now, _ := time.Parse(time.RFC3339, "2010-01-01T01:01:01Z")
	before01Min := now.Add(-1 * time.Minute)
	before20Min := now.Add(-20 * time.Minute)

	// Defaults
	cn := Conditions{
		Count: 100,
		Size:  5 * datasize.MB,
		Time:  10 * time.Minute,
	}
	// Not met
	res, desc := evaluate(cn, now, before01Min, statistics.Value{
		RecordsCount: 50,
		RecordsSize:  1 * datasize.MB,
	})
	assert.False(t, res)
	assert.Equal(t, "no condition met", desc)

	// Count met
	res, desc = evaluate(cn, now, before01Min, statistics.Value{
		RecordsCount: 200,
		RecordsSize:  1 * datasize.MB,
	})
	assert.True(t, res)
	assert.Equal(t, "count threshold met, received: 200 rows, threshold: 100 rows", desc)

	// Size met
	res, desc = evaluate(cn, now, before01Min, statistics.Value{
		RecordsCount: 50,
		RecordsSize:  10 * datasize.MB,
	})
	assert.True(t, res)
	assert.Equal(t, "size threshold met, received: 10.0 MB, threshold: 5.0 MB", desc)

	// Time met
	res, desc = evaluate(cn, now, before20Min, statistics.Value{
		RecordsCount: 50,
		RecordsSize:  1 * datasize.MB,
	})
	assert.True(t, res)
	assert.Equal(t, "time threshold met, opened at: 2010-01-01T00:41:01.000Z, passed: 20m0s threshold: 10m0s", desc)
}
