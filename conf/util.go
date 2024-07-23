package conf

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
)

type Conf interface {
	Id() string
}

func UnmarshalFrom[C Conf](ctx context.Context, server horus.ConfServiceServer, c C) error {
	res, err := server.Get(ctx, horus.ConfById(c.Id()))
	if err != nil {
		s, _ := status.FromError(err)
		if s.Code() == codes.NotFound {
			return nil
		}

		return err
	}
	if err := json.Unmarshal([]byte(res.Value), c); err != nil {
		return fmt.Errorf("invalid value: %w", err)
	}

	return nil
}

func MarshalInto[C Conf](ctx context.Context, c C, server horus.ConfServiceServer) error {
	v, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = server.Create(ctx, &horus.CreateConfRequest{
		Id:    c.Id(),
		Value: string(v),
	})
	if err != nil {
		return err
	}

	return nil
}
