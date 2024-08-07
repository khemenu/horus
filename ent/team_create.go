// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"khepri.dev/horus/ent/membership"
	"khepri.dev/horus/ent/silo"
	"khepri.dev/horus/ent/team"
)

// TeamCreate is the builder for creating a Team entity.
type TeamCreate struct {
	config
	mutation *TeamMutation
	hooks    []Hook
}

// SetDateCreated sets the "date_created" field.
func (tc *TeamCreate) SetDateCreated(t time.Time) *TeamCreate {
	tc.mutation.SetDateCreated(t)
	return tc
}

// SetNillableDateCreated sets the "date_created" field if the given value is not nil.
func (tc *TeamCreate) SetNillableDateCreated(t *time.Time) *TeamCreate {
	if t != nil {
		tc.SetDateCreated(*t)
	}
	return tc
}

// SetAlias sets the "alias" field.
func (tc *TeamCreate) SetAlias(s string) *TeamCreate {
	tc.mutation.SetAlias(s)
	return tc
}

// SetNillableAlias sets the "alias" field if the given value is not nil.
func (tc *TeamCreate) SetNillableAlias(s *string) *TeamCreate {
	if s != nil {
		tc.SetAlias(*s)
	}
	return tc
}

// SetName sets the "name" field.
func (tc *TeamCreate) SetName(s string) *TeamCreate {
	tc.mutation.SetName(s)
	return tc
}

// SetNillableName sets the "name" field if the given value is not nil.
func (tc *TeamCreate) SetNillableName(s *string) *TeamCreate {
	if s != nil {
		tc.SetName(*s)
	}
	return tc
}

// SetDescription sets the "description" field.
func (tc *TeamCreate) SetDescription(s string) *TeamCreate {
	tc.mutation.SetDescription(s)
	return tc
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (tc *TeamCreate) SetNillableDescription(s *string) *TeamCreate {
	if s != nil {
		tc.SetDescription(*s)
	}
	return tc
}

// SetSiloID sets the "silo_id" field.
func (tc *TeamCreate) SetSiloID(u uuid.UUID) *TeamCreate {
	tc.mutation.SetSiloID(u)
	return tc
}

// SetID sets the "id" field.
func (tc *TeamCreate) SetID(u uuid.UUID) *TeamCreate {
	tc.mutation.SetID(u)
	return tc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (tc *TeamCreate) SetNillableID(u *uuid.UUID) *TeamCreate {
	if u != nil {
		tc.SetID(*u)
	}
	return tc
}

// SetSilo sets the "silo" edge to the Silo entity.
func (tc *TeamCreate) SetSilo(s *Silo) *TeamCreate {
	return tc.SetSiloID(s.ID)
}

// AddMemberIDs adds the "members" edge to the Membership entity by IDs.
func (tc *TeamCreate) AddMemberIDs(ids ...uuid.UUID) *TeamCreate {
	tc.mutation.AddMemberIDs(ids...)
	return tc
}

// AddMembers adds the "members" edges to the Membership entity.
func (tc *TeamCreate) AddMembers(m ...*Membership) *TeamCreate {
	ids := make([]uuid.UUID, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return tc.AddMemberIDs(ids...)
}

// Mutation returns the TeamMutation object of the builder.
func (tc *TeamCreate) Mutation() *TeamMutation {
	return tc.mutation
}

// Save creates the Team in the database.
func (tc *TeamCreate) Save(ctx context.Context) (*Team, error) {
	tc.defaults()
	return withHooks(ctx, tc.sqlSave, tc.mutation, tc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TeamCreate) SaveX(ctx context.Context) *Team {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tc *TeamCreate) Exec(ctx context.Context) error {
	_, err := tc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tc *TeamCreate) ExecX(ctx context.Context) {
	if err := tc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tc *TeamCreate) defaults() {
	if _, ok := tc.mutation.DateCreated(); !ok {
		v := team.DefaultDateCreated()
		tc.mutation.SetDateCreated(v)
	}
	if _, ok := tc.mutation.Alias(); !ok {
		v := team.DefaultAlias()
		tc.mutation.SetAlias(v)
	}
	if _, ok := tc.mutation.Name(); !ok {
		v := team.DefaultName
		tc.mutation.SetName(v)
	}
	if _, ok := tc.mutation.Description(); !ok {
		v := team.DefaultDescription
		tc.mutation.SetDescription(v)
	}
	if _, ok := tc.mutation.ID(); !ok {
		v := team.DefaultID()
		tc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tc *TeamCreate) check() error {
	if _, ok := tc.mutation.DateCreated(); !ok {
		return &ValidationError{Name: "date_created", err: errors.New(`ent: missing required field "Team.date_created"`)}
	}
	if _, ok := tc.mutation.Alias(); !ok {
		return &ValidationError{Name: "alias", err: errors.New(`ent: missing required field "Team.alias"`)}
	}
	if v, ok := tc.mutation.Alias(); ok {
		if err := team.AliasValidator(v); err != nil {
			return &ValidationError{Name: "alias", err: fmt.Errorf(`ent: validator failed for field "Team.alias": %w`, err)}
		}
	}
	if _, ok := tc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Team.name"`)}
	}
	if v, ok := tc.mutation.Name(); ok {
		if err := team.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Team.name": %w`, err)}
		}
	}
	if _, ok := tc.mutation.Description(); !ok {
		return &ValidationError{Name: "description", err: errors.New(`ent: missing required field "Team.description"`)}
	}
	if v, ok := tc.mutation.Description(); ok {
		if err := team.DescriptionValidator(v); err != nil {
			return &ValidationError{Name: "description", err: fmt.Errorf(`ent: validator failed for field "Team.description": %w`, err)}
		}
	}
	if _, ok := tc.mutation.SiloID(); !ok {
		return &ValidationError{Name: "silo_id", err: errors.New(`ent: missing required field "Team.silo_id"`)}
	}
	if len(tc.mutation.SiloIDs()) == 0 {
		return &ValidationError{Name: "silo", err: errors.New(`ent: missing required edge "Team.silo"`)}
	}
	return nil
}

func (tc *TeamCreate) sqlSave(ctx context.Context) (*Team, error) {
	if err := tc.check(); err != nil {
		return nil, err
	}
	_node, _spec := tc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	tc.mutation.id = &_node.ID
	tc.mutation.done = true
	return _node, nil
}

func (tc *TeamCreate) createSpec() (*Team, *sqlgraph.CreateSpec) {
	var (
		_node = &Team{config: tc.config}
		_spec = sqlgraph.NewCreateSpec(team.Table, sqlgraph.NewFieldSpec(team.FieldID, field.TypeUUID))
	)
	if id, ok := tc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := tc.mutation.DateCreated(); ok {
		_spec.SetField(team.FieldDateCreated, field.TypeTime, value)
		_node.DateCreated = value
	}
	if value, ok := tc.mutation.Alias(); ok {
		_spec.SetField(team.FieldAlias, field.TypeString, value)
		_node.Alias = value
	}
	if value, ok := tc.mutation.Name(); ok {
		_spec.SetField(team.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := tc.mutation.Description(); ok {
		_spec.SetField(team.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if nodes := tc.mutation.SiloIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   team.SiloTable,
			Columns: []string{team.SiloColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(silo.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.SiloID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.MembersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   team.MembersTable,
			Columns: []string{team.MembersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(membership.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// TeamCreateBulk is the builder for creating many Team entities in bulk.
type TeamCreateBulk struct {
	config
	err      error
	builders []*TeamCreate
}

// Save creates the Team entities in the database.
func (tcb *TeamCreateBulk) Save(ctx context.Context) ([]*Team, error) {
	if tcb.err != nil {
		return nil, tcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(tcb.builders))
	nodes := make([]*Team, len(tcb.builders))
	mutators := make([]Mutator, len(tcb.builders))
	for i := range tcb.builders {
		func(i int, root context.Context) {
			builder := tcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TeamMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, tcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, tcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tcb *TeamCreateBulk) SaveX(ctx context.Context) []*Team {
	v, err := tcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcb *TeamCreateBulk) Exec(ctx context.Context) error {
	_, err := tcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcb *TeamCreateBulk) ExecX(ctx context.Context) {
	if err := tcb.Exec(ctx); err != nil {
		panic(err)
	}
}
