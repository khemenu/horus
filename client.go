package horus

import (
	"google.golang.org/grpc"
)

type Client interface {
	Auth() AuthServiceClient

	User() UserServiceClient
	Account() AccountServiceClient
	Invitation() InvitationServiceClient
	Membership() MembershipServiceClient
	Silo() SiloServiceClient
	Team() TeamServiceClient
	Token() TokenServiceClient
}

type client struct {
	auth AuthServiceClient

	user       UserServiceClient
	account    AccountServiceClient
	invitation InvitationServiceClient
	membership MembershipServiceClient
	silo       SiloServiceClient
	team       TeamServiceClient
	token      TokenServiceClient
}

func NewClient(conn grpc.ClientConnInterface) Client {
	return &client{
		auth: NewAuthServiceClient(conn),

		user:       NewUserServiceClient(conn),
		account:    NewAccountServiceClient(conn),
		invitation: NewInvitationServiceClient(conn),
		membership: NewMembershipServiceClient(conn),
		silo:       NewSiloServiceClient(conn),
		team:       NewTeamServiceClient(conn),
		token:      NewTokenServiceClient(conn),
	}
}

func (c *client) Auth() AuthServiceClient {
	return c.auth
}

func (c *client) User() UserServiceClient {
	return c.user
}

func (c *client) Account() AccountServiceClient {
	return c.account
}

func (c *client) Invitation() InvitationServiceClient {
	return c.invitation
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

// Return error fast if there is no access token.
// func (c *client) preAuthInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
// 	if !strings.HasPrefix(method, "/khepri.horus.AuthService/") && c.access_token == nil {
// 		return status.Error(codes.Unauthenticated, "no access token")
// 	}

// 	return invoker(ctx, method, req, reply, cc, opts...)
// }

// // Set access token if available.
// func (c *client) authInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
// 	switch method {
// 	case AuthService_BasicSignIn_FullMethodName:
// 		if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
// 			return err
// 		}

// 		c.access_token = reply.(*BasicSignInRseponse).GetToken()
// 		return nil

// 	case AuthService_TokenSignIn_FullMethodName:
// 		token := req.(*TokenSignInRequest).GetToken()
// 		if token.Type != horus.TokenTypeAccess {
// 			return status.Errorf(codes.InvalidArgument, "token type must be \"%s\"", horus.TokenTypeAccess)
// 		}

// 		c.access_token = token
// 		reply = &TokenSignInResponse{Token: token}
// 		return nil

// 	case AuthService_SignOut_FullMethodName:
// 		if c.access_token == nil {
// 			reply = &SingOutResponse{}
// 			return nil
// 		}
// 		if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
// 			// TODO: log
// 		}

// 		c.access_token = nil
// 		return nil

// 	default:
// 		return invoker(ctx, method, req, reply, cc, opts...)
// 	}
// }

// // Attach access token to metadata.
// func (c *client) tokenAttachInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
// 	if c.access_token == nil {
// 		return invoker(ctx, method, req, reply, cc, opts...)
// 	}

// 	ctx = metadata.AppendToOutgoingContext(ctx, horus.TokenKeyName, c.access_token.Value)
// 	return invoker(ctx, method, req, reply, cc, opts...)
// }
