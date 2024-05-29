package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Invitation struct {
	grpcSchema
}

func (Invitation) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Unique().
			Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.String("invitee").
			Immutable().
			NotEmpty().
			Annotations(entproto.Field(4)),

		field.String("type").
			Immutable().
			NotEmpty().
			Annotations(entproto.Field(5)),

		field.Time("date_created").
			Immutable().
			Default(time.Now).
			Annotations(entproto.Field(15)),
		field.Time("date_expired").
			Annotations(entproto.Field(14)),
		field.Time("date_accepted").
			Optional().
			Nillable().
			Annotations(entproto.Field(13)),
		field.Time("date_declined").
			Optional().
			Nillable().
			Annotations(entproto.Field(12)),
		field.Time("date_canceled").
			Optional().
			Nillable().
			Annotations(entproto.Field(11)),
	}
}

func (Invitation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("silo", Silo.Type).
			Ref("invitations").
			Immutable().
			Unique().
			Required().
			Annotations(entproto.Field(2)),
		edge.From("inviter", Account.Type).
			Ref("invitations").
			Immutable().
			Unique().
			Required().
			Annotations(entproto.Field(3)),
	}
}
