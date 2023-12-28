package service

import (
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/proto/khepri/horus"
)

type Service interface {
	User() horus.UserServiceServer
	Account() horus.AccountServiceServer
	Membership() horus.MembershipServiceServer
	Silo() horus.SiloServiceServer
	Team() horus.TeamServiceServer
}

type service struct {
	user       horus.UserServiceServer
	account    horus.AccountServiceServer
	membership horus.MembershipServiceServer
	silo       horus.SiloServiceServer
	team       horus.TeamServiceServer
}

func (s *service) User() horus.UserServiceServer {
	return s.user
}

func (s *service) Account() horus.AccountServiceServer {
	return s.account
}

func (s *service) Membership() horus.MembershipServiceServer {
	return s.membership
}

func (s *service) Silo() horus.SiloServiceServer {
	return s.silo
}

func (s *service) Team() horus.TeamServiceServer {
	return s.team
}

type base struct {
	service

	client *ent.Client
	raw    Service
}

func NewService(client *ent.Client) Service {
	s := &base{
		client: client,
		raw: &service{
			user:       horus.NewUserService(client),
			account:    horus.NewAccountService(client),
			membership: horus.NewMembershipService(client),
			silo:       horus.NewSiloService(client),
			team:       horus.NewTeamService(client),
		},
	}
	s.user = &UserService{base: s}
	s.account = &AccountService{base: s}
	s.membership = &MembershipService{base: s}
	s.silo = &SiloService{base: s}
	s.team = &TeamService{base: s}

	return s
}
