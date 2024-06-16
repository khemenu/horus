package server_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

		v, err = s.svc.User().Get(s.ctx, horus.UserByIdV(v.Id))
		s.NoError(err)
		s.Equal(s.me.Actor.ID[:], v.GetParent().GetId())
	})
}

func (t *UserTestSuite) TestGet() {
	t.Run("alias _me returns me", func() {
		v, err := t.svc.User().Get(t.ctx, horus.UserByAlias("_me"))
		t.NoError(err)
		t.Equal(t.me.Actor.ID[:], v.Id)
	})
	t.Run("not found error if user does not exist", func() {
		_, err := t.svc.User().Get(t.ctx, horus.UserByAlias("not exist"))
		t.Error(err)

		s, ok := status.FromError(err)
		t.True(ok)
		t.Equal(codes.NotFound, s.Code())
	})
}
