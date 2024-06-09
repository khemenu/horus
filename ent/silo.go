// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"khepri.dev/horus/ent/silo"
)

// Silo is the model entity for the Silo schema.
type Silo struct {
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
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the SiloQuery when eager-loading is set.
	Edges        SiloEdges `json:"edges"`
	selectValues sql.SelectValues
}

// SiloEdges holds the relations/edges for other nodes in the graph.
type SiloEdges struct {
	// Members holds the value of the members edge.
	Members []*Account `json:"members,omitempty"`
	// Teams holds the value of the teams edge.
	Teams []*Team `json:"teams,omitempty"`
	// Invitations holds the value of the invitations edge.
	Invitations []*Invitation `json:"invitations,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// MembersOrErr returns the Members value or an error if the edge
// was not loaded in eager-loading.
func (e SiloEdges) MembersOrErr() ([]*Account, error) {
	if e.loadedTypes[0] {
		return e.Members, nil
	}
	return nil, &NotLoadedError{edge: "members"}
}

// TeamsOrErr returns the Teams value or an error if the edge
// was not loaded in eager-loading.
func (e SiloEdges) TeamsOrErr() ([]*Team, error) {
	if e.loadedTypes[1] {
		return e.Teams, nil
	}
	return nil, &NotLoadedError{edge: "teams"}
}

// InvitationsOrErr returns the Invitations value or an error if the edge
// was not loaded in eager-loading.
func (e SiloEdges) InvitationsOrErr() ([]*Invitation, error) {
	if e.loadedTypes[2] {
		return e.Invitations, nil
	}
	return nil, &NotLoadedError{edge: "invitations"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Silo) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case silo.FieldAlias, silo.FieldName, silo.FieldDescription:
			values[i] = new(sql.NullString)
		case silo.FieldDateCreated:
			values[i] = new(sql.NullTime)
		case silo.FieldID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Silo fields.
func (s *Silo) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case silo.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				s.ID = *value
			}
		case silo.FieldDateCreated:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field date_created", values[i])
			} else if value.Valid {
				s.DateCreated = value.Time
			}
		case silo.FieldAlias:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field alias", values[i])
			} else if value.Valid {
				s.Alias = value.String
			}
		case silo.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				s.Name = value.String
			}
		case silo.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				s.Description = value.String
			}
		default:
			s.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Silo.
// This includes values selected through modifiers, order, etc.
func (s *Silo) Value(name string) (ent.Value, error) {
	return s.selectValues.Get(name)
}

// QueryMembers queries the "members" edge of the Silo entity.
func (s *Silo) QueryMembers() *AccountQuery {
	return NewSiloClient(s.config).QueryMembers(s)
}

// QueryTeams queries the "teams" edge of the Silo entity.
func (s *Silo) QueryTeams() *TeamQuery {
	return NewSiloClient(s.config).QueryTeams(s)
}

// QueryInvitations queries the "invitations" edge of the Silo entity.
func (s *Silo) QueryInvitations() *InvitationQuery {
	return NewSiloClient(s.config).QueryInvitations(s)
}

// Update returns a builder for updating this Silo.
// Note that you need to call Silo.Unwrap() before calling this method if this Silo
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Silo) Update() *SiloUpdateOne {
	return NewSiloClient(s.config).UpdateOne(s)
}

// Unwrap unwraps the Silo entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (s *Silo) Unwrap() *Silo {
	_tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("ent: Silo is not a transactional entity")
	}
	s.config.driver = _tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *Silo) String() string {
	var builder strings.Builder
	builder.WriteString("Silo(")
	builder.WriteString(fmt.Sprintf("id=%v, ", s.ID))
	builder.WriteString("date_created=")
	builder.WriteString(s.DateCreated.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("alias=")
	builder.WriteString(s.Alias)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(s.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(s.Description)
	builder.WriteByte(')')
	return builder.String()
}

// Silos is a parsable slice of Silo.
type Silos []*Silo
