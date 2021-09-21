package model

import (
	"fmt"
	"sort"
	"sync"
)

// Component https://keboola.docs.apiary.io/#reference/components-and-configurations/get-development-branch-components/get-development-branch-components
type Component struct {
	ComponentKey
	Type      string                 `json:"type" validate:"required"`
	Name      string                 `json:"name" validate:"required"`
	Schema    map[string]interface{} `json:"configurationSchema,omitempty"`
	SchemaRow map[string]interface{} `json:"configurationRowSchema,omitempty"`
}

type ComponentWithConfigs struct {
	BranchId int `json:"branchId" validate:"required"`
	*Component
	Configs []*ConfigWithRows `json:"configurations" validate:"required"`
}

func (c *Component) IsTransformation() bool {
	return c.Type == TransformationType
}

// remoteComponentsProvider - interface for Storage API.
type remoteComponentsProvider interface {
	GetComponent(componentId string) (*Component, error)
}

type ComponentsMap struct {
	mutex          *sync.Mutex
	remoteProvider remoteComponentsProvider
	components     map[string]*Component
}

func NewComponentsMap(remoteProvider remoteComponentsProvider) *ComponentsMap {
	return &ComponentsMap{
		mutex:          &sync.Mutex{},
		remoteProvider: remoteProvider,
		components:     make(map[string]*Component),
	}
}

func (c *ComponentsMap) AllLoaded() []*Component {
	var components []*Component
	for _, c := range c.components {
		components = append(components, c)
	}
	sort.SliceStable(components, func(i, j int) bool {
		return components[i].Id < components[j].Id
	})
	return components
}

func (c *ComponentsMap) Get(key ComponentKey) (*Component, error) {
	// Load component from cache if present
	if component, found := c.doGet(key); found {
		return component, nil
	}

	// Or by API
	if component, err := c.remoteProvider.GetComponent(key.Id); err == nil {
		return component, nil
	} else {
		return nil, err
	}
}

func (c *ComponentsMap) Set(component *Component) {
	if component == nil {
		panic(fmt.Errorf("component is not set"))
	}
	c.doSet(component)
}

func (c *ComponentsMap) doGet(key ComponentKey) (*Component, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	component, found := c.components[key.String()]
	return component, found
}

func (c *ComponentsMap) doSet(component *Component) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.components[component.String()] = component
}