// Code generated by "protoc-gen-entpb". DO NOT EDIT

package bare

import (
	context "context"
	uuid "github.com/google/uuid"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	horus "khepri.dev/horus"
	ent "khepri.dev/horus/ent"
	account "khepri.dev/horus/ent/account"
	predicate "khepri.dev/horus/ent/predicate"
)

type AccountServiceServer struct {
	db *ent.Client
	horus.UnimplementedAccountServiceServer
}

func NewAccountServiceServer(db *ent.Client) *AccountServiceServer {
	return &AccountServiceServer{db: db}
}
func (s *AccountServiceServer) Create(ctx context.Context, req *horus.CreateAccountRequest) (*horus.Account, error) {
	q := s.db.Account.Create()
	if v := req.Alias; v != nil {
		q.SetAlias(*v)
	}
	if v := req.Name; v != nil {
		q.SetName(*v)
	}
	if v := req.Description; v != nil {
		q.SetDescription(*v)
	}
	q.SetRole(toEntRole(req.GetRole()))
	if v, err := GetUserId(ctx, s.db, req.GetOwner()); err != nil {
		return nil, err
	} else {
		q.SetOwnerID(v)
	}
	if v, err := GetSiloId(ctx, s.db, req.GetSilo()); err != nil {
		return nil, err
	} else {
		q.SetSiloID(v)
	}

	res, err := q.Save(ctx)
	if err != nil {
		return nil, ToStatus(err)
	}

	return ToProtoAccount(res), nil
}
func (s *AccountServiceServer) Delete(ctx context.Context, req *horus.DeleteAccountRequest) (*emptypb.Empty, error) {
	q := s.db.Account.Delete()
	if v, err := uuid.FromBytes(req.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
	} else {
		q.Where(account.IDEQ(v))
	}

	_, err := q.Exec(ctx)
	if err != nil {
		return nil, ToStatus(err)
	}

	return &emptypb.Empty{}, nil
}
func (s *AccountServiceServer) Get(ctx context.Context, req *horus.GetAccountRequest) (*horus.Account, error) {
	q := s.db.Account.Query()
	if p, err := GetAccountSpecifier(req); err != nil {
		return nil, err
	} else {
		q.Where(p)
	}

	q.WithOwner(func(q *ent.UserQuery) { q.Select(account.FieldID) })
	q.WithSilo(func(q *ent.SiloQuery) { q.Select(account.FieldID) })
	q.WithMemberships(func(q *ent.MembershipQuery) { q.Select(account.FieldID) })

	res, err := q.Only(ctx)
	if err != nil {
		return nil, ToStatus(err)
	}

	return ToProtoAccount(res), nil
}
func (s *AccountServiceServer) Update(ctx context.Context, req *horus.UpdateAccountRequest) (*horus.Account, error) {
	id, err := uuid.FromBytes(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
	}

	q := s.db.Account.UpdateOneID(id)
	if v := req.Alias; v != nil {
		q.SetAlias(*v)
	}
	if v := req.Name; v != nil {
		q.SetName(*v)
	}
	if v := req.Description; v != nil {
		q.SetDescription(*v)
	}
	if v := req.Role; v != nil {
		q.SetRole(toEntRole(*v))
	}

	res, err := q.Save(ctx)
	if err != nil {
		return nil, ToStatus(err)
	}

	return ToProtoAccount(res), nil
}
func ToProtoAccount(v *ent.Account) *horus.Account {
	m := &horus.Account{}
	m.Id = v.ID[:]
	m.DateCreated = timestamppb.New(v.DateCreated)
	m.Alias = v.Alias
	m.Name = v.Name
	m.Description = v.Description
	m.Role = toPbRole(v.Role)
	if v := v.Edges.Owner; v != nil {
		m.Owner = ToProtoUser(v)
	}
	if v := v.Edges.Silo; v != nil {
		m.Silo = ToProtoSilo(v)
	}
	for _, v := range v.Edges.Memberships {
		m.Memberships = append(m.Memberships, ToProtoMembership(v))
	}
	return m
}
func GetAccountId(ctx context.Context, db *ent.Client, req *horus.GetAccountRequest) (uuid.UUID, error) {
	var r uuid.UUID
	if v, err := uuid.FromBytes(req.GetId()); err != nil {
		return r, status.Errorf(codes.InvalidArgument, "id: %s", err)
	} else {
		return v, nil
	}
}
func GetAccountSpecifier(req *horus.GetAccountRequest) (predicate.Account, error) {
	if v, err := uuid.FromBytes(req.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
	} else {
		return account.IDEQ(v), nil
	}
}
