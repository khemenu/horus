package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/lesomnus/entpb"
)

type Identity struct {
	ent.Schema
}

func (Identity) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin{},
		labelMixin{},
	}
}

func (Identity) Fields() []ent.Field {
	return []ent.Field{
		field.String("kind").
			Annotations(entpb.Field(3)).
			Immutable().
			NotEmpty(),
		field.String("verifier").
			Annotations(entpb.Field(4)).
			NotEmpty(),
	}
}

func (Identity) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Annotations(entpb.Field(2)).
			Ref("identities").
			Immutable().
			Unique().
			Required(),
	}
}
