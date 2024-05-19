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
	// FieldInvitee holds the string denoting the invitee field in the database.
	FieldInvitee = "invitee"
	// FieldCreatedDate holds the string denoting the created_date field in the database.
	FieldCreatedDate = "created_date"
	// FieldExpiredDate holds the string denoting the expired_date field in the database.
	FieldExpiredDate = "expired_date"
	// FieldAcceptedDate holds the string denoting the accepted_date field in the database.
	FieldAcceptedDate = "accepted_date"
	// FieldDeclinedDate holds the string denoting the declined_date field in the database.
	FieldDeclinedDate = "declined_date"
	// FieldCanceledDate holds the string denoting the canceled_date field in the database.
	FieldCanceledDate = "canceled_date"
	// EdgeSilo holds the string denoting the silo edge name in mutations.
	EdgeSilo = "silo"
	// EdgeInviter holds the string denoting the inviter edge name in mutations.
	EdgeInviter = "inviter"
	// Table holds the table name of the invitation in the database.
	Table = "invitations"
	// SiloTable is the table that holds the silo relation/edge.
	SiloTable = "invitations"
	// SiloInverseTable is the table name for the Silo entity.
	// It exists in this package in order to avoid circular dependency with the "silo" package.
	SiloInverseTable = "silos"
	// SiloColumn is the table column denoting the silo relation/edge.
	SiloColumn = "silo_invitations"
	// InviterTable is the table that holds the inviter relation/edge.
	InviterTable = "invitations"
	// InviterInverseTable is the table name for the Account entity.
	// It exists in this package in order to avoid circular dependency with the "account" package.
	InviterInverseTable = "accounts"
	// InviterColumn is the table column denoting the inviter relation/edge.
	InviterColumn = "account_invitations"
)

// Columns holds all SQL columns for invitation fields.
var Columns = []string{
	FieldID,
	FieldInvitee,
	FieldCreatedDate,
	FieldExpiredDate,
	FieldAcceptedDate,
	FieldDeclinedDate,
	FieldCanceledDate,
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
	// InviteeValidator is a validator for the "invitee" field. It is called by the builders before save.
	InviteeValidator func(string) error
	// DefaultCreatedDate holds the default value on creation for the "created_date" field.
	DefaultCreatedDate func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Invitation queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByInvitee orders the results by the invitee field.
func ByInvitee(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldInvitee, opts...).ToFunc()
}

// ByCreatedDate orders the results by the created_date field.
func ByCreatedDate(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedDate, opts...).ToFunc()
}

// ByExpiredDate orders the results by the expired_date field.
func ByExpiredDate(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldExpiredDate, opts...).ToFunc()
}

// ByAcceptedDate orders the results by the accepted_date field.
func ByAcceptedDate(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAcceptedDate, opts...).ToFunc()
}

// ByDeclinedDate orders the results by the declined_date field.
func ByDeclinedDate(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDeclinedDate, opts...).ToFunc()
}

// ByCanceledDate orders the results by the canceled_date field.
func ByCanceledDate(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCanceledDate, opts...).ToFunc()
}

// BySiloField orders the results by silo field.
func BySiloField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSiloStep(), sql.OrderByField(field, opts...))
	}
}

// ByInviterField orders the results by inviter field.
func ByInviterField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newInviterStep(), sql.OrderByField(field, opts...))
	}
}
func newSiloStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SiloInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, SiloTable, SiloColumn),
	)
}
func newInviterStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(InviterInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, InviterTable, InviterColumn),
	)
}
