package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
)

type Silo struct {
	ent.Schema
}

func (Silo) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin{},
		aliasMixin{},
		labelMixin{},
	}
}

func (Silo) Fields() []ent.Field {
	return []ent.Field{}
}

func (Silo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("accounts", Account.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("teams", Team.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("invitations", Invitation.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
