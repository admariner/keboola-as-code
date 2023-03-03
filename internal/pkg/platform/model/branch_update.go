// Code generated by ent, DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/keboola/keboola-as-code/internal/pkg/platform/model/branch"
	"github.com/keboola/keboola-as-code/internal/pkg/platform/model/configuration"
	"github.com/keboola/keboola-as-code/internal/pkg/platform/model/key"
	"github.com/keboola/keboola-as-code/internal/pkg/platform/model/predicate"
)

// BranchUpdate is the builder for updating Branch entities.
type BranchUpdate struct {
	config
	hooks    []Hook
	mutation *BranchMutation
}

// Where appends a list predicates to the BranchUpdate builder.
func (bu *BranchUpdate) Where(ps ...predicate.Branch) *BranchUpdate {
	bu.mutation.Where(ps...)
	return bu
}

// SetName sets the "name" field.
func (bu *BranchUpdate) SetName(s string) *BranchUpdate {
	bu.mutation.SetName(s)
	return bu
}

// SetDescription sets the "description" field.
func (bu *BranchUpdate) SetDescription(s string) *BranchUpdate {
	bu.mutation.SetDescription(s)
	return bu
}

// AddConfigurationIDs adds the "configurations" edge to the Configuration entity by IDs.
func (bu *BranchUpdate) AddConfigurationIDs(ids ...key.ConfigurationKey) *BranchUpdate {
	bu.mutation.AddConfigurationIDs(ids...)
	return bu
}

// AddConfigurations adds the "configurations" edges to the Configuration entity.
func (bu *BranchUpdate) AddConfigurations(c ...*Configuration) *BranchUpdate {
	ids := make([]key.ConfigurationKey, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return bu.AddConfigurationIDs(ids...)
}

// Mutation returns the BranchMutation object of the builder.
func (bu *BranchUpdate) Mutation() *BranchMutation {
	return bu.mutation
}

// ClearConfigurations clears all "configurations" edges to the Configuration entity.
func (bu *BranchUpdate) ClearConfigurations() *BranchUpdate {
	bu.mutation.ClearConfigurations()
	return bu
}

// RemoveConfigurationIDs removes the "configurations" edge to Configuration entities by IDs.
func (bu *BranchUpdate) RemoveConfigurationIDs(ids ...key.ConfigurationKey) *BranchUpdate {
	bu.mutation.RemoveConfigurationIDs(ids...)
	return bu
}

// RemoveConfigurations removes "configurations" edges to Configuration entities.
func (bu *BranchUpdate) RemoveConfigurations(c ...*Configuration) *BranchUpdate {
	ids := make([]key.ConfigurationKey, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return bu.RemoveConfigurationIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (bu *BranchUpdate) Save(ctx context.Context) (int, error) {
	return withHooks[int, BranchMutation](ctx, bu.sqlSave, bu.mutation, bu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (bu *BranchUpdate) SaveX(ctx context.Context) int {
	affected, err := bu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (bu *BranchUpdate) Exec(ctx context.Context) error {
	_, err := bu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bu *BranchUpdate) ExecX(ctx context.Context) {
	if err := bu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (bu *BranchUpdate) check() error {
	if v, ok := bu.mutation.Name(); ok {
		if err := branch.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Branch.name": %w`, err)}
		}
	}
	return nil
}

func (bu *BranchUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := bu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(branch.Table, branch.Columns, sqlgraph.NewFieldSpec(branch.FieldID, field.TypeString))
	if ps := bu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := bu.mutation.Name(); ok {
		_spec.SetField(branch.FieldName, field.TypeString, value)
	}
	if value, ok := bu.mutation.Description(); ok {
		_spec.SetField(branch.FieldDescription, field.TypeString, value)
	}
	if bu.mutation.ConfigurationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   branch.ConfigurationsTable,
			Columns: []string{branch.ConfigurationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: configuration.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.RemovedConfigurationsIDs(); len(nodes) > 0 && !bu.mutation.ConfigurationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   branch.ConfigurationsTable,
			Columns: []string{branch.ConfigurationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: configuration.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.ConfigurationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   branch.ConfigurationsTable,
			Columns: []string{branch.ConfigurationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: configuration.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, bu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{branch.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	bu.mutation.done = true
	return n, nil
}

// BranchUpdateOne is the builder for updating a single Branch entity.
type BranchUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *BranchMutation
}

// SetName sets the "name" field.
func (buo *BranchUpdateOne) SetName(s string) *BranchUpdateOne {
	buo.mutation.SetName(s)
	return buo
}

// SetDescription sets the "description" field.
func (buo *BranchUpdateOne) SetDescription(s string) *BranchUpdateOne {
	buo.mutation.SetDescription(s)
	return buo
}

// AddConfigurationIDs adds the "configurations" edge to the Configuration entity by IDs.
func (buo *BranchUpdateOne) AddConfigurationIDs(ids ...key.ConfigurationKey) *BranchUpdateOne {
	buo.mutation.AddConfigurationIDs(ids...)
	return buo
}

// AddConfigurations adds the "configurations" edges to the Configuration entity.
func (buo *BranchUpdateOne) AddConfigurations(c ...*Configuration) *BranchUpdateOne {
	ids := make([]key.ConfigurationKey, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return buo.AddConfigurationIDs(ids...)
}

// Mutation returns the BranchMutation object of the builder.
func (buo *BranchUpdateOne) Mutation() *BranchMutation {
	return buo.mutation
}

// ClearConfigurations clears all "configurations" edges to the Configuration entity.
func (buo *BranchUpdateOne) ClearConfigurations() *BranchUpdateOne {
	buo.mutation.ClearConfigurations()
	return buo
}

// RemoveConfigurationIDs removes the "configurations" edge to Configuration entities by IDs.
func (buo *BranchUpdateOne) RemoveConfigurationIDs(ids ...key.ConfigurationKey) *BranchUpdateOne {
	buo.mutation.RemoveConfigurationIDs(ids...)
	return buo
}

// RemoveConfigurations removes "configurations" edges to Configuration entities.
func (buo *BranchUpdateOne) RemoveConfigurations(c ...*Configuration) *BranchUpdateOne {
	ids := make([]key.ConfigurationKey, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return buo.RemoveConfigurationIDs(ids...)
}

// Where appends a list predicates to the BranchUpdate builder.
func (buo *BranchUpdateOne) Where(ps ...predicate.Branch) *BranchUpdateOne {
	buo.mutation.Where(ps...)
	return buo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (buo *BranchUpdateOne) Select(field string, fields ...string) *BranchUpdateOne {
	buo.fields = append([]string{field}, fields...)
	return buo
}

// Save executes the query and returns the updated Branch entity.
func (buo *BranchUpdateOne) Save(ctx context.Context) (*Branch, error) {
	return withHooks[*Branch, BranchMutation](ctx, buo.sqlSave, buo.mutation, buo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (buo *BranchUpdateOne) SaveX(ctx context.Context) *Branch {
	node, err := buo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (buo *BranchUpdateOne) Exec(ctx context.Context) error {
	_, err := buo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (buo *BranchUpdateOne) ExecX(ctx context.Context) {
	if err := buo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (buo *BranchUpdateOne) check() error {
	if v, ok := buo.mutation.Name(); ok {
		if err := branch.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Branch.name": %w`, err)}
		}
	}
	return nil
}

func (buo *BranchUpdateOne) sqlSave(ctx context.Context) (_node *Branch, err error) {
	if err := buo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(branch.Table, branch.Columns, sqlgraph.NewFieldSpec(branch.FieldID, field.TypeString))
	id, ok := buo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "Branch.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := buo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, branch.FieldID)
		for _, f := range fields {
			if !branch.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != branch.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := buo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := buo.mutation.Name(); ok {
		_spec.SetField(branch.FieldName, field.TypeString, value)
	}
	if value, ok := buo.mutation.Description(); ok {
		_spec.SetField(branch.FieldDescription, field.TypeString, value)
	}
	if buo.mutation.ConfigurationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   branch.ConfigurationsTable,
			Columns: []string{branch.ConfigurationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: configuration.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.RemovedConfigurationsIDs(); len(nodes) > 0 && !buo.mutation.ConfigurationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   branch.ConfigurationsTable,
			Columns: []string{branch.ConfigurationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: configuration.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.ConfigurationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   branch.ConfigurationsTable,
			Columns: []string{branch.ConfigurationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: configuration.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Branch{config: buo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, buo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{branch.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	buo.mutation.done = true
	return _node, nil
}
