package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Team struct {
	ent.Schema
}

func (Team) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique().Default(uuid.New),
		field.UUID("org_id", uuid.UUID{}).Immutable(),
		field.String("name").NotEmpty().MaxLen(64),
		field.Time("created_at").Immutable().Default(utcNow),
	}
}

func (Team) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("org", Org.Type).Ref("teams").Field("org_id").Immutable().Unique().Required(),
		edge.To("members", Member.Type).Through("memberships", Membership.Type),
	}
}

func (Team) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("org_id", "name").Unique(),
	}
}
