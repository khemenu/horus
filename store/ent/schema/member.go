package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"khepri.dev/horus"
)

type Member struct {
	ent.Schema
}

func (Member) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Unique().Default(uuid.New),
		field.UUID("org_id", uuid.UUID{}).Immutable(),
		field.UUID("user_id", uuid.UUID{}).Immutable(),
		field.Enum("role").GoType(horus.RoleOrg("")),
		field.String("name").Default("").MaxLen(64),
		field.Time("created_at").Immutable().Default(utcNow),
	}
}

func (Member) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("belongs").Field("user_id").Immutable().Unique().Required(),
		edge.From("org", Org.Type).Ref("members").Field("org_id").Immutable().Unique().Required(),
		edge.From("teams", Team.Type).Ref("members").Through("memberships", Membership.Type),
		edge.To("contacts", Identity.Type),
	}
}

func (Member) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("org_id", "user_id").Unique(),
	}
}
