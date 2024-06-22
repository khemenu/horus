package horus

import "khepri.dev/horus/role"

func (r Role) s() role.Role {
	switch r {
	case Role_ROLE_OWNER:
		return role.Owner
	case Role_ROLE_ADMIN:
		return role.Admin
	case Role_ROLE_MEMBER:
		return role.Member
	default:
		return role.Member
	}
}

func (r Role) LowerThan(v role.Role) bool {
	return r.s().LowerThan(v)
}

func (r Role) HigherThan(v role.Role) bool {
	return r.s().HigherThan(v)
}

func RoleFrom(v role.Role) Role {
	switch v {
	case role.Owner:
		return Role_ROLE_OWNER
	case role.Admin:
		return Role_ROLE_ADMIN
	case role.Member:
		return Role_ROLE_MEMBER
	default:
		return Role_ROLE_UNSPECIFIED
	}
}
