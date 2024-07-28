package server_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
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
	t.Run("user is created with the parent as actor", func() {
		v, err := t.svc.User().Create(t.CtxMe(), nil)
		t.NoError(err)

		v, err = t.svc.User().Get(t.CtxMe(), horus.UserByIdV(v.Id))
		t.NoError(err)
		t.Equal(t.me.Actor.ID[:], v.GetParent().GetId())
	})
	t.Run("user cannot be created if the alias is used by another user", func() {
		_, err := t.svc.User().Create(t.CtxMe(), &horus.CreateUserRequest{
			Alias: &t.me.Actor.Alias,
		})
		t.ErrCode(err, codes.AlreadyExists)
		t.ErrorContains(err, "alias")
	})
	t.Run("user cannot be create with explicit parent", func() {
		_, err := t.svc.User().Create(t.CtxMe(), &horus.CreateUserRequest{
			Parent: horus.UserById(t.other.Actor.ID),
		})
		t.ErrCode(err, codes.InvalidArgument)
	})
}

func (t *UserTestSuite) TestGet() {
	t.Run("user owned by the actor can be retrieved with an empty input", func() {
		v, err := t.svc.User().Get(t.CtxMe(), nil)
		t.NoError(err)
		t.Equal(t.me.Actor.ID[:], v.Id)
	})
	t.Run("user owned by the actor can be retrieved using alias Me", func() {
		v, err := t.svc.User().Get(t.CtxMe(), horus.UserByAlias(horus.Me))
		t.NoError(err)
		t.Equal(t.me.Actor.ID[:], v.Id)
	})
	t.Run("user can be retrieved by its parent", func() {
		_, err := t.svc.User().Get(t.CtxMe(), horus.UserById(t.child.Actor.ID))
		t.NoError(err)
	})
	t.Run("user cannot be retrieved if the user is not child of the actor", func() {
		_, err := t.svc.User().Get(t.CtxMe(), horus.UserById(t.other.Actor.ID))
		t.ErrCode(err, codes.PermissionDenied)
	})
	t.Run("user cannot be retrieved if the user does not exist", func() {
		_, err := t.svc.User().Get(t.CtxMe(), horus.UserByAlias("not exist"))
		t.ErrCode(err, codes.NotFound)
	})
}

func (t *UserTestSuite) TestUpdate() {
	t.Run("user can be updated by its owner", func() {
		v, err := t.svc.User().Update(t.CtxMe(), &horus.UpdateUserRequest{
			Key:   horus.UserById(t.me.Actor.ID),
			Alias: fx.Addr("django"),
		})
		t.NoError(err)
		t.Equal("django", v.Alias)
	})
	t.Run("user can be updated by its parent", func() {
		v, err := t.svc.User().Update(t.CtxMe(), &horus.UpdateUserRequest{
			Key:   horus.UserById(t.child.Actor.ID),
			Alias: fx.Addr("django"),
		})
		t.NoError(err)
		t.Equal("django", v.Alias)
	})
	t.Run("user cannot be updated by other", func() {
		_, err := t.svc.User().Update(t.CtxMe(), &horus.UpdateUserRequest{
			Key:   horus.UserById(t.other.Actor.ID),
			Alias: fx.Addr("django"),
		})
		t.ErrCode(err, codes.PermissionDenied)
	})
	t.Run("user cannot be updated if the user does not exist", func() {
		_, err := t.svc.User().Update(t.CtxMe(), &horus.UpdateUserRequest{
			Key:   horus.UserByAlias("not exist"),
			Alias: fx.Addr("django"),
		})
		t.ErrCode(err, codes.NotFound)
	})
}

func (t *UserTestSuite) TestDelete() {
	t.Run("user cannot be deleted who is a root user", func() {
		_, err := t.svc.User().Delete(t.CtxMe(), horus.UserById(t.me.Actor.ID))
		t.ErrCode(err, codes.FailedPrecondition)
	})
	t.Run("user cannot be deleted if the user does not exist", func() {
		_, err := t.svc.User().Delete(t.CtxMe(), horus.UserByAlias("not exists"))
		t.ErrCode(err, codes.NotFound)
	})
	t.Run("user can be deleted by its parent", func() {
		_, err := t.svc.User().Delete(t.CtxChild(), horus.UserById(t.child.Actor.ID))
		t.NoError(err)

		_, err = t.svc.User().Get(t.CtxMe(), horus.UserById(t.child.Actor.ID))
		t.ErrCode(err, codes.NotFound)
	})
	t.Run("user can be deleted by its ancestor", func() {
		grand_child, err := t.svc.User().Create(t.CtxChild(), nil)
		t.NoError(err)

		_, err = t.svc.User().Delete(t.CtxMe(), horus.UserByIdV(grand_child.Id))
		t.NoError(err)

		_, err = t.svc.User().Get(t.CtxMe(), horus.UserByIdV(grand_child.Id))
		t.ErrCode(err, codes.NotFound)
	})
	t.Run("when a user is deleted, the parent of the deleted user becomes the parent of that user's child users", func() {
		grand_child, err := t.svc.User().Create(t.CtxChild(), nil)
		t.NoError(err)

		_, err = t.svc.User().Delete(t.CtxMe(), horus.UserById(t.child.Actor.ID))
		t.NoError(err)

		grand_child, err = t.svc.User().Get(t.CtxMe(), horus.UserByIdV(grand_child.Id))
		t.NoError(err)
		t.Equal(t.me.Actor.ID[:], grand_child.Parent.Id)
	})
	t.Run("user cannot be deleted by the user who is not its ancestor", func() {
		_, err := t.svc.User().Delete(t.CtxOther(), horus.UserById(t.child.Actor.ID))
		t.ErrCode(err, codes.PermissionDenied)
	})
	t.Run("user cannot be delete by their children", func() {
		_, err := t.svc.User().Delete(t.CtxChild(), horus.UserById(t.me.Actor.ID))
		t.ErrCode(err, codes.PermissionDenied)
	})
}
