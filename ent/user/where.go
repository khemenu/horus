// Code generated by ent, DO NOT EDIT.

package user

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"khepri.dev/horus/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.User {
	return predicate.User(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldLTE(FieldID, id))
}

// CreatedDate applies equality check predicate on the "created_date" field. It's identical to CreatedDateEQ.
func CreatedDate(v time.Time) predicate.User {
	return predicate.User(sql.FieldEQ(FieldCreatedDate, v))
}

// CreatedDateEQ applies the EQ predicate on the "created_date" field.
func CreatedDateEQ(v time.Time) predicate.User {
	return predicate.User(sql.FieldEQ(FieldCreatedDate, v))
}

// CreatedDateNEQ applies the NEQ predicate on the "created_date" field.
func CreatedDateNEQ(v time.Time) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldCreatedDate, v))
}

// CreatedDateIn applies the In predicate on the "created_date" field.
func CreatedDateIn(vs ...time.Time) predicate.User {
	return predicate.User(sql.FieldIn(FieldCreatedDate, vs...))
}

// CreatedDateNotIn applies the NotIn predicate on the "created_date" field.
func CreatedDateNotIn(vs ...time.Time) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldCreatedDate, vs...))
}

// CreatedDateGT applies the GT predicate on the "created_date" field.
func CreatedDateGT(v time.Time) predicate.User {
	return predicate.User(sql.FieldGT(FieldCreatedDate, v))
}

// CreatedDateGTE applies the GTE predicate on the "created_date" field.
func CreatedDateGTE(v time.Time) predicate.User {
	return predicate.User(sql.FieldGTE(FieldCreatedDate, v))
}

// CreatedDateLT applies the LT predicate on the "created_date" field.
func CreatedDateLT(v time.Time) predicate.User {
	return predicate.User(sql.FieldLT(FieldCreatedDate, v))
}

// CreatedDateLTE applies the LTE predicate on the "created_date" field.
func CreatedDateLTE(v time.Time) predicate.User {
	return predicate.User(sql.FieldLTE(FieldCreatedDate, v))
}

// HasIdentities applies the HasEdge predicate on the "identities" edge.
func HasIdentities() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, IdentitiesTable, IdentitiesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasIdentitiesWith applies the HasEdge predicate on the "identities" edge with a given conditions (other predicates).
func HasIdentitiesWith(preds ...predicate.Identity) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := newIdentitiesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasAccounts applies the HasEdge predicate on the "accounts" edge.
func HasAccounts() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, AccountsTable, AccountsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasAccountsWith applies the HasEdge predicate on the "accounts" edge with a given conditions (other predicates).
func HasAccountsWith(preds ...predicate.Account) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := newAccountsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasTokens applies the HasEdge predicate on the "tokens" edge.
func HasTokens() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, TokensTable, TokensColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTokensWith applies the HasEdge predicate on the "tokens" edge with a given conditions (other predicates).
func HasTokensWith(preds ...predicate.Token) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := newTokensStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.User) predicate.User {
	return predicate.User(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.User) predicate.User {
	return predicate.User(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.User) predicate.User {
	return predicate.User(sql.NotPredicates(p))
}
