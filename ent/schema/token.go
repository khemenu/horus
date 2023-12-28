package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Token struct {
	grpcSchema
}

func (Token) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable().NotEmpty().Sensitive().Unique().
			Annotations(entproto.Field(1)),

		field.String("type").Immutable().NotEmpty().
			Annotations(entproto.Field(3)),
		field.String("name").Default("").
			Annotations(entproto.Field(4)),

		field.Time("created_at").Immutable().Default(utcNow).
			Annotations(entproto.Field(15)),
		field.Time("expired_at").
			Annotations(entproto.Field(14)),
	}
}

func (Token) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("tokens").Immutable().Unique().Required().
			Annotations(entproto.Field(2)),
	}
}
