package resolve

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus/ent"
)

type AliasedEntity interface {
	GetId() []byte
	GetAlias() string
}

type Query[T ~func(*sql.Selector), U any] interface {
	Where(ps ...T) U
	OnlyID(ctx context.Context) (id uuid.UUID, err error)
}

type Queryer[U Query[T, U], T ~func(*sql.Selector)] interface {
	Query() U
}

func Uuid[T ~func(*sql.Selector), U Queryer[V, T], V Query[T, V]](ctx context.Context, id string, q U, p func(id string) T) (uuid.UUID, error) {
	v, err := uuid.Parse(id)
	if err == nil {
		return v, nil
	}

	v, err = q.Query().Where(p(id)).OnlyID(ctx)
	if err == nil {
		return v, nil
	}

	if ent.IsNotFound(err) {
		err = status.Error(codes.NotFound, "not found")
	}
	return uuid.Nil, fmt.Errorf("query: %w", err)
}
