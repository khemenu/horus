package service

import (
	"context"
	"errors"
	"fmt"

	"khepri.dev/horus"
)

type authService struct {
	stores horus.Stores
}

func (s *authService) SignUp(ctx context.Context, init horus.IdentityInit) (*horus.User, error) {
	if _, err := s.stores.Identities().GetByValue(ctx, init.Value); err == nil {
		return nil, horus.ErrExist
	} else {
		if !errors.Is(err, horus.ErrNotExist) {
			return nil, fmt.Errorf("get identity: %w", err)
		}
	}

	user, err := s.stores.Users().New(ctx)
	if err != nil {
		return nil, fmt.Errorf("new user: %w", err)
	}

	_, err = s.stores.Identities().Create(ctx, &horus.Identity{
		IdentityInit: init,
		OwnerId:      user.Id,
	})
	if err != nil {
		return nil, fmt.Errorf("create identity: %w", err)
	}

	return user, nil
}

// func (s *AuthService) SignIn(ctx context.Context, user_id uuid.UUID) (*Token, error) {
// 	access_token, err := s.Stores.Token.Issue(ctx, user_id, 6*time.Hour)
// 	if err != nil {
// 		return nil, fmt.Errorf("issue access token: %w", err)
// 	}

// 	refresh_token, err := s.Stores.Token.Issue(ctx, user_id, 6*30*24*time.Hour)
// 	if err != nil {

// 	}
// }
