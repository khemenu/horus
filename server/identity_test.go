package server_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/server/frame"
)

type IdentityTestSuite struct {
	Suite
}

func TestIdentity(t *testing.T) {
	s := IdentityTestSuite{
		Suite: NewSuiteWithSqliteStore(),
	}
	suite.Run(t, &s)
}

func (t *IdentityTestSuite) TestCreate() {
	t.Run("identity is created with the actor as an owner", func() {
		v, err := t.svc.Identity().Create(t.CtxMe(), &horus.CreateIdentityRequest{
			Kind:  "DRiVER_LiCENSE",
			Value: "Patrick Star",
		})
		t.NoError(err)

		v, err = t.svc.Identity().Get(t.CtxMe(), horus.IdentityByIdV(v.Id))
		t.NoError(err)
		t.Equal(t.me.Actor.ID[:], v.Owner.Id)
	})
	t.Run("identity can be created with the child as an owner", func() {
		v, err := t.svc.Identity().Create(t.CtxMe(), &horus.CreateIdentityRequest{
			Owner: horus.UserById(t.child.Actor.ID),
			Kind:  "DRiVER_LiCENSE",
			Value: "Patrick Star",
		})
		t.NoError(err)

		v, err = t.svc.Identity().Get(t.CtxChild(), horus.IdentityByIdV(v.Id))
		t.NoError(err)
		t.Equal(t.child.Actor.ID[:], v.Owner.Id)
	})
	t.Run("identity cannot be created by other user", func() {
		_, err := t.svc.Identity().Create(t.CtxOther(), &horus.CreateIdentityRequest{
			Owner: horus.UserById(t.child.Actor.ID),
			Kind:  "DRiVER_LiCENSE",
			Value: "Patrick Star",
		})
		t.ErrCode(err, codes.PermissionDenied)
	})
	t.Run("identity cannot be created if kind-value pair is duplicated for the owner", func() {
		_, err := t.svc.Identity().Create(t.CtxMe(), t.CreateReq_A())
		t.NoError(err)

		_, err = t.svc.Identity().Create(t.CtxMe(), t.CreateReq_B())
		t.NoError(err)

		_, err = t.svc.Identity().Create(t.CtxMe(), t.CreateReq_B())
		t.ErrCode(err, codes.AlreadyExists)
	})
}

func (t *IdentityTestSuite) TestGet() {
	for _, act := range t.UserPermissionActs() {
		t.Run(fmt.Sprintf("identity %s be retrieved by %s", act.Fail, act.Name), func() {
			actor := frame.WithContext(t.ctx, *act.Actor)
			target := frame.WithContext(t.ctx, *act.Target)

			v, err := t.svc.Identity().Create(target, t.CreateReq_A())
			t.NoError(err)

			_, err = t.svc.Identity().Get(actor, horus.IdentityByIdV(v.Id))
			if act.Fail {
				t.ErrCode(err, codes.PermissionDenied)
			} else {
				t.NoError(err)
			}
		})
	}

	t.Run("identity cannot be retrieved if it does not exist", func() {
		_, err := t.svc.Identity().Get(t.CtxOther(), horus.IdentityById(uuid.Nil))
		t.ErrCode(err, codes.NotFound)
	})
}

func (t *IdentityTestSuite) TestList() {
	for _, act := range t.UserPermissionActs() {
		t.Run(fmt.Sprintf("identity %s be listed by %s", act.Fail, act.Name), func() {
			actor := frame.WithContext(t.ctx, *act.Actor)

			_, err := t.svc.Identity().List(actor, &horus.ListIdentityRequest{
				Key: &horus.ListIdentityRequest_Owner{
					Owner: horus.UserById((*act.Target).Actor.ID),
				},
			})
			if act.Fail {
				t.ErrCode(err, codes.PermissionDenied)
			} else {
				t.NoError(err)
			}
		})
	}

	t.Run("identity cannot be listed if the owner does not exist", func() {
		_, err := t.svc.Identity().List(t.CtxOther(), &horus.ListIdentityRequest{Key: &horus.ListIdentityRequest_Owner{
			Owner: horus.UserById(uuid.Nil),
		}})
		t.ErrCode(err, codes.NotFound)
	})
}

func (t *IdentityTestSuite) TestUpdate() {
	for _, act := range t.UserPermissionActs() {
		t.Run(fmt.Sprintf("identity %s be updated by %s", act.Fail, act.Name), func() {
			actor := frame.WithContext(t.ctx, *act.Actor)
			target := frame.WithContext(t.ctx, *act.Target)

			v, err := t.svc.Identity().Create(target, t.CreateReq_A())
			t.NoError(err)
			t.Empty(v.Description)

			_, err = t.svc.Identity().Update(actor, &horus.UpdateIdentityRequest{
				Key:         horus.IdentityByIdV(v.Id),
				Description: fx.Addr("HAIR: PINK"),
			})
			if act.Fail {
				t.ErrCode(err, codes.PermissionDenied)
				return
			}
			t.NoError(err)

			v, err = t.svc.Identity().Get(target, horus.IdentityByIdV(v.Id))
			t.NoError(err)
			t.Equal("HAIR: PINK", v.Description)
		})
	}

	t.Run("identity cannot be listed if it does not exist", func() {
		_, err := t.svc.Identity().Update(t.CtxMe(), &horus.UpdateIdentityRequest{
			Key:         horus.IdentityById(uuid.Nil),
			Description: fx.Addr("HAIR: PINK"),
		})
		t.ErrCode(err, codes.NotFound)
	})
}

func (t *IdentityTestSuite) TestDelete() {
	for _, act := range t.UserPermissionActs() {
		t.Run(fmt.Sprintf("identity %s be deleted by %s", act.Fail, act.Name), func() {
			actor := frame.WithContext(t.ctx, *act.Actor)
			target := frame.WithContext(t.ctx, *act.Target)

			v, err := t.svc.Identity().Create(target, t.CreateReq_A())
			t.NoError(err)

			_, err = t.svc.Identity().Delete(actor, horus.IdentityByIdV(v.Id))
			if act.Fail {
				t.ErrCode(err, codes.PermissionDenied)
				return
			}
			t.NoError(err)

			_, err = t.svc.Identity().Get(t.CtxChild(), horus.IdentityByIdV(v.Id))
			t.ErrCode(err, codes.NotFound)
		})
	}

	t.Run("identity cannot be deleted if it does not exist", func() {
		_, err := t.svc.Identity().Delete(t.CtxMe(), horus.IdentityById(uuid.Nil))
		t.ErrCode(err, codes.NotFound)
	})
}

func (t *IdentityTestSuite) CreateReq_A() *horus.CreateIdentityRequest {
	return &horus.CreateIdentityRequest{
		Kind:  "DRiVER_LiCENSE",
		Value: "Patrick Star",
	}
}

func (t *IdentityTestSuite) CreateReq_B() *horus.CreateIdentityRequest {
	return &horus.CreateIdentityRequest{
		Kind:  "FiSHiNG_LiCENSE",
		Value: "Patrick Star",
	}
}
