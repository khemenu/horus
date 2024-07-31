package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
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
		field.UUID("owner_id", uuid.UUID{}).
			Immutable(),
		field.String("kind").
			Annotations(entpb.Field(3)).
			Immutable().
			NotEmpty(),
		field.String("value").
			Annotations(entpb.Field(4)).
			Immutable(),
		field.String("verifier").
			Annotations(entpb.Field(5, entpb.WithReadOnly())).
			Optional().
			Nillable(),
	}
}

func (Identity) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Annotations(entpb.Field(2)).
			Ref("identities").
			Field("owner_id").
			Immutable().
			Unique().
			Required(),
	}
}

func (Identity) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("owner_id", "kind", "value").Unique(),
	}
}

func (Identity) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entpb.Message(entpb.PathInherit,
			entpb.WithService(entpb.PathInherit,
				&entpb.Rpc{
					Ident: "List",
					Req:   entpb.PbType{Ident: "ListIdentityRequest", Import: entpb.PbThis.Import},
					Res:   entpb.PbType{Ident: "ListIdentityResponse", Import: entpb.PbThis.Import},
				},
			),
		),
	}
}
