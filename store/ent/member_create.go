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
	"khepri.dev/horus"
	"khepri.dev/horus/store/ent/identity"
	"khepri.dev/horus/store/ent/member"
	"khepri.dev/horus/store/ent/org"
	"khepri.dev/horus/store/ent/team"
	"khepri.dev/horus/store/ent/user"
)

// MemberCreate is the builder for creating a Member entity.
type MemberCreate struct {
	config
	mutation *MemberMutation
	hooks    []Hook
}

// SetOrgID sets the "org_id" field.
func (mc *MemberCreate) SetOrgID(u uuid.UUID) *MemberCreate {
	mc.mutation.SetOrgID(u)
	return mc
}

// SetUserID sets the "user_id" field.
func (mc *MemberCreate) SetUserID(u uuid.UUID) *MemberCreate {
	mc.mutation.SetUserID(u)
	return mc
}

// SetRole sets the "role" field.
func (mc *MemberCreate) SetRole(ho horus.RoleOrg) *MemberCreate {
	mc.mutation.SetRole(ho)
	return mc
}

// SetName sets the "name" field.
func (mc *MemberCreate) SetName(s string) *MemberCreate {
	mc.mutation.SetName(s)
	return mc
}

// SetNillableName sets the "name" field if the given value is not nil.
func (mc *MemberCreate) SetNillableName(s *string) *MemberCreate {
	if s != nil {
		mc.SetName(*s)
	}
	return mc
}

// SetCreatedAt sets the "created_at" field.
func (mc *MemberCreate) SetCreatedAt(t time.Time) *MemberCreate {
	mc.mutation.SetCreatedAt(t)
	return mc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (mc *MemberCreate) SetNillableCreatedAt(t *time.Time) *MemberCreate {
	if t != nil {
		mc.SetCreatedAt(*t)
	}
	return mc
}

// SetID sets the "id" field.
func (mc *MemberCreate) SetID(u uuid.UUID) *MemberCreate {
	mc.mutation.SetID(u)
	return mc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (mc *MemberCreate) SetNillableID(u *uuid.UUID) *MemberCreate {
	if u != nil {
		mc.SetID(*u)
	}
	return mc
}

// SetUser sets the "user" edge to the User entity.
func (mc *MemberCreate) SetUser(u *User) *MemberCreate {
	return mc.SetUserID(u.ID)
}

// SetOrg sets the "org" edge to the Org entity.
func (mc *MemberCreate) SetOrg(o *Org) *MemberCreate {
	return mc.SetOrgID(o.ID)
}

// AddTeamIDs adds the "teams" edge to the Team entity by IDs.
func (mc *MemberCreate) AddTeamIDs(ids ...uuid.UUID) *MemberCreate {
	mc.mutation.AddTeamIDs(ids...)
	return mc
}

// AddTeams adds the "teams" edges to the Team entity.
func (mc *MemberCreate) AddTeams(t ...*Team) *MemberCreate {
	ids := make([]uuid.UUID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return mc.AddTeamIDs(ids...)
}

// AddIdentityIDs adds the "identities" edge to the Identity entity by IDs.
func (mc *MemberCreate) AddIdentityIDs(ids ...string) *MemberCreate {
	mc.mutation.AddIdentityIDs(ids...)
	return mc
}

// AddIdentities adds the "identities" edges to the Identity entity.
func (mc *MemberCreate) AddIdentities(i ...*Identity) *MemberCreate {
	ids := make([]string, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return mc.AddIdentityIDs(ids...)
}

// Mutation returns the MemberMutation object of the builder.
func (mc *MemberCreate) Mutation() *MemberMutation {
	return mc.mutation
}

// Save creates the Member in the database.
func (mc *MemberCreate) Save(ctx context.Context) (*Member, error) {
	mc.defaults()
	return withHooks(ctx, mc.sqlSave, mc.mutation, mc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (mc *MemberCreate) SaveX(ctx context.Context) *Member {
	v, err := mc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mc *MemberCreate) Exec(ctx context.Context) error {
	_, err := mc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mc *MemberCreate) ExecX(ctx context.Context) {
	if err := mc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (mc *MemberCreate) defaults() {
	if _, ok := mc.mutation.Name(); !ok {
		v := member.DefaultName
		mc.mutation.SetName(v)
	}
	if _, ok := mc.mutation.CreatedAt(); !ok {
		v := member.DefaultCreatedAt()
		mc.mutation.SetCreatedAt(v)
	}
	if _, ok := mc.mutation.ID(); !ok {
		v := member.DefaultID()
		mc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mc *MemberCreate) check() error {
	if _, ok := mc.mutation.OrgID(); !ok {
		return &ValidationError{Name: "org_id", err: errors.New(`ent: missing required field "Member.org_id"`)}
	}
	if _, ok := mc.mutation.UserID(); !ok {
		return &ValidationError{Name: "user_id", err: errors.New(`ent: missing required field "Member.user_id"`)}
	}
	if _, ok := mc.mutation.Role(); !ok {
		return &ValidationError{Name: "role", err: errors.New(`ent: missing required field "Member.role"`)}
	}
	if v, ok := mc.mutation.Role(); ok {
		if err := member.RoleValidator(v); err != nil {
			return &ValidationError{Name: "role", err: fmt.Errorf(`ent: validator failed for field "Member.role": %w`, err)}
		}
	}
	if _, ok := mc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Member.name"`)}
	}
	if v, ok := mc.mutation.Name(); ok {
		if err := member.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Member.name": %w`, err)}
		}
	}
	if _, ok := mc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Member.created_at"`)}
	}
	if _, ok := mc.mutation.UserID(); !ok {
		return &ValidationError{Name: "user", err: errors.New(`ent: missing required edge "Member.user"`)}
	}
	if _, ok := mc.mutation.OrgID(); !ok {
		return &ValidationError{Name: "org", err: errors.New(`ent: missing required edge "Member.org"`)}
	}
	return nil
}

func (mc *MemberCreate) sqlSave(ctx context.Context) (*Member, error) {
	if err := mc.check(); err != nil {
		return nil, err
	}
	_node, _spec := mc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mc.driver, _spec); err != nil {
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
	mc.mutation.id = &_node.ID
	mc.mutation.done = true
	return _node, nil
}

func (mc *MemberCreate) createSpec() (*Member, *sqlgraph.CreateSpec) {
	var (
		_node = &Member{config: mc.config}
		_spec = sqlgraph.NewCreateSpec(member.Table, sqlgraph.NewFieldSpec(member.FieldID, field.TypeUUID))
	)
	if id, ok := mc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := mc.mutation.Role(); ok {
		_spec.SetField(member.FieldRole, field.TypeEnum, value)
		_node.Role = value
	}
	if value, ok := mc.mutation.Name(); ok {
		_spec.SetField(member.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := mc.mutation.CreatedAt(); ok {
		_spec.SetField(member.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if nodes := mc.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   member.UserTable,
			Columns: []string{member.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.UserID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := mc.mutation.OrgIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   member.OrgTable,
			Columns: []string{member.OrgColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(org.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.OrgID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := mc.mutation.TeamsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   member.TeamsTable,
			Columns: member.TeamsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(team.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		createE := &MembershipCreate{config: mc.config, mutation: newMembershipMutation(mc.config, OpCreate)}
		createE.defaults()
		_, specE := createE.createSpec()
		edge.Target.Fields = specE.Fields
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := mc.mutation.IdentitiesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   member.IdentitiesTable,
			Columns: member.IdentitiesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(identity.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// MemberCreateBulk is the builder for creating many Member entities in bulk.
type MemberCreateBulk struct {
	config
	builders []*MemberCreate
}

// Save creates the Member entities in the database.
func (mcb *MemberCreateBulk) Save(ctx context.Context) ([]*Member, error) {
	specs := make([]*sqlgraph.CreateSpec, len(mcb.builders))
	nodes := make([]*Member, len(mcb.builders))
	mutators := make([]Mutator, len(mcb.builders))
	for i := range mcb.builders {
		func(i int, root context.Context) {
			builder := mcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*MemberMutation)
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
					_, err = mutators[i+1].Mutate(root, mcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, mcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (mcb *MemberCreateBulk) SaveX(ctx context.Context) []*Member {
	v, err := mcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mcb *MemberCreateBulk) Exec(ctx context.Context) error {
	_, err := mcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mcb *MemberCreateBulk) ExecX(ctx context.Context) {
	if err := mcb.Exec(ctx); err != nil {
		panic(err)
	}
}
