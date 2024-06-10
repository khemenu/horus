package server_test

import (
	"slices"
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
			Owner: &horus.User{Id: child.Id},
			Silo:  &horus.Silo{Id: v.Id},
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

		res, err := s.svc.Account().List(s.ctx, nil)
		s.NoError(err)
		s.Len(res.Items, 2)

		slices.SortFunc(res.Items, func(a, b *horus.Account) int {
			return a.DateCreated.AsTime().Compare(b.DateCreated.AsTime())
		})
		s.Equal(v1.Id, res.Items[0].Silo.Id)
		s.Equal(v2.Id, res.Items[1].Silo.Id)
	})
}
