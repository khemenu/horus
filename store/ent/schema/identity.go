package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Identity struct {
	ent.Schema
}

func (Identity) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable().NotEmpty().Sensitive().Unique(),
		field.UUID("owner_id", uuid.UUID{}).Immutable(),
		field.String("name").Default(""),
		field.String("kind").Immutable().NotEmpty(),
		field.String("verified_by").NotEmpty(),
		field.Time("created_at").Immutable().Default(utcNow),
	}
}

func (Identity) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("identities").Field("owner_id").Immutable().Unique().Required(),
	}
}
