package schema

import "golang.org/x/exp/maps"

const (
	VisibilityPrivate string = "PRIVATE"
	VisibilityPublic  string = "PUBLIC"
)

type Visibility string

func (Visibility) Map() map[string]int32 {
	return map[string]int32{
		VisibilityPrivate: 10,
		VisibilityPublic:  20,
	}
}

func (v Visibility) Values() (kinds []string) {
	return maps.Keys(v.Map())
}
