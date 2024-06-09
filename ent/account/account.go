// Code generated by ent, DO NOT EDIT.

package account

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"khepri.dev/horus/role"
)

const (
	// Label holds the string label denoting the account type in the database.
	Label = "account"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldDateCreated holds the string denoting the date_created field in the database.
	FieldDateCreated = "date_created"
	// FieldAlias holds the string denoting the alias field in the database.
	FieldAlias = "alias"
	// FieldOwnerID holds the string denoting the owner_id field in the database.
	FieldOwnerID = "owner_id"
	// FieldSiloID holds the string denoting the silo_id field in the database.
	FieldSiloID = "silo_id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldRole holds the string denoting the role field in the database.
	FieldRole = "role"
	// EdgeOwner holds the string denoting the owner edge name in mutations.
	EdgeOwner = "owner"
	// EdgeSilo holds the string denoting the silo edge name in mutations.
	EdgeSilo = "silo"
	// EdgeMemberships holds the string denoting the memberships edge name in mutations.
	EdgeMemberships = "memberships"
	// EdgeInvitations holds the string denoting the invitations edge name in mutations.
	EdgeInvitations = "invitations"
	// Table holds the table name of the account in the database.
	Table = "accounts"
	// OwnerTable is the table that holds the owner relation/edge.
	OwnerTable = "accounts"
	// OwnerInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	OwnerInverseTable = "users"
	// OwnerColumn is the table column denoting the owner relation/edge.
	OwnerColumn = "owner_id"
	// SiloTable is the table that holds the silo relation/edge.
	SiloTable = "accounts"
	// SiloInverseTable is the table name for the Silo entity.
	// It exists in this package in order to avoid circular dependency with the "silo" package.
	SiloInverseTable = "silos"
	// SiloColumn is the table column denoting the silo relation/edge.
	SiloColumn = "silo_id"
	// MembershipsTable is the table that holds the memberships relation/edge.
	MembershipsTable = "memberships"
	// MembershipsInverseTable is the table name for the Membership entity.
	// It exists in this package in order to avoid circular dependency with the "membership" package.
	MembershipsInverseTable = "memberships"
	// MembershipsColumn is the table column denoting the memberships relation/edge.
	MembershipsColumn = "account_memberships"
	// InvitationsTable is the table that holds the invitations relation/edge.
	InvitationsTable = "invitations"
	// InvitationsInverseTable is the table name for the Invitation entity.
	// It exists in this package in order to avoid circular dependency with the "invitation" package.
	InvitationsInverseTable = "invitations"
	// InvitationsColumn is the table column denoting the invitations relation/edge.
	InvitationsColumn = "account_invitations"
)

// Columns holds all SQL columns for account fields.
var Columns = []string{
	FieldID,
	FieldDateCreated,
	FieldAlias,
	FieldOwnerID,
	FieldSiloID,
	FieldName,
	FieldDescription,
	FieldRole,
}

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
	// DefaultDateCreated holds the default value on creation for the "date_created" field.
	DefaultDateCreated func() time.Time
	// DefaultAlias holds the default value on creation for the "alias" field.
	DefaultAlias func() string
	// AliasValidator is a validator for the "alias" field. It is called by the builders before save.
	AliasValidator func(string) error
	// DefaultName holds the default value on creation for the "name" field.
	DefaultName string
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultDescription holds the default value on creation for the "description" field.
	DefaultDescription string
	// DescriptionValidator is a validator for the "description" field. It is called by the builders before save.
	DescriptionValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// RoleValidator is a validator for the "role" field enum values. It is called by the builders before save.
func RoleValidator(r role.Role) error {
	switch r {
	case "OWNER", "MEMBER":
		return nil
	default:
		return fmt.Errorf("account: invalid enum value for role field: %q", r)
	}
}

// OrderOption defines the ordering options for the Account queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByDateCreated orders the results by the date_created field.
func ByDateCreated(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDateCreated, opts...).ToFunc()
}

// ByAlias orders the results by the alias field.
func ByAlias(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAlias, opts...).ToFunc()
}

// ByOwnerID orders the results by the owner_id field.
func ByOwnerID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOwnerID, opts...).ToFunc()
}

// BySiloID orders the results by the silo_id field.
func BySiloID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSiloID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByRole orders the results by the role field.
func ByRole(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRole, opts...).ToFunc()
}

// ByOwnerField orders the results by owner field.
func ByOwnerField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOwnerStep(), sql.OrderByField(field, opts...))
	}
}

// BySiloField orders the results by silo field.
func BySiloField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSiloStep(), sql.OrderByField(field, opts...))
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

// ByInvitationsCount orders the results by invitations count.
func ByInvitationsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newInvitationsStep(), opts...)
	}
}

// ByInvitations orders the results by invitations terms.
func ByInvitations(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newInvitationsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newOwnerStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OwnerInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, OwnerTable, OwnerColumn),
	)
}
func newSiloStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SiloInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, SiloTable, SiloColumn),
	)
}
func newMembershipsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MembershipsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, MembershipsTable, MembershipsColumn),
	)
}
func newInvitationsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(InvitationsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, InvitationsTable, InvitationsColumn),
	)
}
