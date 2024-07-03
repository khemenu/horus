package server_test

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/role"
	"khepri.dev/horus/server/frame"
)

type AccountTestSuite struct {
	SuiteWithSilo
}

func TestAccount(t *testing.T) {
	s := AccountTestSuite{
		SuiteWithSilo: SuiteWithSilo{
			Suite: NewSuiteWithSqliteStore(),
		},
	}
	suite.Run(t, &s)
}

func (t *AccountTestSuite) TestCreate() {
	type Act struct {
		Actor role.Role
		As    role.Role
		Fail  bool
	}

	for _, act := range []Act{
		{
			Actor: role.Owner,
			As:    role.Owner,
		},
		{
			Actor: role.Owner,
			As:    role.Admin,
		},
		{
			Actor: role.Owner,
			As:    role.Member,
		},
		{
			Actor: role.Admin,
			As:    role.Owner,
			Fail:  true,
		},
		{
			Actor: role.Admin,
			As:    role.Admin,
			Fail:  true,
		},
		{
			Actor: role.Admin,
			As:    role.Member,
		},
		{
			Actor: role.Member,
			As:    role.Owner,
			Fail:  true,
		},
		{
			Actor: role.Member,
			As:    role.Admin,
			Fail:  true,
		},
		{
			Actor: role.Member,
			As:    role.Member,
			Fail:  true,
		},
	} {
		title := "silo " + strings.ToLower(string(act.Actor)) + " "
		title += fx.Cond(act.Fail, "cannot", "can")
		title += " create an account with role " + strings.ToLower(string(act.As)) + " for their child"

		t.Run(title, func() {
			actor := t.silo_admin
			err := t.SetSiloRole(t.silo_owner, actor, act.Actor)
			t.NoError(err)

			ctx := frame.WithContext(t.ctx, actor)
			child, err := t.svc.User().Create(ctx, nil)
			t.NoError(err)

			_, err = t.svc.Account().Create(ctx, &horus.CreateAccountRequest{
				Silo:  horus.SiloByIdV(t.silo.ID[:]),
				Owner: horus.UserByIdV(child.Id),
				Role:  horus.RoleFrom(act.As),
			})
			if act.Fail {
				t.ErrCode(err, codes.PermissionDenied)
			} else {
				t.NoError(err)
			}
		})
	}

	t.Run("account cannot be create if the user is not a child of the actor", func() {
		_, err := t.svc.Account().Create(t.CtxSiloOwner(), &horus.CreateAccountRequest{
			Silo:  horus.SiloByIdV(t.silo.ID[:]),
			Owner: horus.UserByIdV(t.other.Actor.ID[:]),
		})
		t.ErrCode(err, codes.FailedPrecondition)
	})
	t.Run("silo admin cannot create an account if the user is not a child", func() {
		_, err := t.svc.Account().Create(t.CtxSiloAdmin(), &horus.CreateAccountRequest{
			Silo:  horus.SiloByIdV(t.silo.ID[:]),
			Owner: horus.UserByIdV(t.other.Actor.ID[:]),
		})
		t.ErrCode(err, codes.FailedPrecondition)
	})
	t.Run("user cannot have multiple accounts in the silo", func() {
		child, err := t.svc.User().Create(t.CtxSiloOwner(), nil)
		t.NoError(err)

		_, err = t.svc.Account().Create(t.CtxSiloOwner(), &horus.CreateAccountRequest{
			Silo:  horus.SiloByIdV(t.silo.ID[:]),
			Owner: horus.UserByIdV(child.Id),
		})
		t.NoError(err)

		_, err = t.svc.Account().Create(t.CtxSiloOwner(), &horus.CreateAccountRequest{
			Silo:  horus.SiloByIdV(t.silo.ID[:]),
			Owner: horus.UserByIdV(child.Id),
		})
		t.ErrCode(err, codes.AlreadyExists)
	})
}

func (t *AccountTestSuite) TestGet() {
	t.Run("account can be get by a user who has an account in the same silo", func() {
		fs := []*frame.Frame{
			t.silo_owner,
			t.silo_admin,
			t.silo_member,
		}
		for _, f := range fs {
			ctx := frame.WithContext(t.ctx, f)
			for _, f_ := range fs {
				_, err := t.svc.Account().Get(ctx, horus.AccountById(f_.ActingAccount.ID))
				t.NoError(err)
			}
		}
	})
	t.Run("account cannot be get if the account does not exist", func() {
		_, err := t.svc.Account().Get(t.CtxSiloOwner(), horus.AccountByAliasInSilo("not exist", horus.SiloById(t.silo.ID)))
		t.ErrCode(err, codes.NotFound)
	})
	t.Run("account cannot be get if the account is in another silo", func() {
		_, err := t.svc.Account().Get(t.CtxOther(), horus.AccountById(t.silo_owner.ActingAccount.ID))
		t.ErrCode(err, codes.NotFound)
	})
	t.Run("my account cannot be get if the account does not exist", func() {
		_, err := t.svc.Account().Get(t.CtxOther(), horus.AccountByAliasInSilo(horus.Me, horus.SiloById(t.silo.ID)))
		t.ErrCode(err, codes.NotFound)
	})
	t.Run("my account cannot be get if the silo does not exist", func() {
		_, err := t.svc.Account().Get(t.CtxOther(), horus.AccountByAliasInSilo(horus.Me, horus.SiloById(uuid.Nil)))
		t.ErrCode(err, codes.NotFound)
	})
}

func (t *AccountTestSuite) TestUpdate() {
	update := func(actor *frame.Frame, target *frame.Frame) error {
		ctx := frame.WithContext(t.ctx, actor)
		_, err := t.svc.Account().Update(ctx, &horus.UpdateAccountRequest{
			Key:  horus.AccountById(target.ActingAccount.ID),
			Name: fx.Addr("Django"),
		})
		return err
	}

	type Act struct {
		Actor  role.Role
		Target role.Role
		To     role.Role
		Self   bool
		Fail   bool
	}
	prepare := func(act Act) (*frame.Frame, *frame.Frame) {
		actor := t.silo_admin
		err := t.SetSiloRole(t.silo_owner, actor, act.Actor)
		t.NoError(err)

		target := t.silo_member
		if act.Self {
			target = actor
		} else {
			err := t.SetSiloRole(t.silo_owner, target, act.Target)
			t.NoError(err)
		}

		return actor, target
	}

	// Update info.
	for _, act := range []Act{
		// As owner.
		{
			Actor: role.Owner,
			Self:  true,
		},
		{
			Actor:  role.Owner,
			Target: role.Owner,
		},
		{
			Actor:  role.Owner,
			Target: role.Admin,
		},
		{
			Actor:  role.Owner,
			Target: role.Member,
		},
		// As admin.
		{
			Actor: role.Admin,
			Self:  true,
		},
		{
			Actor:  role.Admin,
			Target: role.Owner,
			Fail:   true,
		},
		{
			Actor:  role.Admin,
			Target: role.Admin,
			Fail:   true,
		},
		{
			Actor:  role.Admin,
			Target: role.Member,
		},
		// As member.
		{
			Actor: role.Member,
			Self:  true,
		},
		{
			Actor:  role.Member,
			Target: role.Owner,
			Fail:   true,
		},
		{
			Actor:  role.Member,
			Target: role.Admin,
			Fail:   true,
		},
		{
			Actor:  role.Member,
			Target: role.Member,
			Fail:   true,
		},
	} {
		title := "silo " + strings.ToLower(string(act.Actor)) + " "
		title += fx.Cond(act.Fail, "cannot", "can")
		title += " update "
		title += fx.Cond(act.Self, "itself", "silo "+strings.ToLower(string(act.Target)))

		t.Run(title, func() {
			actor, target := prepare(act)

			err := update(actor, target)
			if act.Fail {
				t.ErrCode(err, codes.PermissionDenied)
			} else {
				t.NoError(err)
			}
		})
	}

	// Update role.
	for _, act := range []Act{
		// As a silo owner.
		{
			Actor: role.Owner,
			To:    role.Admin,
			Self:  true,
		},
		{
			Actor: role.Owner,
			To:    role.Member,
			Self:  true,
		},
		{
			Actor:  role.Owner,
			Target: role.Owner,
			To:     role.Admin,
		},
		{
			Actor:  role.Owner,
			Target: role.Owner,
			To:     role.Member,
		},
		{
			Actor:  role.Owner,
			Target: role.Admin,
			To:     role.Owner,
		},
		{
			Actor:  role.Owner,
			Target: role.Admin,
			To:     role.Member,
		},
		{
			Actor:  role.Owner,
			Target: role.Member,
			To:     role.Owner,
		},
		{
			Actor:  role.Owner,
			Target: role.Member,
			To:     role.Admin,
		},
		// As a silo admin.
		{
			Actor: role.Admin,
			To:    role.Owner,
			Self:  true,
			Fail:  true,
		},
		{
			Actor: role.Admin,
			To:    role.Member,
			Self:  true,
		},
		{
			Actor:  role.Admin,
			Target: role.Owner,
			To:     role.Admin,
			Fail:   true,
		},
		{
			Actor:  role.Admin,
			Target: role.Owner,
			To:     role.Member,
			Fail:   true,
		},
		{
			Actor:  role.Admin,
			Target: role.Admin,
			To:     role.Owner,
			Fail:   true,
		},
		{
			Actor:  role.Admin,
			Target: role.Admin,
			To:     role.Member,
			Fail:   true,
		},
		{
			Actor:  role.Admin,
			Target: role.Member,
			To:     role.Owner,
			Fail:   true,
		},
		{
			Actor:  role.Admin,
			Target: role.Member,
			To:     role.Admin,
			Fail:   true,
		},
		// As a silo member.
		{
			Actor: role.Member,
			To:    role.Owner,
			Self:  true,
			Fail:  true,
		},
		{
			Actor: role.Member,
			To:    role.Admin,
			Self:  true,
			Fail:  true,
		},
		{
			Actor:  role.Member,
			Target: role.Owner,
			To:     role.Owner,
			Fail:   true,
		},
		{
			Actor:  role.Member,
			Target: role.Owner,
			To:     role.Admin,
			Fail:   true,
		},
		{
			Actor:  role.Member,
			Target: role.Admin,
			To:     role.Owner,
			Fail:   true,
		},
		{
			Actor:  role.Member,
			Target: role.Admin,
			To:     role.Member,
			Fail:   true,
		},
		{
			Actor:  role.Member,
			Target: role.Member,
			To:     role.Owner,
			Fail:   true,
		},
		{
			Actor:  role.Member,
			Target: role.Member,
			To:     role.Admin,
			Fail:   true,
		},
	} {
		title := "silo " + strings.ToLower(string(act.Actor)) + " "
		if act.Fail {
			title += "cannot"
		} else {
			title += "can"
		}
		title += " "
		if act.Actor.HigherThan(act.Target) {
			title += "promote"
		} else {
			title += "demote"
		}
		title += " "
		if act.Self {
			title += "itself"
		} else {
			title += "silo " + strings.ToLower(string(act.Target))
		}

		t.Run(title, func() {
			actor, target := prepare(act)

			err := t.SetSiloRole(actor, target, act.To)
			if act.Fail {
				t.ErrCode(err, codes.PermissionDenied)
			} else {
				t.NoError(err)
			}
		})
	}

	t.Run("silo owner cannot demote itself if it is sole owner", func() {
		err := t.SetSiloRole(t.silo_owner, t.silo_owner, role.Admin)
		t.ErrCode(err, codes.FailedPrecondition)
	})
}

func (t *AccountTestSuite) TestDelete() {
	type Act struct {
		Actor  role.Role
		Target role.Role
		Self   bool
		Fail   bool
	}
	prepare := func(act Act) (*frame.Frame, *frame.Frame) {
		actor := t.silo_admin
		err := t.SetSiloRole(t.silo_owner, actor, act.Actor)
		t.NoError(err)

		target := t.silo_member
		if act.Self {
			target = actor
		} else {
			err := t.SetSiloRole(t.silo_owner, target, act.Target)
			t.NoError(err)
		}

		return actor, target
	}

	for _, act := range []Act{
		// As a silo owner.
		{
			Actor: role.Owner,
			Self:  true,
		},
		{
			Actor:  role.Owner,
			Target: role.Owner,
		},
		{
			Actor:  role.Owner,
			Target: role.Admin,
		},
		{
			Actor:  role.Owner,
			Target: role.Member,
		},
		// As a silo admin.
		{
			Actor: role.Admin,
			Self:  true,
		},
		{
			Actor:  role.Admin,
			Target: role.Owner,
			Fail:   true,
		},
		{
			Actor:  role.Admin,
			Target: role.Member,
		},
		// As a silo member.
		{
			Actor: role.Member,
			Self:  true,
		},
		{
			Actor:  role.Member,
			Target: role.Owner,
			Fail:   true,
		},
		{
			Actor:  role.Member,
			Target: role.Admin,
			Fail:   true,
		},
		{
			Actor:  role.Member,
			Target: role.Member,
			Fail:   true,
		},
	} {
		title := "silo " + strings.ToLower(string(act.Actor)) + " "
		if act.Fail {
			title += "cannot"
		} else {
			title += "can"
		}
		title += " delete "
		if act.Self {
			title += "itself"
		} else {
			title += "silo " + strings.ToLower(string(act.Target))
		}

		t.Run(title, func() {
			actor, target := prepare(act)

			ctx := frame.WithContext(t.ctx, actor)
			_, err := t.svc.Account().Delete(ctx, horus.AccountById(target.ActingAccount.ID))
			if act.Fail {
				t.ErrCode(err, codes.PermissionDenied)
			} else {
				t.NoError(err)
			}
		})
	}

	t.Run("silo owner cannot delete itself if it is sole owner", func() {
		_, err := t.svc.Account().Delete(t.CtxSiloOwner(), horus.AccountById(t.silo_owner.ActingAccount.ID))
		t.ErrCode(err, codes.FailedPrecondition)
	})
	t.Run("account cannot be deleted if it does not exist", func() {
		_, err := t.svc.Account().Delete(t.CtxSiloOwner(), horus.AccountById(uuid.Nil))
		t.ErrCode(err, codes.NotFound)
	})
	t.Run("account cannot be deleted if it is in another silo", func() {
		fs := []*frame.Frame{
			t.silo_owner,
			t.silo_admin,
			t.silo_member,
		}
		for _, f := range fs {
			_, err := t.svc.Account().Delete(t.CtxOther(), horus.AccountById(f.ActingAccount.ID))
			t.ErrCode(err, codes.NotFound)
		}
	})
}

func (t *AccountTestSuite) TestList() {
	t.Run("list accounts I owned", func() {
		u, err := t.svc.User().Create(t.CtxMe(), nil)
		t.NoError(err)

		s, err := t.svc.Silo().Create(t.CtxMe(), nil)
		t.NoError(err)

		// Account owned by `u` where the role is a silo owner.
		v1, err := t.svc.Account().Create(t.CtxMe(), &horus.CreateAccountRequest{
			Owner: horus.UserByIdV(u.Id),
			Silo:  horus.SiloByIdV(s.Id),
		})
		t.NoError(err)

		ctx_u := t.AsUser(u.Id)
		s_u, err := t.svc.Silo().Create(ctx_u, nil)
		t.NoError(err)

		// Account owned by `u` where the role is a silo member.
		v2, err := t.svc.Account().Get(ctx_u, horus.AccountByAliasInSilo(horus.Me, horus.SiloByIdV(s_u.Id)))
		t.NoError(err)

		res, err := t.svc.Account().List(ctx_u, &horus.ListAccountRequest{Key: &horus.ListAccountRequest_Mine{}})
		t.NoError(err)
		t.Len(res.Items, 2)

		t.Equal(v2.Id, res.Items[0].Id)
		t.Equal(v1.Id, res.Items[1].Id)
	})

	t.Run("list accounts of silo", func() {
		s, err := t.svc.Silo().Create(t.CtxMe(), nil)
		t.NoError(err)

		v1, err := t.svc.Account().Get(t.CtxMe(), horus.AccountByAliasInSilo(horus.Me, horus.SiloByIdV(s.Id)))
		t.NoError(err)

		child, err := t.svc.User().Create(t.CtxMe(), nil)
		t.NoError(err)

		v2, err := t.svc.Account().Create(t.CtxMe(), &horus.CreateAccountRequest{
			Owner: horus.UserByIdV(child.Id),
			Silo:  horus.SiloByIdV(s.Id),
		})
		t.NoError(err)

		res, err := t.svc.Account().List(t.CtxMe(), &horus.ListAccountRequest{Key: &horus.ListAccountRequest_Silo{
			Silo: horus.SiloByIdV(s.Id),
		}})
		t.NoError(err)
		t.Len(res.Items, 2)

		t.Equal(v2.Id, res.Items[0].Id)
		t.Equal(v1.Id, res.Items[1].Id)
	})
}
