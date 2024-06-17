package server_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
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

func (t *UserTestSuite) TestCreate() {
	t.Run("acting user is set as the new user's parent", func() {
		v, err := t.svc.User().Create(t.ctx, nil)
		t.NoError(err)

		v, err = t.svc.User().Get(t.ctx, horus.UserByIdV(v.Id))
		t.NoError(err)
		t.Equal(t.me.Actor.ID[:], v.GetParent().GetId())
	})
	t.Run("cannot have same alias with other", func() {
		_, err := t.svc.User().Create(t.ctx, &horus.CreateUserRequest{
			Alias: &t.me.Actor.Alias,
		})
		t.ErrCode(err, codes.AlreadyExists)
		t.ErrorContains(err, "alias")
	})
}

func (t *UserTestSuite) TestGet() {
	t.Run("alias _me returns myself", func() {
		v, err := t.svc.User().Get(t.ctx, horus.UserByAlias("_me"))
		t.NoError(err)
		t.Equal(t.me.Actor.ID[:], v.Id)
	})
	t.Run("not found error if the user does not exist", func() {
		_, err := t.svc.User().Get(t.ctx, horus.UserByAlias("not exist"))
		t.ErrCode(err, codes.NotFound)
	})
	t.Run("permission denied error if try to get another existing user info", func() {
		_, err := t.svc.User().Get(t.ctx, horus.UserById(t.other.Actor.ID))
		t.ErrCode(err, codes.PermissionDenied)
	})
	t.Run("user can get their child user info", func() {
		v, err := t.svc.User().Create(t.ctx, nil)
		t.NoError(err)

		_, err = t.svc.User().Get(t.ctx, horus.UserByIdV(v.Id))
		t.NoError(err)
	})
}

func (t *UserTestSuite) TestUpdate() {
	t.Run("user can update info of their child", func() {
		v, err := t.svc.User().Create(t.ctx, nil)
		t.NoError(err)

		alias := "django"
		w, err := t.svc.User().Update(t.ctx, &horus.UpdateUserRequest{
			Key:   horus.UserByIdV(v.Id),
			Alias: &alias,
		})
		t.NoError(err)
		t.Equal(alias, w.Alias)

		w, err = t.svc.User().Get(t.ctx, horus.UserByAlias(alias))
		t.NoError(err)
		t.Equal(v.Id, w.Id)

		_, err = t.svc.User().Get(t.ctx, horus.UserByAlias(v.Alias))
		t.ErrCode(err, codes.NotFound)
	})
}
