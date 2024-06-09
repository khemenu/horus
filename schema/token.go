package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/lesomnus/entpb"
)

type Token struct {
	ent.Schema
}

func (Token) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin{},
	}
}

func (Token) Fields() []ent.Field {
	return []ent.Field{
		field.String("value").
			Annotations(entpb.Field(2)).
			Immutable().
			NotEmpty().
			Unique().
			Sensitive(),

		field.String("type").
			Annotations(entpb.Field(3)).
			Immutable().
			NotEmpty(),

		field.String("name").
			Annotations(entpb.Field(4)).
			Default(""),

		field.Time("date_expired").
			Annotations(entpb.Field(14)),
	}
}

func (Token) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Annotations(entpb.Field(5)).
			Ref("tokens").
			Immutable().
			Unique().
			Required(),
		edge.To("children", Token.Type).
			Annotations(entpb.Field(7)).
			From("parent").
			Annotations(entpb.Field(6)).
			Immutable().
			Unique(),
	}
}
