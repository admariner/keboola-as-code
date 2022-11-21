package description_test

import (
	"testing"

	"github.com/keboola/keboola-as-code/internal/pkg/mapper/description"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/state"
)

func createStateWithMapper(t *testing.T) (*state.State, dependencies.Mocked) {
	t.Helper()
	d := dependencies.NewMockedDeps()
	mockedState := d.MockedState()
	mockedState.Mapper().AddMapper(description.NewMapper())
	return mockedState, d
}
