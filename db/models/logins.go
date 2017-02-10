package models

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/vattle/sqlboiler/strmangle"
	"gopkg.in/nullbio/null.v6"
)

// Login is an object representing the database table.
type Login struct {
	Snowflake  int64     `boil:"snowflake" json:"snowflake" toml:"snowflake" yaml:"snowflake"`
	CreatedAt  time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	DeletedAt  null.Time `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`
	UserID     int64     `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	Type       int       `boil:"type" json:"type" toml:"type" yaml:"type"`
	Data       []byte    `boil:"data" json:"data" toml:"data" yaml:"data"`
	Identifier string    `boil:"identifier" json:"identifier" toml:"identifier" yaml:"identifier"`

	R *loginR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L loginL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// loginR is where relationships are stored.
type loginR struct {
	User *User
}

// loginL is where Load methods for each relationship are stored.
type loginL struct{}

var (
	loginColumns               = []string{"snowflake", "created_at", "deleted_at", "user_id", "type", "data", "identifier"}
	loginColumnsWithoutDefault = []string{"snowflake", "deleted_at", "user_id", "type", "data", "identifier"}
	loginColumnsWithDefault    = []string{"created_at"}
	loginPrimaryKeyColumns     = []string{"snowflake"}
)

type (
	// LoginSlice is an alias for a slice of pointers to Login.
	// This should generally be used opposed to []Login.
	LoginSlice []*Login
	// LoginHook is the signature for custom Login hook methods
	LoginHook func(boil.Executor, *Login) error

	loginQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	loginType                 = reflect.TypeOf(&Login{})
	loginMapping              = queries.MakeStructMapping(loginType)
	loginPrimaryKeyMapping, _ = queries.BindMapping(loginType, loginMapping, loginPrimaryKeyColumns)
	loginInsertCacheMut       sync.RWMutex
	loginInsertCache          = make(map[string]insertCache)
	loginUpdateCacheMut       sync.RWMutex
	loginUpdateCache          = make(map[string]updateCache)
	loginUpsertCacheMut       sync.RWMutex
	loginUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var loginBeforeInsertHooks []LoginHook
var loginBeforeUpdateHooks []LoginHook
var loginBeforeDeleteHooks []LoginHook
var loginBeforeUpsertHooks []LoginHook

var loginAfterInsertHooks []LoginHook
var loginAfterSelectHooks []LoginHook
var loginAfterUpdateHooks []LoginHook
var loginAfterDeleteHooks []LoginHook
var loginAfterUpsertHooks []LoginHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Login) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range loginBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Login) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range loginBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Login) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range loginBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Login) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range loginBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Login) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range loginAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Login) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range loginAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Login) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range loginAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Login) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range loginAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Login) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range loginAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddLoginHook registers your hook function for all future operations.
func AddLoginHook(hookPoint boil.HookPoint, loginHook LoginHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		loginBeforeInsertHooks = append(loginBeforeInsertHooks, loginHook)
	case boil.BeforeUpdateHook:
		loginBeforeUpdateHooks = append(loginBeforeUpdateHooks, loginHook)
	case boil.BeforeDeleteHook:
		loginBeforeDeleteHooks = append(loginBeforeDeleteHooks, loginHook)
	case boil.BeforeUpsertHook:
		loginBeforeUpsertHooks = append(loginBeforeUpsertHooks, loginHook)
	case boil.AfterInsertHook:
		loginAfterInsertHooks = append(loginAfterInsertHooks, loginHook)
	case boil.AfterSelectHook:
		loginAfterSelectHooks = append(loginAfterSelectHooks, loginHook)
	case boil.AfterUpdateHook:
		loginAfterUpdateHooks = append(loginAfterUpdateHooks, loginHook)
	case boil.AfterDeleteHook:
		loginAfterDeleteHooks = append(loginAfterDeleteHooks, loginHook)
	case boil.AfterUpsertHook:
		loginAfterUpsertHooks = append(loginAfterUpsertHooks, loginHook)
	}
}

// OneP returns a single login record from the query, and panics on error.
func (q loginQuery) OneP() *Login {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single login record from the query.
func (q loginQuery) One() (*Login, error) {
	o := &Login{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for logins")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Login records from the query, and panics on error.
func (q loginQuery) AllP() LoginSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Login records from the query.
func (q loginQuery) All() (LoginSlice, error) {
	var o LoginSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Login slice")
	}

	if len(loginAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Login records in the query, and panics on error.
func (q loginQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Login records in the query.
func (q loginQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count logins rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q loginQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q loginQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if logins exists")
	}

	return count > 0, nil
}

// UserG pointed to by the foreign key.
func (o *Login) UserG(mods ...qm.QueryMod) userQuery {
	return o.User(boil.GetDB(), mods...)
}

// User pointed to by the foreign key.
func (o *Login) User(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("snowflake=?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(exec, queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (loginL) LoadUser(e boil.Executor, singular bool, maybeLogin interface{}) error {
	var slice []*Login
	var object *Login

	count := 1
	if singular {
		object = maybeLogin.(*Login)
	} else {
		slice = *maybeLogin.(*LoginSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &loginR{}
		}
		args[0] = object.UserID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &loginR{}
			}
			args[i] = obj.UserID
		}
	}

	query := fmt.Sprintf(
		"select * from \"users\" where \"snowflake\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load User")
	}
	defer results.Close()

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice User")
	}

	if len(loginAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		object.R.User = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.UserID == foreign.Snowflake {
				local.R.User = foreign
				break
			}
		}
	}

	return nil
}

// SetUserG of the login to the related item.
// Sets o.R.User to related.
// Adds o to related.R.Logins.
// Uses the global database handle.
func (o *Login) SetUserG(insert bool, related *User) error {
	return o.SetUser(boil.GetDB(), insert, related)
}

// SetUserP of the login to the related item.
// Sets o.R.User to related.
// Adds o to related.R.Logins.
// Panics on error.
func (o *Login) SetUserP(exec boil.Executor, insert bool, related *User) {
	if err := o.SetUser(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUserGP of the login to the related item.
// Sets o.R.User to related.
// Adds o to related.R.Logins.
// Uses the global database handle and panics on error.
func (o *Login) SetUserGP(insert bool, related *User) {
	if err := o.SetUser(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUser of the login to the related item.
// Sets o.R.User to related.
// Adds o to related.R.Logins.
func (o *Login) SetUser(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"logins\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, loginPrimaryKeyColumns),
	)
	values := []interface{}{related.Snowflake, o.Snowflake}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.UserID = related.Snowflake

	if o.R == nil {
		o.R = &loginR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			Logins: LoginSlice{o},
		}
	} else {
		related.R.Logins = append(related.R.Logins, o)
	}

	return nil
}

// LoginsG retrieves all records.
func LoginsG(mods ...qm.QueryMod) loginQuery {
	return Logins(boil.GetDB(), mods...)
}

// Logins retrieves all the records using an executor.
func Logins(exec boil.Executor, mods ...qm.QueryMod) loginQuery {
	mods = append(mods, qm.From("\"logins\""))
	return loginQuery{NewQuery(exec, mods...)}
}

// FindLoginG retrieves a single record by ID.
func FindLoginG(snowflake int64, selectCols ...string) (*Login, error) {
	return FindLogin(boil.GetDB(), snowflake, selectCols...)
}

// FindLoginGP retrieves a single record by ID, and panics on error.
func FindLoginGP(snowflake int64, selectCols ...string) *Login {
	retobj, err := FindLogin(boil.GetDB(), snowflake, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindLogin retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindLogin(exec boil.Executor, snowflake int64, selectCols ...string) (*Login, error) {
	loginObj := &Login{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"logins\" where \"snowflake\"=$1", sel,
	)

	q := queries.Raw(exec, query, snowflake)

	err := q.Bind(loginObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from logins")
	}

	return loginObj, nil
}

// FindLoginP retrieves a single record by ID with an executor, and panics on error.
func FindLoginP(exec boil.Executor, snowflake int64, selectCols ...string) *Login {
	retobj, err := FindLogin(exec, snowflake, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Login) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Login) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Login) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Login) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no logins provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(loginColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	loginInsertCacheMut.RLock()
	cache, cached := loginInsertCache[key]
	loginInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			loginColumns,
			loginColumnsWithDefault,
			loginColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(loginType, loginMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(loginType, loginMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"logins\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

		if len(cache.retMapping) != 0 {
			cache.query += fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into logins")
	}

	if !cached {
		loginInsertCacheMut.Lock()
		loginInsertCache[key] = cache
		loginInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Login record. See Update for
// whitelist behavior description.
func (o *Login) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Login record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Login) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Login, and panics on error.
// See Update for whitelist behavior description.
func (o *Login) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Login.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Login) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	loginUpdateCacheMut.RLock()
	cache, cached := loginUpdateCache[key]
	loginUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(loginColumns, loginPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update logins, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"logins\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, loginPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(loginType, loginMapping, append(wl, loginPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update logins row")
	}

	if !cached {
		loginUpdateCacheMut.Lock()
		loginUpdateCache[key] = cache
		loginUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q loginQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q loginQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for logins")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o LoginSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o LoginSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o LoginSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o LoginSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), loginPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"logins\" SET %s WHERE (\"snowflake\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(loginPrimaryKeyColumns), len(colNames)+1, len(loginPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in login slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Login) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Login) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Login) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Login) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no logins provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(loginColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
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
	for _, c := range updateColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	loginUpsertCacheMut.RLock()
	cache, cached := loginUpsertCache[key]
	loginUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			loginColumns,
			loginColumnsWithDefault,
			loginColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			loginColumns,
			loginPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert logins, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(loginPrimaryKeyColumns))
			copy(conflict, loginPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"logins\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(loginType, loginMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(loginType, loginMapping, ret)
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

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert logins")
	}

	if !cached {
		loginUpsertCacheMut.Lock()
		loginUpsertCache[key] = cache
		loginUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Login record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Login) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Login record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Login) DeleteG() error {
	if o == nil {
		return errors.New("models: no Login provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Login record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Login) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Login record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Login) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Login provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), loginPrimaryKeyMapping)
	sql := "DELETE FROM \"logins\" WHERE \"snowflake\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from logins")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q loginQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q loginQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no loginQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from logins")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o LoginSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o LoginSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Login slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o LoginSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o LoginSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Login slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(loginBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), loginPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"logins\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, loginPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(loginPrimaryKeyColumns), 1, len(loginPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from login slice")
	}

	if len(loginAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Login) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Login) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Login) ReloadG() error {
	if o == nil {
		return errors.New("models: no Login provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Login) Reload(exec boil.Executor) error {
	ret, err := FindLogin(exec, o.Snowflake)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *LoginSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *LoginSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *LoginSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty LoginSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *LoginSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	logins := LoginSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), loginPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"logins\".* FROM \"logins\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, loginPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(loginPrimaryKeyColumns), 1, len(loginPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&logins)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in LoginSlice")
	}

	*o = logins

	return nil
}

// LoginExists checks if the Login row exists.
func LoginExists(exec boil.Executor, snowflake int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"logins\" where \"snowflake\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, snowflake)
	}

	row := exec.QueryRow(sql, snowflake)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if logins exists")
	}

	return exists, nil
}

// LoginExistsG checks if the Login row exists.
func LoginExistsG(snowflake int64) (bool, error) {
	return LoginExists(boil.GetDB(), snowflake)
}

// LoginExistsGP checks if the Login row exists. Panics on error.
func LoginExistsGP(snowflake int64) bool {
	e, err := LoginExists(boil.GetDB(), snowflake)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// LoginExistsP checks if the Login row exists. Panics on error.
func LoginExistsP(exec boil.Executor, snowflake int64) bool {
	e, err := LoginExists(exec, snowflake)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
