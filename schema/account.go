package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/lesomnus/entpb"
	"khepri.dev/horus/role"
)

type Account struct {
	ent.Schema
}

func (Account) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin{},
		aliasMixin{},
	}
}

func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("owner_id", uuid.UUID{}).
			Immutable(),
		field.UUID("silo_id", uuid.UUID{}).
			Immutable(),

		field.String("name").
			Annotations(entpb.Field(6)).
			Default("").
			MaxLen(64),
		field.String("description").
			Annotations(entpb.Field(7)).
			Default("").
			MaxLen(256),
		field.Enum("role").
			Annotations(entpb.Field(8)).
			GoType(role.Role("")),
	}
}

func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Annotations(entpb.Field(3)).
			Ref("accounts").
			Field("owner_id").
			Immutable().
			Unique().
			Required(),
		edge.From("silo", Silo.Type).
			Annotations(entpb.Field(4)).
			Ref("members").
			Field("silo_id").
			Immutable().
			Unique().
			Required(),
		edge.To("memberships", Membership.Type).
			Annotations(entsql.OnDelete(entsql.Cascade), entpb.Field(5)),
		edge.To("invitations", Invitation.Type).
			Annotations(entsql.OnDelete(entsql.NoAction)),
	}
}

func (Account) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("silo_id", "owner_id").Unique(),
		index.Fields("silo_id", "alias").Unique(),
	}
}
