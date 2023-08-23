package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Token struct {
	ent.Schema
}

func (Token) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("value").Immutable().NotEmpty().Sensitive().Unique(),
		field.UUID("owner_id", uuid.UUID{}).Immutable(),
		field.String("type").Immutable().NotEmpty(),
		field.String("name").Default(""),
		field.Time("created_at").Immutable().Default(utcNow),
		field.Time("expired_at"),
	}
}

func (Token) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("tokens").Field("owner_id").Immutable().Unique().Required(),
	}
}
