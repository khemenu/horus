package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/horus/server/frame"
	"khepri.dev/horus/pb"
)

type GrpcServerConfig struct {
	*horus.Config
}

func (c *GrpcServerConfig) Normalize() error {
	errs := []error{}

	if c.Config == nil {
		c.Config = &horus.Config{}
	}
	if err := c.Config.Normalize(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("invalid config: %w", errors.Join(errs...))
	}

	return nil
}

type GrpcServer interface {
	pb.HorusServer
	Register(server grpc.ServiceRegistrar)
	UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
}

type grpcServer struct {
	horus.Horus
	pb.UnimplementedHorusServer

	conf *GrpcServerConfig
}

func NewGrpcServer(h horus.Horus, conf *GrpcServerConfig) (GrpcServer, error) {
	if conf == nil {
		conf = &GrpcServerConfig{}
	}
	if conf.Config == nil {
		conf.Config = h.Config()
	}
	if err := conf.Normalize(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &grpcServer{
		Horus: h,
		conf:  conf,
	}, nil
}

func (s *grpcServer) Register(server grpc.ServiceRegistrar) {
	pb.RegisterHorusServer(server, s)
}

func (s *grpcServer) UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "missing metadata")
	}

	values, ok := md[horus.CookieNameAccessToken]
	if !ok || len(values) != 1 {
		return nil, status.Error(codes.InvalidArgument, "no access token")
	}

	access_token, err := s.Tokens().GetByValue(ctx, values[0], horus.AccessToken)
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		return nil, status.Error(codes.Internal, "failed to get token details")
	}

	f := frame.NewFrame(s.Horus, access_token)
	if _, err := f.User(ctx); err != nil {
		return nil, status.Error(codes.Internal, "get user details")
	}

	ctx = frame.WithCtx(ctx, f)
	return handler(ctx, req)
}

func (s *grpcServer) mustUser(ctx context.Context) *horus.User {
	frame := frame.MustFromCtx(ctx)

	user, err := frame.User(ctx)
	if err != nil {
		panic(err)
	}

	return user
}

func grpcInternalErr(ctx context.Context, err error) error {
	return status.Error(codes.Internal, fmt.Sprintf("%s: %s", codes.Internal.String(), err.Error()))
}

func grpcStatusWithCode(code codes.Code) error {
	return status.Error(code, code.String())
}

func toPbRoleOrg(v horus.RoleOrg) pb.RoleOrg {
	switch v {
	case horus.RoleOrgOwner:
		return pb.RoleOrg_ROLE_ORG_OWNER
	case horus.RoleOrgMember:
		return pb.RoleOrg_ROLE_ORG_MEMBER
	}

	return pb.RoleOrg_ROLE_ORG_UNSPECIFIED
}

func fromPbRoleOrg(v pb.RoleOrg) horus.RoleOrg {
	switch v {
	case pb.RoleOrg_ROLE_ORG_OWNER:
		return horus.RoleOrgOwner
	case pb.RoleOrg_ROLE_ORG_MEMBER:
		return horus.RoleOrgMember
	}

	return ""
}

func parseUuidArg[T ~[16]byte](v []byte, name string) (T, error) {
	rst, err := uuid.FromBytes(v)
	if err != nil {
		desc := ""
		if name == "" {
			desc = "invalid ID"
		} else {
			desc = fmt.Sprintf("invalid %s ID", name)
		}

		return T{}, status.Errorf(codes.InvalidArgument, desc)
	}

	return T(rst), nil
}

func parseOrgId(v []byte) (horus.OrgId, error) {
	return parseUuidArg[horus.OrgId](v, "organization")
}

func parseTeamId(v []byte) (horus.TeamId, error) {
	return parseUuidArg[horus.TeamId](v, "team")
}

func parseMemberId(v []byte) (horus.MemberId, error) {
	return parseUuidArg[horus.MemberId](v, "member")
}
