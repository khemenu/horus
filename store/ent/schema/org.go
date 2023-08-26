package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Org struct {
	ent.Schema
}

func (Org) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique().Default(uuid.New),
		field.String("name").Default("").MaxLen(64),
		field.Time("created_at").Immutable().Default(utcNow),
	}
}

func (Org) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("teams", Team.Type),
		edge.To("members", Member.Type),
	}
}
