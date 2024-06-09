package role

const (
	Owner  Role = "OWNER"
	Member Role = "MEMBER"
)

type Role string

func (r Role) Values() []string {
	return []string{
		string(Owner),
		string(Member),
	}
}
