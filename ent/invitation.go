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
	"khepri.dev/horus/ent/invitation"
	"khepri.dev/horus/ent/silo"
)

// Invitation is the model entity for the Invitation schema.
type Invitation struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Invitee holds the value of the "invitee" field.
	Invitee string `json:"invitee,omitempty"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
	// DateCreated holds the value of the "date_created" field.
	DateCreated time.Time `json:"date_created,omitempty"`
	// DateExpired holds the value of the "date_expired" field.
	DateExpired time.Time `json:"date_expired,omitempty"`
	// DateAccepted holds the value of the "date_accepted" field.
	DateAccepted *time.Time `json:"date_accepted,omitempty"`
	// DateDeclined holds the value of the "date_declined" field.
	DateDeclined *time.Time `json:"date_declined,omitempty"`
	// DateCanceled holds the value of the "date_canceled" field.
	DateCanceled *time.Time `json:"date_canceled,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the InvitationQuery when eager-loading is set.
	Edges               InvitationEdges `json:"edges"`
	account_invitations *uuid.UUID
	silo_invitations    *uuid.UUID
	selectValues        sql.SelectValues
}

// InvitationEdges holds the relations/edges for other nodes in the graph.
type InvitationEdges struct {
	// Silo holds the value of the silo edge.
	Silo *Silo `json:"silo,omitempty"`
	// Inviter holds the value of the inviter edge.
	Inviter *Account `json:"inviter,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// SiloOrErr returns the Silo value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e InvitationEdges) SiloOrErr() (*Silo, error) {
	if e.loadedTypes[0] {
		if e.Silo == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: silo.Label}
		}
		return e.Silo, nil
	}
	return nil, &NotLoadedError{edge: "silo"}
}

// InviterOrErr returns the Inviter value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e InvitationEdges) InviterOrErr() (*Account, error) {
	if e.loadedTypes[1] {
		if e.Inviter == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: account.Label}
		}
		return e.Inviter, nil
	}
	return nil, &NotLoadedError{edge: "inviter"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Invitation) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case invitation.FieldInvitee, invitation.FieldType:
			values[i] = new(sql.NullString)
		case invitation.FieldDateCreated, invitation.FieldDateExpired, invitation.FieldDateAccepted, invitation.FieldDateDeclined, invitation.FieldDateCanceled:
			values[i] = new(sql.NullTime)
		case invitation.FieldID:
			values[i] = new(uuid.UUID)
		case invitation.ForeignKeys[0]: // account_invitations
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		case invitation.ForeignKeys[1]: // silo_invitations
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Invitation fields.
func (i *Invitation) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for j := range columns {
		switch columns[j] {
		case invitation.FieldID:
			if value, ok := values[j].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[j])
			} else if value != nil {
				i.ID = *value
			}
		case invitation.FieldInvitee:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field invitee", values[j])
			} else if value.Valid {
				i.Invitee = value.String
			}
		case invitation.FieldType:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[j])
			} else if value.Valid {
				i.Type = value.String
			}
		case invitation.FieldDateCreated:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field date_created", values[j])
			} else if value.Valid {
				i.DateCreated = value.Time
			}
		case invitation.FieldDateExpired:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field date_expired", values[j])
			} else if value.Valid {
				i.DateExpired = value.Time
			}
		case invitation.FieldDateAccepted:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field date_accepted", values[j])
			} else if value.Valid {
				i.DateAccepted = new(time.Time)
				*i.DateAccepted = value.Time
			}
		case invitation.FieldDateDeclined:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field date_declined", values[j])
			} else if value.Valid {
				i.DateDeclined = new(time.Time)
				*i.DateDeclined = value.Time
			}
		case invitation.FieldDateCanceled:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field date_canceled", values[j])
			} else if value.Valid {
				i.DateCanceled = new(time.Time)
				*i.DateCanceled = value.Time
			}
		case invitation.ForeignKeys[0]:
			if value, ok := values[j].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field account_invitations", values[j])
			} else if value.Valid {
				i.account_invitations = new(uuid.UUID)
				*i.account_invitations = *value.S.(*uuid.UUID)
			}
		case invitation.ForeignKeys[1]:
			if value, ok := values[j].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field silo_invitations", values[j])
			} else if value.Valid {
				i.silo_invitations = new(uuid.UUID)
				*i.silo_invitations = *value.S.(*uuid.UUID)
			}
		default:
			i.selectValues.Set(columns[j], values[j])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Invitation.
// This includes values selected through modifiers, order, etc.
func (i *Invitation) Value(name string) (ent.Value, error) {
	return i.selectValues.Get(name)
}

// QuerySilo queries the "silo" edge of the Invitation entity.
func (i *Invitation) QuerySilo() *SiloQuery {
	return NewInvitationClient(i.config).QuerySilo(i)
}

// QueryInviter queries the "inviter" edge of the Invitation entity.
func (i *Invitation) QueryInviter() *AccountQuery {
	return NewInvitationClient(i.config).QueryInviter(i)
}

// Update returns a builder for updating this Invitation.
// Note that you need to call Invitation.Unwrap() before calling this method if this Invitation
// was returned from a transaction, and the transaction was committed or rolled back.
func (i *Invitation) Update() *InvitationUpdateOne {
	return NewInvitationClient(i.config).UpdateOne(i)
}

// Unwrap unwraps the Invitation entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (i *Invitation) Unwrap() *Invitation {
	_tx, ok := i.config.driver.(*txDriver)
	if !ok {
		panic("ent: Invitation is not a transactional entity")
	}
	i.config.driver = _tx.drv
	return i
}

// String implements the fmt.Stringer.
func (i *Invitation) String() string {
	var builder strings.Builder
	builder.WriteString("Invitation(")
	builder.WriteString(fmt.Sprintf("id=%v, ", i.ID))
	builder.WriteString("invitee=")
	builder.WriteString(i.Invitee)
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(i.Type)
	builder.WriteString(", ")
	builder.WriteString("date_created=")
	builder.WriteString(i.DateCreated.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("date_expired=")
	builder.WriteString(i.DateExpired.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := i.DateAccepted; v != nil {
		builder.WriteString("date_accepted=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := i.DateDeclined; v != nil {
		builder.WriteString("date_declined=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := i.DateCanceled; v != nil {
		builder.WriteString("date_canceled=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteByte(')')
	return builder.String()
}

// Invitations is a parsable slice of Invitation.
type Invitations []*Invitation
