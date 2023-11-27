// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Monitor is an object representing the database table.
type Monitor struct {
	ID            string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	AccountID     string    `boil:"account_id" json:"account_id" toml:"account_id" yaml:"account_id"`
	EndpointURL   string    `boil:"endpoint_url" json:"endpoint_url" toml:"endpoint_url" yaml:"endpoint_url"`
	IsEndpointUp  bool      `boil:"is_endpoint_up" json:"is_endpoint_up" toml:"is_endpoint_up" yaml:"is_endpoint_up"`
	CreatedAt     time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	LastCheckedAt null.Time `boil:"last_checked_at" json:"last_checked_at,omitempty" toml:"last_checked_at" yaml:"last_checked_at,omitempty"`

	R *monitorR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L monitorL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var MonitorColumns = struct {
	ID            string
	AccountID     string
	EndpointURL   string
	IsEndpointUp  string
	CreatedAt     string
	LastCheckedAt string
}{
	ID:            "id",
	AccountID:     "account_id",
	EndpointURL:   "endpoint_url",
	IsEndpointUp:  "is_endpoint_up",
	CreatedAt:     "created_at",
	LastCheckedAt: "last_checked_at",
}

var MonitorTableColumns = struct {
	ID            string
	AccountID     string
	EndpointURL   string
	IsEndpointUp  string
	CreatedAt     string
	LastCheckedAt string
}{
	ID:            "monitors.id",
	AccountID:     "monitors.account_id",
	EndpointURL:   "monitors.endpoint_url",
	IsEndpointUp:  "monitors.is_endpoint_up",
	CreatedAt:     "monitors.created_at",
	LastCheckedAt: "monitors.last_checked_at",
}

// Generated where

type whereHelperbool struct{ field string }

func (w whereHelperbool) EQ(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperbool) NEQ(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperbool) LT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperbool) LTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperbool) GT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperbool) GTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }

var MonitorWhere = struct {
	ID            whereHelperstring
	AccountID     whereHelperstring
	EndpointURL   whereHelperstring
	IsEndpointUp  whereHelperbool
	CreatedAt     whereHelpertime_Time
	LastCheckedAt whereHelpernull_Time
}{
	ID:            whereHelperstring{field: "\"monitors\".\"id\""},
	AccountID:     whereHelperstring{field: "\"monitors\".\"account_id\""},
	EndpointURL:   whereHelperstring{field: "\"monitors\".\"endpoint_url\""},
	IsEndpointUp:  whereHelperbool{field: "\"monitors\".\"is_endpoint_up\""},
	CreatedAt:     whereHelpertime_Time{field: "\"monitors\".\"created_at\""},
	LastCheckedAt: whereHelpernull_Time{field: "\"monitors\".\"last_checked_at\""},
}

// MonitorRels is where relationship names are stored.
var MonitorRels = struct {
	Account string
}{
	Account: "Account",
}

// monitorR is where relationships are stored.
type monitorR struct {
	Account *Account `boil:"Account" json:"Account" toml:"Account" yaml:"Account"`
}

// NewStruct creates a new relationship struct
func (*monitorR) NewStruct() *monitorR {
	return &monitorR{}
}

func (r *monitorR) GetAccount() *Account {
	if r == nil {
		return nil
	}
	return r.Account
}

// monitorL is where Load methods for each relationship are stored.
type monitorL struct{}

var (
	monitorAllColumns            = []string{"id", "account_id", "endpoint_url", "is_endpoint_up", "created_at", "last_checked_at"}
	monitorColumnsWithoutDefault = []string{"id", "account_id", "endpoint_url"}
	monitorColumnsWithDefault    = []string{"is_endpoint_up", "created_at", "last_checked_at"}
	monitorPrimaryKeyColumns     = []string{"id"}
	monitorGeneratedColumns      = []string{}
)

type (
	// MonitorSlice is an alias for a slice of pointers to Monitor.
	// This should almost always be used instead of []Monitor.
	MonitorSlice []*Monitor
	// MonitorHook is the signature for custom Monitor hook methods
	MonitorHook func(context.Context, boil.ContextExecutor, *Monitor) error

	monitorQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	monitorType                 = reflect.TypeOf(&Monitor{})
	monitorMapping              = queries.MakeStructMapping(monitorType)
	monitorPrimaryKeyMapping, _ = queries.BindMapping(monitorType, monitorMapping, monitorPrimaryKeyColumns)
	monitorInsertCacheMut       sync.RWMutex
	monitorInsertCache          = make(map[string]insertCache)
	monitorUpdateCacheMut       sync.RWMutex
	monitorUpdateCache          = make(map[string]updateCache)
	monitorUpsertCacheMut       sync.RWMutex
	monitorUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var monitorAfterSelectHooks []MonitorHook

var monitorBeforeInsertHooks []MonitorHook
var monitorAfterInsertHooks []MonitorHook

var monitorBeforeUpdateHooks []MonitorHook
var monitorAfterUpdateHooks []MonitorHook

var monitorBeforeDeleteHooks []MonitorHook
var monitorAfterDeleteHooks []MonitorHook

var monitorBeforeUpsertHooks []MonitorHook
var monitorAfterUpsertHooks []MonitorHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Monitor) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range monitorAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Monitor) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range monitorBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Monitor) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range monitorAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Monitor) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range monitorBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Monitor) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range monitorAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Monitor) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range monitorBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Monitor) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range monitorAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Monitor) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range monitorBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Monitor) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range monitorAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddMonitorHook registers your hook function for all future operations.
func AddMonitorHook(hookPoint boil.HookPoint, monitorHook MonitorHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		monitorAfterSelectHooks = append(monitorAfterSelectHooks, monitorHook)
	case boil.BeforeInsertHook:
		monitorBeforeInsertHooks = append(monitorBeforeInsertHooks, monitorHook)
	case boil.AfterInsertHook:
		monitorAfterInsertHooks = append(monitorAfterInsertHooks, monitorHook)
	case boil.BeforeUpdateHook:
		monitorBeforeUpdateHooks = append(monitorBeforeUpdateHooks, monitorHook)
	case boil.AfterUpdateHook:
		monitorAfterUpdateHooks = append(monitorAfterUpdateHooks, monitorHook)
	case boil.BeforeDeleteHook:
		monitorBeforeDeleteHooks = append(monitorBeforeDeleteHooks, monitorHook)
	case boil.AfterDeleteHook:
		monitorAfterDeleteHooks = append(monitorAfterDeleteHooks, monitorHook)
	case boil.BeforeUpsertHook:
		monitorBeforeUpsertHooks = append(monitorBeforeUpsertHooks, monitorHook)
	case boil.AfterUpsertHook:
		monitorAfterUpsertHooks = append(monitorAfterUpsertHooks, monitorHook)
	}
}

// One returns a single monitor record from the query.
func (q monitorQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Monitor, error) {
	o := &Monitor{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for monitors")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Monitor records from the query.
func (q monitorQuery) All(ctx context.Context, exec boil.ContextExecutor) (MonitorSlice, error) {
	var o []*Monitor

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Monitor slice")
	}

	if len(monitorAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Monitor records in the query.
func (q monitorQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count monitors rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q monitorQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if monitors exists")
	}

	return count > 0, nil
}

// Account pointed to by the foreign key.
func (o *Monitor) Account(mods ...qm.QueryMod) accountQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.AccountID),
	}

	queryMods = append(queryMods, mods...)

	return Accounts(queryMods...)
}

// LoadAccount allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (monitorL) LoadAccount(ctx context.Context, e boil.ContextExecutor, singular bool, maybeMonitor interface{}, mods queries.Applicator) error {
	var slice []*Monitor
	var object *Monitor

	if singular {
		var ok bool
		object, ok = maybeMonitor.(*Monitor)
		if !ok {
			object = new(Monitor)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeMonitor)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeMonitor))
			}
		}
	} else {
		s, ok := maybeMonitor.(*[]*Monitor)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeMonitor)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeMonitor))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &monitorR{}
		}
		args = append(args, object.AccountID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &monitorR{}
			}

			for _, a := range args {
				if a == obj.AccountID {
					continue Outer
				}
			}

			args = append(args, obj.AccountID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`accounts`),
		qm.WhereIn(`accounts.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Account")
	}

	var resultSlice []*Account
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Account")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for accounts")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for accounts")
	}

	if len(accountAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Account = foreign
		if foreign.R == nil {
			foreign.R = &accountR{}
		}
		foreign.R.Monitors = append(foreign.R.Monitors, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.AccountID == foreign.ID {
				local.R.Account = foreign
				if foreign.R == nil {
					foreign.R = &accountR{}
				}
				foreign.R.Monitors = append(foreign.R.Monitors, local)
				break
			}
		}
	}

	return nil
}

// SetAccount of the monitor to the related item.
// Sets o.R.Account to related.
// Adds o to related.R.Monitors.
func (o *Monitor) SetAccount(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Account) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"monitors\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"account_id"}),
		strmangle.WhereClause("\"", "\"", 2, monitorPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.AccountID = related.ID
	if o.R == nil {
		o.R = &monitorR{
			Account: related,
		}
	} else {
		o.R.Account = related
	}

	if related.R == nil {
		related.R = &accountR{
			Monitors: MonitorSlice{o},
		}
	} else {
		related.R.Monitors = append(related.R.Monitors, o)
	}

	return nil
}

// Monitors retrieves all the records using an executor.
func Monitors(mods ...qm.QueryMod) monitorQuery {
	mods = append(mods, qm.From("\"monitors\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"monitors\".*"})
	}

	return monitorQuery{q}
}

// FindMonitor retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindMonitor(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*Monitor, error) {
	monitorObj := &Monitor{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"monitors\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, monitorObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from monitors")
	}

	if err = monitorObj.doAfterSelectHooks(ctx, exec); err != nil {
		return monitorObj, err
	}

	return monitorObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Monitor) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no monitors provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(monitorColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	monitorInsertCacheMut.RLock()
	cache, cached := monitorInsertCache[key]
	monitorInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			monitorAllColumns,
			monitorColumnsWithDefault,
			monitorColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(monitorType, monitorMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(monitorType, monitorMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"monitors\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"monitors\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into monitors")
	}

	if !cached {
		monitorInsertCacheMut.Lock()
		monitorInsertCache[key] = cache
		monitorInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Monitor.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Monitor) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	monitorUpdateCacheMut.RLock()
	cache, cached := monitorUpdateCache[key]
	monitorUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			monitorAllColumns,
			monitorPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update monitors, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"monitors\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, monitorPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(monitorType, monitorMapping, append(wl, monitorPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update monitors row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for monitors")
	}

	if !cached {
		monitorUpdateCacheMut.Lock()
		monitorUpdateCache[key] = cache
		monitorUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q monitorQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for monitors")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for monitors")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o MonitorSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), monitorPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"monitors\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, monitorPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in monitor slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all monitor")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Monitor) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no monitors provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(monitorColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	monitorUpsertCacheMut.RLock()
	cache, cached := monitorUpsertCache[key]
	monitorUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			monitorAllColumns,
			monitorColumnsWithDefault,
			monitorColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			monitorAllColumns,
			monitorPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert monitors, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(monitorPrimaryKeyColumns))
			copy(conflict, monitorPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"monitors\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(monitorType, monitorMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(monitorType, monitorMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert monitors")
	}

	if !cached {
		monitorUpsertCacheMut.Lock()
		monitorUpsertCache[key] = cache
		monitorUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Monitor record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Monitor) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Monitor provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), monitorPrimaryKeyMapping)
	sql := "DELETE FROM \"monitors\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from monitors")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for monitors")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q monitorQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no monitorQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from monitors")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for monitors")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o MonitorSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(monitorBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), monitorPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"monitors\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, monitorPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from monitor slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for monitors")
	}

	if len(monitorAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Monitor) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindMonitor(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *MonitorSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := MonitorSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), monitorPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"monitors\".* FROM \"monitors\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, monitorPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in MonitorSlice")
	}

	*o = slice

	return nil
}

// MonitorExists checks if the Monitor row exists.
func MonitorExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"monitors\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if monitors exists")
	}

	return exists, nil
}

// Exists checks if the Monitor row exists.
func (o *Monitor) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return MonitorExists(ctx, exec, o.ID)
}