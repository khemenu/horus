package server_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/server/frame"
)

type TokenTestSuite struct {
	Suite
}

func TestToken(t *testing.T) {
	s := TokenTestSuite{
		Suite: NewSuiteWithSqliteStore(),
	}
	suite.Run(t, &s)
}

func (s *TokenTestSuite) TestCreate() {
	s.Run("user cannot create tokens for other user", func() {
		_, err := s.svc.Token().Create(s.ctx, &horus.CreateTokenRequest{
			Type:  horus.TokenTypeRefresh,
			Owner: horus.UserById(s.other.Actor.ID),
		})
		s.ErrorContains(err, "Permission")
	})

	s.Run("parent user can create tokens for their child user", func() {
		child, err := s.svc.User().Create(s.ctx, &horus.CreateUserRequest{})
		s.NoError(err)

		v, err := s.svc.Token().Create(s.ctx, &horus.CreateTokenRequest{
			Type:  horus.TokenTypeRefresh,
			Owner: horus.UserByIdV(child.Id),
		})
		s.NoError(err)

		_, err = s.svc.Token().Get(s.ctx, &horus.GetTokenRequest{Key: &horus.GetTokenRequest_Id{
			Id: v.Id,
		}})
		s.ErrorContains(err, "not found")

		_, err = s.svc.Token().Get(
			frame.WithContext(s.ctx, &frame.Frame{
				Actor: &ent.User{ID: uuid.UUID(child.Id)},
			}),
			&horus.GetTokenRequest{Key: &horus.GetTokenRequest_Id{
				Id: v.Id,
			}},
		)
		s.NoError(err)
	})
}
