package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/log"
	"khepri.dev/horus/store/ent"
	"khepri.dev/horus/store/ent/member"
)

func member_(v *ent.Member) *horus.Member {
	return &horus.Member{
		Id:     horus.MemberId(v.ID),
		OrgId:  horus.OrgId(v.OrgID),
		UserId: horus.UserId(v.UserID),
		Name:   v.Name,
		Role:   v.Role,

		CreatedAt: v.CreatedAt,
	}
}

type memberStore struct {
	client *ent.Client
}

func NewMemberStore(client *ent.Client) (horus.MemberStore, error) {
	s := &memberStore{
		client: client,
	}

	return s, nil
}

func newMember(ctx context.Context, client *ent.Client, init horus.MemberInit) (*horus.Member, error) {
	res, err := client.Member.Create().
		SetOrgID(uuid.UUID(init.OrgId)).
		SetUserID(uuid.UUID(init.UserId)).
		SetRole(init.Role).
		SetName(init.Name).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			if strings.Contains(err.Error(), "FOREIGN KEY") {
				// User or Org does not exist.
				return nil, errors.Join(horus.ErrNotExist, err)
			} else if strings.Contains(err.Error(), "UNIQUE") {
				// Add a user who already a member.
				return nil, errors.Join(horus.ErrExist, err)
			}
		}

		return nil, fmt.Errorf("save: %w", err)
	}

	log.FromCtx(ctx).Info("new user", "id", res.ID)
	return member_(res), nil
}

func (s *memberStore) New(ctx context.Context, init horus.MemberInit) (*horus.Member, error) {
	return newMember(ctx, s.client, init)
}

func (s *memberStore) GetByUserIdFromOrg(ctx context.Context, org_id horus.OrgId, user_id horus.UserId) (*horus.Member, error) {
	res, err := s.client.Member.Query().
		Where(member.And(
			member.OrgID(uuid.UUID(org_id)),
			member.UserID(uuid.UUID(user_id)),
		)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return member_(res), nil
}

func (s *memberStore) GetAllByOrgId(ctx context.Context, org_id horus.OrgId) ([]*horus.Member, error) {
	res, err := s.client.Member.Query().
		Where(member.OrgID(uuid.UUID(org_id))).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	if len(res) == 0 {
		// Org does not exist.
		// There must be at least one member who is owner of the Org.
		return nil, horus.ErrNotExist
	}

	return fx.MapV(res, member_), nil
}
