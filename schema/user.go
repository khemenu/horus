package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/lesomnus/entpb"
)

type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin{},
		aliasMixin{},
	}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("parent_id", uuid.UUID{}).
			Optional().
			Nillable(),

		field.Uint("sign_in_attempt_count").
			Default(0),

		field.Time("date_unlocked").
			Comment("For users created by other users, this value is initially NULL.").
			Annotations(entpb.Field(14, entpb.WithReadOnly())).
			Optional().
			Nillable(),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", User.Type).
			Annotations(entpb.Field(4)).
			From("parent").
			Annotations(entpb.Field(3, entpb.WithWritable())).
			Field("parent_id").
			Unique(),
		edge.To("identities", Identity.Type).
			Annotations(entpb.Field(5)),
		edge.To("accounts", Account.Type).
			Annotations(entpb.Field(6)),
		edge.To("tokens", Token.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
