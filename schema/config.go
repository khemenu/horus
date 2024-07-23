package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/lesomnus/entpb"
)

type Conf struct {
	ent.Schema
}

type baseMixinWithStringKey struct {
	baseMixin
}

func (m baseMixinWithStringKey) Fields() []ent.Field {
	fields := m.baseMixin.Fields()
	for i, f := range fields {
		if f.Descriptor().Name != "id" {
			continue
		}

		fields[i] = field.String("id").
			Annotations(entpb.Field(1, entpb.WithReadOnly())).
			Unique().
			Immutable()
	}

	return fields
}

func (Conf) Mixin() []ent.Mixin {
	return []ent.Mixin{
		baseMixinWithStringKey{},
	}
}

func (Conf) Fields() []ent.Field {
	return []ent.Field{
		field.String("value").
			Annotations(entpb.Field(2)).
			NotEmpty(),

		field.Time("date_updated").
			Annotations(entpb.Field(14, entpb.WithReadOnly())).
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}
