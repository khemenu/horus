// Code generated by ent, DO NOT EDIT.

package team

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the team type in the database.
	Label = "team"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldSiloID holds the string denoting the silo_id field in the database.
	FieldSiloID = "silo_id"
	// FieldAlias holds the string denoting the alias field in the database.
	FieldAlias = "alias"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldInterVisibility holds the string denoting the inter_visibility field in the database.
	FieldInterVisibility = "inter_visibility"
	// FieldIntraVisibility holds the string denoting the intra_visibility field in the database.
	FieldIntraVisibility = "intra_visibility"
	// FieldCreatedDate holds the string denoting the created_date field in the database.
	FieldCreatedDate = "created_date"
	// EdgeSilo holds the string denoting the silo edge name in mutations.
	EdgeSilo = "silo"
	// EdgeMembers holds the string denoting the members edge name in mutations.
	EdgeMembers = "members"
	// Table holds the table name of the team in the database.
	Table = "teams"
	// SiloTable is the table that holds the silo relation/edge.
	SiloTable = "teams"
	// SiloInverseTable is the table name for the Silo entity.
	// It exists in this package in order to avoid circular dependency with the "silo" package.
	SiloInverseTable = "silos"
	// SiloColumn is the table column denoting the silo relation/edge.
	SiloColumn = "silo_id"
	// MembersTable is the table that holds the members relation/edge.
	MembersTable = "memberships"
	// MembersInverseTable is the table name for the Membership entity.
	// It exists in this package in order to avoid circular dependency with the "membership" package.
	MembersInverseTable = "memberships"
	// MembersColumn is the table column denoting the members relation/edge.
	MembersColumn = "team_members"
)

// Columns holds all SQL columns for team fields.
var Columns = []string{
	FieldID,
	FieldSiloID,
	FieldAlias,
	FieldName,
	FieldDescription,
	FieldInterVisibility,
	FieldIntraVisibility,
	FieldCreatedDate,
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
	// DefaultAlias holds the default value on creation for the "alias" field.
	DefaultAlias func() string
	// AliasValidator is a validator for the "alias" field. It is called by the builders before save.
	AliasValidator func(string) error
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultDescription holds the default value on creation for the "description" field.
	DefaultDescription string
	// DescriptionValidator is a validator for the "description" field. It is called by the builders before save.
	DescriptionValidator func(string) error
	// DefaultCreatedDate holds the default value on creation for the "created_date" field.
	DefaultCreatedDate func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// InterVisibility defines the type for the "inter_visibility" enum field.
type InterVisibility string

// InterVisibility values.
const (
	InterVisibilityPUBLIC  InterVisibility = "PUBLIC"
	InterVisibilityPRIVATE InterVisibility = "PRIVATE"
)

func (iv InterVisibility) String() string {
	return string(iv)
}

// InterVisibilityValidator is a validator for the "inter_visibility" field enum values. It is called by the builders before save.
func InterVisibilityValidator(iv InterVisibility) error {
	switch iv {
	case InterVisibilityPUBLIC, InterVisibilityPRIVATE:
		return nil
	default:
		return fmt.Errorf("team: invalid enum value for inter_visibility field: %q", iv)
	}
}

// IntraVisibility defines the type for the "intra_visibility" enum field.
type IntraVisibility string

// IntraVisibility values.
const (
	IntraVisibilityPRIVATE IntraVisibility = "PRIVATE"
	IntraVisibilityPUBLIC  IntraVisibility = "PUBLIC"
)

func (iv IntraVisibility) String() string {
	return string(iv)
}

// IntraVisibilityValidator is a validator for the "intra_visibility" field enum values. It is called by the builders before save.
func IntraVisibilityValidator(iv IntraVisibility) error {
	switch iv {
	case IntraVisibilityPRIVATE, IntraVisibilityPUBLIC:
		return nil
	default:
		return fmt.Errorf("team: invalid enum value for intra_visibility field: %q", iv)
	}
}

// OrderOption defines the ordering options for the Team queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// BySiloID orders the results by the silo_id field.
func BySiloID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSiloID, opts...).ToFunc()
}

// ByAlias orders the results by the alias field.
func ByAlias(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAlias, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByInterVisibility orders the results by the inter_visibility field.
func ByInterVisibility(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldInterVisibility, opts...).ToFunc()
}

// ByIntraVisibility orders the results by the intra_visibility field.
func ByIntraVisibility(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIntraVisibility, opts...).ToFunc()
}

// ByCreatedDate orders the results by the created_date field.
func ByCreatedDate(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedDate, opts...).ToFunc()
}

// BySiloField orders the results by silo field.
func BySiloField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSiloStep(), sql.OrderByField(field, opts...))
	}
}

// ByMembersCount orders the results by members count.
func ByMembersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newMembersStep(), opts...)
	}
}

// ByMembers orders the results by members terms.
func ByMembers(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMembersStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newSiloStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SiloInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, SiloTable, SiloColumn),
	)
}
func newMembersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MembersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, MembersTable, MembersColumn),
	)
}
