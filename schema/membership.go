package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/lesomnus/entpb"
	"khepri.dev/horus/role"
)

type Membership struct {
	ent.Schema
}

func (Membership) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin{},
	}
}
func (Membership) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("role").
			Annotations(entpb.Field(6)).
			GoType(role.Role("")),
	}
}

func (Membership) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Annotations(entpb.Field(2)).
			Ref("memberships").
			Immutable().
			Unique().
			Required(),
		edge.From("team", Team.Type).
			Annotations(entpb.Field(3)).
			Ref("members").
			Immutable().
			Unique().
			Required(),
	}
}
