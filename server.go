package horus

import (
	"google.golang.org/grpc"
)

type Server interface {
	Auth() AuthServiceServer
	Store
}

type Store interface {
	Conf() ConfServiceServer
	User() UserServiceServer
	Identity() IdentityServiceServer
	Account() AccountServiceServer
	Invitation() InvitationServiceServer
	Membership() MembershipServiceServer
	Silo() SiloServiceServer
	Team() TeamServiceServer
	Token() TokenServiceServer
}

func GrpcRegister(s *grpc.Server, svr Server) {
	RegisterAuthServiceServer(s, svr.Auth())
	RegisterConfServiceServer(s, svr.Conf())
	RegisterUserServiceServer(s, svr.User())
	RegisterIdentityServiceServer(s, svr.Identity())
	RegisterAccountServiceServer(s, svr.Account())
	RegisterInvitationServiceServer(s, svr.Invitation())
	RegisterMembershipServiceServer(s, svr.Membership())
	RegisterSiloServiceServer(s, svr.Silo())
	RegisterTeamServiceServer(s, svr.Team())
	RegisterTokenServiceServer(s, svr.Token())
}
