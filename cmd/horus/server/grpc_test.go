package server_test

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/horus/server"
	"khepri.dev/horus/pb"
)

type horusGrpc struct {
	horus.Horus

	require *require.Assertions
	ctx     context.Context
	client  pb.HorusClient

	user     *horus.User
	identity *horus.Identity
}

func (h *horusGrpc) WithNewIdentity(ctx context.Context, init *horus.IdentityInit) *horusGrpc {
	if init.VerifiedBy == "" {
		init.VerifiedBy = "horus"
	}

	identity, err := h.Identities().New(ctx, init)
	h.require.NoError(err)

	user, err := h.Users().GetById(ctx, identity.OwnerId)
	h.require.NoError(err)

	access_token, err := h.Tokens().Issue(ctx, horus.TokenInit{
		OwnerId:  identity.OwnerId,
		Type:     horus.AccessToken,
		Duration: time.Hour,
	})
	h.require.NoError(err)

	return &horusGrpc{
		Horus: h.Horus,

		require: h.require,
		ctx:     metadata.NewOutgoingContext(ctx, metadata.Pairs(horus.CookieNameAccessToken, access_token.Value)),
		client:  h.client,

		user:     user,
		identity: identity,
	}
}

func WithHorusGrpc(conf *server.GrpcServerConfig, f func(require *require.Assertions, ctx context.Context, h *horusGrpc)) func(t *testing.T) {
	if conf == nil {
		conf = &server.GrpcServerConfig{}
	}
	return WithHorus(conf.Config, func(require *require.Assertions, h horus.Horus) {
		ctx := context.Background()

		horus_server, err := server.NewGrpcServer(h, conf)
		require.NoError(err)

		lis := bufconn.Listen(2 << 20)
		dialer := func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}

		conn, err := grpc.DialContext(ctx, "bufnet",
			grpc.WithContextDialer(dialer),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		require.NoError(err)

		s := grpc.NewServer(
			grpc.UnaryInterceptor(horus_server.UnaryInterceptor),
		)
		horus_server.Register(s)

		var wg sync.WaitGroup
		wg.Add(1)
		defer wg.Wait()

		go func() {
			defer wg.Done()
			err := s.Serve(lis)
			require.NoError(err)
		}()

		defer s.GracefulStop()

		hg := &horusGrpc{
			Horus: h,

			require: require,
			ctx:     ctx,
			client:  pb.NewHorusClient(conn),
		}

		hg = hg.WithNewIdentity(ctx, &horus.IdentityInit{
			Kind:       horus.IdentityMail,
			Value:      "ra@khepri.dev",
			VerifiedBy: "horus",
		})

		f(require, hg.ctx, hg)
	})
}

// func TestGrpcInterceptor(t *testing.T) {
// 	t.Run("without access token", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h horus.Horus, client pb.HorusClient) {
// 		ctx = metadata.NewOutgoingContext(ctx, metadata.MD{})

// 		_, err := client.Status(ctx, &pb.StatusReq{})
// 		require.Equal(codes.InvalidArgument, status.Code(err))
// 	}))

// 	t.Run("with invalid access token", WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h horus.Horus, client pb.HorusClient) {
// 		ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
// 			horus.CookieNameAccessToken: []string{"invalid"},
// 		})

// 		_, err := client.Status(ctx, &pb.StatusReq{})
// 		require.Equal(codes.Unauthenticated, status.Code(err))
// 	}))
// }

// func TestGrpcStatus(t *testing.T) {
// 	WithHorusGrpc(nil, func(require *require.Assertions, ctx context.Context, h horus.Horus, client pb.HorusClient) {
// 		res, err := client.Status(ctx, &pb.StatusReq{})
// 		require.NoError(err)
// 		require.NotEmpty(res.UserAlias)

// 		expired_at, err := time.Parse(time.RFC3339, res.SessionExpiredAt)
// 		require.NoError(err)
// 		require.True(time.Now().Before(expired_at))
// 	})(t)
// }
