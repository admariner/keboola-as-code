package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIf_Evaluate(t *testing.T) {
	t.Parallel()

	// Empty
	result, err := If("").Evaluate(map[string]interface{}{})
	assert.True(t, result)
	assert.NoError(t, err)

	// Simple
	result, err = If("true").Evaluate(map[string]interface{}{})
	assert.True(t, result)
	assert.NoError(t, err)

	// Parameter - true
	result, err = If("[my-param]").Evaluate(map[string]interface{}{"my-param": true})
	assert.True(t, result)
	assert.NoError(t, err)

	// Parameter - false
	result, err = If("[my-param]").Evaluate(map[string]interface{}{"my-param": false})
	assert.False(t, result)
	assert.NoError(t, err)

	// Parameter - not found
	result, err = If("[my-param]").Evaluate(map[string]interface{}{})
	assert.False(t, result)
	assert.Error(t, err)
	assert.Equal(t, "cannot evaluate condition:\n- expression: [my-param]\n- error: No parameter 'my-param' found.", err.Error())

	// Invalid expression
	result, err = If(">>>>>").Evaluate(map[string]interface{}{})
	assert.False(t, result)
	assert.Error(t, err)
	assert.Equal(t, "cannot compile condition:\n- expression: >>>>>\n- error: Invalid token: '>>>>>'", err.Error())
}
