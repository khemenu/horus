package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Token struct {
	grpcSchema
}

func (Token) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Unique().
			Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.String("value").
			Immutable().
			NotEmpty().
			Unique().
			Sensitive().
			Annotations(entproto.Field(2)),

		field.String("type").
			Immutable().
			NotEmpty().
			Annotations(entproto.Field(4)),
		field.String("name").
			Default("").
			Annotations(entproto.Field(5)),

		field.Time("create_date").
			Immutable().
			Default(utcNow).
			Annotations(entproto.Field(15)),
		field.Time("expired_date").
			Annotations(entproto.Field(14)),
	}
}

func (Token) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("tokens").Immutable().Unique().Required().
			Annotations(entproto.Field(3)),
	}
}
