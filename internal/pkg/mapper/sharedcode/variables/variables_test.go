package variables_test

import (
	"testing"

	"github.com/keboola/keboola-as-code/internal/pkg/mapper/sharedcode/variables"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/state"
)

func createStateWithMapper(t *testing.T) (*state.State, dependencies.Mocked) {
	t.Helper()
	d := dependencies.NewMockedDeps(t)
	mockedState := d.MockedState()
	mockedState.Mapper().AddMapper(variables.NewMapper(mockedState))
	return mockedState, d
}
