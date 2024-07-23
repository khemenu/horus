package bare

import (
	horus "khepri.dev/horus"
	ent "khepri.dev/horus/ent"
)

func NewStore(db *ent.Client) horus.Store {
	return &store{db}
}

type store struct {
	db *ent.Client
}

func (s *store) Conf() horus.ConfServiceServer {
	return NewConfServiceServer(s.db)
}

func (s *store) User() horus.UserServiceServer {
	return NewUserServiceServer(s.db)
}

func (s *store) Account() horus.AccountServiceServer {
	return NewAccountServiceServer(s.db)
}

func (s *store) Invitation() horus.InvitationServiceServer {
	return NewInvitationServiceServer(s.db)
}

func (s *store) Membership() horus.MembershipServiceServer {
	return NewMembershipServiceServer(s.db)
}

func (s *store) Silo() horus.SiloServiceServer {
	return NewSiloServiceServer(s.db)
}

func (s *store) Team() horus.TeamServiceServer {
	return NewTeamServiceServer(s.db)
}

func (s *store) Token() horus.TokenServiceServer {
	return NewTokenServiceServer(s.db)
}
