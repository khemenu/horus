// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"khepri.dev/horus/ent/conf"
)

// Conf is the model entity for the Conf schema.
type Conf struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// DateCreated holds the value of the "date_created" field.
	DateCreated time.Time `json:"date_created,omitempty"`
	// Value holds the value of the "value" field.
	Value string `json:"value,omitempty"`
	// DateUpdated holds the value of the "date_updated" field.
	DateUpdated  time.Time `json:"date_updated,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Conf) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case conf.FieldID, conf.FieldValue:
			values[i] = new(sql.NullString)
		case conf.FieldDateCreated, conf.FieldDateUpdated:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Conf fields.
func (c *Conf) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case conf.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				c.ID = value.String
			}
		case conf.FieldDateCreated:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field date_created", values[i])
			} else if value.Valid {
				c.DateCreated = value.Time
			}
		case conf.FieldValue:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field value", values[i])
			} else if value.Valid {
				c.Value = value.String
			}
		case conf.FieldDateUpdated:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field date_updated", values[i])
			} else if value.Valid {
				c.DateUpdated = value.Time
			}
		default:
			c.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// GetValue returns the ent.Value that was dynamically selected and assigned to the Conf.
// This includes values selected through modifiers, order, etc.
func (c *Conf) GetValue(name string) (ent.Value, error) {
	return c.selectValues.Get(name)
}

// Update returns a builder for updating this Conf.
// Note that you need to call Conf.Unwrap() before calling this method if this Conf
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Conf) Update() *ConfUpdateOne {
	return NewConfClient(c.config).UpdateOne(c)
}

// Unwrap unwraps the Conf entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Conf) Unwrap() *Conf {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Conf is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Conf) String() string {
	var builder strings.Builder
	builder.WriteString("Conf(")
	builder.WriteString(fmt.Sprintf("id=%v, ", c.ID))
	builder.WriteString("date_created=")
	builder.WriteString(c.DateCreated.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("value=")
	builder.WriteString(c.Value)
	builder.WriteString(", ")
	builder.WriteString("date_updated=")
	builder.WriteString(c.DateUpdated.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Confs is a parsable slice of Conf.
type Confs []*Conf
