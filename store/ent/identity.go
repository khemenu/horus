// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"khepri.dev/horus"
	"khepri.dev/horus/store/ent/identity"
	"khepri.dev/horus/store/ent/user"
)

// Identity is the model entity for the Identity schema.
type Identity struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// OwnerID holds the value of the "owner_id" field.
	OwnerID uuid.UUID `json:"owner_id,omitempty"`
	// Kind holds the value of the "kind" field.
	Kind horus.IdentityKind `json:"kind,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// VerifiedBy holds the value of the "verified_by" field.
	VerifiedBy horus.Verifier `json:"verified_by,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the IdentityQuery when eager-loading is set.
	Edges        IdentityEdges `json:"edges"`
	selectValues sql.SelectValues
}

// IdentityEdges holds the relations/edges for other nodes in the graph.
type IdentityEdges struct {
	// Owner holds the value of the owner edge.
	Owner *User `json:"owner,omitempty"`
	// Member holds the value of the member edge.
	Member []*Member `json:"member,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// OwnerOrErr returns the Owner value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e IdentityEdges) OwnerOrErr() (*User, error) {
	if e.loadedTypes[0] {
		if e.Owner == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.Owner, nil
	}
	return nil, &NotLoadedError{edge: "owner"}
}

// MemberOrErr returns the Member value or an error if the edge
// was not loaded in eager-loading.
func (e IdentityEdges) MemberOrErr() ([]*Member, error) {
	if e.loadedTypes[1] {
		return e.Member, nil
	}
	return nil, &NotLoadedError{edge: "member"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Identity) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case identity.FieldID, identity.FieldKind, identity.FieldName, identity.FieldVerifiedBy:
			values[i] = new(sql.NullString)
		case identity.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		case identity.FieldOwnerID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Identity fields.
func (i *Identity) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for j := range columns {
		switch columns[j] {
		case identity.FieldID:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[j])
			} else if value.Valid {
				i.ID = value.String
			}
		case identity.FieldOwnerID:
			if value, ok := values[j].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field owner_id", values[j])
			} else if value != nil {
				i.OwnerID = *value
			}
		case identity.FieldKind:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field kind", values[j])
			} else if value.Valid {
				i.Kind = horus.IdentityKind(value.String)
			}
		case identity.FieldName:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[j])
			} else if value.Valid {
				i.Name = value.String
			}
		case identity.FieldVerifiedBy:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field verified_by", values[j])
			} else if value.Valid {
				i.VerifiedBy = horus.Verifier(value.String)
			}
		case identity.FieldCreatedAt:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[j])
			} else if value.Valid {
				i.CreatedAt = value.Time
			}
		default:
			i.selectValues.Set(columns[j], values[j])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Identity.
// This includes values selected through modifiers, order, etc.
func (i *Identity) Value(name string) (ent.Value, error) {
	return i.selectValues.Get(name)
}

// QueryOwner queries the "owner" edge of the Identity entity.
func (i *Identity) QueryOwner() *UserQuery {
	return NewIdentityClient(i.config).QueryOwner(i)
}

// QueryMember queries the "member" edge of the Identity entity.
func (i *Identity) QueryMember() *MemberQuery {
	return NewIdentityClient(i.config).QueryMember(i)
}

// Update returns a builder for updating this Identity.
// Note that you need to call Identity.Unwrap() before calling this method if this Identity
// was returned from a transaction, and the transaction was committed or rolled back.
func (i *Identity) Update() *IdentityUpdateOne {
	return NewIdentityClient(i.config).UpdateOne(i)
}

// Unwrap unwraps the Identity entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (i *Identity) Unwrap() *Identity {
	_tx, ok := i.config.driver.(*txDriver)
	if !ok {
		panic("ent: Identity is not a transactional entity")
	}
	i.config.driver = _tx.drv
	return i
}

// String implements the fmt.Stringer.
func (i *Identity) String() string {
	var builder strings.Builder
	builder.WriteString("Identity(")
	builder.WriteString(fmt.Sprintf("id=%v, ", i.ID))
	builder.WriteString("owner_id=")
	builder.WriteString(fmt.Sprintf("%v", i.OwnerID))
	builder.WriteString(", ")
	builder.WriteString("kind=")
	builder.WriteString(fmt.Sprintf("%v", i.Kind))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(i.Name)
	builder.WriteString(", ")
	builder.WriteString("verified_by=")
	builder.WriteString(fmt.Sprintf("%v", i.VerifiedBy))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(i.CreatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Identities is a parsable slice of Identity.
type Identities []*Identity
