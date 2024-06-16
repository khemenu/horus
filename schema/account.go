package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/lesomnus/entpb"
	"khepri.dev/horus/role"
)

type Account struct {
	ent.Schema
}

func (Account) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixin{},
		aliasMixin{IsCommon: true},
		labelMixin{},
	}
}

func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("owner_id", uuid.UUID{}).
			Immutable(),
		field.UUID("silo_id", uuid.UUID{}).
			Immutable(),

		field.Enum("role").
			Annotations(entpb.Field(6)).
			GoType(role.Role("")),
	}
}

func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Annotations(entpb.Field(3)).
			Ref("accounts").
			Field("owner_id").
			Immutable().
			Unique().
			Required(),
		edge.From("silo", Silo.Type).
			Annotations(entpb.Field(4)).
			Ref("members").
			Field("silo_id").
			Immutable().
			Unique().
			Required(),
		edge.To("memberships", Membership.Type).
			Annotations(entsql.OnDelete(entsql.Cascade), entpb.Field(5)),
		edge.To("invitations", Invitation.Type).
			Annotations(entsql.OnDelete(entsql.NoAction)),
	}
}

func (Account) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("silo_id", "owner_id").Unique(),
		index.Fields("silo_id", "alias").Unique(),
	}
}

func (Account) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entpb.Message(entpb.PathInherit,
			entpb.WithService(entpb.PathInherit,
				&entpb.Rpc{
					Ident: "List",
					Req:   entpb.PbType{Ident: "ListAccountRequest", Import: "khepri/horus/extend.proto"},
					Res:   entpb.PbType{Ident: "ListAccountResponse", Import: "khepri/horus/extend.proto"},
				},
			),
		),
	}
}
