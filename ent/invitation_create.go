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
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/invitation"
	"khepri.dev/horus/ent/silo"
)

// InvitationCreate is the builder for creating a Invitation entity.
type InvitationCreate struct {
	config
	mutation *InvitationMutation
	hooks    []Hook
}

// SetInvitee sets the "invitee" field.
func (ic *InvitationCreate) SetInvitee(s string) *InvitationCreate {
	ic.mutation.SetInvitee(s)
	return ic
}

// SetCreatedDate sets the "created_date" field.
func (ic *InvitationCreate) SetCreatedDate(t time.Time) *InvitationCreate {
	ic.mutation.SetCreatedDate(t)
	return ic
}

// SetNillableCreatedDate sets the "created_date" field if the given value is not nil.
func (ic *InvitationCreate) SetNillableCreatedDate(t *time.Time) *InvitationCreate {
	if t != nil {
		ic.SetCreatedDate(*t)
	}
	return ic
}

// SetExpiredDate sets the "expired_date" field.
func (ic *InvitationCreate) SetExpiredDate(t time.Time) *InvitationCreate {
	ic.mutation.SetExpiredDate(t)
	return ic
}

// SetAcceptedDate sets the "accepted_date" field.
func (ic *InvitationCreate) SetAcceptedDate(t time.Time) *InvitationCreate {
	ic.mutation.SetAcceptedDate(t)
	return ic
}

// SetDeclinedDate sets the "declined_date" field.
func (ic *InvitationCreate) SetDeclinedDate(t time.Time) *InvitationCreate {
	ic.mutation.SetDeclinedDate(t)
	return ic
}

// SetCanceledDate sets the "canceled_date" field.
func (ic *InvitationCreate) SetCanceledDate(t time.Time) *InvitationCreate {
	ic.mutation.SetCanceledDate(t)
	return ic
}

// SetID sets the "id" field.
func (ic *InvitationCreate) SetID(u uuid.UUID) *InvitationCreate {
	ic.mutation.SetID(u)
	return ic
}

// SetNillableID sets the "id" field if the given value is not nil.
func (ic *InvitationCreate) SetNillableID(u *uuid.UUID) *InvitationCreate {
	if u != nil {
		ic.SetID(*u)
	}
	return ic
}

// SetSiloID sets the "silo" edge to the Silo entity by ID.
func (ic *InvitationCreate) SetSiloID(id uuid.UUID) *InvitationCreate {
	ic.mutation.SetSiloID(id)
	return ic
}

// SetSilo sets the "silo" edge to the Silo entity.
func (ic *InvitationCreate) SetSilo(s *Silo) *InvitationCreate {
	return ic.SetSiloID(s.ID)
}

// SetInviterID sets the "inviter" edge to the Account entity by ID.
func (ic *InvitationCreate) SetInviterID(id uuid.UUID) *InvitationCreate {
	ic.mutation.SetInviterID(id)
	return ic
}

// SetInviter sets the "inviter" edge to the Account entity.
func (ic *InvitationCreate) SetInviter(a *Account) *InvitationCreate {
	return ic.SetInviterID(a.ID)
}

// Mutation returns the InvitationMutation object of the builder.
func (ic *InvitationCreate) Mutation() *InvitationMutation {
	return ic.mutation
}

// Save creates the Invitation in the database.
func (ic *InvitationCreate) Save(ctx context.Context) (*Invitation, error) {
	ic.defaults()
	return withHooks(ctx, ic.sqlSave, ic.mutation, ic.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ic *InvitationCreate) SaveX(ctx context.Context) *Invitation {
	v, err := ic.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ic *InvitationCreate) Exec(ctx context.Context) error {
	_, err := ic.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ic *InvitationCreate) ExecX(ctx context.Context) {
	if err := ic.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ic *InvitationCreate) defaults() {
	if _, ok := ic.mutation.CreatedDate(); !ok {
		v := invitation.DefaultCreatedDate()
		ic.mutation.SetCreatedDate(v)
	}
	if _, ok := ic.mutation.ID(); !ok {
		v := invitation.DefaultID()
		ic.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ic *InvitationCreate) check() error {
	if _, ok := ic.mutation.Invitee(); !ok {
		return &ValidationError{Name: "invitee", err: errors.New(`ent: missing required field "Invitation.invitee"`)}
	}
	if v, ok := ic.mutation.Invitee(); ok {
		if err := invitation.InviteeValidator(v); err != nil {
			return &ValidationError{Name: "invitee", err: fmt.Errorf(`ent: validator failed for field "Invitation.invitee": %w`, err)}
		}
	}
	if _, ok := ic.mutation.CreatedDate(); !ok {
		return &ValidationError{Name: "created_date", err: errors.New(`ent: missing required field "Invitation.created_date"`)}
	}
	if _, ok := ic.mutation.ExpiredDate(); !ok {
		return &ValidationError{Name: "expired_date", err: errors.New(`ent: missing required field "Invitation.expired_date"`)}
	}
	if _, ok := ic.mutation.AcceptedDate(); !ok {
		return &ValidationError{Name: "accepted_date", err: errors.New(`ent: missing required field "Invitation.accepted_date"`)}
	}
	if _, ok := ic.mutation.DeclinedDate(); !ok {
		return &ValidationError{Name: "declined_date", err: errors.New(`ent: missing required field "Invitation.declined_date"`)}
	}
	if _, ok := ic.mutation.CanceledDate(); !ok {
		return &ValidationError{Name: "canceled_date", err: errors.New(`ent: missing required field "Invitation.canceled_date"`)}
	}
	if _, ok := ic.mutation.SiloID(); !ok {
		return &ValidationError{Name: "silo", err: errors.New(`ent: missing required edge "Invitation.silo"`)}
	}
	if _, ok := ic.mutation.InviterID(); !ok {
		return &ValidationError{Name: "inviter", err: errors.New(`ent: missing required edge "Invitation.inviter"`)}
	}
	return nil
}

func (ic *InvitationCreate) sqlSave(ctx context.Context) (*Invitation, error) {
	if err := ic.check(); err != nil {
		return nil, err
	}
	_node, _spec := ic.createSpec()
	if err := sqlgraph.CreateNode(ctx, ic.driver, _spec); err != nil {
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
	ic.mutation.id = &_node.ID
	ic.mutation.done = true
	return _node, nil
}

func (ic *InvitationCreate) createSpec() (*Invitation, *sqlgraph.CreateSpec) {
	var (
		_node = &Invitation{config: ic.config}
		_spec = sqlgraph.NewCreateSpec(invitation.Table, sqlgraph.NewFieldSpec(invitation.FieldID, field.TypeUUID))
	)
	if id, ok := ic.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := ic.mutation.Invitee(); ok {
		_spec.SetField(invitation.FieldInvitee, field.TypeString, value)
		_node.Invitee = value
	}
	if value, ok := ic.mutation.CreatedDate(); ok {
		_spec.SetField(invitation.FieldCreatedDate, field.TypeTime, value)
		_node.CreatedDate = value
	}
	if value, ok := ic.mutation.ExpiredDate(); ok {
		_spec.SetField(invitation.FieldExpiredDate, field.TypeTime, value)
		_node.ExpiredDate = &value
	}
	if value, ok := ic.mutation.AcceptedDate(); ok {
		_spec.SetField(invitation.FieldAcceptedDate, field.TypeTime, value)
		_node.AcceptedDate = &value
	}
	if value, ok := ic.mutation.DeclinedDate(); ok {
		_spec.SetField(invitation.FieldDeclinedDate, field.TypeTime, value)
		_node.DeclinedDate = &value
	}
	if value, ok := ic.mutation.CanceledDate(); ok {
		_spec.SetField(invitation.FieldCanceledDate, field.TypeTime, value)
		_node.CanceledDate = &value
	}
	if nodes := ic.mutation.SiloIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   invitation.SiloTable,
			Columns: []string{invitation.SiloColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(silo.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.silo_invitations = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ic.mutation.InviterIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   invitation.InviterTable,
			Columns: []string{invitation.InviterColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(account.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.account_invitations = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// InvitationCreateBulk is the builder for creating many Invitation entities in bulk.
type InvitationCreateBulk struct {
	config
	err      error
	builders []*InvitationCreate
}

// Save creates the Invitation entities in the database.
func (icb *InvitationCreateBulk) Save(ctx context.Context) ([]*Invitation, error) {
	if icb.err != nil {
		return nil, icb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(icb.builders))
	nodes := make([]*Invitation, len(icb.builders))
	mutators := make([]Mutator, len(icb.builders))
	for i := range icb.builders {
		func(i int, root context.Context) {
			builder := icb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*InvitationMutation)
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
					_, err = mutators[i+1].Mutate(root, icb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, icb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, icb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (icb *InvitationCreateBulk) SaveX(ctx context.Context) []*Invitation {
	v, err := icb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (icb *InvitationCreateBulk) Exec(ctx context.Context) error {
	_, err := icb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (icb *InvitationCreateBulk) ExecX(ctx context.Context) {
	if err := icb.Exec(ctx); err != nil {
		panic(err)
	}
}