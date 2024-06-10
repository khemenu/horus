// Code generated by "protoc-gen-entpb". DO NOT EDIT

package bare

import (
	context "context"
	uuid "github.com/google/uuid"
	runtime "github.com/lesomnus/entpb/cmd/protoc-gen-entpb/runtime"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	horus "khepri.dev/horus"
	ent "khepri.dev/horus/ent"
	token "khepri.dev/horus/ent/token"
)

type TokenServiceServer struct {
	db *ent.Client
	horus.UnimplementedTokenServiceServer
}

func NewTokenServiceServer(db *ent.Client) *TokenServiceServer {
	return &TokenServiceServer{db: db}
}
func (s *TokenServiceServer) Create(ctx context.Context, req *horus.CreateTokenRequest) (*horus.Token, error) {
	q := s.db.Token.Create()
	q.SetValue(req.GetValue())
	q.SetType(req.GetType())
	if v := req.Name; v != nil {
		q.SetName(*v)
	}
	if v := req.GetDateExpired(); v != nil {
		w := v.AsTime()
		q.SetDateExpired(w)
	}
	if v := req.GetOwner().GetId(); v == nil {
		return nil, status.Errorf(codes.InvalidArgument, "field \"owner\" not provided")
	} else {
		if w, err := uuid.FromBytes(v); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "owner: %s", err)
		} else {
			q.SetOwnerID(w)
		}
	}

	res, err := q.Save(ctx)
	if err != nil {
		return nil, runtime.EntErrorToStatus(err)
	}

	return ToProtoToken(res), nil
}
func (s *TokenServiceServer) Delete(ctx context.Context, req *horus.DeleteTokenRequest) (*emptypb.Empty, error) {
	q := s.db.Token.Delete()
	switch t := req.GetKey().(type) {
	case *horus.DeleteTokenRequest_Id:
		if v, err := uuid.FromBytes(t.Id); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
		} else {
			q.Where(token.IDEQ(v))
		}
	case *horus.DeleteTokenRequest_Value:
		q.Where(token.ValueEQ(t.Value))
	default:
		return nil, status.Errorf(codes.InvalidArgument, "key not provided")
	}

	_, err := q.Exec(ctx)
	if err != nil {
		return nil, runtime.EntErrorToStatus(err)
	}

	return &emptypb.Empty{}, nil
}
func (s *TokenServiceServer) Get(ctx context.Context, req *horus.GetTokenRequest) (*horus.Token, error) {
	q := s.db.Token.Query()
	switch t := req.GetKey().(type) {
	case *horus.GetTokenRequest_Id:
		if v, err := uuid.FromBytes(t.Id); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
		} else {
			q.Where(token.IDEQ(v))
		}
	case *horus.GetTokenRequest_Value:
		q.Where(token.ValueEQ(t.Value))
	default:
		return nil, status.Errorf(codes.InvalidArgument, "key not provided")
	}

	q.WithOwner(func(q *ent.UserQuery) { q.Select(token.FieldID) })
	q.WithParent(func(q *ent.TokenQuery) { q.Select(token.FieldID) })
	q.WithChildren(func(q *ent.TokenQuery) { q.Select(token.FieldID) })

	res, err := q.Only(ctx)
	if err != nil {
		return nil, runtime.EntErrorToStatus(err)
	}

	return ToProtoToken(res), nil
}
func (s *TokenServiceServer) Update(ctx context.Context, req *horus.UpdateTokenRequest) (*horus.Token, error) {
	id, err := uuid.FromBytes(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
	}

	q := s.db.Token.UpdateOneID(id)
	if v := req.Name; v != nil {
		q.SetName(*v)
	}
	if v := req.DateExpired; v != nil {
		w := v.AsTime()
		q.SetDateExpired(w)
	}

	res, err := q.Save(ctx)
	if err != nil {
		return nil, runtime.EntErrorToStatus(err)
	}

	return ToProtoToken(res), nil
}
func ToProtoToken(v *ent.Token) *horus.Token {
	m := &horus.Token{}
	m.Id = v.ID[:]
	m.DateCreated = timestamppb.New(v.DateCreated)
	m.Value = v.Value
	m.Type = v.Type
	m.Name = v.Name
	m.DateExpired = timestamppb.New(v.DateExpired)
	if v := v.Edges.Owner; v != nil {
		m.Owner = &horus.User{Id: v.ID[:]}
	}
	if v := v.Edges.Parent; v != nil {
		m.Parent = &horus.Token{Id: v.ID[:]}
	}
	for _, v := range v.Edges.Children {
		m.Children = append(m.Children, &horus.Token{Id: v.ID[:]})
	}
	return m
}
