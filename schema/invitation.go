package schema

import (
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
		field.UUID("id", uuid.UUID{}).Unique().Default(uuid.New).
			Annotations(entproto.Field(1)),
		// field.UUID("silo_id", uuid.UUID{}).Immutable().
		// 	Annotations(entproto.Field(1)),
		// field.UUID("inviter_id", uuid.UUID{}).Immutable().
		// 	Annotations(entproto.Field(1)),
		field.String("invitee").Immutable().NotEmpty().
			Annotations(entproto.Field(4)),

		field.Time("created_date").Immutable().Default(utcNow).
			Annotations(entproto.Field(15)),
		field.Time("expired_date").Nillable().
			Annotations(entproto.Field(14)),
		field.Time("accepted_date").Nillable().
			Annotations(entproto.Field(13)),
		field.Time("declined_date").Nillable().
			Annotations(entproto.Field(12)),
		field.Time("canceled_date").Nillable().
			Annotations(entproto.Field(11)),
	}
}

func (Invitation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("silo", Silo.Type).Ref("invitations").Immutable().Unique().Required().
			Annotations(entproto.Field(2)),
		edge.From("inviter", Account.Type).Ref("invitations").Immutable().Unique().Required().
			Annotations(entproto.Field(3)),
	}
}
