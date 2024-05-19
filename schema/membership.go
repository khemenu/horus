package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Membership struct {
	grpcSchema
}

func (Membership) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique().Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.Enum("role").
			Values(RoleTeam("").Values()...).
			Annotations(
				entproto.Field(4),
				entproto.Enum(RoleTeam("").Map()),
			),

		field.Time("created_date").Immutable().Default(utcNow).
			Annotations(entproto.Field(15)),
	}
}

func (Membership) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).Ref("memberships").
			Immutable().Unique().Required().
			Annotations(entproto.Field(2)),
		edge.From("team", Team.Type).Ref("members").
			Immutable().Unique().Required().
			Annotations(entproto.Field(3)),
		// edge.To("team", Team.Type).Unique().Immutable().Required().Field("team_id").
		// 	Annotations(entsql.OnDelete(entsql.Cascade), entproto.Field(2)),
		// edge.To("member", Account.Type).Unique().Immutable().Required().Field("account_id").
		// 	Annotations(entsql.OnDelete(entsql.Cascade), entproto.Field(3)),
	}
}
