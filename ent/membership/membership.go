// Code generated by ent, DO NOT EDIT.

package membership

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"khepri.dev/horus/role"
)

const (
	// Label holds the string label denoting the membership type in the database.
	Label = "membership"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldDateCreated holds the string denoting the date_created field in the database.
	FieldDateCreated = "date_created"
	// FieldRole holds the string denoting the role field in the database.
	FieldRole = "role"
	// EdgeAccount holds the string denoting the account edge name in mutations.
	EdgeAccount = "account"
	// EdgeTeam holds the string denoting the team edge name in mutations.
	EdgeTeam = "team"
	// Table holds the table name of the membership in the database.
	Table = "memberships"
	// AccountTable is the table that holds the account relation/edge.
	AccountTable = "memberships"
	// AccountInverseTable is the table name for the Account entity.
	// It exists in this package in order to avoid circular dependency with the "account" package.
	AccountInverseTable = "accounts"
	// AccountColumn is the table column denoting the account relation/edge.
	AccountColumn = "account_memberships"
	// TeamTable is the table that holds the team relation/edge.
	TeamTable = "memberships"
	// TeamInverseTable is the table name for the Team entity.
	// It exists in this package in order to avoid circular dependency with the "team" package.
	TeamInverseTable = "teams"
	// TeamColumn is the table column denoting the team relation/edge.
	TeamColumn = "team_members"
)

// Columns holds all SQL columns for membership fields.
var Columns = []string{
	FieldID,
	FieldDateCreated,
	FieldRole,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "memberships"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"account_memberships",
	"team_members",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultDateCreated holds the default value on creation for the "date_created" field.
	DefaultDateCreated func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// RoleValidator is a validator for the "role" field enum values. It is called by the builders before save.
func RoleValidator(r role.Role) error {
	switch r {
	case "OWNER", "ADMIN", "MEMBER":
		return nil
	default:
		return fmt.Errorf("membership: invalid enum value for role field: %q", r)
	}
}

// OrderOption defines the ordering options for the Membership queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByDateCreated orders the results by the date_created field.
func ByDateCreated(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDateCreated, opts...).ToFunc()
}

// ByRole orders the results by the role field.
func ByRole(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRole, opts...).ToFunc()
}

// ByAccountField orders the results by account field.
func ByAccountField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newAccountStep(), sql.OrderByField(field, opts...))
	}
}

// ByTeamField orders the results by team field.
func ByTeamField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newTeamStep(), sql.OrderByField(field, opts...))
	}
}
func newAccountStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(AccountInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, AccountTable, AccountColumn),
	)
}
func newTeamStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(TeamInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, TeamTable, TeamColumn),
	)
}
