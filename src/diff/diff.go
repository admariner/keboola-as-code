package diff

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"keboola-as-code/src/model"
	"keboola-as-code/src/utils"
	"reflect"
	"strings"
)

type typeName string

type structField struct {
	name    string
	reflect reflect.StructField
}

type Differ struct {
	state     *model.State
	results   []*Result
	typeCache map[typeName][]structField
	error     *utils.Error
}

type ResultState int

const (
	ResultNotSet ResultState = iota
	ResultNotEqual
	ResultEqual
	ResultOnlyInRemote
	ResultOnlyInLocal
)

type Result struct {
	model.ObjectState
	State         ResultState
	ChangedFields []string
	Differences   map[string]string
}

type Results struct {
	Results []*Result
}

func NewDiffer(state *model.State) *Differ {
	return &Differ{
		state:     state,
		typeCache: make(map[typeName][]structField),
	}
}

func (d *Differ) Diff() (*Results, error) {
	// Diff all states
	d.results = []*Result{}
	d.error = &utils.Error{}
	for _, objectState := range d.state.All() {
		result, err := d.doDiff(objectState)
		if err != nil {
			d.error.Add(err)
		} else {
			d.results = append(d.results, result)
		}
	}
	// Check errors
	var err error
	if d.error.Len() > 0 {
		err = fmt.Errorf("%s", d.error)
	}

	return &Results{d.results}, err
}

func (d *Differ) doDiff(state model.ObjectState) (*Result, error) {
	result := &Result{ObjectState: state}
	remoteState := state.RemoteState()
	localState := state.LocalState()
	remoteType := reflect.TypeOf(remoteState).Elem()
	localType := reflect.TypeOf(localState).Elem()
	remoteValues := reflect.ValueOf(remoteState)
	localValues := reflect.ValueOf(localState)

	// Types must be same
	if remoteType.String() != localType.String() {
		panic(fmt.Errorf("local(%s) and remote(%s) states must have same data type", remoteType, localType))
	}

	// Get available fields for diff
	diffFields := d.getDiffFields(remoteType)
	if len(diffFields) == 0 {
		return nil, fmt.Errorf(`no field with tag "diff:true" in struct "%s"`, remoteType.String())
	}

	// Check values
	result.ChangedFields = make([]string, 0)
	result.Differences = make(map[string]string)
	if remoteValues.IsNil() && localValues.IsNil() {
		panic(fmt.Errorf("both local and remote state are not set"))
	}
	if remoteValues.IsNil() {
		result.State = ResultOnlyInLocal
		return result, nil
	}
	if localValues.IsNil() {
		result.State = ResultOnlyInRemote
		return result, nil
	}

	// Get pointer value
	if remoteValues.Type().Kind() == reflect.Ptr {
		remoteValues = remoteValues.Elem()
	}
	if localValues.Type().Kind() == reflect.Ptr {
		localValues = localValues.Elem()
	}

	// Diff
	for _, field := range diffFields {
		difference := cmp.Diff(
			remoteValues.FieldByName(field.reflect.Name).Interface(),
			localValues.FieldByName(field.reflect.Name).Interface(),
		)
		if len(difference) > 0 {
			result.ChangedFields = append(result.ChangedFields, field.name)
			result.Differences[field.name] = difference
		}
	}

	if len(result.ChangedFields) > 0 {
		result.State = ResultNotEqual
	} else {
		result.State = ResultEqual
	}

	return result, nil
}

func (d *Differ) getDiffFields(t reflect.Type) []structField {
	if v, ok := d.typeCache[typeName(t.Name())]; ok {
		return v
	} else {
		diffFields := make([]structField, 0)
		numFields := t.NumField()
		for i := 0; i < numFields; i++ {
			fieldType := t.Field(i)

			// Use JSON name if present
			name := fieldType.Name
			jsonName := strings.Split(fieldType.Tag.Get("json"), ",")[0]
			if jsonName != "" {
				name = jsonName
			}

			// Field must be marked with tag `diff:"true"`
			tag := fieldType.Tag.Get("diff")
			if tag == "true" {
				diffFields = append(diffFields, structField{name, fieldType})
			}
		}
		name := typeName(t.Name())
		d.typeCache[name] = diffFields
		return diffFields
	}
}
