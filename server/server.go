package server

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/user"
	"khepri.dev/horus/server/bare"
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
	conf       horus.ConfServiceServer
	user       horus.UserServiceServer
	identity   horus.IdentityServiceServer
	account    horus.AccountServiceServer
	invitation horus.InvitationServiceServer
	membership horus.MembershipServiceServer
	silo       horus.SiloServiceServer
	team       horus.TeamServiceServer
	token      horus.TokenServiceServer
}

func (s *store) Conf() horus.ConfServiceServer {
	return s.conf
}

func (s *store) User() horus.UserServiceServer {
	return s.user
}

func (s *store) Identity() horus.IdentityServiceServer {
	return s.identity
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
	b_.bare = bare.NewStore(client)
	return &b_
}

func (b *base) isAncestor(ctx context.Context, actor *ent.User, target *ent.User) (bool, error) {
	if target.ParentID == nil {
		return false, nil
	}

	parent_id := target.ParentID
	for {
		if parent_id == nil || *parent_id == actor.ID {
			return true, nil
		}

		u, err := b.db.User.Query().
			Select(user.FieldParentID).
			Where(user.ID(*parent_id)).
			Only(ctx)
		if err != nil {
			return false, bare.ToStatus(fmt.Errorf("query: %w", err))
		}
		if u.ParentID == nil {
			return false, nil
		}

		parent_id = u.ParentID
	}
}

func (b *base) isAncestorOrMe(ctx context.Context, actor *ent.User, target *ent.User) (bool, error) {
	if actor.ID == target.ID || actor.Alias == target.Alias {
		return true, nil
	}

	return b.isAncestor(ctx, actor, target)
}

func (b *base) isAncestorQ(ctx context.Context, actor *ent.User, req *horus.GetUserRequest) (*ent.User, error) {
	p, err := bare.GetUserSpecifier(req)
	if err != nil {
		return nil, err
	}

	u, err := b.db.User.Query().
		Select(user.FieldParentID).
		Where(p).
		Only(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}

	ok, err := b.isAncestor(ctx, actor, u)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return u, nil
}

func (b *base) isAncestorOrMeQ(ctx context.Context, actor *ent.User, req *horus.GetUserRequest) (*ent.User, error) {
	if req == nil || req.GetKey() == nil {
		return actor, nil
	}

	switch k := req.Key.(type) {
	case *horus.GetUserRequest_Id:
		if bytes.Equal(actor.ID[:], k.Id) {
			return actor, nil
		}
	case *horus.GetUserRequest_Alias:
		if k.Alias == actor.Alias || k.Alias == horus.Me {
			return actor, nil
		}
	case *horus.GetUserRequest_Query:
		q := k.Query

		if v, ok := strings.CutPrefix(q, "@"); ok {
			return b.isAncestorOrMeQ(ctx, actor, horus.UserByAlias(v))

		}
		v, err := uuid.Parse(q)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid query string: %s", err)
		}
		return b.isAncestorOrMeQ(ctx, actor, horus.UserById(v))
	}

	return b.isAncestorQ(ctx, actor, req)
}

func (s *base) hasPermission(ctx context.Context, actor *ent.User, req *horus.GetUserRequest) (*ent.User, error) {
	v, err := s.isAncestorOrMeQ(ctx, actor, req)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
	}

	return v, nil
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
			conf:       &ConfServiceServer{base: b},
			user:       &UserServiceServer{base: b},
			identity:   &IdentityServiceServer{base: b},
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

		return handler(ctx, req)
	}
}
