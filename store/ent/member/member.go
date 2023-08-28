// Code generated by ent, DO NOT EDIT.

package member

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"khepri.dev/horus"
)

const (
	// Label holds the string label denoting the member type in the database.
	Label = "member"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldOrgID holds the string denoting the org_id field in the database.
	FieldOrgID = "org_id"
	// FieldUserID holds the string denoting the user_id field in the database.
	FieldUserID = "user_id"
	// FieldRole holds the string denoting the role field in the database.
	FieldRole = "role"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// EdgeOrg holds the string denoting the org edge name in mutations.
	EdgeOrg = "org"
	// EdgeTeams holds the string denoting the teams edge name in mutations.
	EdgeTeams = "teams"
	// EdgeIdentities holds the string denoting the identities edge name in mutations.
	EdgeIdentities = "identities"
	// EdgeMemberships holds the string denoting the memberships edge name in mutations.
	EdgeMemberships = "memberships"
	// Table holds the table name of the member in the database.
	Table = "members"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "members"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "user_id"
	// OrgTable is the table that holds the org relation/edge.
	OrgTable = "members"
	// OrgInverseTable is the table name for the Org entity.
	// It exists in this package in order to avoid circular dependency with the "org" package.
	OrgInverseTable = "orgs"
	// OrgColumn is the table column denoting the org relation/edge.
	OrgColumn = "org_id"
	// TeamsTable is the table that holds the teams relation/edge. The primary key declared below.
	TeamsTable = "memberships"
	// TeamsInverseTable is the table name for the Team entity.
	// It exists in this package in order to avoid circular dependency with the "team" package.
	TeamsInverseTable = "teams"
	// IdentitiesTable is the table that holds the identities relation/edge. The primary key declared below.
	IdentitiesTable = "member_identities"
	// IdentitiesInverseTable is the table name for the Identity entity.
	// It exists in this package in order to avoid circular dependency with the "identity" package.
	IdentitiesInverseTable = "identities"
	// MembershipsTable is the table that holds the memberships relation/edge.
	MembershipsTable = "memberships"
	// MembershipsInverseTable is the table name for the Membership entity.
	// It exists in this package in order to avoid circular dependency with the "membership" package.
	MembershipsInverseTable = "memberships"
	// MembershipsColumn is the table column denoting the memberships relation/edge.
	MembershipsColumn = "member_id"
)

// Columns holds all SQL columns for member fields.
var Columns = []string{
	FieldID,
	FieldOrgID,
	FieldUserID,
	FieldRole,
	FieldName,
	FieldCreatedAt,
}

var (
	// TeamsPrimaryKey and TeamsColumn2 are the table columns denoting the
	// primary key for the teams relation (M2M).
	TeamsPrimaryKey = []string{"team_id", "member_id"}
	// IdentitiesPrimaryKey and IdentitiesColumn2 are the table columns denoting the
	// primary key for the identities relation (M2M).
	IdentitiesPrimaryKey = []string{"member_id", "identity_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultName holds the default value on creation for the "name" field.
	DefaultName string
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// RoleValidator is a validator for the "role" field enum values. It is called by the builders before save.
func RoleValidator(r horus.RoleOrg) error {
	switch r {
	case "owner", "member", "invitee":
		return nil
	default:
		return fmt.Errorf("member: invalid enum value for role field: %q", r)
	}
}

// OrderOption defines the ordering options for the Member queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByOrgID orders the results by the org_id field.
func ByOrgID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOrgID, opts...).ToFunc()
}

// ByUserID orders the results by the user_id field.
func ByUserID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUserID, opts...).ToFunc()
}

// ByRole orders the results by the role field.
func ByRole(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRole, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUserField orders the results by user field.
func ByUserField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newUserStep(), sql.OrderByField(field, opts...))
	}
}

// ByOrgField orders the results by org field.
func ByOrgField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOrgStep(), sql.OrderByField(field, opts...))
	}
}

// ByTeamsCount orders the results by teams count.
func ByTeamsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newTeamsStep(), opts...)
	}
}

// ByTeams orders the results by teams terms.
func ByTeams(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newTeamsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByIdentitiesCount orders the results by identities count.
func ByIdentitiesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newIdentitiesStep(), opts...)
	}
}

// ByIdentities orders the results by identities terms.
func ByIdentities(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newIdentitiesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByMembershipsCount orders the results by memberships count.
func ByMembershipsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newMembershipsStep(), opts...)
	}
}

// ByMemberships orders the results by memberships terms.
func ByMemberships(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMembershipsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newUserStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(UserInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, UserTable, UserColumn),
	)
}
func newOrgStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OrgInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, OrgTable, OrgColumn),
	)
}
func newTeamsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(TeamsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, TeamsTable, TeamsPrimaryKey...),
	)
}
func newIdentitiesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(IdentitiesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, IdentitiesTable, IdentitiesPrimaryKey...),
	)
}
func newMembershipsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MembershipsInverseTable, MembershipsColumn),
		sqlgraph.Edge(sqlgraph.O2M, true, MembershipsTable, MembershipsColumn),
	)
}
