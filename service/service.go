package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/service/bare"
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
	client  *ent.Client
	bare    horus.Store
	service horus.Service

	keyer tokens.Keyer
}

func NewService(client *ent.Client) horus.Service {
	b := &base{
		client: client,
		bare:   bare.NewStore(client),

		keyer: tokens.NewArgon2i(&tokens.Argon2State{
			Parallelism: 4,
			TagLength:   32,
			MemorySize:  32 * (1 << 10), // 32 MiB
			Iterations:  3,
		}),
	}

	svc := &service{
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

	b.service = svc
	return svc
}

func UnaryInterceptor(svc horus.Service, db *ent.Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if strings.HasPrefix(info.FullMethod, "/khepri.horus.AuthService/") {
			return handler(ctx, req)
		}

		token := horus.Must(ctx)

		user, err := db.User.Get(ctx, uuid.UUID(token.Owner.Id))
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to get user details")
		}

		ctx = frame.WithContext(ctx, &frame.Frame{
			Actor: user,
		})
		return handler(ctx, req)
	}
}
