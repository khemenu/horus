// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"khepri.dev/horus/ent/user"
)

// User is the model entity for the User schema.
type User struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Alias holds the value of the "alias" field.
	Alias string `json:"alias,omitempty"`
	// DateCreated holds the value of the "date_created" field.
	DateCreated time.Time `json:"date_created,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the UserQuery when eager-loading is set.
	Edges         UserEdges `json:"edges"`
	user_children *uuid.UUID
	selectValues  sql.SelectValues
}

// UserEdges holds the relations/edges for other nodes in the graph.
type UserEdges struct {
	// Parent holds the value of the parent edge.
	Parent *User `json:"parent,omitempty"`
	// Children holds the value of the children edge.
	Children []*User `json:"children,omitempty"`
	// Identities holds the value of the identities edge.
	Identities []*Identity `json:"identities,omitempty"`
	// Accounts holds the value of the accounts edge.
	Accounts []*Account `json:"accounts,omitempty"`
	// Tokens holds the value of the tokens edge.
	Tokens []*Token `json:"tokens,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [5]bool
}

// ParentOrErr returns the Parent value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e UserEdges) ParentOrErr() (*User, error) {
	if e.loadedTypes[0] {
		if e.Parent == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.Parent, nil
	}
	return nil, &NotLoadedError{edge: "parent"}
}

// ChildrenOrErr returns the Children value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) ChildrenOrErr() ([]*User, error) {
	if e.loadedTypes[1] {
		return e.Children, nil
	}
	return nil, &NotLoadedError{edge: "children"}
}

// IdentitiesOrErr returns the Identities value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) IdentitiesOrErr() ([]*Identity, error) {
	if e.loadedTypes[2] {
		return e.Identities, nil
	}
	return nil, &NotLoadedError{edge: "identities"}
}

// AccountsOrErr returns the Accounts value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) AccountsOrErr() ([]*Account, error) {
	if e.loadedTypes[3] {
		return e.Accounts, nil
	}
	return nil, &NotLoadedError{edge: "accounts"}
}

// TokensOrErr returns the Tokens value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) TokensOrErr() ([]*Token, error) {
	if e.loadedTypes[4] {
		return e.Tokens, nil
	}
	return nil, &NotLoadedError{edge: "tokens"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*User) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case user.FieldAlias:
			values[i] = new(sql.NullString)
		case user.FieldDateCreated:
			values[i] = new(sql.NullTime)
		case user.FieldID:
			values[i] = new(uuid.UUID)
		case user.ForeignKeys[0]: // user_children
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the User fields.
func (u *User) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case user.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				u.ID = *value
			}
		case user.FieldAlias:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field alias", values[i])
			} else if value.Valid {
				u.Alias = value.String
			}
		case user.FieldDateCreated:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field date_created", values[i])
			} else if value.Valid {
				u.DateCreated = value.Time
			}
		case user.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field user_children", values[i])
			} else if value.Valid {
				u.user_children = new(uuid.UUID)
				*u.user_children = *value.S.(*uuid.UUID)
			}
		default:
			u.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the User.
// This includes values selected through modifiers, order, etc.
func (u *User) Value(name string) (ent.Value, error) {
	return u.selectValues.Get(name)
}

// QueryParent queries the "parent" edge of the User entity.
func (u *User) QueryParent() *UserQuery {
	return NewUserClient(u.config).QueryParent(u)
}

// QueryChildren queries the "children" edge of the User entity.
func (u *User) QueryChildren() *UserQuery {
	return NewUserClient(u.config).QueryChildren(u)
}

// QueryIdentities queries the "identities" edge of the User entity.
func (u *User) QueryIdentities() *IdentityQuery {
	return NewUserClient(u.config).QueryIdentities(u)
}

// QueryAccounts queries the "accounts" edge of the User entity.
func (u *User) QueryAccounts() *AccountQuery {
	return NewUserClient(u.config).QueryAccounts(u)
}

// QueryTokens queries the "tokens" edge of the User entity.
func (u *User) QueryTokens() *TokenQuery {
	return NewUserClient(u.config).QueryTokens(u)
}

// Update returns a builder for updating this User.
// Note that you need to call User.Unwrap() before calling this method if this User
// was returned from a transaction, and the transaction was committed or rolled back.
func (u *User) Update() *UserUpdateOne {
	return NewUserClient(u.config).UpdateOne(u)
}

// Unwrap unwraps the User entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (u *User) Unwrap() *User {
	_tx, ok := u.config.driver.(*txDriver)
	if !ok {
		panic("ent: User is not a transactional entity")
	}
	u.config.driver = _tx.drv
	return u
}

// String implements the fmt.Stringer.
func (u *User) String() string {
	var builder strings.Builder
	builder.WriteString("User(")
	builder.WriteString(fmt.Sprintf("id=%v, ", u.ID))
	builder.WriteString("alias=")
	builder.WriteString(u.Alias)
	builder.WriteString(", ")
	builder.WriteString("date_created=")
	builder.WriteString(u.DateCreated.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Users is a parsable slice of User.
type Users []*User
