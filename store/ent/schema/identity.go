package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"khepri.dev/horus"
)

type Identity struct {
	ent.Schema
}

func (Identity) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable().NotEmpty().Sensitive().Unique(),
		field.UUID("owner_id", uuid.UUID{}).Immutable(),
		field.String("kind").Immutable().NotEmpty().GoType(horus.IdentityKind("")),
		field.String("name").Default(""),
		field.String("verified_by").NotEmpty().GoType(horus.Verifier("")),
		field.Time("created_at").Immutable().Default(utcNow),
	}
}

func (Identity) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("identities").Field("owner_id").Immutable().Unique().Required(),
		edge.From("member", Member.Type).Ref("identities"),
	}
}
