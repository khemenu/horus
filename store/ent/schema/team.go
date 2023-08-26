package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Team struct {
	ent.Schema
}

func (Team) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique().Default(uuid.New),
		field.UUID("org_id", uuid.UUID{}).Immutable(),
		field.String("name").Default("").MaxLen(64),
		field.Time("created_at").Immutable().Default(utcNow),
	}
}

func (Team) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("org", Org.Type).Ref("teams").Field("org_id").Immutable().Unique().Required(),
		edge.To("members", Member.Type).Through("memberships", Membership.Type),
	}
}
