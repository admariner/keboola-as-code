// Package use represents the process of replacing of values when applying a template.
package use

import (
	"context"
	"fmt"
	"sync"

	jsonnetLib "github.com/google/go-jsonnet"
	"github.com/keboola/go-client/pkg/storageapi"
	"github.com/keboola/go-utils/pkg/orderedmap"

	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/jsonnet"
	"github.com/keboola/keboola-as-code/internal/pkg/jsonnet/fsimporter"
	"github.com/keboola/keboola-as-code/internal/pkg/mapper/template/jsonnetfiles"
	"github.com/keboola/keboola-as-code/internal/pkg/mapper/template/metadata"
	"github.com/keboola/keboola-as-code/internal/pkg/mapper/template/replacevalues"
	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/template"
	"github.com/keboola/keboola-as-code/internal/pkg/template/jsonnet/function"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/strhelper"
)

// Context represents the process of the replacing values when applying a template.
//
// Process description:
//  1. There is some template.
//     - It contains objects IDs defined by functions, for example: ConfigId("my-config-id"), ConfigRowId("my-row-id")
//  2. When loading JsonNet files, functions are called.
//     - A placeholder is generated for each unique value.
//     - For example, each ConfigId("my-config-id") is replaced by "<<~~func:ticket:1~~>>".
//     - This is because we do not know in advance how many new IDs will need to be generated.
//     - Function call can contain an expression, for example ConfigId("my-config-" + tableName), and this prevents forward analysis.
//     - Functions are defined in Context.registerJsonNetFunctions().
//  3. When the entire template is loaded, the placeholders are replaced with new IDs.
//     - For example, each "<<~~func:ticket:1~~>>" is replaced by "3496482342".
//     - Replacements are defined by Context.Replacements().
//     - Values are replaced by "internal/pkg/mapper/template/replacevalues".
//  4. Then the objects are copied to the project,
//     - See "pkg/lib/operation/project/local/template/use/operation.go".
//     - A new path is generated for each new object, according to the project naming.
//
// Context.JsonNetContext() returns JsonNet functions.
// Context.Replacements() returns placeholders for new IDs.
type Context struct {
	_context
	templateRef       model.TemplateRef
	instanceId        string
	instanceIdShort   string
	jsonNetCtx        *jsonnet.Context
	replacements      *replacevalues.Values
	inputsValues      map[string]template.InputValue
	tickets           *storageapi.TicketProvider
	components        *model.ComponentsMap
	placeholdersCount int
	ticketsResolved   bool

	lock          *sync.Mutex
	placeholders  PlaceholdersMap
	objectIds     metadata.ObjectIdsMap
	inputsUsage   *metadata.InputsUsage
	inputsDefsMap map[string]*template.Input
}

type _context context.Context

// PlaceholdersMap -  original template value -> placeholder.
type PlaceholdersMap map[interface{}]Placeholder

type Placeholder struct {
	asString string      // placeholder as string for use in Json file, eg. string("<<~~placeholder:1~~>>)
	asValue  interface{} // eg. ConfigId, RowId, eg. ConfigId("<<~~placeholder:1~~>>)
}

type PlaceholderResolver func(p Placeholder, cb ResolveCallback)

type ResolveCallback func(newId interface{})

type inputUsageNotifier struct {
	*Context
	ctx context.Context
}

const (
	placeholderStart      = "<<~~"
	placeholderEnd        = "~~>>"
	instanceIdShortLength = 8
)

func NewContext(ctx context.Context, templateRef model.TemplateRef, objectsRoot filesystem.Fs, instanceId string, targetBranch model.BranchKey, inputsValues template.InputsValues, inputsDefsMap map[string]*template.Input, tickets *storageapi.TicketProvider, components *model.ComponentsMap) *Context {
	ctx = template.NewContext(ctx)
	c := &Context{
		_context:        ctx,
		templateRef:     templateRef,
		instanceId:      instanceId,
		instanceIdShort: strhelper.FirstN(instanceId, instanceIdShortLength),
		jsonNetCtx:      jsonnet.NewContext().WithCtx(ctx).WithImporter(fsimporter.New(objectsRoot)),
		replacements:    replacevalues.NewValues(),
		inputsValues:    make(map[string]template.InputValue),
		tickets:         tickets,
		components:      components,
		lock:            &sync.Mutex{},
		placeholders:    make(PlaceholdersMap),
		objectIds:       make(metadata.ObjectIdsMap),
		inputsUsage:     metadata.NewInputsUsage(),
		inputsDefsMap:   inputsDefsMap,
	}

	// Convert inputsValues to map
	for _, input := range inputsValues {
		c.inputsValues[input.Id] = input
	}

	// Replace BranchId, in template all objects have BranchId = 0
	c.replacements.AddKey(model.BranchKey{Id: 0}, targetBranch)

	// Register JsonNet functions
	c.registerJsonNetFunctions()

	// Let's see where the inputs were used
	c.registerInputsUsageNotifier()

	return c
}

func (c *Context) TemplateRef() model.TemplateRef {
	return c.templateRef
}

func (c *Context) InstanceId() string {
	return c.instanceId
}

func (c *Context) JsonNetContext() *jsonnet.Context {
	return c.jsonNetCtx
}

func (c *Context) Replacements() (*replacevalues.Values, error) {
	// Generate new IDs
	if !c.ticketsResolved {
		if err := c.tickets.Resolve(); err != nil {
			return nil, err
		}
		c.ticketsResolved = true
	}
	return c.replacements, nil
}

func (c *Context) RemoteObjectsFilter() model.ObjectsFilter {
	return model.NoFilter()
}

func (c *Context) LocalObjectsFilter() model.ObjectsFilter {
	return model.NoFilter()
}

func (c *Context) ObjectIds() metadata.ObjectIdsMap {
	return c.objectIds
}

func (c *Context) InputsUsage() *metadata.InputsUsage {
	return c.inputsUsage
}

// RegisterPlaceholder for an object oldId, it can be resolved later/async.
func (c *Context) RegisterPlaceholder(oldId interface{}, fn PlaceholderResolver) Placeholder {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, found := c.placeholders[oldId]; !found {
		// Generate placeholder, it will be later replaced by a new ID
		c.placeholdersCount++
		p := Placeholder{asString: fmt.Sprintf("%splaceholder:%d%s", placeholderStart, c.placeholdersCount, placeholderEnd)}

		// Convert string to an ID value
		switch oldId.(type) {
		case storageapi.ConfigID:
			p.asValue = storageapi.ConfigID(p.asString)
		case storageapi.RowID:
			p.asValue = storageapi.RowID(p.asString)
		default:
			panic(fmt.Errorf("unexpected ID type"))
		}

		// Store oldId -> placeholder
		c.placeholders[oldId] = p

		// Resolve newId async by provider function
		fn(p, func(newId interface{}) {
			c.replacements.AddId(p.asValue, newId)
			c.objectIds[newId] = oldId
		})
	}
	return c.placeholders[oldId]
}

func (c *Context) registerJsonNetFunctions() {
	c.jsonNetCtx.NativeFunctionWithAlias(function.ConfigId(c.mapId))
	c.jsonNetCtx.NativeFunctionWithAlias(function.ConfigRowId(c.mapId))
	c.jsonNetCtx.NativeFunctionWithAlias(function.Input(c.inputValue))
	c.jsonNetCtx.NativeFunctionWithAlias(function.InputIsAvailable(c.inputValue))
	c.jsonNetCtx.NativeFunctionWithAlias(function.InstanceId(c.instanceId))
	c.jsonNetCtx.NativeFunctionWithAlias(function.InstanceIdShort(c.instanceIdShort))
	c.jsonNetCtx.NativeFunctionWithAlias(function.ComponentIsAvailable(c.components))
	c.jsonNetCtx.NativeFunctionWithAlias(function.SnowflakeWriterComponentId(c.components))
}

// mapId maps ConfigId/ConfigRowId in JsonNet files to a <<~~ticket:123~~>> placeholder.
// When all JsonNet files are processed, new IDs are generated in parallel.
func (c *Context) mapId(oldId interface{}) string {
	p := c.RegisterPlaceholder(oldId, func(p Placeholder, cb ResolveCallback) {
		// Placeholder -> new ID
		var newId interface{}
		c.tickets.Request(func(ticket *storageapi.Ticket) {
			switch p.asValue.(type) {
			case storageapi.ConfigID:
				newId = storageapi.ConfigID(ticket.ID)
			case storageapi.RowID:
				newId = storageapi.RowID(ticket.ID)
			default:
				panic(fmt.Errorf("unexpected ID type"))
			}
			cb(newId)
		})
	})
	return p.asString
}

func (c *Context) inputValue(inputId string) (template.InputValue, bool) {
	v, ok := c.inputsValues[inputId]
	return v, ok
}

func (c *Context) registerInputsUsageNotifier() {
	c.jsonNetCtx.NotifierFactory(func(ctx context.Context) jsonnetLib.Notifier {
		return &inputUsageNotifier{Context: c, ctx: ctx}
	})
}

func (n *inputUsageNotifier) OnGeneratedValue(fnName string, args []interface{}, partial bool, partialValue, _ interface{}, steps []interface{}) {
	// Only for Input function
	if fnName != "Input" {
		return
	}

	// One argument expected
	if len(args) != 1 {
		return
	}

	// Argument is input name
	inputName, ok := args[0].(string)
	if !ok {
		return
	}

	// Check if input exists and has been filled in by user
	if input, found := n.inputsValues[inputName]; !found || input.Skipped {
		return
	}

	// Convert steps to orderedmap format
	var mappedSteps []orderedmap.Step
	for _, step := range steps {
		switch v := step.(type) {
		case jsonnetLib.ObjectFieldStep:
			mappedSteps = append(mappedSteps, orderedmap.MapStep(v.Field))
		case jsonnetLib.ArrayIndexStep:
			mappedSteps = append(mappedSteps, orderedmap.SliceStep(v.Index))
		default:
			panic(fmt.Errorf(`unexpected type "%T"`, v))
		}
	}

	// Get file definition
	fileDef, _ := n.ctx.Value(jsonnetfiles.FileDefCtxKey).(*filesystem.FileDef)
	if fileDef == nil {
		return
	}

	// Get key of the parent object
	objectKey, ok := fileDef.MetadataOrNil(filesystem.ObjectKeyMetadata).(model.Key)
	if !ok {
		return
	}

	// We are only interested in the inputs used in the configuration.
	if !fileDef.HasTag(model.FileKindObjectConfig) {
		return
	}

	// Replace tickets in object key
	objectKeyRaw, err := n.replacements.Replace(objectKey)
	if err != nil {
		panic(err)
	}

	// Store
	objectKey = objectKeyRaw.(model.Key)
	n.lock.Lock()
	defer n.lock.Unlock()
	if !partial {
		// Values has been generated by the Input function, store input usage
		n.inputsUsage.Values[objectKey] = append(n.inputsUsage.Values[objectKey], metadata.InputUsage{
			Name:    inputName,
			JsonKey: mappedSteps,
			Def:     n.inputsDefsMap[inputName],
		})
	} else if jsonObject, ok := partialValue.(map[string]any); ok && len(jsonObject) > 0 {
		// Get JSON keys
		var keys []string
		for jsonKey := range jsonObject {
			keys = append(keys, jsonKey)
		}

		// Part of the object has been generated by the Input function, store input usage
		n.inputsUsage.Values[objectKey] = append(n.inputsUsage.Values[objectKey], metadata.InputUsage{
			Name:       inputName,
			JsonKey:    mappedSteps,
			Def:        n.inputsDefsMap[inputName],
			ObjectKeys: keys,
		})
	}
}