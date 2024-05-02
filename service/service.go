package service

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/proto/khepri/horus"
	"khepri.dev/horus/service/frame"
	"khepri.dev/horus/tokens"
)

type Store interface {
	User() horus.UserServiceServer
	Account() horus.AccountServiceServer
	Membership() horus.MembershipServiceServer
	Silo() horus.SiloServiceServer
	Team() horus.TeamServiceServer
	Token() horus.TokenServiceServer
}

type Service interface {
	Auth() horus.AuthServiceServer
	Store
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

type service struct {
	auth horus.AuthServiceServer
	store
}

func (s *service) Auth() horus.AuthServiceServer {
	return s.auth
}

type base struct {
	service

	client *ent.Client
	store  Store

	keyer tokens.Keyer
}

func NewService(client *ent.Client) Service {
	s := &base{
		client: client,
		store: &store{
			user:       horus.NewUserService(client),
			account:    horus.NewAccountService(client),
			membership: horus.NewMembershipService(client),
			silo:       horus.NewSiloService(client),
			team:       horus.NewTeamService(client),
			token:      horus.NewTokenService(client),
		},
	}
	s.auth = &AuthService{base: s}
	s.user = &UserService{base: s}
	s.account = &AccountService{base: s}
	s.membership = &MembershipService{base: s}
	s.silo = &SiloService{base: s}
	s.team = &TeamService{base: s}
	s.token = &TokenService{base: s}

	return s
}

func GrpcUnaryInterceptor(svc Service, db *ent.Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "missing metadata")
		}

		values, ok := md[tokens.CookieName]
		if !ok || len(values) != 1 {
			return nil, status.Error(codes.InvalidArgument, "no access token")
		}

		token, err := svc.Token().Get(ctx, &horus.GetTokenRequest{
			Id:   values[0],
			View: horus.GetTokenRequest_WITH_EDGE_IDS,
		})
		if err != nil {
			switch status.Code(err) {
			case codes.NotFound:
				return nil, status.Error(codes.Unauthenticated, "invalid token")
			default:
				return nil, status.Error(codes.Internal, "failed to get token details")
			}
		}

		user, err := db.User.Get(ctx, uuid.UUID(token.Owner.Id))
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to get user details")
		}

		f := frame.New()
		f.Actor = user
		ctx = frame.WithContext(ctx, f)
		return handler(ctx, req)
	}
}

func GrpcRegisterStoreService(svc Service, s *grpc.Server) {
	horus.RegisterUserServiceServer(s, svc.User())
	horus.RegisterAccountServiceServer(s, svc.Account())
	horus.RegisterMembershipServiceServer(s, svc.Membership())
	horus.RegisterSiloServiceServer(s, svc.Silo())
	horus.RegisterTeamServiceServer(s, svc.Team())
	horus.RegisterTokenServiceServer(s, svc.Token())
}
