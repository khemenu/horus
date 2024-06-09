package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/lesomnus/entpb"
)

type Silo struct {
	ent.Schema
}

func (Silo) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin{},
		aliasMixin{},
	}
}

func (Silo) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Annotations(entpb.Field(3)).
			Default("").
			MaxLen(64),
		field.String("description").
			Annotations(entpb.Field(4)).
			Default("").
			MaxLen(256),
	}
}

func (Silo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("members", Account.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("teams", Team.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("invitations", Invitation.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
