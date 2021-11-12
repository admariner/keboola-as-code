package helper

import (
	"fmt"

	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
)

// SharedCodeHelper gets some values from shared codes.
type SharedCodeHelper struct {
	state  *model.State
	naming *model.Naming
}

func New(state *model.State, naming *model.Naming) *SharedCodeHelper {
	return &SharedCodeHelper{state: state, naming: naming}
}

func (h *SharedCodeHelper) IsSharedCodeKey(key model.Key) (bool, error) {
	// Is config?
	configKey, ok := key.(model.ConfigKey)
	if !ok {
		return false, nil
	}

	// Is shared code?
	component, err := h.state.Components().Get(configKey.ComponentKey())
	if err != nil || !component.IsSharedCode() {
		return false, err
	}
	return true, nil
}

func (h *SharedCodeHelper) IsSharedCodeRowKey(key model.Key) (bool, error) {
	// Is config row?
	configRowKey, ok := key.(model.ConfigRowKey)
	if !ok {
		return false, nil
	}

	// Is shared code?
	component, err := h.state.Components().Get(configRowKey.ComponentKey())
	if err != nil || !component.IsSharedCode() {
		return false, err
	}
	return true, nil
}

func (h *SharedCodeHelper) IsTransformation(key model.Key) (bool, error) {
	// Is config?
	configKey, ok := key.(model.ConfigKey)
	if !ok {
		return false, nil
	}

	// Is shared code?
	component, err := h.state.Components().Get(configKey.ComponentKey())
	if err != nil || !component.IsTransformation() {
		return false, err
	}
	return true, nil
}

// GetTargetComponentId returns the component for which the shared code is intended.
func (h *SharedCodeHelper) GetTargetComponentId(sharedCodeConfig *model.Config) (string, error) {
	componentIdRaw, found := sharedCodeConfig.Content.Get(model.SharedCodeComponentIdContentKey)
	if !found {
		return "", fmt.Errorf(`missing "%s" in %s`, model.SharedCodeComponentIdContentKey, sharedCodeConfig.Desc())
	}

	componentId, ok := componentIdRaw.(string)
	if !ok {
		return "", fmt.Errorf(`key "%s" must be string, found %T, in %s`, model.SharedCodeComponentIdContentKey, componentIdRaw, sharedCodeConfig.Desc())
	}

	return componentId, nil
}

func (h *SharedCodeHelper) GetSharedCodePath(object model.Object) (*model.Config, string, error) {
	// Shared code is used by transformation
	ok, err := h.IsTransformation(object.Key())
	if err != nil || !ok {
		return nil, "", err
	}
	transformation := object.(*model.Config)

	// Path is stored in configuration
	sharedCodePathRaw, found := transformation.Content.Get(model.SharedCodePathContentKey)
	if !found {
		return nil, "", nil
	}

	// Path must be string
	sharedCodePath, ok := sharedCodePathRaw.(string)
	if !ok {
		return nil, "", fmt.Errorf(`key "%s" must be string, found %T, in %s`, model.SharedCodePathContentKey, sharedCodePathRaw, object.Desc())
	}
	return transformation, sharedCodePath, nil
}

func (h *SharedCodeHelper) GetSharedCodeKey(object model.Object) (*model.Config, model.Key, error) {
	// Shared code is used by transformation
	ok, err := h.IsTransformation(object.Key())
	if err != nil || !ok {
		return nil, nil, err
	}
	transformation := object.(*model.Config)

	// ID is stored in configuration
	sharedCodeIdRaw, found := transformation.Content.Get(model.SharedCodeIdContentKey)
	if !found {
		return nil, nil, nil
	}

	// ID must be string
	sharedCodeConfigId, ok := sharedCodeIdRaw.(string)
	if !ok {
		return nil, nil, fmt.Errorf(`key "%s" must be string, found %T, in %s`, model.SharedCodeIdContentKey, sharedCodeIdRaw, object.Desc())
	}

	// Id -> key
	sharedCodeKey := model.ConfigKey{
		BranchId:    transformation.BranchId, // same branch
		ComponentId: model.SharedCodeComponentId,
		Id:          sharedCodeConfigId,
	}
	return transformation, sharedCodeKey, nil
}

func (h *SharedCodeHelper) GetSharedCodeByPath(branchKey model.BranchKey, path string) *model.ConfigState {
	// Get branch
	branch, found := h.state.Get(branchKey)
	if !found {
		return nil
	}

	// Get key by path
	path = filesystem.Join(branch.Path(), path)
	keyRaw, found := h.naming.FindByPath(path)
	if !found {
		return nil
	}

	// Is config?
	key, ok := keyRaw.(model.ConfigKey)
	if !ok {
		return nil
	}

	// Is from right parent?
	if branchKey != key.BranchKey() {
		return nil
	}

	// Shared code?
	if key.ComponentId != model.SharedCodeComponentId {
		return nil
	}

	// Ok
	return h.state.MustGet(key).(*model.ConfigState)
}

func (h *SharedCodeHelper) GetSharedCodeRowByPath(sharedCode *model.ConfigState, path string) *model.ConfigRowState {
	// Get key by path
	path = filesystem.Join(sharedCode.Path(), path)
	keyRaw, found := h.naming.FindByPath(path)
	if !found {
		return nil
	}

	// Is config row?
	key, ok := keyRaw.(model.ConfigRowKey)
	if !ok {
		return nil
	}

	// Is from parent?
	if sharedCode.Key() != key.ConfigKey() {
		return nil
	}

	// Ok
	return h.state.MustGet(key).(*model.ConfigRowState)
}

func (h *SharedCodeHelper) GetSharedCodeVariablesId(configRow *model.ConfigRow) (string, bool) {
	// Variables ID is stored in configuration
	variablesIdRaw, found := configRow.Content.Get(model.SharedCodeVariablesIdContentKey)
	if !found {
		return "", false
	}

	// Variables ID must be string
	variablesId, ok := variablesIdRaw.(string)
	if !ok {
		return "", false
	}

	return variablesId, true
}