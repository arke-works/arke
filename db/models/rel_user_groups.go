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

// RelUserGroup is an object representing the database table.
type RelUserGroup struct {
	UserID    int64     `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	GroupID   int64     `boil:"group_id" json:"group_id" toml:"group_id" yaml:"group_id"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	DeletedAt null.Time `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`

	R *relUserGroupR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L relUserGroupL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// relUserGroupR is where relationships are stored.
type relUserGroupR struct {
	User  *User
	Group *Group
}

// relUserGroupL is where Load methods for each relationship are stored.
type relUserGroupL struct{}

var (
	relUserGroupColumns               = []string{"user_id", "group_id", "created_at", "deleted_at"}
	relUserGroupColumnsWithoutDefault = []string{"user_id", "group_id", "deleted_at"}
	relUserGroupColumnsWithDefault    = []string{"created_at"}
	relUserGroupPrimaryKeyColumns     = []string{"user_id", "group_id"}
)

type (
	// RelUserGroupSlice is an alias for a slice of pointers to RelUserGroup.
	// This should generally be used opposed to []RelUserGroup.
	RelUserGroupSlice []*RelUserGroup
	// RelUserGroupHook is the signature for custom RelUserGroup hook methods
	RelUserGroupHook func(boil.Executor, *RelUserGroup) error

	relUserGroupQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	relUserGroupType                 = reflect.TypeOf(&RelUserGroup{})
	relUserGroupMapping              = queries.MakeStructMapping(relUserGroupType)
	relUserGroupPrimaryKeyMapping, _ = queries.BindMapping(relUserGroupType, relUserGroupMapping, relUserGroupPrimaryKeyColumns)
	relUserGroupInsertCacheMut       sync.RWMutex
	relUserGroupInsertCache          = make(map[string]insertCache)
	relUserGroupUpdateCacheMut       sync.RWMutex
	relUserGroupUpdateCache          = make(map[string]updateCache)
	relUserGroupUpsertCacheMut       sync.RWMutex
	relUserGroupUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var relUserGroupBeforeInsertHooks []RelUserGroupHook
var relUserGroupBeforeUpdateHooks []RelUserGroupHook
var relUserGroupBeforeDeleteHooks []RelUserGroupHook
var relUserGroupBeforeUpsertHooks []RelUserGroupHook

var relUserGroupAfterInsertHooks []RelUserGroupHook
var relUserGroupAfterSelectHooks []RelUserGroupHook
var relUserGroupAfterUpdateHooks []RelUserGroupHook
var relUserGroupAfterDeleteHooks []RelUserGroupHook
var relUserGroupAfterUpsertHooks []RelUserGroupHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *RelUserGroup) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range relUserGroupBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *RelUserGroup) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range relUserGroupBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *RelUserGroup) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range relUserGroupBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *RelUserGroup) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range relUserGroupBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *RelUserGroup) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range relUserGroupAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *RelUserGroup) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range relUserGroupAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *RelUserGroup) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range relUserGroupAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *RelUserGroup) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range relUserGroupAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *RelUserGroup) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range relUserGroupAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddRelUserGroupHook registers your hook function for all future operations.
func AddRelUserGroupHook(hookPoint boil.HookPoint, relUserGroupHook RelUserGroupHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		relUserGroupBeforeInsertHooks = append(relUserGroupBeforeInsertHooks, relUserGroupHook)
	case boil.BeforeUpdateHook:
		relUserGroupBeforeUpdateHooks = append(relUserGroupBeforeUpdateHooks, relUserGroupHook)
	case boil.BeforeDeleteHook:
		relUserGroupBeforeDeleteHooks = append(relUserGroupBeforeDeleteHooks, relUserGroupHook)
	case boil.BeforeUpsertHook:
		relUserGroupBeforeUpsertHooks = append(relUserGroupBeforeUpsertHooks, relUserGroupHook)
	case boil.AfterInsertHook:
		relUserGroupAfterInsertHooks = append(relUserGroupAfterInsertHooks, relUserGroupHook)
	case boil.AfterSelectHook:
		relUserGroupAfterSelectHooks = append(relUserGroupAfterSelectHooks, relUserGroupHook)
	case boil.AfterUpdateHook:
		relUserGroupAfterUpdateHooks = append(relUserGroupAfterUpdateHooks, relUserGroupHook)
	case boil.AfterDeleteHook:
		relUserGroupAfterDeleteHooks = append(relUserGroupAfterDeleteHooks, relUserGroupHook)
	case boil.AfterUpsertHook:
		relUserGroupAfterUpsertHooks = append(relUserGroupAfterUpsertHooks, relUserGroupHook)
	}
}

// OneP returns a single relUserGroup record from the query, and panics on error.
func (q relUserGroupQuery) OneP() *RelUserGroup {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single relUserGroup record from the query.
func (q relUserGroupQuery) One() (*RelUserGroup, error) {
	o := &RelUserGroup{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for rel_user_groups")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all RelUserGroup records from the query, and panics on error.
func (q relUserGroupQuery) AllP() RelUserGroupSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all RelUserGroup records from the query.
func (q relUserGroupQuery) All() (RelUserGroupSlice, error) {
	var o RelUserGroupSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to RelUserGroup slice")
	}

	if len(relUserGroupAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all RelUserGroup records in the query, and panics on error.
func (q relUserGroupQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all RelUserGroup records in the query.
func (q relUserGroupQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count rel_user_groups rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q relUserGroupQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q relUserGroupQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if rel_user_groups exists")
	}

	return count > 0, nil
}

// UserG pointed to by the foreign key.
func (o *RelUserGroup) UserG(mods ...qm.QueryMod) userQuery {
	return o.User(boil.GetDB(), mods...)
}

// User pointed to by the foreign key.
func (o *RelUserGroup) User(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("snowflake=?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(exec, queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// GroupG pointed to by the foreign key.
func (o *RelUserGroup) GroupG(mods ...qm.QueryMod) groupQuery {
	return o.Group(boil.GetDB(), mods...)
}

// Group pointed to by the foreign key.
func (o *RelUserGroup) Group(exec boil.Executor, mods ...qm.QueryMod) groupQuery {
	queryMods := []qm.QueryMod{
		qm.Where("snowflake=?", o.GroupID),
	}

	queryMods = append(queryMods, mods...)

	query := Groups(exec, queryMods...)
	queries.SetFrom(query.Query, "\"groups\"")

	return query
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (relUserGroupL) LoadUser(e boil.Executor, singular bool, maybeRelUserGroup interface{}) error {
	var slice []*RelUserGroup
	var object *RelUserGroup

	count := 1
	if singular {
		object = maybeRelUserGroup.(*RelUserGroup)
	} else {
		slice = *maybeRelUserGroup.(*RelUserGroupSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &relUserGroupR{}
		}
		args[0] = object.UserID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &relUserGroupR{}
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

	if len(relUserGroupAfterSelectHooks) != 0 {
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

// LoadGroup allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (relUserGroupL) LoadGroup(e boil.Executor, singular bool, maybeRelUserGroup interface{}) error {
	var slice []*RelUserGroup
	var object *RelUserGroup

	count := 1
	if singular {
		object = maybeRelUserGroup.(*RelUserGroup)
	} else {
		slice = *maybeRelUserGroup.(*RelUserGroupSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &relUserGroupR{}
		}
		args[0] = object.GroupID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &relUserGroupR{}
			}
			args[i] = obj.GroupID
		}
	}

	query := fmt.Sprintf(
		"select * from \"groups\" where \"snowflake\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Group")
	}
	defer results.Close()

	var resultSlice []*Group
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Group")
	}

	if len(relUserGroupAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		object.R.Group = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.GroupID == foreign.Snowflake {
				local.R.Group = foreign
				break
			}
		}
	}

	return nil
}

// SetUserG of the rel_user_group to the related item.
// Sets o.R.User to related.
// Adds o to related.R.RelUserGroups.
// Uses the global database handle.
func (o *RelUserGroup) SetUserG(insert bool, related *User) error {
	return o.SetUser(boil.GetDB(), insert, related)
}

// SetUserP of the rel_user_group to the related item.
// Sets o.R.User to related.
// Adds o to related.R.RelUserGroups.
// Panics on error.
func (o *RelUserGroup) SetUserP(exec boil.Executor, insert bool, related *User) {
	if err := o.SetUser(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUserGP of the rel_user_group to the related item.
// Sets o.R.User to related.
// Adds o to related.R.RelUserGroups.
// Uses the global database handle and panics on error.
func (o *RelUserGroup) SetUserGP(insert bool, related *User) {
	if err := o.SetUser(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUser of the rel_user_group to the related item.
// Sets o.R.User to related.
// Adds o to related.R.RelUserGroups.
func (o *RelUserGroup) SetUser(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"rel_user_groups\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, relUserGroupPrimaryKeyColumns),
	)
	values := []interface{}{related.Snowflake, o.UserID, o.GroupID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.UserID = related.Snowflake

	if o.R == nil {
		o.R = &relUserGroupR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			RelUserGroups: RelUserGroupSlice{o},
		}
	} else {
		related.R.RelUserGroups = append(related.R.RelUserGroups, o)
	}

	return nil
}

// SetGroupG of the rel_user_group to the related item.
// Sets o.R.Group to related.
// Adds o to related.R.RelUserGroups.
// Uses the global database handle.
func (o *RelUserGroup) SetGroupG(insert bool, related *Group) error {
	return o.SetGroup(boil.GetDB(), insert, related)
}

// SetGroupP of the rel_user_group to the related item.
// Sets o.R.Group to related.
// Adds o to related.R.RelUserGroups.
// Panics on error.
func (o *RelUserGroup) SetGroupP(exec boil.Executor, insert bool, related *Group) {
	if err := o.SetGroup(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetGroupGP of the rel_user_group to the related item.
// Sets o.R.Group to related.
// Adds o to related.R.RelUserGroups.
// Uses the global database handle and panics on error.
func (o *RelUserGroup) SetGroupGP(insert bool, related *Group) {
	if err := o.SetGroup(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetGroup of the rel_user_group to the related item.
// Sets o.R.Group to related.
// Adds o to related.R.RelUserGroups.
func (o *RelUserGroup) SetGroup(exec boil.Executor, insert bool, related *Group) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"rel_user_groups\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"group_id"}),
		strmangle.WhereClause("\"", "\"", 2, relUserGroupPrimaryKeyColumns),
	)
	values := []interface{}{related.Snowflake, o.UserID, o.GroupID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.GroupID = related.Snowflake

	if o.R == nil {
		o.R = &relUserGroupR{
			Group: related,
		}
	} else {
		o.R.Group = related
	}

	if related.R == nil {
		related.R = &groupR{
			RelUserGroups: RelUserGroupSlice{o},
		}
	} else {
		related.R.RelUserGroups = append(related.R.RelUserGroups, o)
	}

	return nil
}

// RelUserGroupsG retrieves all records.
func RelUserGroupsG(mods ...qm.QueryMod) relUserGroupQuery {
	return RelUserGroups(boil.GetDB(), mods...)
}

// RelUserGroups retrieves all the records using an executor.
func RelUserGroups(exec boil.Executor, mods ...qm.QueryMod) relUserGroupQuery {
	mods = append(mods, qm.From("\"rel_user_groups\""))
	return relUserGroupQuery{NewQuery(exec, mods...)}
}

// FindRelUserGroupG retrieves a single record by ID.
func FindRelUserGroupG(userID int64, groupID int64, selectCols ...string) (*RelUserGroup, error) {
	return FindRelUserGroup(boil.GetDB(), userID, groupID, selectCols...)
}

// FindRelUserGroupGP retrieves a single record by ID, and panics on error.
func FindRelUserGroupGP(userID int64, groupID int64, selectCols ...string) *RelUserGroup {
	retobj, err := FindRelUserGroup(boil.GetDB(), userID, groupID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindRelUserGroup retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindRelUserGroup(exec boil.Executor, userID int64, groupID int64, selectCols ...string) (*RelUserGroup, error) {
	relUserGroupObj := &RelUserGroup{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"rel_user_groups\" where \"user_id\"=$1 AND \"group_id\"=$2", sel,
	)

	q := queries.Raw(exec, query, userID, groupID)

	err := q.Bind(relUserGroupObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from rel_user_groups")
	}

	return relUserGroupObj, nil
}

// FindRelUserGroupP retrieves a single record by ID with an executor, and panics on error.
func FindRelUserGroupP(exec boil.Executor, userID int64, groupID int64, selectCols ...string) *RelUserGroup {
	retobj, err := FindRelUserGroup(exec, userID, groupID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *RelUserGroup) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *RelUserGroup) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *RelUserGroup) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *RelUserGroup) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no rel_user_groups provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(relUserGroupColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	relUserGroupInsertCacheMut.RLock()
	cache, cached := relUserGroupInsertCache[key]
	relUserGroupInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			relUserGroupColumns,
			relUserGroupColumnsWithDefault,
			relUserGroupColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(relUserGroupType, relUserGroupMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(relUserGroupType, relUserGroupMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"rel_user_groups\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into rel_user_groups")
	}

	if !cached {
		relUserGroupInsertCacheMut.Lock()
		relUserGroupInsertCache[key] = cache
		relUserGroupInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single RelUserGroup record. See Update for
// whitelist behavior description.
func (o *RelUserGroup) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single RelUserGroup record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *RelUserGroup) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the RelUserGroup, and panics on error.
// See Update for whitelist behavior description.
func (o *RelUserGroup) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the RelUserGroup.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *RelUserGroup) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	relUserGroupUpdateCacheMut.RLock()
	cache, cached := relUserGroupUpdateCache[key]
	relUserGroupUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(relUserGroupColumns, relUserGroupPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update rel_user_groups, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"rel_user_groups\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, relUserGroupPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(relUserGroupType, relUserGroupMapping, append(wl, relUserGroupPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update rel_user_groups row")
	}

	if !cached {
		relUserGroupUpdateCacheMut.Lock()
		relUserGroupUpdateCache[key] = cache
		relUserGroupUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q relUserGroupQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q relUserGroupQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for rel_user_groups")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o RelUserGroupSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o RelUserGroupSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o RelUserGroupSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o RelUserGroupSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), relUserGroupPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"rel_user_groups\" SET %s WHERE (\"user_id\",\"group_id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(relUserGroupPrimaryKeyColumns), len(colNames)+1, len(relUserGroupPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in relUserGroup slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *RelUserGroup) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *RelUserGroup) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *RelUserGroup) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *RelUserGroup) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no rel_user_groups provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(relUserGroupColumnsWithDefault, o)

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

	relUserGroupUpsertCacheMut.RLock()
	cache, cached := relUserGroupUpsertCache[key]
	relUserGroupUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			relUserGroupColumns,
			relUserGroupColumnsWithDefault,
			relUserGroupColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			relUserGroupColumns,
			relUserGroupPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert rel_user_groups, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(relUserGroupPrimaryKeyColumns))
			copy(conflict, relUserGroupPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"rel_user_groups\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(relUserGroupType, relUserGroupMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(relUserGroupType, relUserGroupMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert rel_user_groups")
	}

	if !cached {
		relUserGroupUpsertCacheMut.Lock()
		relUserGroupUpsertCache[key] = cache
		relUserGroupUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single RelUserGroup record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *RelUserGroup) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single RelUserGroup record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *RelUserGroup) DeleteG() error {
	if o == nil {
		return errors.New("models: no RelUserGroup provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single RelUserGroup record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *RelUserGroup) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single RelUserGroup record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *RelUserGroup) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no RelUserGroup provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), relUserGroupPrimaryKeyMapping)
	sql := "DELETE FROM \"rel_user_groups\" WHERE \"user_id\"=$1 AND \"group_id\"=$2"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from rel_user_groups")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q relUserGroupQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q relUserGroupQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no relUserGroupQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from rel_user_groups")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o RelUserGroupSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o RelUserGroupSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no RelUserGroup slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o RelUserGroupSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o RelUserGroupSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no RelUserGroup slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(relUserGroupBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), relUserGroupPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"rel_user_groups\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, relUserGroupPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(relUserGroupPrimaryKeyColumns), 1, len(relUserGroupPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from relUserGroup slice")
	}

	if len(relUserGroupAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *RelUserGroup) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *RelUserGroup) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *RelUserGroup) ReloadG() error {
	if o == nil {
		return errors.New("models: no RelUserGroup provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *RelUserGroup) Reload(exec boil.Executor) error {
	ret, err := FindRelUserGroup(exec, o.UserID, o.GroupID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *RelUserGroupSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *RelUserGroupSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *RelUserGroupSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty RelUserGroupSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *RelUserGroupSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	relUserGroups := RelUserGroupSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), relUserGroupPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"rel_user_groups\".* FROM \"rel_user_groups\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, relUserGroupPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(relUserGroupPrimaryKeyColumns), 1, len(relUserGroupPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&relUserGroups)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in RelUserGroupSlice")
	}

	*o = relUserGroups

	return nil
}

// RelUserGroupExists checks if the RelUserGroup row exists.
func RelUserGroupExists(exec boil.Executor, userID int64, groupID int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"rel_user_groups\" where \"user_id\"=$1 AND \"group_id\"=$2 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, userID, groupID)
	}

	row := exec.QueryRow(sql, userID, groupID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if rel_user_groups exists")
	}

	return exists, nil
}

// RelUserGroupExistsG checks if the RelUserGroup row exists.
func RelUserGroupExistsG(userID int64, groupID int64) (bool, error) {
	return RelUserGroupExists(boil.GetDB(), userID, groupID)
}

// RelUserGroupExistsGP checks if the RelUserGroup row exists. Panics on error.
func RelUserGroupExistsGP(userID int64, groupID int64) bool {
	e, err := RelUserGroupExists(boil.GetDB(), userID, groupID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// RelUserGroupExistsP checks if the RelUserGroup row exists. Panics on error.
func RelUserGroupExistsP(exec boil.Executor, userID int64, groupID int64) bool {
	e, err := RelUserGroupExists(exec, userID, groupID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
