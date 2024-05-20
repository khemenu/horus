// Code generated by protoc-gen-entgrpc. DO NOT EDIT.
package bare

import (
	context "context"
	runtime "entgo.io/contrib/entproto/runtime"
	sqlgraph "entgo.io/ent/dialect/sql/sqlgraph"
	uuid "github.com/google/uuid"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	horus "khepri.dev/horus"
	ent "khepri.dev/horus/ent"
	account "khepri.dev/horus/ent/account"
	membership "khepri.dev/horus/ent/membership"
	team "khepri.dev/horus/ent/team"
	regexp "regexp"
	strings "strings"
)

// MembershipService implements MembershipServiceServer
type MembershipService struct {
	client *ent.Client
	horus.UnimplementedMembershipServiceServer
}

// NewMembershipService returns a new MembershipService
func NewMembershipService(client *ent.Client) *MembershipService {
	return &MembershipService{
		client: client,
	}
}

var protoIdentNormalizeRegexpMembership_Role = regexp.MustCompile(`[^a-zA-Z0-9_]+`)

func protoIdentNormalizeMembership_Role(e string) string {
	return protoIdentNormalizeRegexpMembership_Role.ReplaceAllString(e, "_")
}

func toProtoMembership_Role(e membership.Role) horus.Membership_Role {
	if v, ok := horus.Membership_Role_value[strings.ToUpper("ROLE_"+protoIdentNormalizeMembership_Role(string(e)))]; ok {
		return horus.Membership_Role(v)
	}
	return horus.Membership_Role(0)
}

func toEntMembership_Role(e horus.Membership_Role) membership.Role {
	if v, ok := horus.Membership_Role_name[int32(e)]; ok {
		entVal := map[string]string{
			"ROLE_OWNER":  "OWNER",
			"ROLE_MEMBER": "MEMBER",
		}[v]
		return membership.Role(entVal)
	}
	return ""
}

// toProtoMembership transforms the ent type to the pb type
func toProtoMembership(e *ent.Membership) (*horus.Membership, error) {
	v := &horus.Membership{}
	created_date := timestamppb.New(e.CreatedDate)
	v.CreatedDate = created_date
	id, err := e.ID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	v.Id = id
	role := toProtoMembership_Role(e.Role)
	v.Role = role

	if edg := e.Edges.Account; edg != nil {
		id, err := edg.ID.MarshalBinary()
		if err != nil {
			return nil, err
		}
		v.Account = &horus.Account{
			Id: id,
		}
	}

	if edg := e.Edges.Team; edg != nil {
		id, err := edg.ID.MarshalBinary()
		if err != nil {
			return nil, err
		}
		v.Team = &horus.Team{
			Id: id,
		}
	}
	return v, nil
}

// Create implements MembershipServiceServer.Create
func (svc *MembershipService) Create(ctx context.Context, req *horus.CreateMembershipRequest) (*horus.Membership, error) {
	membership := req.GetMembership()
	m, err := svc.createBuilder(membership)
	if err != nil {
		return nil, err
	}
	res, err := m.Save(ctx)
	switch {
	case err == nil:
		proto, err := toProtoMembership(res)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal error: %s", err)
		}
		return proto, nil
	case sqlgraph.IsUniqueConstraintError(err):
		return nil, status.Errorf(codes.AlreadyExists, "already exists: %s", err)
	case ent.IsConstraintError(err):
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// Get implements MembershipServiceServer.Get
func (svc *MembershipService) Get(ctx context.Context, req *horus.GetMembershipRequest) (*horus.Membership, error) {

	var (
		err error
		get *ent.Membership
	)
	var id uuid.UUID
	if err := (&id).UnmarshalBinary(req.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	}
	switch req.GetView() {
	case horus.GetMembershipRequest_VIEW_UNSPECIFIED:
		fallthrough
	case horus.GetMembershipRequest_BASIC:
		get, err = svc.client.Membership.Get(ctx, id)
	case horus.GetMembershipRequest_WITH_EDGE_IDS:
		get, err = svc.client.Membership.Query().
			Where(membership.ID(id)).
			WithAccount(func(query *ent.AccountQuery) {
				query.Select(account.FieldID)
			}).
			WithTeam(func(query *ent.TeamQuery) {
				query.Select(team.FieldID)
			}).
			Only(ctx)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid argument: unknown view")
	}
	switch {
	case err == nil:
		return toProtoMembership(get)
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// Update implements MembershipServiceServer.Update
func (svc *MembershipService) Update(ctx context.Context, req *horus.UpdateMembershipRequest) (*horus.Membership, error) {
	membership := req.GetMembership()
	var membershipID uuid.UUID
	if err := (&membershipID).UnmarshalBinary(membership.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	}
	m := svc.client.Membership.UpdateOneID(membershipID)
	membershipRole := toEntMembership_Role(membership.GetRole())
	m.SetRole(membershipRole)

	res, err := m.Save(ctx)
	switch {
	case err == nil:
		proto, err := toProtoMembership(res)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal error: %s", err)
		}
		return proto, nil
	case sqlgraph.IsUniqueConstraintError(err):
		return nil, status.Errorf(codes.AlreadyExists, "already exists: %s", err)
	case ent.IsConstraintError(err):
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// Delete implements MembershipServiceServer.Delete
func (svc *MembershipService) Delete(ctx context.Context, req *horus.DeleteMembershipRequest) (*emptypb.Empty, error) {
	var err error
	var id uuid.UUID
	if err := (&id).UnmarshalBinary(req.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	}
	err = svc.client.Membership.DeleteOneID(id).Exec(ctx)
	switch {
	case err == nil:
		return &emptypb.Empty{}, nil
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

func (svc *MembershipService) createBuilder(membership *horus.Membership) (*ent.MembershipCreate, error) {
	m := svc.client.Membership.Create()
	membershipCreatedDate := runtime.ExtractTime(membership.GetCreatedDate())
	m.SetCreatedDate(membershipCreatedDate)
	membershipRole := toEntMembership_Role(membership.GetRole())
	m.SetRole(membershipRole)
	if membership.GetAccount() != nil {
		var membershipAccount uuid.UUID
		if err := (&membershipAccount).UnmarshalBinary(membership.GetAccount().GetId()); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
		}
		m.SetAccountID(membershipAccount)
	}
	if membership.GetTeam() != nil {
		var membershipTeam uuid.UUID
		if err := (&membershipTeam).UnmarshalBinary(membership.GetTeam().GetId()); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
		}
		m.SetTeamID(membershipTeam)
	}
	return m, nil
}