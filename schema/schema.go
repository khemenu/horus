package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
)

type grpcSchema struct {
	ent.Schema
}

func (grpcSchema) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(
			entproto.PackageName("khepri.horus"),
		),
		entproto.Service(entproto.Methods(entproto.MethodCreate |
			entproto.MethodGet |
			entproto.MethodUpdate |
			entproto.MethodDelete,
		)),
	}
}

type grpcSchemaWithList struct {
	ent.Schema
}

func (grpcSchemaWithList) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(
			entproto.PackageName("khepri.horus"),
		),
		entproto.Service(entproto.Methods(entproto.MethodCreate |
			entproto.MethodGet |
			entproto.MethodList |
			entproto.MethodUpdate |
			entproto.MethodDelete,
		)),
	}
}
