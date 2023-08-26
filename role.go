package horus

import "khepri.dev/horus/internal/fx"

type RoleOrg string
type RoleTeam string

const (
	RoleOrgOwner  RoleOrg = "owner"
	RoleOrgMember RoleOrg = "member"

	RoleTeamOwner  RoleTeam = "owner"
	RoleTeamMember RoleTeam = "member"
)

func (RoleOrg) Values() []string {
	return fx.MapV([]RoleOrg{
		RoleOrgOwner,
		RoleOrgMember,
	}, func(v RoleOrg) string {
		return string(v)
	})
}

func (RoleTeam) Values() []string {
	return fx.MapV([]RoleTeam{
		RoleTeamOwner,
		RoleTeamMember,
	}, func(v RoleTeam) string {
		return string(v)
	})
}
