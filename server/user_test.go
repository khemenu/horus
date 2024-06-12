package server_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
)

type UserTestSuite struct {
	Suite
}

func TestUser(t *testing.T) {
	s := UserTestSuite{
		Suite: NewSuiteWithSqliteStore(),
	}
	suite.Run(t, &s)
}

func (s *UserTestSuite) TestCreate() {
	s.Run("current user is set as the new user's parent", func() {
		v, err := s.svc.User().Create(s.ctx, nil)
		s.NoError(err)

		v, err = s.svc.User().Get(s.ctx, &horus.GetUserRequest{Key: &horus.GetUserRequest_Id{
			Id: v.Id,
		}})
		s.NoError(err)
		s.Equal(s.me.Actor.ID[:], v.GetParent().GetId())
	})
}
