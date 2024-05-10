package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/service/frame"
	"khepri.dev/horus/tokens"
)

type service struct {
	auth horus.AuthServiceServer
	store
}

func (s *service) Auth() horus.AuthServiceServer {
	return s.auth
}

type store struct {
	user       horus.UserServiceServer
	account    horus.AccountServiceServer
	membership horus.MembershipServiceServer
	silo       horus.SiloServiceServer
	team       horus.TeamServiceServer
	token      horus.TokenServiceServer
}

func (s *store) User() horus.UserServiceServer {
	return s.user
}

func (s *store) Account() horus.AccountServiceServer {
	return s.account
}

func (s *store) Membership() horus.MembershipServiceServer {
	return s.membership
}

func (s *store) Silo() horus.SiloServiceServer {
	return s.silo
}

func (s *store) Team() horus.TeamServiceServer {
	return s.team
}

func (s *store) Token() horus.TokenServiceServer {
	return s.token
}

type base struct {
	client *ent.Client
	bare   horus.Store

	keyer tokens.Keyer
}

func NewService(client *ent.Client) horus.Service {
	b := &base{
		client: client,
		keyer: tokens.NewArgon2iKeyer(tokens.Argon2iKeyerInit{
			Time:    3,
			Memory:  32 * (1 << 10),
			Threads: 4,
			KeyLen:  32,
		}),
	}

	return &service{
		auth: &AuthService{base: b},
		store: store{
			user:       &UserService{base: b},
			account:    &AccountService{base: b},
			membership: &MembershipService{base: b},
			silo:       &SiloService{base: b},
			team:       &TeamService{base: b},
			token:      &TokenService{base: b},
		},
	}
}

func GrpcUnaryInterceptor(svc horus.Service, db *ent.Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "missing metadata")
		}

		var token string
		if entry := md.Get("cookie"); len(entry) > 0 {
			prefix := fmt.Sprintf("%s=", tokens.CookieName)
			for _, cookie := range strings.Split(entry[0], "; ") {
				if !strings.HasPrefix(cookie, prefix) {
					continue
				}

				kv := strings.SplitN(cookie, "=", 2)
				if len(kv) != 2 {
					break
				}

				token = kv[1]
				break
			}
		}
		if token == "" {
			return nil, status.Error(codes.Unauthenticated, "no access token")
		}

		res, err := svc.Auth().TokenSignIn(ctx, &horus.TokenSignInRequest{Token: &horus.Token{
			Value: token,
		}})
		if err != nil {
			switch status.Code(err) {
			case codes.Unauthenticated:
				return nil, err
			default:
				return nil, status.Error(codes.Internal, "failed to get token details")
			}
		}

		user, err := db.User.Get(ctx, uuid.UUID(res.Token.Owner.Id))
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to get user details")
		}

		f := frame.New()
		f.Actor = user
		ctx = frame.WithContext(ctx, f)
		return handler(ctx, req)
	}
}
