package server

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/server/bare"
	"khepri.dev/horus/server/frame"
	"khepri.dev/horus/tokens"
)

type server struct {
	auth horus.AuthServiceServer
	store
}

func (s *server) Auth() horus.AuthServiceServer {
	return s.auth
}

type store struct {
	user       horus.UserServiceServer
	account    horus.AccountServiceServer
	invitation horus.InvitationServiceServer
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

func (s *store) Invitation() horus.InvitationServiceServer {
	return s.invitation
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
	db      *ent.Client
	bare    horus.Store
	covered horus.Server

	keyer tokens.Keyer
}

func (b *base) withClient(client *ent.Client) *base {
	b_ := *b
	b_.db = client
	return &b_
}

func NewServer(client *ent.Client) horus.Server {
	b := &base{
		db:   client,
		bare: bare.NewStore(client),

		keyer: tokens.NewArgon2i(&tokens.Argon2State{
			Parallelism: 4,
			TagLength:   32,
			MemorySize:  32 * (1 << 10), // 32 MiB
			Iterations:  3,
		}),
	}

	svc := &server{
		auth: &AuthService{base: b},
		store: store{
			user:       &UserServiceServer{base: b},
			account:    &AccountServiceServer{base: b},
			invitation: &InvitationServiceServer{base: b},
			membership: &MembershipServiceServer{base: b},
			silo:       &SiloServiceServer{base: b},
			team:       &TeamServiceServer{base: b},
			token:      &TokenServiceServer{base: b},
		},
	}

	b.covered = svc
	return svc
}

func UnaryInterceptor(svc horus.Server, db *ent.Client) grpc.UnaryServerInterceptor {
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
