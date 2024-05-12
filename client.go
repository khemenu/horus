package horus

import (
	"google.golang.org/grpc"
)

type Client interface {
	User() UserServiceClient
	Account() AccountServiceClient
	Membership() MembershipServiceClient
	Silo() SiloServiceClient
	Team() TeamServiceClient
	Token() TokenServiceClient
}

type client struct {
	user       UserServiceClient
	account    AccountServiceClient
	membership MembershipServiceClient
	silo       SiloServiceClient
	team       TeamServiceClient
	token      TokenServiceClient
}

func NewClient(conn grpc.ClientConnInterface) Client {
	return &client{
		user:       NewUserServiceClient(conn),
		account:    NewAccountServiceClient(conn),
		membership: NewMembershipServiceClient(conn),
		silo:       NewSiloServiceClient(conn),
		team:       NewTeamServiceClient(conn),
		token:      NewTokenServiceClient(conn),
	}
}

func (c *client) User() UserServiceClient {
	return c.user
}

func (c *client) Account() AccountServiceClient {
	return c.account
}

func (c *client) Membership() MembershipServiceClient {
	return c.membership
}

func (c *client) Silo() SiloServiceClient {
	return c.silo
}

func (c *client) Team() TeamServiceClient {
	return c.team
}

func (c *client) Token() TokenServiceClient {
	return c.token
}
