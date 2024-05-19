package horus

import (
	"google.golang.org/grpc"
)

type Service interface {
	Auth() AuthServiceServer
	Store
}

type Store interface {
	User() UserServiceServer
	Account() AccountServiceServer
	Membership() MembershipServiceServer
	Silo() SiloServiceServer
	Team() TeamServiceServer
	Token() TokenServiceServer
}

func GrpcRegister(svc Service, s *grpc.Server) {
	RegisterAuthServiceServer(s, svc.Auth())
	RegisterUserServiceServer(s, svc.User())
	RegisterAccountServiceServer(s, svc.Account())
	RegisterMembershipServiceServer(s, svc.Membership())
	RegisterSiloServiceServer(s, svc.Silo())
	RegisterTeamServiceServer(s, svc.Team())
	RegisterTokenServiceServer(s, svc.Token())
}
