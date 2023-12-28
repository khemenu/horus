package frame

import (
	"context"
	"fmt"
	"time"

	"khepri.dev/horus"
	"khepri.dev/horus/log"
)

type Frame interface {
	ExpiredAt() time.Time
	User(ctx context.Context) (*horus.User, error)
}

type frame struct {
	horus horus.Horus

	access_token *horus.Token
	user         *horus.User
}

func NewFrame(h horus.Horus, access_token *horus.Token) Frame {
	if access_token.Type != horus.AccessToken {
		panic("token must be an access token")
	}

	return &frame{horus: h, access_token: access_token}
}

func (f *frame) ExpiredAt() time.Time {
	return f.access_token.ExpiredAt
}

func (f *frame) User(ctx context.Context) (*horus.User, error) {
	if f.user != nil {
		return f.user, nil
	}

	user, err := f.horus.Users().GetById(ctx, f.access_token.OwnerId)
	if err != nil {
		return nil, fmt.Errorf("get user details: %w", err)
	}

	log.FromCtx(ctx).Info("verified", "user_id", user.Id)

	f.user = user
	return user, nil
}
