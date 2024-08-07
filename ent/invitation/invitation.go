// Code generated by ent, DO NOT EDIT.

package invitation

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the invitation type in the database.
	Label = "invitation"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldDateCreated holds the string denoting the date_created field in the database.
	FieldDateCreated = "date_created"
	// FieldInvitee holds the string denoting the invitee field in the database.
	FieldInvitee = "invitee"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldDateExpired holds the string denoting the date_expired field in the database.
	FieldDateExpired = "date_expired"
	// FieldDateAccepted holds the string denoting the date_accepted field in the database.
	FieldDateAccepted = "date_accepted"
	// FieldDateDeclined holds the string denoting the date_declined field in the database.
	FieldDateDeclined = "date_declined"
	// FieldDateCanceled holds the string denoting the date_canceled field in the database.
	FieldDateCanceled = "date_canceled"
	// EdgeInviter holds the string denoting the inviter edge name in mutations.
	EdgeInviter = "inviter"
	// EdgeSilo holds the string denoting the silo edge name in mutations.
	EdgeSilo = "silo"
	// Table holds the table name of the invitation in the database.
	Table = "invitations"
	// InviterTable is the table that holds the inviter relation/edge.
	InviterTable = "invitations"
	// InviterInverseTable is the table name for the Account entity.
	// It exists in this package in order to avoid circular dependency with the "account" package.
	InviterInverseTable = "accounts"
	// InviterColumn is the table column denoting the inviter relation/edge.
	InviterColumn = "account_invitations"
	// SiloTable is the table that holds the silo relation/edge.
	SiloTable = "invitations"
	// SiloInverseTable is the table name for the Silo entity.
	// It exists in this package in order to avoid circular dependency with the "silo" package.
	SiloInverseTable = "silos"
	// SiloColumn is the table column denoting the silo relation/edge.
	SiloColumn = "silo_invitations"
)

// Columns holds all SQL columns for invitation fields.
var Columns = []string{
	FieldID,
	FieldDateCreated,
	FieldInvitee,
	FieldType,
	FieldDateExpired,
	FieldDateAccepted,
	FieldDateDeclined,
	FieldDateCanceled,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "invitations"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"account_invitations",
	"silo_invitations",
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
	// InviteeValidator is a validator for the "invitee" field. It is called by the builders before save.
	InviteeValidator func(string) error
	// TypeValidator is a validator for the "type" field. It is called by the builders before save.
	TypeValidator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Invitation queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByDateCreated orders the results by the date_created field.
func ByDateCreated(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDateCreated, opts...).ToFunc()
}

// ByInvitee orders the results by the invitee field.
func ByInvitee(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldInvitee, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// ByDateExpired orders the results by the date_expired field.
func ByDateExpired(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDateExpired, opts...).ToFunc()
}

// ByDateAccepted orders the results by the date_accepted field.
func ByDateAccepted(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDateAccepted, opts...).ToFunc()
}

// ByDateDeclined orders the results by the date_declined field.
func ByDateDeclined(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDateDeclined, opts...).ToFunc()
}

// ByDateCanceled orders the results by the date_canceled field.
func ByDateCanceled(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDateCanceled, opts...).ToFunc()
}

// ByInviterField orders the results by inviter field.
func ByInviterField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newInviterStep(), sql.OrderByField(field, opts...))
	}
}

// BySiloField orders the results by silo field.
func BySiloField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSiloStep(), sql.OrderByField(field, opts...))
	}
}
func newInviterStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(InviterInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, InviterTable, InviterColumn),
	)
}
func newSiloStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SiloInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, SiloTable, SiloColumn),
	)
}
