package server_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/server/frame"
)

type TokenTestSuite struct {
	Suite
}

func TestToken(t *testing.T) {
	s := TokenTestSuite{
		Suite: NewSuiteWithSqliteStore(),
	}
	suite.Run(t, &s)
}

func (t *TokenTestSuite) TestCreate() {
	t.Run("password", func() {
		pw := "very secure"
		v, err := t.svc.Token().Create(t.ctx, &horus.CreateTokenRequest{
			Value: pw,
			Type:  horus.TokenTypeBasic,
		})
		t.NoError(err)

		res, err := t.svc.Auth().BasicSignIn(t.ctx, &horus.BasicSignInRequest{
			Username: t.me.Actor.Alias,
			Password: pw,
		})
		t.NoError(err)
		t.Equal(horus.TokenTypeAccess, res.Token.Type)

		v2, err := t.svc.Token().Get(t.ctx, horus.TokenByIdV(res.Token.Id))
		t.NoError(err)
		t.Equal(v.Id, v2.GetParent().GetId())
	})
	t.Run("user can create tokens for their child", func() {
		child, err := t.svc.User().Create(t.ctx, nil)
		t.NoError(err)

		v, err := t.svc.Token().Create(t.ctx, &horus.CreateTokenRequest{
			Type:  horus.TokenTypeRefresh,
			Owner: horus.UserByIdV(child.Id),
		})
		t.NoError(err)

		_, err = t.svc.Token().Get(t.ctx, horus.TokenByIdV(v.Id))
		t.ErrCode(err, codes.NotFound)

		child_ctx := frame.WithContext(t.ctx, &frame.Frame{
			Actor: &ent.User{ID: uuid.UUID(child.Id)},
		})
		_, err = t.svc.Token().Get(child_ctx, horus.TokenByIdV(v.Id))
		t.NoError(err)
	})
	t.Run("user cannot create tokens for other user", func() {
		_, err := t.svc.Token().Create(t.ctx, &horus.CreateTokenRequest{
			Type:  horus.TokenTypeRefresh,
			Owner: horus.UserById(t.other.Actor.ID),
		})
		t.ErrCode(err, codes.NotFound)
	})
}
