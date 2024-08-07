package server_test

import (
	"encoding/base64"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"khepri.dev/horus"
	"khepri.dev/horus/conf"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/server/frame"
)

type AuthTestSuite struct {
	Suite
}

func TestAuth(t *testing.T) {
	s := AuthTestSuite{
		Suite: NewSuiteWithSqliteStore(),
	}
	suite.Run(t, &s)
}

func (t *AuthTestSuite) TestBasicSignIn() {
	pw := "bigboobz"

	t.Run("user can sign in using their username and password", func() {
		_, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Value: pw,
			Type:  horus.TokenTypePassword,
		})
		t.NoError(err)

		v, err := t.svc.Auth().BasicSignIn(t.ctx, &horus.BasicSignInRequest{
			Username: t.me.Actor.Alias,
			Password: pw,
		})
		t.NoError(err)
		t.Equal(horus.TokenTypeAccess, v.Token.Type)
	})
	t.Run("user cannot sign in using an old password", func() {
		_, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Value: pw,
			Type:  horus.TokenTypePassword,
		})
		t.NoError(err)

		v, err := t.svc.Auth().BasicSignIn(t.ctx, &horus.BasicSignInRequest{
			Username: t.me.Actor.Alias,
			Password: pw,
		})
		t.NoError(err)
		t.Equal(horus.TokenTypeAccess, v.Token.Type)

		_, err = t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Value: "pam",
			Type:  horus.TokenTypePassword,
		})
		t.NoError(err)

		_, err = t.svc.Auth().BasicSignIn(t.ctx, &horus.BasicSignInRequest{
			Username: t.me.Actor.Alias,
			Password: pw,
		})
		t.ErrCode(err, codes.Unauthenticated)
	})
	t.Run("user cannot sign in as another user with their password", func() {
		_, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Value: pw,
			Type:  horus.TokenTypePassword,
		})
		t.NoError(err)

		_, err = t.svc.Auth().BasicSignIn(t.ctx, &horus.BasicSignInRequest{
			Username: t.other.Actor.Alias,
			Password: pw,
		})
		t.ErrCode(err, codes.Unauthenticated)
	})
	t.Run("user cannot sign in if it is locked out", func() {
		c := conf.ConfSignInLockout{
			Enabled:      true,
			Count:        1,
			LockedPeriod: 60,
		}
		err := conf.MarshalInto(t.ctx, &c, t.bare.Conf())
		t.NoError(err)

		_, err = t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Value: pw,
			Type:  horus.TokenTypePassword,
		})
		t.NoError(err)

		_, err = t.svc.Auth().BasicSignIn(t.ctx, &horus.BasicSignInRequest{
			Username: t.me.Actor.Alias,
			Password: "0000 0000",
		})
		t.ErrCode(err, codes.Unauthenticated)

		_, err = t.svc.Auth().BasicSignIn(t.ctx, &horus.BasicSignInRequest{
			Username: t.me.Actor.Alias,
			Password: "0000 0001",
		})
		t.ErrCode(err, codes.FailedPrecondition)
	})
	t.Run("user can try sign in again after locked period of time", func() {
		c := conf.ConfSignInLockout{
			Enabled:      true,
			Count:        1,
			LockedPeriod: 0,
		}
		err := conf.MarshalInto(t.ctx, c, t.bare.Conf())
		t.NoError(err)

		_, err = t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Value: pw,
			Type:  horus.TokenTypePassword,
		})
		t.NoError(err)

		_, err = t.svc.Auth().BasicSignIn(t.ctx, &horus.BasicSignInRequest{
			Username: t.me.Actor.Alias,
			Password: "0000 0000",
		})
		t.ErrCode(err, codes.Unauthenticated)

		_, err = t.svc.Auth().BasicSignIn(t.ctx, &horus.BasicSignInRequest{
			Username: t.me.Actor.Alias,
			Password: "0000 0001",
		})
		t.ErrCode(err, codes.Unauthenticated)
	})
	t.Run("user cannot sign in with credentials that have reached the use count limit", func() {
		_, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Value: pw,
			Type:  horus.TokenTypePassword,

			UseCountLimit: fx.Addr(uint64(1)),
		})
		t.NoError(err)

		v, err := t.svc.Auth().BasicSignIn(t.ctx, &horus.BasicSignInRequest{
			Username: t.me.Actor.Alias,
			Password: pw,
		})
		t.NoError(err)
		t.Equal(horus.TokenTypeAccess, v.Token.Type)

		_, err = t.svc.Auth().BasicSignIn(t.ctx, &horus.BasicSignInRequest{
			Username: t.me.Actor.Alias,
			Password: pw,
		})
		t.ErrCode(err, codes.Unauthenticated)
	})
}

func (t *AuthTestSuite) TestTokenSignIn() {
	t.Run("user can sign in using their access token", func() {
		v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Type: horus.TokenTypeAccess,
		})
		t.NoError(err)

		ctx := frame.WithContext(t.ctx, &frame.Frame{})
		w, err := t.svc.Auth().TokenSignIn(ctx, &horus.TokenSignInRequest{
			Token: v.Value,
		})
		t.NoError(err)
		t.Equal(v.Id, w.Token.Id)

		f, ok := frame.Get(ctx)
		t.True(ok)
		t.Equal(t.me.Actor.ID, f.Actor.ID)
	})
	t.Run("user cannot sign in using a deleted access token", func() {
		v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Type: horus.TokenTypeAccess,
		})
		t.NoError(err)

		_, err = t.svc.Token().Delete(t.CtxMe(), horus.TokenByIdV(v.Id))
		t.NoError(err)

		_, err = t.svc.Auth().TokenSignIn(t.ctx, &horus.TokenSignInRequest{
			Token: v.Value,
		})
		t.ErrCode(err, codes.Unauthenticated)
	})
	t.Run("user cannot sign in using an access token that does not exist", func() {
		v := append(uuid.Nil[:], []byte("not exist")...)
		_, err := t.svc.Auth().TokenSignIn(t.ctx, &horus.TokenSignInRequest{
			Token: base64.RawStdEncoding.EncodeToString(v),
		})
		t.ErrCode(err, codes.Unauthenticated)
	})
	t.Run("user cannot sign in using their refresh token", func() {
		v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Type: horus.TokenTypeRefresh,
		})
		t.NoError(err)

		_, err = t.svc.Auth().TokenSignIn(t.ctx, &horus.TokenSignInRequest{
			Token: v.Value,
		})
		t.ErrCode(err, codes.Unauthenticated)
	})
}

func (t *AuthTestSuite) TestRefresh() {
	t.Run("user can create an access token using their refresh token", func() {
		v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Type: horus.TokenTypeRefresh,
		})
		t.NoError(err)

		w, err := t.svc.Auth().Refresh(t.ctx, &horus.RefreshRequest{
			Token: v.Value,
		})
		t.NoError(err)
		t.Equal(horus.TokenTypeAccess, w.Token.Type)

		_, err = t.svc.Auth().TokenSignIn(t.ctx, &horus.TokenSignInRequest{
			Token: w.Token.Value,
		})
		t.NoError(err)
	})
	t.Run("user cannot sign in using an old access token created by the refresh token", func() {
		v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Type: horus.TokenTypeRefresh,
		})
		t.NoError(err)

		w, err := t.svc.Auth().Refresh(t.ctx, &horus.RefreshRequest{
			Token: v.Value,
		})
		t.NoError(err)

		x, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Type: horus.TokenTypeAccess,
		})
		t.NoError(err)

		_, err = t.svc.Auth().Refresh(t.ctx, &horus.RefreshRequest{
			Token: v.Value,
		})
		t.NoError(err)

		// `w`, created by the refresh token, is expired after it is refreshed.
		_, err = t.svc.Auth().TokenSignIn(t.ctx, &horus.TokenSignInRequest{
			Token: w.Token.Value,
		})
		t.ErrCode(err, codes.Unauthenticated)

		// `x`, created using Token service, is not expired after the refresh.
		_, err = t.svc.Auth().TokenSignIn(t.ctx, &horus.TokenSignInRequest{
			Token: x.Value,
		})
		t.NoError(err)
	})
}
