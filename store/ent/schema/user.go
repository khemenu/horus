package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique().Default(uuid.New),
		field.String("alias").Unique(),
		field.Time("created_at").Immutable().Default(utcNow).Annotations(entsql.Default("CURRENT_TIMESTAMP")),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tokens", Token.Type),
		edge.To("identities", Identity.Type),
		edge.To("authorizer", Authorizer.Type).Unique(),
	}
}
