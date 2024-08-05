// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/invitation"
	"khepri.dev/horus/ent/predicate"
	"khepri.dev/horus/ent/silo"
	"khepri.dev/horus/ent/team"
)

// SiloQuery is the builder for querying Silo entities.
type SiloQuery struct {
	config
	ctx             *QueryContext
	order           []silo.OrderOption
	inters          []Interceptor
	predicates      []predicate.Silo
	withAccounts    *AccountQuery
	withTeams       *TeamQuery
	withInvitations *InvitationQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the SiloQuery builder.
func (sq *SiloQuery) Where(ps ...predicate.Silo) *SiloQuery {
	sq.predicates = append(sq.predicates, ps...)
	return sq
}

// Limit the number of records to be returned by this query.
func (sq *SiloQuery) Limit(limit int) *SiloQuery {
	sq.ctx.Limit = &limit
	return sq
}

// Offset to start from.
func (sq *SiloQuery) Offset(offset int) *SiloQuery {
	sq.ctx.Offset = &offset
	return sq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (sq *SiloQuery) Unique(unique bool) *SiloQuery {
	sq.ctx.Unique = &unique
	return sq
}

// Order specifies how the records should be ordered.
func (sq *SiloQuery) Order(o ...silo.OrderOption) *SiloQuery {
	sq.order = append(sq.order, o...)
	return sq
}

// QueryAccounts chains the current query on the "accounts" edge.
func (sq *SiloQuery) QueryAccounts() *AccountQuery {
	query := (&AccountClient{config: sq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := sq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(silo.Table, silo.FieldID, selector),
			sqlgraph.To(account.Table, account.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, silo.AccountsTable, silo.AccountsColumn),
		)
		fromU = sqlgraph.SetNeighbors(sq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryTeams chains the current query on the "teams" edge.
func (sq *SiloQuery) QueryTeams() *TeamQuery {
	query := (&TeamClient{config: sq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := sq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(silo.Table, silo.FieldID, selector),
			sqlgraph.To(team.Table, team.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, silo.TeamsTable, silo.TeamsColumn),
		)
		fromU = sqlgraph.SetNeighbors(sq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryInvitations chains the current query on the "invitations" edge.
func (sq *SiloQuery) QueryInvitations() *InvitationQuery {
	query := (&InvitationClient{config: sq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := sq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(silo.Table, silo.FieldID, selector),
			sqlgraph.To(invitation.Table, invitation.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, silo.InvitationsTable, silo.InvitationsColumn),
		)
		fromU = sqlgraph.SetNeighbors(sq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Silo entity from the query.
// Returns a *NotFoundError when no Silo was found.
func (sq *SiloQuery) First(ctx context.Context) (*Silo, error) {
	nodes, err := sq.Limit(1).All(setContextOp(ctx, sq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{silo.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (sq *SiloQuery) FirstX(ctx context.Context) *Silo {
	node, err := sq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Silo ID from the query.
// Returns a *NotFoundError when no Silo ID was found.
func (sq *SiloQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = sq.Limit(1).IDs(setContextOp(ctx, sq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{silo.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (sq *SiloQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := sq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Silo entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Silo entity is found.
// Returns a *NotFoundError when no Silo entities are found.
func (sq *SiloQuery) Only(ctx context.Context) (*Silo, error) {
	nodes, err := sq.Limit(2).All(setContextOp(ctx, sq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{silo.Label}
	default:
		return nil, &NotSingularError{silo.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (sq *SiloQuery) OnlyX(ctx context.Context) *Silo {
	node, err := sq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Silo ID in the query.
// Returns a *NotSingularError when more than one Silo ID is found.
// Returns a *NotFoundError when no entities are found.
func (sq *SiloQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = sq.Limit(2).IDs(setContextOp(ctx, sq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{silo.Label}
	default:
		err = &NotSingularError{silo.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (sq *SiloQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := sq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Silos.
func (sq *SiloQuery) All(ctx context.Context) ([]*Silo, error) {
	ctx = setContextOp(ctx, sq.ctx, ent.OpQueryAll)
	if err := sq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Silo, *SiloQuery]()
	return withInterceptors[[]*Silo](ctx, sq, qr, sq.inters)
}

// AllX is like All, but panics if an error occurs.
func (sq *SiloQuery) AllX(ctx context.Context) []*Silo {
	nodes, err := sq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Silo IDs.
func (sq *SiloQuery) IDs(ctx context.Context) (ids []uuid.UUID, err error) {
	if sq.ctx.Unique == nil && sq.path != nil {
		sq.Unique(true)
	}
	ctx = setContextOp(ctx, sq.ctx, ent.OpQueryIDs)
	if err = sq.Select(silo.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (sq *SiloQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := sq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (sq *SiloQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, sq.ctx, ent.OpQueryCount)
	if err := sq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, sq, querierCount[*SiloQuery](), sq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (sq *SiloQuery) CountX(ctx context.Context) int {
	count, err := sq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (sq *SiloQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, sq.ctx, ent.OpQueryExist)
	switch _, err := sq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (sq *SiloQuery) ExistX(ctx context.Context) bool {
	exist, err := sq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the SiloQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (sq *SiloQuery) Clone() *SiloQuery {
	if sq == nil {
		return nil
	}
	return &SiloQuery{
		config:          sq.config,
		ctx:             sq.ctx.Clone(),
		order:           append([]silo.OrderOption{}, sq.order...),
		inters:          append([]Interceptor{}, sq.inters...),
		predicates:      append([]predicate.Silo{}, sq.predicates...),
		withAccounts:    sq.withAccounts.Clone(),
		withTeams:       sq.withTeams.Clone(),
		withInvitations: sq.withInvitations.Clone(),
		// clone intermediate query.
		sql:  sq.sql.Clone(),
		path: sq.path,
	}
}

// WithAccounts tells the query-builder to eager-load the nodes that are connected to
// the "accounts" edge. The optional arguments are used to configure the query builder of the edge.
func (sq *SiloQuery) WithAccounts(opts ...func(*AccountQuery)) *SiloQuery {
	query := (&AccountClient{config: sq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	sq.withAccounts = query
	return sq
}

// WithTeams tells the query-builder to eager-load the nodes that are connected to
// the "teams" edge. The optional arguments are used to configure the query builder of the edge.
func (sq *SiloQuery) WithTeams(opts ...func(*TeamQuery)) *SiloQuery {
	query := (&TeamClient{config: sq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	sq.withTeams = query
	return sq
}

// WithInvitations tells the query-builder to eager-load the nodes that are connected to
// the "invitations" edge. The optional arguments are used to configure the query builder of the edge.
func (sq *SiloQuery) WithInvitations(opts ...func(*InvitationQuery)) *SiloQuery {
	query := (&InvitationClient{config: sq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	sq.withInvitations = query
	return sq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		DateCreated time.Time `json:"date_created,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Silo.Query().
//		GroupBy(silo.FieldDateCreated).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (sq *SiloQuery) GroupBy(field string, fields ...string) *SiloGroupBy {
	sq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &SiloGroupBy{build: sq}
	grbuild.flds = &sq.ctx.Fields
	grbuild.label = silo.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		DateCreated time.Time `json:"date_created,omitempty"`
//	}
//
//	client.Silo.Query().
//		Select(silo.FieldDateCreated).
//		Scan(ctx, &v)
func (sq *SiloQuery) Select(fields ...string) *SiloSelect {
	sq.ctx.Fields = append(sq.ctx.Fields, fields...)
	sbuild := &SiloSelect{SiloQuery: sq}
	sbuild.label = silo.Label
	sbuild.flds, sbuild.scan = &sq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a SiloSelect configured with the given aggregations.
func (sq *SiloQuery) Aggregate(fns ...AggregateFunc) *SiloSelect {
	return sq.Select().Aggregate(fns...)
}

func (sq *SiloQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range sq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, sq); err != nil {
				return err
			}
		}
	}
	for _, f := range sq.ctx.Fields {
		if !silo.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if sq.path != nil {
		prev, err := sq.path(ctx)
		if err != nil {
			return err
		}
		sq.sql = prev
	}
	return nil
}

func (sq *SiloQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Silo, error) {
	var (
		nodes       = []*Silo{}
		_spec       = sq.querySpec()
		loadedTypes = [3]bool{
			sq.withAccounts != nil,
			sq.withTeams != nil,
			sq.withInvitations != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Silo).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Silo{config: sq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, sq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := sq.withAccounts; query != nil {
		if err := sq.loadAccounts(ctx, query, nodes,
			func(n *Silo) { n.Edges.Accounts = []*Account{} },
			func(n *Silo, e *Account) { n.Edges.Accounts = append(n.Edges.Accounts, e) }); err != nil {
			return nil, err
		}
	}
	if query := sq.withTeams; query != nil {
		if err := sq.loadTeams(ctx, query, nodes,
			func(n *Silo) { n.Edges.Teams = []*Team{} },
			func(n *Silo, e *Team) { n.Edges.Teams = append(n.Edges.Teams, e) }); err != nil {
			return nil, err
		}
	}
	if query := sq.withInvitations; query != nil {
		if err := sq.loadInvitations(ctx, query, nodes,
			func(n *Silo) { n.Edges.Invitations = []*Invitation{} },
			func(n *Silo, e *Invitation) { n.Edges.Invitations = append(n.Edges.Invitations, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (sq *SiloQuery) loadAccounts(ctx context.Context, query *AccountQuery, nodes []*Silo, init func(*Silo), assign func(*Silo, *Account)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*Silo)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(account.FieldSiloID)
	}
	query.Where(predicate.Account(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(silo.AccountsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.SiloID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "silo_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (sq *SiloQuery) loadTeams(ctx context.Context, query *TeamQuery, nodes []*Silo, init func(*Silo), assign func(*Silo, *Team)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*Silo)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(team.FieldSiloID)
	}
	query.Where(predicate.Team(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(silo.TeamsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.SiloID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "silo_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (sq *SiloQuery) loadInvitations(ctx context.Context, query *InvitationQuery, nodes []*Silo, init func(*Silo), assign func(*Silo, *Invitation)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*Silo)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Invitation(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(silo.InvitationsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.silo_invitations
		if fk == nil {
			return fmt.Errorf(`foreign-key "silo_invitations" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "silo_invitations" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (sq *SiloQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := sq.querySpec()
	_spec.Node.Columns = sq.ctx.Fields
	if len(sq.ctx.Fields) > 0 {
		_spec.Unique = sq.ctx.Unique != nil && *sq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, sq.driver, _spec)
}

func (sq *SiloQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(silo.Table, silo.Columns, sqlgraph.NewFieldSpec(silo.FieldID, field.TypeUUID))
	_spec.From = sq.sql
	if unique := sq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if sq.path != nil {
		_spec.Unique = true
	}
	if fields := sq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, silo.FieldID)
		for i := range fields {
			if fields[i] != silo.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := sq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := sq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := sq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := sq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (sq *SiloQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(sq.driver.Dialect())
	t1 := builder.Table(silo.Table)
	columns := sq.ctx.Fields
	if len(columns) == 0 {
		columns = silo.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if sq.sql != nil {
		selector = sq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if sq.ctx.Unique != nil && *sq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range sq.predicates {
		p(selector)
	}
	for _, p := range sq.order {
		p(selector)
	}
	if offset := sq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := sq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// SiloGroupBy is the group-by builder for Silo entities.
type SiloGroupBy struct {
	selector
	build *SiloQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (sgb *SiloGroupBy) Aggregate(fns ...AggregateFunc) *SiloGroupBy {
	sgb.fns = append(sgb.fns, fns...)
	return sgb
}

// Scan applies the selector query and scans the result into the given value.
func (sgb *SiloGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, sgb.build.ctx, ent.OpQueryGroupBy)
	if err := sgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*SiloQuery, *SiloGroupBy](ctx, sgb.build, sgb, sgb.build.inters, v)
}

func (sgb *SiloGroupBy) sqlScan(ctx context.Context, root *SiloQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(sgb.fns))
	for _, fn := range sgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*sgb.flds)+len(sgb.fns))
		for _, f := range *sgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*sgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := sgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// SiloSelect is the builder for selecting fields of Silo entities.
type SiloSelect struct {
	*SiloQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ss *SiloSelect) Aggregate(fns ...AggregateFunc) *SiloSelect {
	ss.fns = append(ss.fns, fns...)
	return ss
}

// Scan applies the selector query and scans the result into the given value.
func (ss *SiloSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ss.ctx, ent.OpQuerySelect)
	if err := ss.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*SiloQuery, *SiloSelect](ctx, ss.SiloQuery, ss, ss.inters, v)
}

func (ss *SiloSelect) sqlScan(ctx context.Context, root *SiloQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ss.fns))
	for _, fn := range ss.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ss.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
