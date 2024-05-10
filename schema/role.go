package schema

import "golang.org/x/exp/maps"

const (
	RoleOwner  string = "OWNER"
	RoleMember string = "MEMBER"
)

type RoleSilo string

func (RoleSilo) Map() map[string]int32 {
	return map[string]int32{
		RoleOwner:  10,
		RoleMember: 20,
	}
}

func (r RoleSilo) Values() (kinds []string) {
	return maps.Keys(r.Map())
}

type RoleTeam string

func (RoleTeam) Map() map[string]int32 {
	return map[string]int32{
		RoleOwner:  10,
		RoleMember: 20,
	}
}

func (r RoleTeam) Values() (kinds []string) {
	return maps.Keys(r.Map())
}
