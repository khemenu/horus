package horus

import (
	"google.golang.org/grpc"
)

type Server interface {
	Auth() AuthServiceServer
	Store
}

type Store interface {
	User() UserServiceServer
	Account() AccountServiceServer
	Invitation() InvitationServiceServer
	Membership() MembershipServiceServer
	Silo() SiloServiceServer
	Team() TeamServiceServer
	Token() TokenServiceServer
}

func GrpcRegister(s *grpc.Server, svr Server) {
	RegisterAuthServiceServer(s, svr.Auth())
	RegisterUserServiceServer(s, svr.User())
	RegisterAccountServiceServer(s, svr.Account())
	RegisterInvitationServiceServer(s, svr.Invitation())
	RegisterMembershipServiceServer(s, svr.Membership())
	RegisterSiloServiceServer(s, svr.Silo())
	RegisterTeamServiceServer(s, svr.Team())
	RegisterTokenServiceServer(s, svr.Token())
}
