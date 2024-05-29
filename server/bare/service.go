package bare

import (
	horus "khepri.dev/horus"
	ent "khepri.dev/horus/ent"
)

func NewStore(client *ent.Client) horus.Store {
	return &store{client}
}

type store struct {
	client *ent.Client
}

func (s *store) User() horus.UserServiceServer {
	return NewUserService(s.client)
}

func (s *store) Account() horus.AccountServiceServer {
	return NewAccountService(s.client)
}

func (s *store) Invitation() horus.InvitationServiceServer {
	return NewInvitationService(s.client)
}

func (s *store) Membership() horus.MembershipServiceServer {
	return NewMembershipService(s.client)
}

func (s *store) Silo() horus.SiloServiceServer {
	return NewSiloService(s.client)
}

func (s *store) Team() horus.TeamServiceServer {
	return NewTeamService(s.client)
}

func (s *store) Token() horus.TokenServiceServer {
	return NewTokenService(s.client)
}
