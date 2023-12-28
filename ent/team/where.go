// Code generated by ent, DO NOT EDIT.

package team

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"khepri.dev/horus/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.Team {
	return predicate.Team(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.Team {
	return predicate.Team(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.Team {
	return predicate.Team(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.Team {
	return predicate.Team(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.Team {
	return predicate.Team(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.Team {
	return predicate.Team(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.Team {
	return predicate.Team(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.Team {
	return predicate.Team(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.Team {
	return predicate.Team(sql.FieldLTE(FieldID, id))
}

// SiloID applies equality check predicate on the "silo_id" field. It's identical to SiloIDEQ.
func SiloID(v uuid.UUID) predicate.Team {
	return predicate.Team(sql.FieldEQ(FieldSiloID, v))
}

// Alias applies equality check predicate on the "alias" field. It's identical to AliasEQ.
func Alias(v string) predicate.Team {
	return predicate.Team(sql.FieldEQ(FieldAlias, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Team {
	return predicate.Team(sql.FieldEQ(FieldName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.Team {
	return predicate.Team(sql.FieldEQ(FieldDescription, v))
}

// CreatedDate applies equality check predicate on the "created_date" field. It's identical to CreatedDateEQ.
func CreatedDate(v time.Time) predicate.Team {
	return predicate.Team(sql.FieldEQ(FieldCreatedDate, v))
}

// SiloIDEQ applies the EQ predicate on the "silo_id" field.
func SiloIDEQ(v uuid.UUID) predicate.Team {
	return predicate.Team(sql.FieldEQ(FieldSiloID, v))
}

// SiloIDNEQ applies the NEQ predicate on the "silo_id" field.
func SiloIDNEQ(v uuid.UUID) predicate.Team {
	return predicate.Team(sql.FieldNEQ(FieldSiloID, v))
}

// SiloIDIn applies the In predicate on the "silo_id" field.
func SiloIDIn(vs ...uuid.UUID) predicate.Team {
	return predicate.Team(sql.FieldIn(FieldSiloID, vs...))
}

// SiloIDNotIn applies the NotIn predicate on the "silo_id" field.
func SiloIDNotIn(vs ...uuid.UUID) predicate.Team {
	return predicate.Team(sql.FieldNotIn(FieldSiloID, vs...))
}

// AliasEQ applies the EQ predicate on the "alias" field.
func AliasEQ(v string) predicate.Team {
	return predicate.Team(sql.FieldEQ(FieldAlias, v))
}

// AliasNEQ applies the NEQ predicate on the "alias" field.
func AliasNEQ(v string) predicate.Team {
	return predicate.Team(sql.FieldNEQ(FieldAlias, v))
}

// AliasIn applies the In predicate on the "alias" field.
func AliasIn(vs ...string) predicate.Team {
	return predicate.Team(sql.FieldIn(FieldAlias, vs...))
}

// AliasNotIn applies the NotIn predicate on the "alias" field.
func AliasNotIn(vs ...string) predicate.Team {
	return predicate.Team(sql.FieldNotIn(FieldAlias, vs...))
}

// AliasGT applies the GT predicate on the "alias" field.
func AliasGT(v string) predicate.Team {
	return predicate.Team(sql.FieldGT(FieldAlias, v))
}

// AliasGTE applies the GTE predicate on the "alias" field.
func AliasGTE(v string) predicate.Team {
	return predicate.Team(sql.FieldGTE(FieldAlias, v))
}

// AliasLT applies the LT predicate on the "alias" field.
func AliasLT(v string) predicate.Team {
	return predicate.Team(sql.FieldLT(FieldAlias, v))
}

// AliasLTE applies the LTE predicate on the "alias" field.
func AliasLTE(v string) predicate.Team {
	return predicate.Team(sql.FieldLTE(FieldAlias, v))
}

// AliasContains applies the Contains predicate on the "alias" field.
func AliasContains(v string) predicate.Team {
	return predicate.Team(sql.FieldContains(FieldAlias, v))
}

// AliasHasPrefix applies the HasPrefix predicate on the "alias" field.
func AliasHasPrefix(v string) predicate.Team {
	return predicate.Team(sql.FieldHasPrefix(FieldAlias, v))
}

// AliasHasSuffix applies the HasSuffix predicate on the "alias" field.
func AliasHasSuffix(v string) predicate.Team {
	return predicate.Team(sql.FieldHasSuffix(FieldAlias, v))
}

// AliasEqualFold applies the EqualFold predicate on the "alias" field.
func AliasEqualFold(v string) predicate.Team {
	return predicate.Team(sql.FieldEqualFold(FieldAlias, v))
}

// AliasContainsFold applies the ContainsFold predicate on the "alias" field.
func AliasContainsFold(v string) predicate.Team {
	return predicate.Team(sql.FieldContainsFold(FieldAlias, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Team {
	return predicate.Team(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Team {
	return predicate.Team(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Team {
	return predicate.Team(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Team {
	return predicate.Team(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Team {
	return predicate.Team(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Team {
	return predicate.Team(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Team {
	return predicate.Team(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Team {
	return predicate.Team(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Team {
	return predicate.Team(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Team {
	return predicate.Team(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Team {
	return predicate.Team(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Team {
	return predicate.Team(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Team {
	return predicate.Team(sql.FieldContainsFold(FieldName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.Team {
	return predicate.Team(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.Team {
	return predicate.Team(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.Team {
	return predicate.Team(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.Team {
	return predicate.Team(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.Team {
	return predicate.Team(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.Team {
	return predicate.Team(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.Team {
	return predicate.Team(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.Team {
	return predicate.Team(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.Team {
	return predicate.Team(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.Team {
	return predicate.Team(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.Team {
	return predicate.Team(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.Team {
	return predicate.Team(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.Team {
	return predicate.Team(sql.FieldContainsFold(FieldDescription, v))
}

// InterVisibilityEQ applies the EQ predicate on the "inter_visibility" field.
func InterVisibilityEQ(v InterVisibility) predicate.Team {
	return predicate.Team(sql.FieldEQ(FieldInterVisibility, v))
}

// InterVisibilityNEQ applies the NEQ predicate on the "inter_visibility" field.
func InterVisibilityNEQ(v InterVisibility) predicate.Team {
	return predicate.Team(sql.FieldNEQ(FieldInterVisibility, v))
}

// InterVisibilityIn applies the In predicate on the "inter_visibility" field.
func InterVisibilityIn(vs ...InterVisibility) predicate.Team {
	return predicate.Team(sql.FieldIn(FieldInterVisibility, vs...))
}

// InterVisibilityNotIn applies the NotIn predicate on the "inter_visibility" field.
func InterVisibilityNotIn(vs ...InterVisibility) predicate.Team {
	return predicate.Team(sql.FieldNotIn(FieldInterVisibility, vs...))
}

// IntraVisibilityEQ applies the EQ predicate on the "intra_visibility" field.
func IntraVisibilityEQ(v IntraVisibility) predicate.Team {
	return predicate.Team(sql.FieldEQ(FieldIntraVisibility, v))
}

// IntraVisibilityNEQ applies the NEQ predicate on the "intra_visibility" field.
func IntraVisibilityNEQ(v IntraVisibility) predicate.Team {
	return predicate.Team(sql.FieldNEQ(FieldIntraVisibility, v))
}

// IntraVisibilityIn applies the In predicate on the "intra_visibility" field.
func IntraVisibilityIn(vs ...IntraVisibility) predicate.Team {
	return predicate.Team(sql.FieldIn(FieldIntraVisibility, vs...))
}

// IntraVisibilityNotIn applies the NotIn predicate on the "intra_visibility" field.
func IntraVisibilityNotIn(vs ...IntraVisibility) predicate.Team {
	return predicate.Team(sql.FieldNotIn(FieldIntraVisibility, vs...))
}

// CreatedDateEQ applies the EQ predicate on the "created_date" field.
func CreatedDateEQ(v time.Time) predicate.Team {
	return predicate.Team(sql.FieldEQ(FieldCreatedDate, v))
}

// CreatedDateNEQ applies the NEQ predicate on the "created_date" field.
func CreatedDateNEQ(v time.Time) predicate.Team {
	return predicate.Team(sql.FieldNEQ(FieldCreatedDate, v))
}

// CreatedDateIn applies the In predicate on the "created_date" field.
func CreatedDateIn(vs ...time.Time) predicate.Team {
	return predicate.Team(sql.FieldIn(FieldCreatedDate, vs...))
}

// CreatedDateNotIn applies the NotIn predicate on the "created_date" field.
func CreatedDateNotIn(vs ...time.Time) predicate.Team {
	return predicate.Team(sql.FieldNotIn(FieldCreatedDate, vs...))
}

// CreatedDateGT applies the GT predicate on the "created_date" field.
func CreatedDateGT(v time.Time) predicate.Team {
	return predicate.Team(sql.FieldGT(FieldCreatedDate, v))
}

// CreatedDateGTE applies the GTE predicate on the "created_date" field.
func CreatedDateGTE(v time.Time) predicate.Team {
	return predicate.Team(sql.FieldGTE(FieldCreatedDate, v))
}

// CreatedDateLT applies the LT predicate on the "created_date" field.
func CreatedDateLT(v time.Time) predicate.Team {
	return predicate.Team(sql.FieldLT(FieldCreatedDate, v))
}

// CreatedDateLTE applies the LTE predicate on the "created_date" field.
func CreatedDateLTE(v time.Time) predicate.Team {
	return predicate.Team(sql.FieldLTE(FieldCreatedDate, v))
}

// HasSilo applies the HasEdge predicate on the "silo" edge.
func HasSilo() predicate.Team {
	return predicate.Team(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, SiloTable, SiloColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSiloWith applies the HasEdge predicate on the "silo" edge with a given conditions (other predicates).
func HasSiloWith(preds ...predicate.Silo) predicate.Team {
	return predicate.Team(func(s *sql.Selector) {
		step := newSiloStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasMembers applies the HasEdge predicate on the "members" edge.
func HasMembers() predicate.Team {
	return predicate.Team(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, MembersTable, MembersColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMembersWith applies the HasEdge predicate on the "members" edge with a given conditions (other predicates).
func HasMembersWith(preds ...predicate.Membership) predicate.Team {
	return predicate.Team(func(s *sql.Selector) {
		step := newMembersStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Team) predicate.Team {
	return predicate.Team(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Team) predicate.Team {
	return predicate.Team(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Team) predicate.Team {
	return predicate.Team(sql.NotPredicates(p))
}
