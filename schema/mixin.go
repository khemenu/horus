package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
	"github.com/lesomnus/entpb"
	"khepri.dev/horus/alias"
)

type baseMixin struct {
	mixin.Schema
}

func (baseMixin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Annotations(entpb.Field(1, entpb.WithReadOnly())).
			Unique().
			Immutable().
			Default(uuid.New),

		field.Time("date_created").
			Annotations(entpb.Field(15, entpb.WithReadOnly())).
			Immutable().
			Default(time.Now),
	}
}

func (baseMixin) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entpb.Message("khepri/horus/common.proto",
			entpb.WithService("khepri/horus/store.proto",
				entpb.RpcEntCreate(),
				entpb.RpcEntGet(),
				entpb.RpcEntUpdate(),
				entpb.RpcEntDelete(),
			),
		),
	}
}

type aliasMixin struct {
	mixin.Schema
}

func (aliasMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("alias").
			Annotations(entpb.Field(2)).
			Unique().
			NotEmpty().
			MaxLen(32).
			DefaultFunc(alias.New).
			Validate(alias.ValidateE),
	}
}
