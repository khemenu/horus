package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"github.com/lesomnus/entpb"
)

type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin{},
		aliasMixin{},
	}
}

func (User) Fields() []ent.Field {
	return []ent.Field{}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", User.Type).
			Annotations(entpb.Field(4)).
			From("parent").
			Annotations(entpb.Field(3, entpb.WithWritable())).
			Unique(),
		edge.To("identities", Identity.Type).
			Annotations(entpb.Field(5)),
		edge.To("accounts", Account.Type).
			Annotations(entpb.Field(6)),
		edge.To("tokens", Token.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
