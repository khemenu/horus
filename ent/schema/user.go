package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type User struct {
	grpcSchema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique().Default(uuid.New).
			Annotations(entproto.Field(1)),

		field.Time("created_date").Immutable().Default(utcNow).
			Annotations(entproto.Field(15)),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("identities", Identity.Type).
			Annotations(entproto.Field(2)),
		edge.To("accounts", Account.Type).
			Annotations(entproto.Field(3)),
		edge.To("tokens", Token.Type).
			Annotations(entsql.OnDelete(entsql.Cascade), entproto.Skip()),
	}
}
