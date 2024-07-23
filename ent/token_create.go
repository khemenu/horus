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
	"khepri.dev/horus/ent/token"
	"khepri.dev/horus/ent/user"
)

// TokenCreate is the builder for creating a Token entity.
type TokenCreate struct {
	config
	mutation *TokenMutation
	hooks    []Hook
}

// SetDateCreated sets the "date_created" field.
func (tc *TokenCreate) SetDateCreated(t time.Time) *TokenCreate {
	tc.mutation.SetDateCreated(t)
	return tc
}

// SetNillableDateCreated sets the "date_created" field if the given value is not nil.
func (tc *TokenCreate) SetNillableDateCreated(t *time.Time) *TokenCreate {
	if t != nil {
		tc.SetDateCreated(*t)
	}
	return tc
}

// SetValue sets the "value" field.
func (tc *TokenCreate) SetValue(s string) *TokenCreate {
	tc.mutation.SetValue(s)
	return tc
}

// SetType sets the "type" field.
func (tc *TokenCreate) SetType(s string) *TokenCreate {
	tc.mutation.SetType(s)
	return tc
}

// SetName sets the "name" field.
func (tc *TokenCreate) SetName(s string) *TokenCreate {
	tc.mutation.SetName(s)
	return tc
}

// SetNillableName sets the "name" field if the given value is not nil.
func (tc *TokenCreate) SetNillableName(s *string) *TokenCreate {
	if s != nil {
		tc.SetName(*s)
	}
	return tc
}

// SetUseCountLimit sets the "use_count_limit" field.
func (tc *TokenCreate) SetUseCountLimit(u uint64) *TokenCreate {
	tc.mutation.SetUseCountLimit(u)
	return tc
}

// SetNillableUseCountLimit sets the "use_count_limit" field if the given value is not nil.
func (tc *TokenCreate) SetNillableUseCountLimit(u *uint64) *TokenCreate {
	if u != nil {
		tc.SetUseCountLimit(*u)
	}
	return tc
}

// SetDateExpired sets the "date_expired" field.
func (tc *TokenCreate) SetDateExpired(t time.Time) *TokenCreate {
	tc.mutation.SetDateExpired(t)
	return tc
}

// SetID sets the "id" field.
func (tc *TokenCreate) SetID(u uuid.UUID) *TokenCreate {
	tc.mutation.SetID(u)
	return tc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (tc *TokenCreate) SetNillableID(u *uuid.UUID) *TokenCreate {
	if u != nil {
		tc.SetID(*u)
	}
	return tc
}

// SetOwnerID sets the "owner" edge to the User entity by ID.
func (tc *TokenCreate) SetOwnerID(id uuid.UUID) *TokenCreate {
	tc.mutation.SetOwnerID(id)
	return tc
}

// SetOwner sets the "owner" edge to the User entity.
func (tc *TokenCreate) SetOwner(u *User) *TokenCreate {
	return tc.SetOwnerID(u.ID)
}

// SetParentID sets the "parent" edge to the Token entity by ID.
func (tc *TokenCreate) SetParentID(id uuid.UUID) *TokenCreate {
	tc.mutation.SetParentID(id)
	return tc
}

// SetNillableParentID sets the "parent" edge to the Token entity by ID if the given value is not nil.
func (tc *TokenCreate) SetNillableParentID(id *uuid.UUID) *TokenCreate {
	if id != nil {
		tc = tc.SetParentID(*id)
	}
	return tc
}

// SetParent sets the "parent" edge to the Token entity.
func (tc *TokenCreate) SetParent(t *Token) *TokenCreate {
	return tc.SetParentID(t.ID)
}

// AddChildIDs adds the "children" edge to the Token entity by IDs.
func (tc *TokenCreate) AddChildIDs(ids ...uuid.UUID) *TokenCreate {
	tc.mutation.AddChildIDs(ids...)
	return tc
}

// AddChildren adds the "children" edges to the Token entity.
func (tc *TokenCreate) AddChildren(t ...*Token) *TokenCreate {
	ids := make([]uuid.UUID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tc.AddChildIDs(ids...)
}

// Mutation returns the TokenMutation object of the builder.
func (tc *TokenCreate) Mutation() *TokenMutation {
	return tc.mutation
}

// Save creates the Token in the database.
func (tc *TokenCreate) Save(ctx context.Context) (*Token, error) {
	tc.defaults()
	return withHooks(ctx, tc.sqlSave, tc.mutation, tc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TokenCreate) SaveX(ctx context.Context) *Token {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tc *TokenCreate) Exec(ctx context.Context) error {
	_, err := tc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tc *TokenCreate) ExecX(ctx context.Context) {
	if err := tc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tc *TokenCreate) defaults() {
	if _, ok := tc.mutation.DateCreated(); !ok {
		v := token.DefaultDateCreated()
		tc.mutation.SetDateCreated(v)
	}
	if _, ok := tc.mutation.Name(); !ok {
		v := token.DefaultName
		tc.mutation.SetName(v)
	}
	if _, ok := tc.mutation.UseCountLimit(); !ok {
		v := token.DefaultUseCountLimit
		tc.mutation.SetUseCountLimit(v)
	}
	if _, ok := tc.mutation.ID(); !ok {
		v := token.DefaultID()
		tc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tc *TokenCreate) check() error {
	if _, ok := tc.mutation.DateCreated(); !ok {
		return &ValidationError{Name: "date_created", err: errors.New(`ent: missing required field "Token.date_created"`)}
	}
	if _, ok := tc.mutation.Value(); !ok {
		return &ValidationError{Name: "value", err: errors.New(`ent: missing required field "Token.value"`)}
	}
	if v, ok := tc.mutation.Value(); ok {
		if err := token.ValueValidator(v); err != nil {
			return &ValidationError{Name: "value", err: fmt.Errorf(`ent: validator failed for field "Token.value": %w`, err)}
		}
	}
	if _, ok := tc.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "Token.type"`)}
	}
	if v, ok := tc.mutation.GetType(); ok {
		if err := token.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Token.type": %w`, err)}
		}
	}
	if _, ok := tc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Token.name"`)}
	}
	if _, ok := tc.mutation.UseCountLimit(); !ok {
		return &ValidationError{Name: "use_count_limit", err: errors.New(`ent: missing required field "Token.use_count_limit"`)}
	}
	if _, ok := tc.mutation.DateExpired(); !ok {
		return &ValidationError{Name: "date_expired", err: errors.New(`ent: missing required field "Token.date_expired"`)}
	}
	if _, ok := tc.mutation.OwnerID(); !ok {
		return &ValidationError{Name: "owner", err: errors.New(`ent: missing required edge "Token.owner"`)}
	}
	return nil
}

func (tc *TokenCreate) sqlSave(ctx context.Context) (*Token, error) {
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

func (tc *TokenCreate) createSpec() (*Token, *sqlgraph.CreateSpec) {
	var (
		_node = &Token{config: tc.config}
		_spec = sqlgraph.NewCreateSpec(token.Table, sqlgraph.NewFieldSpec(token.FieldID, field.TypeUUID))
	)
	if id, ok := tc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := tc.mutation.DateCreated(); ok {
		_spec.SetField(token.FieldDateCreated, field.TypeTime, value)
		_node.DateCreated = value
	}
	if value, ok := tc.mutation.Value(); ok {
		_spec.SetField(token.FieldValue, field.TypeString, value)
		_node.Value = value
	}
	if value, ok := tc.mutation.GetType(); ok {
		_spec.SetField(token.FieldType, field.TypeString, value)
		_node.Type = value
	}
	if value, ok := tc.mutation.Name(); ok {
		_spec.SetField(token.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := tc.mutation.UseCountLimit(); ok {
		_spec.SetField(token.FieldUseCountLimit, field.TypeUint64, value)
		_node.UseCountLimit = value
	}
	if value, ok := tc.mutation.DateExpired(); ok {
		_spec.SetField(token.FieldDateExpired, field.TypeTime, value)
		_node.DateExpired = value
	}
	if nodes := tc.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   token.OwnerTable,
			Columns: []string{token.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.user_tokens = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   token.ParentTable,
			Columns: []string{token.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(token.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.token_children = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.ChildrenIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   token.ChildrenTable,
			Columns: []string{token.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(token.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// TokenCreateBulk is the builder for creating many Token entities in bulk.
type TokenCreateBulk struct {
	config
	err      error
	builders []*TokenCreate
}

// Save creates the Token entities in the database.
func (tcb *TokenCreateBulk) Save(ctx context.Context) ([]*Token, error) {
	if tcb.err != nil {
		return nil, tcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(tcb.builders))
	nodes := make([]*Token, len(tcb.builders))
	mutators := make([]Mutator, len(tcb.builders))
	for i := range tcb.builders {
		func(i int, root context.Context) {
			builder := tcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TokenMutation)
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
func (tcb *TokenCreateBulk) SaveX(ctx context.Context) []*Token {
	v, err := tcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcb *TokenCreateBulk) Exec(ctx context.Context) error {
	_, err := tcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcb *TokenCreateBulk) ExecX(ctx context.Context) {
	if err := tcb.Exec(ctx); err != nil {
		panic(err)
	}
}
