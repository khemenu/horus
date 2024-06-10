// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/silo"
	"khepri.dev/horus/ent/user"
	"khepri.dev/horus/role"
)

// Account is the model entity for the Account schema.
type Account struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// DateCreated holds the value of the "date_created" field.
	DateCreated time.Time `json:"date_created,omitempty"`
	// Alias holds the value of the "alias" field.
	Alias string `json:"alias,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// OwnerID holds the value of the "owner_id" field.
	OwnerID uuid.UUID `json:"owner_id,omitempty"`
	// SiloID holds the value of the "silo_id" field.
	SiloID uuid.UUID `json:"silo_id,omitempty"`
	// Role holds the value of the "role" field.
	Role role.Role `json:"role,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the AccountQuery when eager-loading is set.
	Edges        AccountEdges `json:"edges"`
	selectValues sql.SelectValues
}

// AccountEdges holds the relations/edges for other nodes in the graph.
type AccountEdges struct {
	// Owner holds the value of the owner edge.
	Owner *User `json:"owner,omitempty"`
	// Silo holds the value of the silo edge.
	Silo *Silo `json:"silo,omitempty"`
	// Memberships holds the value of the memberships edge.
	Memberships []*Membership `json:"memberships,omitempty"`
	// Invitations holds the value of the invitations edge.
	Invitations []*Invitation `json:"invitations,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// OwnerOrErr returns the Owner value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AccountEdges) OwnerOrErr() (*User, error) {
	if e.Owner != nil {
		return e.Owner, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: user.Label}
	}
	return nil, &NotLoadedError{edge: "owner"}
}

// SiloOrErr returns the Silo value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AccountEdges) SiloOrErr() (*Silo, error) {
	if e.Silo != nil {
		return e.Silo, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: silo.Label}
	}
	return nil, &NotLoadedError{edge: "silo"}
}

// MembershipsOrErr returns the Memberships value or an error if the edge
// was not loaded in eager-loading.
func (e AccountEdges) MembershipsOrErr() ([]*Membership, error) {
	if e.loadedTypes[2] {
		return e.Memberships, nil
	}
	return nil, &NotLoadedError{edge: "memberships"}
}

// InvitationsOrErr returns the Invitations value or an error if the edge
// was not loaded in eager-loading.
func (e AccountEdges) InvitationsOrErr() ([]*Invitation, error) {
	if e.loadedTypes[3] {
		return e.Invitations, nil
	}
	return nil, &NotLoadedError{edge: "invitations"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Account) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case account.FieldAlias, account.FieldName, account.FieldDescription, account.FieldRole:
			values[i] = new(sql.NullString)
		case account.FieldDateCreated:
			values[i] = new(sql.NullTime)
		case account.FieldID, account.FieldOwnerID, account.FieldSiloID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Account fields.
func (a *Account) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case account.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				a.ID = *value
			}
		case account.FieldDateCreated:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field date_created", values[i])
			} else if value.Valid {
				a.DateCreated = value.Time
			}
		case account.FieldAlias:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field alias", values[i])
			} else if value.Valid {
				a.Alias = value.String
			}
		case account.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				a.Name = value.String
			}
		case account.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				a.Description = value.String
			}
		case account.FieldOwnerID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field owner_id", values[i])
			} else if value != nil {
				a.OwnerID = *value
			}
		case account.FieldSiloID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field silo_id", values[i])
			} else if value != nil {
				a.SiloID = *value
			}
		case account.FieldRole:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field role", values[i])
			} else if value.Valid {
				a.Role = role.Role(value.String)
			}
		default:
			a.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Account.
// This includes values selected through modifiers, order, etc.
func (a *Account) Value(name string) (ent.Value, error) {
	return a.selectValues.Get(name)
}

// QueryOwner queries the "owner" edge of the Account entity.
func (a *Account) QueryOwner() *UserQuery {
	return NewAccountClient(a.config).QueryOwner(a)
}

// QuerySilo queries the "silo" edge of the Account entity.
func (a *Account) QuerySilo() *SiloQuery {
	return NewAccountClient(a.config).QuerySilo(a)
}

// QueryMemberships queries the "memberships" edge of the Account entity.
func (a *Account) QueryMemberships() *MembershipQuery {
	return NewAccountClient(a.config).QueryMemberships(a)
}

// QueryInvitations queries the "invitations" edge of the Account entity.
func (a *Account) QueryInvitations() *InvitationQuery {
	return NewAccountClient(a.config).QueryInvitations(a)
}

// Update returns a builder for updating this Account.
// Note that you need to call Account.Unwrap() before calling this method if this Account
// was returned from a transaction, and the transaction was committed or rolled back.
func (a *Account) Update() *AccountUpdateOne {
	return NewAccountClient(a.config).UpdateOne(a)
}

// Unwrap unwraps the Account entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (a *Account) Unwrap() *Account {
	_tx, ok := a.config.driver.(*txDriver)
	if !ok {
		panic("ent: Account is not a transactional entity")
	}
	a.config.driver = _tx.drv
	return a
}

// String implements the fmt.Stringer.
func (a *Account) String() string {
	var builder strings.Builder
	builder.WriteString("Account(")
	builder.WriteString(fmt.Sprintf("id=%v, ", a.ID))
	builder.WriteString("date_created=")
	builder.WriteString(a.DateCreated.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("alias=")
	builder.WriteString(a.Alias)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(a.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(a.Description)
	builder.WriteString(", ")
	builder.WriteString("owner_id=")
	builder.WriteString(fmt.Sprintf("%v", a.OwnerID))
	builder.WriteString(", ")
	builder.WriteString("silo_id=")
	builder.WriteString(fmt.Sprintf("%v", a.SiloID))
	builder.WriteString(", ")
	builder.WriteString("role=")
	builder.WriteString(fmt.Sprintf("%v", a.Role))
	builder.WriteByte(')')
	return builder.String()
}

// Accounts is a parsable slice of Account.
type Accounts []*Account
