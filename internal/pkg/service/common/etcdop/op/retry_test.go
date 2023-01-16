package op

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/assert"
)

func TestAtomicUpdateBackoff(t *testing.T) {
	t.Parallel()

	b := newBackoff()
	b.RandomizationFactor = 0

	clk := clock.NewMock()
	b.Clock = clk
	b.Reset()

	// Get all delays without sleep
	var delays []time.Duration
	for {
		delay := b.NextBackOff()
		if delay == backoff.Stop {
			break
		}
		delays = append(delays, delay)
		clk.Add(delay)
	}

	// Assert
	assert.Equal(t, []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		40 * time.Millisecond,
		80 * time.Millisecond,
		160 * time.Millisecond,
		320 * time.Millisecond,
		640 * time.Millisecond,
		1280 * time.Millisecond,
		2000 * time.Millisecond,
		2000 * time.Millisecond,
		2000 * time.Millisecond,
		2000 * time.Millisecond,
		2000 * time.Millisecond,
		2000 * time.Millisecond,
	}, delays)
}
