package server_test

import (
	"fmt"
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
		SiloRole   role.Role
		TargetRole role.Role

		Fail Fail
		Code codes.Code
	}

	acts := []Act{}
	for _, silo_role := range role.Values() {
		for _, target_role := range role.Values() {
			acts = append(acts, Act{
				SiloRole:   silo_role,
				TargetRole: target_role,

				Fail: Fail(!fx.Or(
					silo_role == role.Owner,
					silo_role.HigherThan(target_role),
				)),
				Code: codes.PermissionDenied,
			})
		}
	}
	for _, act := range acts {
		t.Run(fmt.Sprintf("account with %s role %s be created by the silo %s for their child user", act.TargetRole, act.Fail, act.SiloRole), func() {
			actor := t.silo_admin
			ctx := frame.WithContext(t.ctx, actor)

			err := t.SetSiloRole(t.silo_owner, actor, act.SiloRole)
			t.NoError(err)

			child, err := t.svc.User().Create(ctx, nil)
			t.NoError(err)

			_, err = t.svc.Account().Create(ctx, &horus.CreateAccountRequest{
				Silo:  horus.SiloByIdV(t.silo.ID[:]),
				Owner: horus.UserByIdV(child.Id),
				Role:  fx.Addr(horus.RoleFrom(act.TargetRole)),
			})
			if act.Fail {
				t.ErrCode(err, codes.PermissionDenied)
			} else {
				t.NoError(err)
			}
		})
	}

	acts = []Act{}
	for _, silo_role := range role.Values() {
		acts = append(acts, Act{
			SiloRole: silo_role,

			Fail: true,
			Code: fx.Cond(
				fx.Or(
					silo_role == role.Owner,
					silo_role == role.Admin,
				),
				codes.FailedPrecondition,
				codes.PermissionDenied, // Member does not allowed to create an account.
			),
		})
	}
	for _, act := range acts {
		t.Run(fmt.Sprintf("account cannot be created if the user is not a child of the silo %s", act.SiloRole), func() {
			actor := t.silo_admin
			ctx := frame.WithContext(t.ctx, actor)

			err := t.SetSiloRole(t.silo_owner, actor, act.SiloRole)
			t.NoError(err)

			_, err = t.svc.Account().Create(ctx, &horus.CreateAccountRequest{
				Silo:  horus.SiloById(t.silo.ID),
				Owner: horus.UserById(t.other.Actor.ID),
			})
			t.ErrCode(err, act.Code)
		})
	}

	t.Run("account cannot be created if the owner already has an account in the silo", func() {
		child, err := t.svc.User().Create(t.CtxSiloOwner(), nil)
		t.NoError(err)

		_, err = t.svc.Account().Create(t.CtxSiloOwner(), &horus.CreateAccountRequest{
			Silo:  horus.SiloById(t.silo.ID),
			Owner: horus.UserByIdV(child.Id),
		})
		t.NoError(err)

		_, err = t.svc.Account().Create(t.CtxSiloOwner(), &horus.CreateAccountRequest{
			Silo:  horus.SiloById(t.silo.ID),
			Owner: horus.UserByIdV(child.Id),
		})
		t.ErrCode(err, codes.AlreadyExists)
	})
	t.Run("account is created with new user as an owner if the owner is not provided", func() {
		a, err := t.svc.Account().Create(t.CtxSiloOwner(), &horus.CreateAccountRequest{
			Silo: horus.SiloById(t.silo.ID),
		})
		t.NoError(err)

		u, err := t.svc.User().Get(t.CtxSiloOwner(), horus.UserByIdV(a.Owner.Id))
		t.NoError(err)
		t.Equal(t.silo_owner.Actor.ID[:], u.Parent.Id)
	})
}

func (t *AccountTestSuite) TestGet() {
	type Act struct {
		SiloRole role.Role

		OtherSilo bool

		Fail Fail
		Code codes.Code
	}

	acts := []Act{}
	for _, silo_role := range role.Values() {
		acts = append(acts, Act{
			SiloRole: silo_role,

			Fail: false,
		})
	}
	for _, silo_role := range role.Values() {
		acts = append(acts, Act{
			SiloRole: silo_role,

			OtherSilo: true,

			Fail: true,
			Code: codes.NotFound,
		})
	}

	for _, silo_role := range role.Values() {
		t.Run(fmt.Sprintf("account can be retrieved by its owner who is silo %s", silo_role), func() {
			actor := t.silo_admin
			err := t.SetSiloRole(t.silo_owner, actor, silo_role)
			t.NoError(err)

			ctx := frame.WithContext(t.ctx, actor)
			a, err := t.svc.Account().Get(ctx, horus.AccountById(actor.ActingAccount.ID))
			t.NoError(err)
			t.Equal(actor.Actor.ID[:], a.Owner.Id)
		})
	}
	for _, act := range acts {
		for _, target_silo_role := range role.Values() {
			title := fmt.Sprintf("account owned by silo %s %s be retrieved by silo %s", target_silo_role, act.Fail, act.SiloRole)
			if act.OtherSilo {
				title += " who is in another silo"
			}

			t.Run(title, func() {
				actor := t.silo_admin
				err := t.SetSiloRole(t.silo_owner, actor, act.SiloRole)
				t.NoError(err)

				var target *frame.Frame
				if act.OtherSilo {
					target = t.other_silo_admin
					err := t.SetSiloRole(t.other_silo_owner, target, target_silo_role)
					t.NoError(err)
				} else {
					target = t.silo_member
					err := t.SetSiloRole(t.silo_owner, target, target_silo_role)
					t.NoError(err)
				}

				ctx := frame.WithContext(t.ctx, actor)
				a, err := t.svc.Account().Get(ctx, horus.AccountById(target.ActingAccount.ID))
				if act.Fail {
					t.ErrCode(err, act.Code)
				} else {
					t.NoError(err)
					t.Equal(target.Actor.ID[:], a.Owner.Id)
				}
			})
		}
	}

	t.Run("account can be retrieved by a user who has an account in the same silo", func() {
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
	t.Run("account cannot be retrieved if the account does not exist", func() {
		_, err := t.svc.Account().Get(t.CtxSiloOwner(), horus.AccountByAliasInSilo("not exist", horus.SiloById(t.silo.ID)))
		t.ErrCode(err, codes.NotFound)
	})
	t.Run("account cannot be retrieved if the account is in another silo", func() {
		_, err := t.svc.Account().Get(t.CtxOther(), horus.AccountById(t.silo_owner.ActingAccount.ID))
		t.ErrCode(err, codes.NotFound)
	})
	t.Run("account of the actor can be retrieved using Me alias", func() {
		v, err := t.svc.Account().Get(t.CtxSiloMember(), horus.AccountByAliasInSilo(horus.Me, horus.SiloById(t.silo.ID)))
		t.NoError(err)
		t.Equal(t.silo_member.ActingAccount.ID[:], v.Id)
	})
	t.Run("account of the actor cannot be retrieved using Me alias if the account does not exist", func() {
		_, err := t.svc.Account().Get(t.CtxOther(), horus.AccountByAliasInSilo(horus.Me, horus.SiloById(t.silo.ID)))
		t.ErrCode(err, codes.NotFound)
	})
	t.Run("account of the actor cannot be retrieved using Me alias if the silo does not exist", func() {
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

		Fail Fail
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

	acts := []Act{}
	for _, actor_role := range role.Values() {
		acts = append(acts, Act{
			Actor: actor_role,
			Self:  true,
		})
	}
	for _, actor_role := range role.Values() {
		for _, target_role := range role.Values() {
			acts = append(acts, Act{
				Actor:  actor_role,
				Target: target_role,

				Fail: Fail(actor_role != role.Owner && !actor_role.HigherThan(target_role)),
			})
		}
	}

	// Update info.
	for _, act := range acts {
		title := fmt.Sprintf("account owned by silo %s %s be updated by ", act.Target, act.Fail)
		if act.Self {
			title += "account owner themselves"
		} else {
			title += fmt.Sprintf("silo %s", act.Actor)
		}

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

	acts = []Act{}
	for _, actor_role := range role.Values() {
		for _, role_to := range role.Values() {
			if actor_role == role_to {
				continue
			}

			acts = append(acts, Act{
				Actor: actor_role,
				To:    role_to,
				Self:  true,

				Fail: Fail(actor_role.LowerThan(role_to)),
			})
		}
	}
	for _, actor_role := range role.Values() {
		for _, target_role := range role.Values() {
			for _, role_to := range role.Values() {
				if target_role == role_to {
					continue
				}

				higher_than_actor := !actor_role.HigherThan(target_role) || !actor_role.HigherThan(role_to)
				acts = append(acts, Act{
					Actor:  actor_role,
					Target: target_role,
					To:     role_to,

					Fail: Fail(actor_role != role.Owner && higher_than_actor),
				})
			}
		}
	}
	for _, act := range acts {
		dir := fx.Cond(act.Target.HigherThan(act.To), "demoted", "promoted")
		title := fmt.Sprintf("silo %s %s %s to %s by ", act.Target, act.Fail, dir, act.To)
		if act.Self {
			title += "themselves"
		} else {
			title += fmt.Sprintf("silo %s", act.Actor)
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

	t.Run("account owned by silo owner cannot be demoted themselves if the actor is sole owner of the silo", func() {
		err := t.SetSiloRole(t.silo_owner, t.silo_owner, role.Admin)
		t.ErrCode(err, codes.FailedPrecondition)
	})
}

func (t *AccountTestSuite) TestDelete() {
	type Act struct {
		Actor  role.Role
		Target role.Role
		Self   bool

		Fail Fail
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

	acts := []Act{}
	for _, actor_role := range role.Values() {
		acts = append(acts, Act{
			Actor: actor_role,
			Self:  true,
		})
	}
	for _, actor_role := range role.Values() {
		for _, target_role := range role.Values() {
			acts = append(acts, Act{
				Actor:  actor_role,
				Target: target_role,

				Fail: Fail(actor_role != role.Owner && !actor_role.HigherThan(target_role)),
			})
		}
	}

	for _, act := range acts {
		title := fmt.Sprintf("account owned by %s %s deleted by ", act.Target, act.Fail)
		if act.Self {
			title += "themselves"
		} else {
			title += fmt.Sprintf("silo %s", act.Actor)
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

	t.Run("account owned by silo owner cannot be deleted themselves if the actor is sole owner in the silo", func() {
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
