package horus

import "khepri.dev/horus/internal/fx"

type RoleOrg string
type RoleTeam string

const (
	RoleOrgOwner   RoleOrg = "owner"
	RoleOrgMember  RoleOrg = "member"
	RoleOrgInvitee RoleOrg = "invitee"

	RoleTeamOwner   RoleTeam = "owner"
	RoleTeamMember  RoleTeam = "member"
	RoleTeamInvitee RoleTeam = "invitee"
)

func (RoleOrg) Values() []string {
	return fx.MapV([]RoleOrg{
		RoleOrgOwner,
		RoleOrgMember,
		RoleOrgInvitee,
	}, func(v RoleOrg) string {
		return string(v)
	})
}

func (RoleTeam) Values() []string {
	return fx.MapV([]RoleTeam{
		RoleTeamOwner,
		RoleTeamMember,
		RoleTeamInvitee,
	}, func(v RoleTeam) string {
		return string(v)
	})
}
