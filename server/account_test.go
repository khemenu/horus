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

		_, err = s.svc.Account().Create(s.ctx, &horus.CreateAccountRequest{Account: &horus.Account{
			Owner: &horus.User{Id: child.Id},
			Silo:  &horus.Silo{Id: v.Id},
		}})
		s.NoError(err)
	})
}
