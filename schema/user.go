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

type User struct {
	grpcSchema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Unique().
			Immutable().
			Default(uuid.New).
			Annotations(entproto.Field(1)),
		field.String("alias").
			Unique().
			NotEmpty().
			MaxLen(32).
			DefaultFunc(alias.New).
			Validate(alias.ValidateE).
			Annotations(entproto.Field(2)),

		field.Time("date_created").
			Immutable().
			Default(time.Now).
			Annotations(entproto.Field(15)),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", User.Type).
			Annotations(entproto.Field(4)).
			From("parent").
			Unique().
			Annotations(entproto.Field(3)),
		edge.To("identities", Identity.Type).
			Annotations(entproto.Field(5)),
		edge.To("accounts", Account.Type).
			Annotations(entproto.Field(6)),
		edge.To("tokens", Token.Type).
			Annotations(entsql.OnDelete(entsql.Cascade), entproto.Skip()),
	}
}
