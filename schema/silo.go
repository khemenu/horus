package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"khepri.dev/horus/alias"
)

type Silo struct {
	grpcSchema
}

func (Silo) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Unique().
			Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.String("alias").
			Unique().
			NotEmpty().
			DefaultFunc(alias.New).
			Validate(alias.ValidateE).
			Annotations(entproto.Field(2)),

		field.String("name").
			Default("").
			MaxLen(64).
			Annotations(entproto.Field(3)),
		field.String("description").
			Default("").
			MaxLen(256).
			Annotations(entproto.Field(4)),

		field.Time("date_created").
			Immutable().
			Default(time.Now).
			Annotations(entproto.Field(15)),
	}
}

func (Silo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("members", Account.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entproto.Skip(),
			),
		edge.To("teams", Team.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entproto.Skip(),
			),
		edge.To("invitations", Invitation.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
				entproto.Skip(),
			),
	}
}
