package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Authorizer struct {
	ent.Schema
}

func (Authorizer) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("owner_id", uuid.UUID{}).Immutable().Unique(),
		field.String("primary_id").NotEmpty(),
		// field.Strings("secondaries"), Not implemented

	}
}

func (Authorizer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("primary", Identity.Type).Field("primary_id").Required().Unique(),
		edge.From("owner", User.Type).Ref("authorizer").Field("owner_id").Immutable().Unique().Required(),
	}
}
