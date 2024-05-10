package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"khepri.dev/horus/alias"
)

type Team struct {
	grpcSchema
}

func (Team) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Unique().
			Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.UUID("silo_id", uuid.UUID{}).
			Immutable().
			Annotations(entproto.Skip()),
		field.String("alias").
			Unique().NotEmpty().
			DefaultFunc(alias.New).Validate(alias.ValidateE).
			Annotations(entproto.Field(2)),

		field.String("name").
			NotEmpty().MaxLen(64).
			Annotations(entproto.Field(5)),
		field.String("description").
			MaxLen(256).
			Default("").
			Annotations(entproto.Field(6)),

		field.Enum("inter_visibility").
			Comment("Team visibility to members who does not have a membership").
			Values(Visibility("").Values()...).
			Annotations(
				entproto.Field(7),
				entproto.Enum(Visibility("").Map()),
			),
		field.Enum("intra_visibility").
			Comment("Member visibility among members in the same team").
			Values(Visibility("").Values()...).
			Annotations(
				entproto.Field(8),
				entproto.Enum(Visibility("").Map()),
			),

		field.Time("created_date").
			Immutable().
			Default(utcNow).
			Annotations(entproto.Field(15)),
	}
}

func (Team) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("silo", Silo.Type).Ref("teams").Field("silo_id").Immutable().Unique().Required().
			Annotations(entproto.Field(3)),
		edge.To("members", Membership.Type).
			Annotations(entsql.OnDelete(entsql.Cascade), entproto.Skip()),
	}
}

func (Team) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("silo_id", "alias").Unique(),
	}
}
