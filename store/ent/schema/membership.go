package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"khepri.dev/horus"
)

type Membership struct {
	ent.Schema
}

func (Membership) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("team_id", uuid.UUID{}).Immutable(),
		field.UUID("member_id", uuid.UUID{}).Immutable(),
		field.Enum("role").GoType(horus.RoleTeam("")),
		field.Time("created_at").Immutable().Default(utcNow),
	}
}

func (Membership) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("team", Team.Type).Unique().Immutable().Required().Field("team_id"),
		edge.To("member", Member.Type).Unique().Immutable().Required().Field("member_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Membership) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("team_id", "member_id"),
	}
}
