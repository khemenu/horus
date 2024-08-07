package role

import (
	"khepri.dev/horus/internal/fx"
)

const (
	Owner  Role = "OWNER"
	Admin  Role = "ADMIN"
	Member Role = "MEMBER"
)

func Values() []Role {
	return []Role{
		Owner,
		Admin,
		Member,
	}
}

type Role string

func (r Role) V() int {
	switch r {
	case Owner:
		return 30
	case Admin:
		return 20
	case Member:
		return 10
	default:
		return 0
	}
}
func (r Role) IsNil() bool {
	return r.V() == 0
}

func (r Role) Values() []string {
	return fx.MapV(Values(), func(v Role) string {
		return string(v)
	})
}

// -1 if x has less permissions than y,
// .0 if x equals y,
// +1 if x has more permissions than y,
func Compare(x, y Role) int {
	v := x.V() - y.V()
	if v < 0 {
		return -1
	}
	if v > 0 {
		return 1
	}
	return 0
}

func Lower(x, y Role) bool {
	return Compare(x, y) < 0
}

func Higher(x, y Role) bool {
	return Compare(x, y) > 0
}

func (r Role) LowerThan(v Role) bool {
	return Lower(r, v)
}

func (r Role) HigherThan(v Role) bool {
	return Higher(r, v)
}
