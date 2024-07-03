package server_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
)

type SiloTestSuite struct {
	SuiteWithSilo
}

func TestSilo(t *testing.T) {
	s := SiloTestSuite{
		SuiteWithSilo: SuiteWithSilo{
			Suite: NewSuiteWithSqliteStore(),
		},
	}
	suite.Run(t, &s)
}

func (t *SiloTestSuite) TestCreate() {
	t.Run("owner account is created when the silo is created", func() {
		v, err := t.svc.Silo().Create(t.CtxMe(), &horus.CreateSiloRequest{
			Alias: fx.Addr("horus"),
			Name:  fx.Addr("Horus"),
		})
		t.NoError(err)
		t.Equal("horus", v.Alias)
		t.Equal("Horus", v.Name)

		v, err = t.svc.Silo().Get(t.CtxMe(), horus.SiloByIdV(v.Id))
		t.NoError(err)
		t.Equal("horus", v.Alias)
		t.Equal("Horus", v.Name)

		w, err := t.svc.Account().Get(t.CtxMe(), horus.AccountByAliasInSilo("founder", horus.SiloByIdV(v.Id)))
		t.NoError(err)
		t.Equal(horus.Role_ROLE_OWNER, w.Role)
	})
	t.Run("silo alias cannot be duplicated", func() {
		_, err := t.svc.Silo().Create(t.CtxMe(), &horus.CreateSiloRequest{
			Alias: &t.silo.Alias,
			Name:  fx.Addr("Horus"),
		})
		t.ErrCode(err, codes.AlreadyExists)
	})
}

func (t *SiloTestSuite) TestGet() {
	t.Run("as a silo owner", func() {
		v, err := t.svc.Silo().Get(t.CtxSiloOwner(), horus.SiloById(t.silo.ID))
		t.NoError(err)
		t.Equal(t.silo.Alias, v.Alias)
	})
	t.Run("as a silo admin", func() {
		v, err := t.svc.Silo().Get(t.CtxSiloAdmin(), horus.SiloById(t.silo.ID))
		t.NoError(err)
		t.Equal(t.silo.Alias, v.Alias)
	})
	t.Run("as a silo member", func() {
		v, err := t.svc.Silo().Get(t.CtxSiloMember(), horus.SiloById(t.silo.ID))
		t.NoError(err)
		t.Equal(t.silo.Alias, v.Alias)
	})
	t.Run("not found error if the user is an outsider", func() {
		_, err := t.svc.Silo().Get(t.CtxOther(), horus.SiloById(t.silo.ID))
		t.Equal(codes.NotFound, status.Code(err))
	})
	t.Run("not found error if the silo does not exist", func() {
		_, err := t.svc.Silo().Get(t.CtxOther(), horus.SiloByAlias("not exist"))
		t.Equal(codes.NotFound, status.Code(err))
	})
}

func (t *SiloTestSuite) TestUpdate() {
	t.Run("as a silo owner", func() {
		_, err := t.svc.Silo().Update(t.CtxSiloOwner(), &horus.UpdateSiloRequest{
			Key:  horus.SiloById(t.silo.ID),
			Name: fx.Addr("Khepri"),
		})
		t.NoError(err)

		v, err := t.svc.Silo().Get(t.CtxSiloOwner(), horus.SiloById(t.silo.ID))
		t.NoError(err)
		t.Equal("Khepri", v.Name)
	})
	t.Run("permission denied error if the user is not the owner of the silo", func() {
		_, err := t.svc.Silo().Update(t.CtxSiloAdmin(), &horus.UpdateSiloRequest{
			Key:  horus.SiloById(t.silo.ID),
			Name: fx.Addr("Khepri"),
		})
		t.Equal(codes.PermissionDenied, status.Code(err))

		_, err = t.svc.Silo().Update(t.CtxSiloMember(), &horus.UpdateSiloRequest{
			Key:  horus.SiloById(t.silo.ID),
			Name: fx.Addr("Khepri"),
		})
		t.Equal(codes.PermissionDenied, status.Code(err))
	})
	t.Run("not found error if the user is an outsider", func() {
		_, err := t.svc.Silo().Update(t.CtxOther(), &horus.UpdateSiloRequest{
			Key:  horus.SiloById(t.silo.ID),
			Name: fx.Addr("Khepri"),
		})
		t.Equal(codes.NotFound, status.Code(err))
	})
}

func (t *SiloTestSuite) TestDelete() {
	t.Run("not found error if the silo is deleted", func() {
		_, err := t.svc.Silo().Delete(t.CtxSiloOwner(), horus.SiloById(t.silo.ID))
		t.NoError(err)

		_, err = t.svc.Silo().Get(t.CtxSiloOwner(), horus.SiloById(t.silo.ID))
		t.ErrCode(err, codes.NotFound)
	})
	t.Run("all accounts in silo are deleted", func() {
		_, err := t.svc.Account().Get(t.CtxSiloOwner(), horus.AccountByAliasInSilo(t.silo_admin.ActingAccount.Alias, horus.SiloById(t.silo.ID)))
		t.NoError(err)

		_, err = t.svc.Silo().Delete(t.CtxSiloOwner(), horus.SiloById(t.silo.ID))
		t.NoError(err)

		_, err = t.svc.Account().Get(t.CtxSiloOwner(), horus.AccountByAliasInSilo(t.silo_admin.ActingAccount.Alias, horus.SiloById(t.silo.ID)))
		t.ErrCode(err, codes.NotFound)
	})
	t.Run("permission denied error if the silo member tries to delete the silo", func() {
		_, err := t.svc.Silo().Delete(t.CtxSiloAdmin(), horus.SiloById(t.silo.ID))
		t.ErrCode(err, codes.PermissionDenied)
	})
}
