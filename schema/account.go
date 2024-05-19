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

type Account struct {
	grpcSchemaWithList
}

func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Unique().
			Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.UUID("silo_id", uuid.UUID{}).
			Immutable().
			Annotations(entproto.Skip()),
		field.String("alias").
			NotEmpty().
			DefaultFunc(alias.New).Validate(alias.ValidateE).
			Annotations(entproto.Field(2)),

		field.String("name").
			NotEmpty().MaxLen(64).
			Annotations(entproto.Field(6)),
		field.String("description").
			MaxLen(256).
			Default("").
			Annotations(entproto.Field(7)),
		field.Enum("role").
			Values(RoleSilo("").Values()...).
			Annotations(
				entproto.Field(8),
				entproto.Enum(RoleSilo("").Map()),
			),

		field.Time("created_date").
			Immutable().
			Default(utcNow).
			Annotations(entproto.Field(15)),
	}
}

func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("accounts").Immutable().Unique().Required().
			Annotations(entproto.Field(3)),
		edge.From("silo", Silo.Type).Ref("members").Field("silo_id").Immutable().Unique().Required().
			Annotations(entproto.Field(4)),
		edge.To("memberships", Membership.Type).
			Annotations(entsql.OnDelete(entsql.Cascade), entproto.Field(5)),
		edge.To("invitations", Invitation.Type).
			Annotations(entsql.OnDelete(entsql.NoAction), entproto.Skip()),
	}
}

func (Account) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("silo_id", "alias").Unique(),
	}
}
