package server_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"khepri.dev/horus"
)

type AccountTestSuite struct {
	Suite
}

func TestAccount(t *testing.T) {
	s := AccountTestSuite{
		Suite: NewSuiteWithSqliteStore(),
	}
	suite.Run(t, &s)
}

func (s *AccountTestSuite) TestCreate() {
	s.Run("silo owners can create account for their child user", func() {
		v, err := s.svc.Silo().Create(s.ctx, nil)
		s.NoError(err)

		child, err := s.svc.User().Create(s.ctx, nil)
		s.NoError(err)

		_, err = s.svc.Account().Create(s.ctx, &horus.CreateAccountRequest{
			Owner: &horus.GetUserRequest{Key: &horus.GetUserRequest_Id{Id: child.Id}},
			Silo:  &horus.GetSiloRequest{Key: &horus.GetSiloRequest_Id{Id: v.Id}},
		})
		s.NoError(err)
	})
}

func (s *AccountTestSuite) TestList() {
	s.Run("list accounts I owned", func() {
		v1, err := s.svc.Silo().Create(s.ctx, nil)
		s.NoError(err)

		v2, err := s.svc.Silo().Create(s.ctx, nil)
		s.NoError(err)

		res, err := s.svc.Account().List(s.ctx, &horus.ListAccountRequest{Key: &horus.ListAccountRequest_Mine{}})
		s.NoError(err)
		s.Len(res.Items, 2)

		s.Equal(v2.Id, res.Items[0].Silo.Id)
		s.Equal(v1.Id, res.Items[1].Silo.Id)
	})

	s.Run("list accounts of silo", func() {
		g, err := s.svc.Silo().Create(s.ctx, nil)
		s.NoError(err)

		u, err := s.svc.User().Create(s.ctx, nil)
		s.NoError(err)

		_, err = s.svc.Account().Create(s.ctx, &horus.CreateAccountRequest{
			Silo:  &horus.GetSiloRequest{Key: &horus.GetSiloRequest_Id{Id: g.Id}},
			Owner: &horus.GetUserRequest{Key: &horus.GetUserRequest_Id{Id: u.Id}},
		})
		s.NoError(err)

		res, err := s.svc.Account().List(s.ctx, &horus.ListAccountRequest{Key: &horus.ListAccountRequest_SiloId{
			SiloId: g.Id,
		}})
		s.NoError(err)
		s.Len(res.Items, 2)

		s.Equal(u.Id, res.Items[0].Owner.Id)
		s.Equal(s.me.Actor.ID[:], res.Items[1].Owner.Id)
	})
}
