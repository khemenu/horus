package schema

import (
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Identity struct {
	grpcSchema
}

func (Identity) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable().NotEmpty().Unique().
			Annotations(entproto.Field(1)),
		// field.UUID("owner_id", uuid.UUID{}).Immutable().
		// 	Annotations(entproto.Field(2)),
		field.String("kind").Immutable().NotEmpty().
			Annotations(entproto.Field(3)),
		field.String("verifier").NotEmpty().
			Annotations(entproto.Field(4)),
		field.String("name").Default("").
			Annotations(entproto.Field(5)),

		field.Time("created_date").Immutable().Default(time.Now).
			Annotations(entproto.Field(15)),
	}
}

func (Identity) Edges() []ent.Edge {
	return []ent.Edge{
		// edge.To("owner", User.Type).Unique().Immutable().Required().
		// 	Annotations(entproto.Field(2)),
		edge.From("owner", User.Type).Ref("identities").Immutable().Unique().Required().
			Annotations(entproto.Field(2)),
	}
}
