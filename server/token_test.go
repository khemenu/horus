package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
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
	pw := "0118 999 881 999 119 725 3"

	for _, req := range []*horus.CreateTokenRequest{
		{
			Value: pw,
			Type:  horus.TokenTypePassword,
		},
		{
			Type: horus.TokenTypeRefresh,
		},
		{
			Type: horus.TokenTypeAccess,
		},
	} {
		t.Run(fmt.Sprintf("token of type %s is created with the actor as its owner", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
			})
			t.NoError(err)
			t.Equal(t.me.Actor.ID[:], v.Owner.Id)
		})
		t.Run(fmt.Sprintf("token of type %s is created with actor's token as its parent", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
			})
			t.NoError(err)

			v, err = t.svc.Token().Get(t.CtxMe(), horus.TokenByIdV(v.Id))
			t.NoError(err)
			t.Equal(t.me.Token.ID[:], v.Parent.Id)
		})
		t.Run(fmt.Sprintf("token of type %s can be created with the actor's child as its owner", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
				Owner: horus.UserById(t.child.Actor.ID),
			})
			t.NoError(err)
			t.Equal(t.child.Actor.ID[:], v.Owner.Id)
		})
		t.Run(fmt.Sprintf("token of type %s cannot be created with another user as its owner", req.Type), func() {
			_, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
				Owner: horus.UserById(t.other.Actor.ID),
			})
			t.ErrCode(err, codes.NotFound)
		})
		t.Run(fmt.Sprintf("token of type %s cannot be created with a user that does not exist as its owner", req.Type), func() {
			_, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
				Owner: horus.UserByAlias("not exist"),
			})
			t.ErrCode(err, codes.NotFound)
		})
	}
	t.Run("token of password type does not reveal its value", func() {
		v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Value: pw,
			Type:  horus.TokenTypePassword,
		})
		t.NoError(err)
		t.Empty(v.Value)
	})
	t.Run("token of password type is salted and hashed", func() {
		v1, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Value: pw,
			Type:  horus.TokenTypePassword,
		})
		t.NoError(err)

		w1, err := t.db.Token.Get(t.ctx, uuid.UUID(v1.Id))
		t.NoError(err)
		t.NotEqual(pw, w1.Value)

		v2, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Value: pw,
			Type:  horus.TokenTypePassword,
		})
		t.NoError(err)

		w2, err := t.db.Token.Get(t.ctx, uuid.UUID(v2.Id))
		t.NoError(err)
		t.NotEqual(pw, w2.Value)
		t.NotEqual(w1.Value, w2.Value)
	})
	t.Run("token of password type is only one that exist", func() {
		v1, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Value: pw,
			Type:  horus.TokenTypePassword,
		})
		t.NoError(err)

		_, err = t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Value: pw,
			Type:  horus.TokenTypePassword,
		})
		t.NoError(err)

		_, err = t.svc.Token().Get(t.CtxMe(), horus.TokenByIdV(v1.Id))
		t.ErrCode(err, codes.NotFound)
	})
	t.Run("token of password type cannot be created without value", func() {
		_, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
			Type: horus.TokenTypePassword,
		})
		t.ErrCode(err, codes.InvalidArgument)
	})
}

func (t *TokenTestSuite) TestGet() {
	pw := "666"

	for _, req := range []*horus.CreateTokenRequest{
		{
			Value: pw,
			Type:  horus.TokenTypePassword,
		},
		{
			Type: horus.TokenTypeRefresh,
		},
		{
			Type: horus.TokenTypeAccess,
		},
	} {
		t.Run(fmt.Sprintf("token of type %s can be retrieved by its owner", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
			})
			t.NoError(err)

			w, err := t.svc.Token().Get(t.CtxMe(), horus.TokenByIdV(v.Id))
			t.NoError(err)
			t.Equal(v.Id, w.Id)
		})
		t.Run(fmt.Sprintf("token of type %s cannot be retrieved by its owner's parent", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxChild(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
			})
			t.NoError(err)

			_, err = t.svc.Token().Get(t.CtxMe(), horus.TokenByIdV(v.Id))
			t.ErrCode(err, codes.NotFound)
		})
		t.Run(fmt.Sprintf("token of type %s cannot be retrieved by another user", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
			})
			t.NoError(err)

			_, err = t.svc.Token().Get(t.CtxOther(), horus.TokenByIdV(v.Id))
			t.ErrCode(err, codes.NotFound)
		})
		t.Run(fmt.Sprintf("token of type %s cannot be retrieved with its value", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
			})
			t.NoError(err)

			w, err := t.svc.Token().Get(t.CtxMe(), horus.TokenByIdV(v.Id))
			t.NoError(err)
			t.Empty(w.Value)
		})
	}

	for _, req := range []*horus.CreateTokenRequest{
		{
			Type: horus.TokenTypeRefresh,
		},
		{
			Type: horus.TokenTypeAccess,
		},
	} {
		t.Run(fmt.Sprintf("token of type %s cannot be created with value", req.Type), func() {
			_, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: "Vincent",
				Type:  req.Type,
			})
			t.ErrCode(err, codes.InvalidArgument)
		})
	}
}

func (t *TokenTestSuite) TestUpdate() {
	pw := "0000"

	for _, req := range []*horus.CreateTokenRequest{
		{
			Value: pw,
			Type:  horus.TokenTypePassword,
		},
		{
			Type: horus.TokenTypeRefresh,
		},
		{
			Type: horus.TokenTypeAccess,
		},
	} {
		t.Run(fmt.Sprintf("token of type %s can be updated by its owner", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
			})
			t.NoError(err)

			w, err := t.svc.Token().Update(t.CtxMe(), &horus.UpdateTokenRequest{
				Key:  horus.TokenByIdV(v.Id),
				Name: fx.Addr("Moreau"),
			})
			t.NoError(err)
			t.Equal("Moreau", w.Name)
		})
		t.Run(fmt.Sprintf("token of type %s cannot be retrieved with its value by updating it", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
			})
			t.NoError(err)

			w, err := t.svc.Token().Update(t.CtxMe(), &horus.UpdateTokenRequest{
				Key:  horus.TokenByIdV(v.Id),
				Name: fx.Addr("Moreau"),
			})
			t.NoError(err)
			t.Empty(w.Value)
		})
		t.Run(fmt.Sprintf("token of type %s cannot be updated by its owner's parent", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxChild(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
			})
			t.NoError(err)

			_, err = t.svc.Token().Update(t.CtxMe(), &horus.UpdateTokenRequest{
				Key:  horus.TokenByIdV(v.Id),
				Name: fx.Addr("Moreau"),
			})
			t.ErrCode(err, codes.NotFound)
		})
		t.Run(fmt.Sprintf("token of type %s cannot be updated by another user", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
			})
			t.NoError(err)

			_, err = t.svc.Token().Update(t.CtxOther(), &horus.UpdateTokenRequest{
				Key:  horus.TokenByIdV(v.Id),
				Name: fx.Addr("Moreau"),
			})
			t.ErrCode(err, codes.NotFound)
		})
		t.Run(fmt.Sprintf("token of type %s cannot be updated if it does not exist", req.Type), func() {
			_, err := t.svc.Token().Update(t.CtxMe(), &horus.UpdateTokenRequest{
				Key:  horus.TokenById(uuid.Nil),
				Name: fx.Addr("Moreau"),
			})
			t.ErrCode(err, codes.NotFound)
		})
		t.Run(fmt.Sprintf("token of type %s cannot be updated if it is expired", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
			})
			t.NoError(err)

			_, err = t.svc.Token().Update(t.CtxMe(), &horus.UpdateTokenRequest{
				Key:         horus.TokenByIdV(v.Id),
				DateExpired: timestamppb.New(time.Now().Add(-time.Hour)),
			})
			t.NoError(err)

			_, err = t.svc.Token().Update(t.CtxMe(), &horus.UpdateTokenRequest{
				Key:  horus.TokenByIdV(v.Id),
				Name: fx.Addr("Moreau"),
			})
			t.ErrCode(err, codes.NotFound)
		})
	}
}

func (t *TokenTestSuite) TestDelete() {
	pw := "SHER"

	for _, req := range []*horus.CreateTokenRequest{
		{
			Value: pw,
			Type:  horus.TokenTypePassword,
		},
		{
			Type: horus.TokenTypeRefresh,
		},
		{
			Type: horus.TokenTypeAccess,
		},
	} {
		t.Run(fmt.Sprintf("token of type %s can be deleted by its owner", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
			})
			t.NoError(err)

			_, err = t.svc.Token().Delete(t.CtxMe(), horus.TokenByIdV(v.Id))
			t.NoError(err)
		})
		t.Run(fmt.Sprintf("token of type %s cannot be deleted by its owner's parent", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxChild(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
			})
			t.NoError(err)

			_, err = t.svc.Token().Delete(t.CtxMe(), horus.TokenByIdV(v.Id))
			t.ErrCode(err, codes.NotFound)
		})
		t.Run(fmt.Sprintf("token of type %s cannot be deleted by another user", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
			})
			t.NoError(err)

			_, err = t.svc.Token().Delete(t.CtxOther(), horus.TokenByIdV(v.Id))
			t.ErrCode(err, codes.NotFound)
		})
		t.Run(fmt.Sprintf("token of type %s cannot be deleted if it is expired", req.Type), func() {
			v, err := t.svc.Token().Create(t.CtxMe(), &horus.CreateTokenRequest{
				Value: req.Value,
				Type:  req.Type,
			})
			t.NoError(err)

			_, err = t.svc.Token().Update(t.CtxMe(), &horus.UpdateTokenRequest{
				Key:         horus.TokenByIdV(v.Id),
				DateExpired: timestamppb.New(time.Now().Add(-time.Hour)),
			})
			t.NoError(err)

			_, err = t.svc.Token().Delete(t.CtxMe(), horus.TokenByIdV(v.Id))
			t.ErrCode(err, codes.NotFound)
		})
	}
}
