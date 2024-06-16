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
	predicate "khepri.dev/horus/ent/predicate"
	team "khepri.dev/horus/ent/team"
)

type TeamServiceServer struct {
	db *ent.Client
	horus.UnimplementedTeamServiceServer
}

func NewTeamServiceServer(db *ent.Client) *TeamServiceServer {
	return &TeamServiceServer{db: db}
}
func (s *TeamServiceServer) Create(ctx context.Context, req *horus.CreateTeamRequest) (*horus.Team, error) {
	q := s.db.Team.Create()
	if v := req.Alias; v != nil {
		q.SetAlias(*v)
	}
	if v := req.Name; v != nil {
		q.SetName(*v)
	}
	if v := req.Description; v != nil {
		q.SetDescription(*v)
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

	return ToProtoTeam(res), nil
}
func (s *TeamServiceServer) Delete(ctx context.Context, req *horus.DeleteTeamRequest) (*emptypb.Empty, error) {
	q := s.db.Team.Delete()
	if v, err := uuid.FromBytes(req.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
	} else {
		q.Where(team.IDEQ(v))
	}

	_, err := q.Exec(ctx)
	if err != nil {
		return nil, ToStatus(err)
	}

	return &emptypb.Empty{}, nil
}
func (s *TeamServiceServer) Get(ctx context.Context, req *horus.GetTeamRequest) (*horus.Team, error) {
	q := s.db.Team.Query()
	if p, err := GetTeamSpecifier(req); err != nil {
		return nil, err
	} else {
		q.Where(p)
	}

	q.WithSilo(func(q *ent.SiloQuery) { q.Select(team.FieldID) })

	res, err := q.Only(ctx)
	if err != nil {
		return nil, ToStatus(err)
	}

	return ToProtoTeam(res), nil
}
func (s *TeamServiceServer) Update(ctx context.Context, req *horus.UpdateTeamRequest) (*horus.Team, error) {
	id, err := uuid.FromBytes(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
	}

	q := s.db.Team.UpdateOneID(id)
	if v := req.Alias; v != nil {
		q.SetAlias(*v)
	}
	if v := req.Name; v != nil {
		q.SetName(*v)
	}
	if v := req.Description; v != nil {
		q.SetDescription(*v)
	}

	res, err := q.Save(ctx)
	if err != nil {
		return nil, ToStatus(err)
	}

	return ToProtoTeam(res), nil
}
func ToProtoTeam(v *ent.Team) *horus.Team {
	m := &horus.Team{}
	m.Id = v.ID[:]
	m.DateCreated = timestamppb.New(v.DateCreated)
	m.Alias = v.Alias
	m.Name = v.Name
	m.Description = v.Description
	if v := v.Edges.Silo; v != nil {
		m.Silo = ToProtoSilo(v)
	}
	return m
}
func GetTeamId(ctx context.Context, db *ent.Client, req *horus.GetTeamRequest) (uuid.UUID, error) {
	var r uuid.UUID
	if v, err := uuid.FromBytes(req.GetId()); err != nil {
		return r, status.Errorf(codes.InvalidArgument, "id: %s", err)
	} else {
		return v, nil
	}
}
func GetTeamSpecifier(req *horus.GetTeamRequest) (predicate.Team, error) {
	if v, err := uuid.FromBytes(req.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
	} else {
		return team.IDEQ(v), nil
	}
}
