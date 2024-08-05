package server_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/server/frame"
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
	t.Run("user is created with the actor as a parent", func() {
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
	t.Run("user cannot be created with explicit parent", func() {
		_, err := t.svc.User().Create(t.CtxMe(), &horus.CreateUserRequest{
			Parent: horus.UserById(t.other.Actor.ID),
		})
		t.ErrCode(err, codes.InvalidArgument)
	})
}

func (t *UserTestSuite) TestGet() {
	for _, act := range t.UserPermissionActs() {
		t.Run(fmt.Sprintf("user %s be retrieved by %s", act.Fail, act.Name), func() {
			actor := frame.WithContext(t.ctx, *act.Actor)
			target := (*act.Target).Actor

			v, err := t.svc.User().Get(actor, horus.UserById(target.ID))
			if act.Fail {
				t.ErrCode(err, codes.PermissionDenied)
				return
			}
			t.NoError(err)
			t.Equal(target.ID[:], v.Id)
		})
	}

	t.Run("user owned by the actor can be retrieved with an empty input", func() {
		v, err := t.svc.User().Get(t.CtxMe(), nil)
		t.NoError(err)
		t.Equal(t.me.Actor.ID[:], v.Id)
	})
	t.Run("user owned by the actor can be retrieved using an alias Me", func() {
		v, err := t.svc.User().Get(t.CtxMe(), horus.UserByAlias(horus.Me))
		t.NoError(err)
		t.Equal(t.me.Actor.ID[:], v.Id)
	})
	t.Run("user cannot be retrieved if the user does not exist", func() {
		_, err := t.svc.User().Get(t.CtxMe(), horus.UserByAlias("not exist"))
		t.ErrCode(err, codes.NotFound)
	})
}

func (t *UserTestSuite) TestUpdate() {
	for _, act := range t.UserPermissionActs() {
		t.Run(fmt.Sprintf("user %s be updated by %s", act.Fail, act.Name), func() {
			actor := frame.WithContext(t.ctx, *act.Actor)
			target := (*act.Target).Actor

			v, err := t.svc.User().Update(actor, &horus.UpdateUserRequest{
				Key:   horus.UserById(target.ID),
				Alias: fx.Addr("django"),
			})
			if act.Fail {
				t.ErrCode(err, codes.PermissionDenied)
				return
			}
			t.NoError(err)
			t.Equal(target.ID[:], v.Id)
		})
	}

	t.Run("user cannot be updated if the user does not exist", func() {
		_, err := t.svc.User().Update(t.CtxMe(), &horus.UpdateUserRequest{
			Key:   horus.UserByAlias("not exist"),
			Alias: fx.Addr("django"),
		})
		t.ErrCode(err, codes.NotFound)
	})
}

func (t *UserTestSuite) TestDelete() {
	for _, act := range t.UserPermissionActs() {
		t.Run(fmt.Sprintf("user %s be deleted by %s", act.Fail, act.Name), func() {
			actor := frame.WithContext(t.ctx, *act.Actor)
			target := (*act.Target).Actor

			_, err := t.svc.User().Delete(actor, horus.UserById(target.ID))
			if act.Fail {
				t.ErrCode(err, codes.PermissionDenied)
				return
			}
			t.NoError(err)

			_, err = t.svc.User().Get(actor, horus.UserById(target.ID))
			t.ErrCode(err, codes.NotFound)
		})
	}

	t.Run("user cannot be deleted if the user is root", func() {
		_, err := t.svc.User().Delete(t.CtxRoot(), horus.UserById(t.root.Actor.ID))
		t.ErrCode(err, codes.FailedPrecondition)
	})
	t.Run("user cannot be deleted if the user does not exist", func() {
		_, err := t.svc.User().Delete(t.CtxMe(), horus.UserByAlias("not exists"))
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
}

type UserPermissionAct struct {
	Name   string
	Actor  **frame.Frame
	Target **frame.Frame
	Fail   Fail
}

func (t *Suite) UserPermissionActs() []UserPermissionAct {
	return []UserPermissionAct{
		{
			Name:   "its owner's ancestor",
			Actor:  &t.parent,
			Target: &t.child,
		},
		{
			Name:   "its owner's parent",
			Actor:  &t.parent,
			Target: &t.me,
		},
		{
			Name:   "its owner",
			Actor:  &t.me,
			Target: &t.me,
		},
		{
			Name:   "its owner's child",
			Actor:  &t.child,
			Target: &t.me,
			Fail:   true,
		},
		{
			Name:   "its owner's descendant",
			Actor:  &t.child,
			Target: &t.parent,
			Fail:   true,
		},
		{
			Name:   "other user",
			Actor:  &t.other,
			Target: &t.me,
			Fail:   true,
		},
	}
}
