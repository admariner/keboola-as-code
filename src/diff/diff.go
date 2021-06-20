package diff

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"keboola-as-code/src/model"
	"keboola-as-code/src/utils"
	"reflect"
)

type typeName string

type Differ struct {
	state     *model.State
	results   []*Result
	typeCache map[typeName][]reflect.StructField
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
	State   ResultState
	Changes map[string]string
}

type Results struct {
	Results []*Result
}

func NewDiffer(state *model.State) *Differ {
	return &Differ{
		state:     state,
		typeCache: make(map[typeName][]reflect.StructField),
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
	// Validate
	result := &Result{ObjectState: state}
	remoteState := state.RemoteState()
	localState := state.LocalState()
	if remoteState == nil && localState == nil {
		panic(fmt.Errorf("both local and remote state are not set"))
	}
	if remoteState == nil {
		result.State = ResultOnlyInLocal
		return result, nil
	}
	if localState == nil {
		result.State = ResultOnlyInLocal
		return result, nil
	}

	// Types must be same
	remoteType := reflect.TypeOf(remoteState).Elem()
	localType := reflect.TypeOf(localState).Elem()
	if remoteType.String() != localType.String() {
		panic(fmt.Errorf("local(%s) and remote(%s) states must have same data type", remoteType, localType))
	}

	// Get available fields for diff
	diffFields := d.getDiffFields(remoteType)
	if len(diffFields) == 0 {
		return nil, fmt.Errorf(`no field with tag "diff:true" in struct "%s"`, remoteType.String())
	}

	// Diff all diffFields
	result.Changes = make(map[string]string)
	remoteValues := reflect.ValueOf(remoteState).Elem()
	localValues := reflect.ValueOf(remoteState).Elem()
	for _, field := range diffFields {
		difference := cmp.Diff(
			remoteValues.FieldByName(field.Name).Interface(),
			localValues.FieldByName(field.Name).Interface(),
		)
		if len(difference) > 0 {
			result.Changes[field.Name] = difference
		}
	}
	if len(result.Changes) > 0 {
		result.State = ResultEqual
	} else {
		result.State = ResultNotEqual
	}

	return result, nil
}

func (d *Differ) getDiffFields(t reflect.Type) []reflect.StructField {
	if v, ok := d.typeCache[typeName(t.Name())]; ok {
		return v
	} else {
		var diffFields []reflect.StructField
		numFields := t.NumField()
		for i := 0; i < numFields; i++ {
			fieldType := t.Field(i)
			tag := fieldType.Tag.Get("diff")
			if tag == "true" {
				diffFields = append(diffFields, fieldType)
			}
		}
		d.typeCache[typeName(t.Name())] = diffFields
		return diffFields
	}
}
