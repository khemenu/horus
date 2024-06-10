package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/lesomnus/entpb"
)

type Team struct {
	ent.Schema
}

func (Team) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin{},
		aliasMixin{IsCommon: true},
		labelMixin{},
	}
}

func (Team) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("silo_id", uuid.UUID{}).
			Immutable(),
	}
}

func (Team) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("silo", Silo.Type).
			Annotations(entpb.Field(3)).
			Ref("teams").
			Field("silo_id").
			Immutable().
			Unique().
			Required(),
		edge.To("members", Membership.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Team) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("silo_id", "alias").Unique(),
	}
}
