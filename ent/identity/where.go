// Code generated by ent, DO NOT EDIT.

package identity

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"khepri.dev/horus/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.Identity {
	return predicate.Identity(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.Identity {
	return predicate.Identity(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.Identity {
	return predicate.Identity(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.Identity {
	return predicate.Identity(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.Identity {
	return predicate.Identity(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.Identity {
	return predicate.Identity(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.Identity {
	return predicate.Identity(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.Identity {
	return predicate.Identity(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.Identity {
	return predicate.Identity(sql.FieldLTE(FieldID, id))
}

// DateCreated applies equality check predicate on the "date_created" field. It's identical to DateCreatedEQ.
func DateCreated(v time.Time) predicate.Identity {
	return predicate.Identity(sql.FieldEQ(FieldDateCreated, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Identity {
	return predicate.Identity(sql.FieldEQ(FieldName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.Identity {
	return predicate.Identity(sql.FieldEQ(FieldDescription, v))
}

// Kind applies equality check predicate on the "kind" field. It's identical to KindEQ.
func Kind(v string) predicate.Identity {
	return predicate.Identity(sql.FieldEQ(FieldKind, v))
}

// Verifier applies equality check predicate on the "verifier" field. It's identical to VerifierEQ.
func Verifier(v string) predicate.Identity {
	return predicate.Identity(sql.FieldEQ(FieldVerifier, v))
}

// DateCreatedEQ applies the EQ predicate on the "date_created" field.
func DateCreatedEQ(v time.Time) predicate.Identity {
	return predicate.Identity(sql.FieldEQ(FieldDateCreated, v))
}

// DateCreatedNEQ applies the NEQ predicate on the "date_created" field.
func DateCreatedNEQ(v time.Time) predicate.Identity {
	return predicate.Identity(sql.FieldNEQ(FieldDateCreated, v))
}

// DateCreatedIn applies the In predicate on the "date_created" field.
func DateCreatedIn(vs ...time.Time) predicate.Identity {
	return predicate.Identity(sql.FieldIn(FieldDateCreated, vs...))
}

// DateCreatedNotIn applies the NotIn predicate on the "date_created" field.
func DateCreatedNotIn(vs ...time.Time) predicate.Identity {
	return predicate.Identity(sql.FieldNotIn(FieldDateCreated, vs...))
}

// DateCreatedGT applies the GT predicate on the "date_created" field.
func DateCreatedGT(v time.Time) predicate.Identity {
	return predicate.Identity(sql.FieldGT(FieldDateCreated, v))
}

// DateCreatedGTE applies the GTE predicate on the "date_created" field.
func DateCreatedGTE(v time.Time) predicate.Identity {
	return predicate.Identity(sql.FieldGTE(FieldDateCreated, v))
}

// DateCreatedLT applies the LT predicate on the "date_created" field.
func DateCreatedLT(v time.Time) predicate.Identity {
	return predicate.Identity(sql.FieldLT(FieldDateCreated, v))
}

// DateCreatedLTE applies the LTE predicate on the "date_created" field.
func DateCreatedLTE(v time.Time) predicate.Identity {
	return predicate.Identity(sql.FieldLTE(FieldDateCreated, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Identity {
	return predicate.Identity(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Identity {
	return predicate.Identity(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Identity {
	return predicate.Identity(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Identity {
	return predicate.Identity(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Identity {
	return predicate.Identity(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Identity {
	return predicate.Identity(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Identity {
	return predicate.Identity(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Identity {
	return predicate.Identity(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Identity {
	return predicate.Identity(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Identity {
	return predicate.Identity(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Identity {
	return predicate.Identity(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Identity {
	return predicate.Identity(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Identity {
	return predicate.Identity(sql.FieldContainsFold(FieldName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.Identity {
	return predicate.Identity(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.Identity {
	return predicate.Identity(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.Identity {
	return predicate.Identity(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.Identity {
	return predicate.Identity(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.Identity {
	return predicate.Identity(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.Identity {
	return predicate.Identity(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.Identity {
	return predicate.Identity(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.Identity {
	return predicate.Identity(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.Identity {
	return predicate.Identity(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.Identity {
	return predicate.Identity(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.Identity {
	return predicate.Identity(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.Identity {
	return predicate.Identity(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.Identity {
	return predicate.Identity(sql.FieldContainsFold(FieldDescription, v))
}

// KindEQ applies the EQ predicate on the "kind" field.
func KindEQ(v string) predicate.Identity {
	return predicate.Identity(sql.FieldEQ(FieldKind, v))
}

// KindNEQ applies the NEQ predicate on the "kind" field.
func KindNEQ(v string) predicate.Identity {
	return predicate.Identity(sql.FieldNEQ(FieldKind, v))
}

// KindIn applies the In predicate on the "kind" field.
func KindIn(vs ...string) predicate.Identity {
	return predicate.Identity(sql.FieldIn(FieldKind, vs...))
}

// KindNotIn applies the NotIn predicate on the "kind" field.
func KindNotIn(vs ...string) predicate.Identity {
	return predicate.Identity(sql.FieldNotIn(FieldKind, vs...))
}

// KindGT applies the GT predicate on the "kind" field.
func KindGT(v string) predicate.Identity {
	return predicate.Identity(sql.FieldGT(FieldKind, v))
}

// KindGTE applies the GTE predicate on the "kind" field.
func KindGTE(v string) predicate.Identity {
	return predicate.Identity(sql.FieldGTE(FieldKind, v))
}

// KindLT applies the LT predicate on the "kind" field.
func KindLT(v string) predicate.Identity {
	return predicate.Identity(sql.FieldLT(FieldKind, v))
}

// KindLTE applies the LTE predicate on the "kind" field.
func KindLTE(v string) predicate.Identity {
	return predicate.Identity(sql.FieldLTE(FieldKind, v))
}

// KindContains applies the Contains predicate on the "kind" field.
func KindContains(v string) predicate.Identity {
	return predicate.Identity(sql.FieldContains(FieldKind, v))
}

// KindHasPrefix applies the HasPrefix predicate on the "kind" field.
func KindHasPrefix(v string) predicate.Identity {
	return predicate.Identity(sql.FieldHasPrefix(FieldKind, v))
}

// KindHasSuffix applies the HasSuffix predicate on the "kind" field.
func KindHasSuffix(v string) predicate.Identity {
	return predicate.Identity(sql.FieldHasSuffix(FieldKind, v))
}

// KindEqualFold applies the EqualFold predicate on the "kind" field.
func KindEqualFold(v string) predicate.Identity {
	return predicate.Identity(sql.FieldEqualFold(FieldKind, v))
}

// KindContainsFold applies the ContainsFold predicate on the "kind" field.
func KindContainsFold(v string) predicate.Identity {
	return predicate.Identity(sql.FieldContainsFold(FieldKind, v))
}

// VerifierEQ applies the EQ predicate on the "verifier" field.
func VerifierEQ(v string) predicate.Identity {
	return predicate.Identity(sql.FieldEQ(FieldVerifier, v))
}

// VerifierNEQ applies the NEQ predicate on the "verifier" field.
func VerifierNEQ(v string) predicate.Identity {
	return predicate.Identity(sql.FieldNEQ(FieldVerifier, v))
}

// VerifierIn applies the In predicate on the "verifier" field.
func VerifierIn(vs ...string) predicate.Identity {
	return predicate.Identity(sql.FieldIn(FieldVerifier, vs...))
}

// VerifierNotIn applies the NotIn predicate on the "verifier" field.
func VerifierNotIn(vs ...string) predicate.Identity {
	return predicate.Identity(sql.FieldNotIn(FieldVerifier, vs...))
}

// VerifierGT applies the GT predicate on the "verifier" field.
func VerifierGT(v string) predicate.Identity {
	return predicate.Identity(sql.FieldGT(FieldVerifier, v))
}

// VerifierGTE applies the GTE predicate on the "verifier" field.
func VerifierGTE(v string) predicate.Identity {
	return predicate.Identity(sql.FieldGTE(FieldVerifier, v))
}

// VerifierLT applies the LT predicate on the "verifier" field.
func VerifierLT(v string) predicate.Identity {
	return predicate.Identity(sql.FieldLT(FieldVerifier, v))
}

// VerifierLTE applies the LTE predicate on the "verifier" field.
func VerifierLTE(v string) predicate.Identity {
	return predicate.Identity(sql.FieldLTE(FieldVerifier, v))
}

// VerifierContains applies the Contains predicate on the "verifier" field.
func VerifierContains(v string) predicate.Identity {
	return predicate.Identity(sql.FieldContains(FieldVerifier, v))
}

// VerifierHasPrefix applies the HasPrefix predicate on the "verifier" field.
func VerifierHasPrefix(v string) predicate.Identity {
	return predicate.Identity(sql.FieldHasPrefix(FieldVerifier, v))
}

// VerifierHasSuffix applies the HasSuffix predicate on the "verifier" field.
func VerifierHasSuffix(v string) predicate.Identity {
	return predicate.Identity(sql.FieldHasSuffix(FieldVerifier, v))
}

// VerifierEqualFold applies the EqualFold predicate on the "verifier" field.
func VerifierEqualFold(v string) predicate.Identity {
	return predicate.Identity(sql.FieldEqualFold(FieldVerifier, v))
}

// VerifierContainsFold applies the ContainsFold predicate on the "verifier" field.
func VerifierContainsFold(v string) predicate.Identity {
	return predicate.Identity(sql.FieldContainsFold(FieldVerifier, v))
}

// HasOwner applies the HasEdge predicate on the "owner" edge.
func HasOwner() predicate.Identity {
	return predicate.Identity(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, OwnerTable, OwnerColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasOwnerWith applies the HasEdge predicate on the "owner" edge with a given conditions (other predicates).
func HasOwnerWith(preds ...predicate.User) predicate.Identity {
	return predicate.Identity(func(s *sql.Selector) {
		step := newOwnerStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Identity) predicate.Identity {
	return predicate.Identity(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Identity) predicate.Identity {
	return predicate.Identity(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Identity) predicate.Identity {
	return predicate.Identity(sql.NotPredicates(p))
}
