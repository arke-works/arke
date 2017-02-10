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

// Group is an object representing the database table.
type Group struct {
	Snowflake  int64      `boil:"snowflake" json:"snowflake" toml:"snowflake" yaml:"snowflake"`
	CreatedAt  time.Time  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	DeletedAt  null.Time  `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`
	Name       string     `boil:"name" json:"name" toml:"name" yaml:"name"`
	Permission null.Bytes `boil:"permission" json:"permission,omitempty" toml:"permission" yaml:"permission,omitempty"`
	ParentID   null.Int64 `boil:"parent_id" json:"parent_id,omitempty" toml:"parent_id" yaml:"parent_id,omitempty"`

	R *groupR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L groupL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// groupR is where relationships are stored.
type groupR struct {
	Parent        *Group
	ParentGroups  GroupSlice
	RelUserGroups RelUserGroupSlice
}

// groupL is where Load methods for each relationship are stored.
type groupL struct{}

var (
	groupColumns               = []string{"snowflake", "created_at", "deleted_at", "name", "permission", "parent_id"}
	groupColumnsWithoutDefault = []string{"snowflake", "deleted_at", "name", "permission", "parent_id"}
	groupColumnsWithDefault    = []string{"created_at"}
	groupPrimaryKeyColumns     = []string{"snowflake"}
)

type (
	// GroupSlice is an alias for a slice of pointers to Group.
	// This should generally be used opposed to []Group.
	GroupSlice []*Group
	// GroupHook is the signature for custom Group hook methods
	GroupHook func(boil.Executor, *Group) error

	groupQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	groupType                 = reflect.TypeOf(&Group{})
	groupMapping              = queries.MakeStructMapping(groupType)
	groupPrimaryKeyMapping, _ = queries.BindMapping(groupType, groupMapping, groupPrimaryKeyColumns)
	groupInsertCacheMut       sync.RWMutex
	groupInsertCache          = make(map[string]insertCache)
	groupUpdateCacheMut       sync.RWMutex
	groupUpdateCache          = make(map[string]updateCache)
	groupUpsertCacheMut       sync.RWMutex
	groupUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var groupBeforeInsertHooks []GroupHook
var groupBeforeUpdateHooks []GroupHook
var groupBeforeDeleteHooks []GroupHook
var groupBeforeUpsertHooks []GroupHook

var groupAfterInsertHooks []GroupHook
var groupAfterSelectHooks []GroupHook
var groupAfterUpdateHooks []GroupHook
var groupAfterDeleteHooks []GroupHook
var groupAfterUpsertHooks []GroupHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Group) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range groupBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Group) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range groupBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Group) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range groupBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Group) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range groupBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Group) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range groupAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Group) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range groupAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Group) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range groupAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Group) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range groupAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Group) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range groupAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddGroupHook registers your hook function for all future operations.
func AddGroupHook(hookPoint boil.HookPoint, groupHook GroupHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		groupBeforeInsertHooks = append(groupBeforeInsertHooks, groupHook)
	case boil.BeforeUpdateHook:
		groupBeforeUpdateHooks = append(groupBeforeUpdateHooks, groupHook)
	case boil.BeforeDeleteHook:
		groupBeforeDeleteHooks = append(groupBeforeDeleteHooks, groupHook)
	case boil.BeforeUpsertHook:
		groupBeforeUpsertHooks = append(groupBeforeUpsertHooks, groupHook)
	case boil.AfterInsertHook:
		groupAfterInsertHooks = append(groupAfterInsertHooks, groupHook)
	case boil.AfterSelectHook:
		groupAfterSelectHooks = append(groupAfterSelectHooks, groupHook)
	case boil.AfterUpdateHook:
		groupAfterUpdateHooks = append(groupAfterUpdateHooks, groupHook)
	case boil.AfterDeleteHook:
		groupAfterDeleteHooks = append(groupAfterDeleteHooks, groupHook)
	case boil.AfterUpsertHook:
		groupAfterUpsertHooks = append(groupAfterUpsertHooks, groupHook)
	}
}

// OneP returns a single group record from the query, and panics on error.
func (q groupQuery) OneP() *Group {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single group record from the query.
func (q groupQuery) One() (*Group, error) {
	o := &Group{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for groups")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Group records from the query, and panics on error.
func (q groupQuery) AllP() GroupSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Group records from the query.
func (q groupQuery) All() (GroupSlice, error) {
	var o GroupSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Group slice")
	}

	if len(groupAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Group records in the query, and panics on error.
func (q groupQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Group records in the query.
func (q groupQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count groups rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q groupQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q groupQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if groups exists")
	}

	return count > 0, nil
}

// ParentG pointed to by the foreign key.
func (o *Group) ParentG(mods ...qm.QueryMod) groupQuery {
	return o.Parent(boil.GetDB(), mods...)
}

// Parent pointed to by the foreign key.
func (o *Group) Parent(exec boil.Executor, mods ...qm.QueryMod) groupQuery {
	queryMods := []qm.QueryMod{
		qm.Where("snowflake=?", o.ParentID),
	}

	queryMods = append(queryMods, mods...)

	query := Groups(exec, queryMods...)
	queries.SetFrom(query.Query, "\"groups\"")

	return query
}

// ParentGroupsG retrieves all the group's groups via parent_id column.
func (o *Group) ParentGroupsG(mods ...qm.QueryMod) groupQuery {
	return o.ParentGroups(boil.GetDB(), mods...)
}

// ParentGroups retrieves all the group's groups with an executor via parent_id column.
func (o *Group) ParentGroups(exec boil.Executor, mods ...qm.QueryMod) groupQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"parent_id\"=?", o.Snowflake),
	)

	query := Groups(exec, queryMods...)
	queries.SetFrom(query.Query, "\"groups\" as \"a\"")
	return query
}

// RelUserGroupsG retrieves all the rel_user_group's rel user groups.
func (o *Group) RelUserGroupsG(mods ...qm.QueryMod) relUserGroupQuery {
	return o.RelUserGroups(boil.GetDB(), mods...)
}

// RelUserGroups retrieves all the rel_user_group's rel user groups with an executor.
func (o *Group) RelUserGroups(exec boil.Executor, mods ...qm.QueryMod) relUserGroupQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"group_id\"=?", o.Snowflake),
	)

	query := RelUserGroups(exec, queryMods...)
	queries.SetFrom(query.Query, "\"rel_user_groups\" as \"a\"")
	return query
}

// LoadParent allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (groupL) LoadParent(e boil.Executor, singular bool, maybeGroup interface{}) error {
	var slice []*Group
	var object *Group

	count := 1
	if singular {
		object = maybeGroup.(*Group)
	} else {
		slice = *maybeGroup.(*GroupSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &groupR{}
		}
		args[0] = object.ParentID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &groupR{}
			}
			args[i] = obj.ParentID
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

	if len(groupAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		object.R.Parent = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ParentID.Int64 == foreign.Snowflake {
				local.R.Parent = foreign
				break
			}
		}
	}

	return nil
}

// LoadParentGroups allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (groupL) LoadParentGroups(e boil.Executor, singular bool, maybeGroup interface{}) error {
	var slice []*Group
	var object *Group

	count := 1
	if singular {
		object = maybeGroup.(*Group)
	} else {
		slice = *maybeGroup.(*GroupSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &groupR{}
		}
		args[0] = object.Snowflake
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &groupR{}
			}
			args[i] = obj.Snowflake
		}
	}

	query := fmt.Sprintf(
		"select * from \"groups\" where \"parent_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load groups")
	}
	defer results.Close()

	var resultSlice []*Group
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice groups")
	}

	if len(groupAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.ParentGroups = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.Snowflake == foreign.ParentID.Int64 {
				local.R.ParentGroups = append(local.R.ParentGroups, foreign)
				break
			}
		}
	}

	return nil
}

// LoadRelUserGroups allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (groupL) LoadRelUserGroups(e boil.Executor, singular bool, maybeGroup interface{}) error {
	var slice []*Group
	var object *Group

	count := 1
	if singular {
		object = maybeGroup.(*Group)
	} else {
		slice = *maybeGroup.(*GroupSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &groupR{}
		}
		args[0] = object.Snowflake
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &groupR{}
			}
			args[i] = obj.Snowflake
		}
	}

	query := fmt.Sprintf(
		"select * from \"rel_user_groups\" where \"group_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load rel_user_groups")
	}
	defer results.Close()

	var resultSlice []*RelUserGroup
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice rel_user_groups")
	}

	if len(relUserGroupAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.RelUserGroups = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.Snowflake == foreign.GroupID {
				local.R.RelUserGroups = append(local.R.RelUserGroups, foreign)
				break
			}
		}
	}

	return nil
}

// SetParentG of the group to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentGroups.
// Uses the global database handle.
func (o *Group) SetParentG(insert bool, related *Group) error {
	return o.SetParent(boil.GetDB(), insert, related)
}

// SetParentP of the group to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentGroups.
// Panics on error.
func (o *Group) SetParentP(exec boil.Executor, insert bool, related *Group) {
	if err := o.SetParent(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParentGP of the group to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentGroups.
// Uses the global database handle and panics on error.
func (o *Group) SetParentGP(insert bool, related *Group) {
	if err := o.SetParent(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParent of the group to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentGroups.
func (o *Group) SetParent(exec boil.Executor, insert bool, related *Group) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"groups\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"parent_id"}),
		strmangle.WhereClause("\"", "\"", 2, groupPrimaryKeyColumns),
	)
	values := []interface{}{related.Snowflake, o.Snowflake}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ParentID.Int64 = related.Snowflake
	o.ParentID.Valid = true

	if o.R == nil {
		o.R = &groupR{
			Parent: related,
		}
	} else {
		o.R.Parent = related
	}

	if related.R == nil {
		related.R = &groupR{
			ParentGroups: GroupSlice{o},
		}
	} else {
		related.R.ParentGroups = append(related.R.ParentGroups, o)
	}

	return nil
}

// RemoveParentG relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle.
func (o *Group) RemoveParentG(related *Group) error {
	return o.RemoveParent(boil.GetDB(), related)
}

// RemoveParentP relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Panics on error.
func (o *Group) RemoveParentP(exec boil.Executor, related *Group) {
	if err := o.RemoveParent(exec, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParentGP relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle and panics on error.
func (o *Group) RemoveParentGP(related *Group) {
	if err := o.RemoveParent(boil.GetDB(), related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParent relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *Group) RemoveParent(exec boil.Executor, related *Group) error {
	var err error

	o.ParentID.Valid = false
	if err = o.Update(exec, "parent_id"); err != nil {
		o.ParentID.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.Parent = nil
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.ParentGroups {
		if o.ParentID.Int64 != ri.ParentID.Int64 {
			continue
		}

		ln := len(related.R.ParentGroups)
		if ln > 1 && i < ln-1 {
			related.R.ParentGroups[i] = related.R.ParentGroups[ln-1]
		}
		related.R.ParentGroups = related.R.ParentGroups[:ln-1]
		break
	}
	return nil
}

// AddParentGroupsG adds the given related objects to the existing relationships
// of the group, optionally inserting them as new records.
// Appends related to o.R.ParentGroups.
// Sets related.R.Parent appropriately.
// Uses the global database handle.
func (o *Group) AddParentGroupsG(insert bool, related ...*Group) error {
	return o.AddParentGroups(boil.GetDB(), insert, related...)
}

// AddParentGroupsP adds the given related objects to the existing relationships
// of the group, optionally inserting them as new records.
// Appends related to o.R.ParentGroups.
// Sets related.R.Parent appropriately.
// Panics on error.
func (o *Group) AddParentGroupsP(exec boil.Executor, insert bool, related ...*Group) {
	if err := o.AddParentGroups(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddParentGroupsGP adds the given related objects to the existing relationships
// of the group, optionally inserting them as new records.
// Appends related to o.R.ParentGroups.
// Sets related.R.Parent appropriately.
// Uses the global database handle and panics on error.
func (o *Group) AddParentGroupsGP(insert bool, related ...*Group) {
	if err := o.AddParentGroups(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddParentGroups adds the given related objects to the existing relationships
// of the group, optionally inserting them as new records.
// Appends related to o.R.ParentGroups.
// Sets related.R.Parent appropriately.
func (o *Group) AddParentGroups(exec boil.Executor, insert bool, related ...*Group) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.ParentID.Int64 = o.Snowflake
			rel.ParentID.Valid = true
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"groups\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"parent_id"}),
				strmangle.WhereClause("\"", "\"", 2, groupPrimaryKeyColumns),
			)
			values := []interface{}{o.Snowflake, rel.Snowflake}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.ParentID.Int64 = o.Snowflake
			rel.ParentID.Valid = true
		}
	}

	if o.R == nil {
		o.R = &groupR{
			ParentGroups: related,
		}
	} else {
		o.R.ParentGroups = append(o.R.ParentGroups, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &groupR{
				Parent: o,
			}
		} else {
			rel.R.Parent = o
		}
	}
	return nil
}

// SetParentGroupsG removes all previously related items of the
// group replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentGroups accordingly.
// Replaces o.R.ParentGroups with related.
// Sets related.R.Parent's ParentGroups accordingly.
// Uses the global database handle.
func (o *Group) SetParentGroupsG(insert bool, related ...*Group) error {
	return o.SetParentGroups(boil.GetDB(), insert, related...)
}

// SetParentGroupsP removes all previously related items of the
// group replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentGroups accordingly.
// Replaces o.R.ParentGroups with related.
// Sets related.R.Parent's ParentGroups accordingly.
// Panics on error.
func (o *Group) SetParentGroupsP(exec boil.Executor, insert bool, related ...*Group) {
	if err := o.SetParentGroups(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParentGroupsGP removes all previously related items of the
// group replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentGroups accordingly.
// Replaces o.R.ParentGroups with related.
// Sets related.R.Parent's ParentGroups accordingly.
// Uses the global database handle and panics on error.
func (o *Group) SetParentGroupsGP(insert bool, related ...*Group) {
	if err := o.SetParentGroups(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParentGroups removes all previously related items of the
// group replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentGroups accordingly.
// Replaces o.R.ParentGroups with related.
// Sets related.R.Parent's ParentGroups accordingly.
func (o *Group) SetParentGroups(exec boil.Executor, insert bool, related ...*Group) error {
	query := "update \"groups\" set \"parent_id\" = null where \"parent_id\" = $1"
	values := []interface{}{o.Snowflake}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.ParentGroups {
			rel.ParentID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Parent = nil
		}

		o.R.ParentGroups = nil
	}
	return o.AddParentGroups(exec, insert, related...)
}

// RemoveParentGroupsG relationships from objects passed in.
// Removes related items from R.ParentGroups (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
// Uses the global database handle.
func (o *Group) RemoveParentGroupsG(related ...*Group) error {
	return o.RemoveParentGroups(boil.GetDB(), related...)
}

// RemoveParentGroupsP relationships from objects passed in.
// Removes related items from R.ParentGroups (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
// Panics on error.
func (o *Group) RemoveParentGroupsP(exec boil.Executor, related ...*Group) {
	if err := o.RemoveParentGroups(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParentGroupsGP relationships from objects passed in.
// Removes related items from R.ParentGroups (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
// Uses the global database handle and panics on error.
func (o *Group) RemoveParentGroupsGP(related ...*Group) {
	if err := o.RemoveParentGroups(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParentGroups relationships from objects passed in.
// Removes related items from R.ParentGroups (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
func (o *Group) RemoveParentGroups(exec boil.Executor, related ...*Group) error {
	var err error
	for _, rel := range related {
		rel.ParentID.Valid = false
		if rel.R != nil {
			rel.R.Parent = nil
		}
		if err = rel.Update(exec, "parent_id"); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.ParentGroups {
			if rel != ri {
				continue
			}

			ln := len(o.R.ParentGroups)
			if ln > 1 && i < ln-1 {
				o.R.ParentGroups[i] = o.R.ParentGroups[ln-1]
			}
			o.R.ParentGroups = o.R.ParentGroups[:ln-1]
			break
		}
	}

	return nil
}

// AddRelUserGroupsG adds the given related objects to the existing relationships
// of the group, optionally inserting them as new records.
// Appends related to o.R.RelUserGroups.
// Sets related.R.Group appropriately.
// Uses the global database handle.
func (o *Group) AddRelUserGroupsG(insert bool, related ...*RelUserGroup) error {
	return o.AddRelUserGroups(boil.GetDB(), insert, related...)
}

// AddRelUserGroupsP adds the given related objects to the existing relationships
// of the group, optionally inserting them as new records.
// Appends related to o.R.RelUserGroups.
// Sets related.R.Group appropriately.
// Panics on error.
func (o *Group) AddRelUserGroupsP(exec boil.Executor, insert bool, related ...*RelUserGroup) {
	if err := o.AddRelUserGroups(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddRelUserGroupsGP adds the given related objects to the existing relationships
// of the group, optionally inserting them as new records.
// Appends related to o.R.RelUserGroups.
// Sets related.R.Group appropriately.
// Uses the global database handle and panics on error.
func (o *Group) AddRelUserGroupsGP(insert bool, related ...*RelUserGroup) {
	if err := o.AddRelUserGroups(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddRelUserGroups adds the given related objects to the existing relationships
// of the group, optionally inserting them as new records.
// Appends related to o.R.RelUserGroups.
// Sets related.R.Group appropriately.
func (o *Group) AddRelUserGroups(exec boil.Executor, insert bool, related ...*RelUserGroup) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.GroupID = o.Snowflake
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"rel_user_groups\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"group_id"}),
				strmangle.WhereClause("\"", "\"", 2, relUserGroupPrimaryKeyColumns),
			)
			values := []interface{}{o.Snowflake, rel.UserID, rel.GroupID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.GroupID = o.Snowflake
		}
	}

	if o.R == nil {
		o.R = &groupR{
			RelUserGroups: related,
		}
	} else {
		o.R.RelUserGroups = append(o.R.RelUserGroups, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &relUserGroupR{
				Group: o,
			}
		} else {
			rel.R.Group = o
		}
	}
	return nil
}

// GroupsG retrieves all records.
func GroupsG(mods ...qm.QueryMod) groupQuery {
	return Groups(boil.GetDB(), mods...)
}

// Groups retrieves all the records using an executor.
func Groups(exec boil.Executor, mods ...qm.QueryMod) groupQuery {
	mods = append(mods, qm.From("\"groups\""))
	return groupQuery{NewQuery(exec, mods...)}
}

// FindGroupG retrieves a single record by ID.
func FindGroupG(snowflake int64, selectCols ...string) (*Group, error) {
	return FindGroup(boil.GetDB(), snowflake, selectCols...)
}

// FindGroupGP retrieves a single record by ID, and panics on error.
func FindGroupGP(snowflake int64, selectCols ...string) *Group {
	retobj, err := FindGroup(boil.GetDB(), snowflake, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindGroup retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindGroup(exec boil.Executor, snowflake int64, selectCols ...string) (*Group, error) {
	groupObj := &Group{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"groups\" where \"snowflake\"=$1", sel,
	)

	q := queries.Raw(exec, query, snowflake)

	err := q.Bind(groupObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from groups")
	}

	return groupObj, nil
}

// FindGroupP retrieves a single record by ID with an executor, and panics on error.
func FindGroupP(exec boil.Executor, snowflake int64, selectCols ...string) *Group {
	retobj, err := FindGroup(exec, snowflake, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Group) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Group) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Group) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Group) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no groups provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(groupColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	groupInsertCacheMut.RLock()
	cache, cached := groupInsertCache[key]
	groupInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			groupColumns,
			groupColumnsWithDefault,
			groupColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(groupType, groupMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(groupType, groupMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"groups\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into groups")
	}

	if !cached {
		groupInsertCacheMut.Lock()
		groupInsertCache[key] = cache
		groupInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Group record. See Update for
// whitelist behavior description.
func (o *Group) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Group record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Group) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Group, and panics on error.
// See Update for whitelist behavior description.
func (o *Group) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Group.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Group) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	groupUpdateCacheMut.RLock()
	cache, cached := groupUpdateCache[key]
	groupUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(groupColumns, groupPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update groups, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"groups\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, groupPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(groupType, groupMapping, append(wl, groupPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update groups row")
	}

	if !cached {
		groupUpdateCacheMut.Lock()
		groupUpdateCache[key] = cache
		groupUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q groupQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q groupQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for groups")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o GroupSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o GroupSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o GroupSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o GroupSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), groupPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"groups\" SET %s WHERE (\"snowflake\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(groupPrimaryKeyColumns), len(colNames)+1, len(groupPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in group slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Group) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Group) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Group) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Group) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no groups provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(groupColumnsWithDefault, o)

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

	groupUpsertCacheMut.RLock()
	cache, cached := groupUpsertCache[key]
	groupUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			groupColumns,
			groupColumnsWithDefault,
			groupColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			groupColumns,
			groupPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert groups, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(groupPrimaryKeyColumns))
			copy(conflict, groupPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"groups\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(groupType, groupMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(groupType, groupMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert groups")
	}

	if !cached {
		groupUpsertCacheMut.Lock()
		groupUpsertCache[key] = cache
		groupUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Group record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Group) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Group record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Group) DeleteG() error {
	if o == nil {
		return errors.New("models: no Group provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Group record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Group) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Group record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Group) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Group provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), groupPrimaryKeyMapping)
	sql := "DELETE FROM \"groups\" WHERE \"snowflake\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from groups")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q groupQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q groupQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no groupQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from groups")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o GroupSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o GroupSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Group slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o GroupSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o GroupSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Group slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(groupBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), groupPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"groups\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, groupPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(groupPrimaryKeyColumns), 1, len(groupPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from group slice")
	}

	if len(groupAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Group) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Group) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Group) ReloadG() error {
	if o == nil {
		return errors.New("models: no Group provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Group) Reload(exec boil.Executor) error {
	ret, err := FindGroup(exec, o.Snowflake)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *GroupSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *GroupSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GroupSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty GroupSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *GroupSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	groups := GroupSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), groupPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"groups\".* FROM \"groups\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, groupPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(groupPrimaryKeyColumns), 1, len(groupPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&groups)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in GroupSlice")
	}

	*o = groups

	return nil
}

// GroupExists checks if the Group row exists.
func GroupExists(exec boil.Executor, snowflake int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"groups\" where \"snowflake\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, snowflake)
	}

	row := exec.QueryRow(sql, snowflake)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if groups exists")
	}

	return exists, nil
}

// GroupExistsG checks if the Group row exists.
func GroupExistsG(snowflake int64) (bool, error) {
	return GroupExists(boil.GetDB(), snowflake)
}

// GroupExistsGP checks if the Group row exists. Panics on error.
func GroupExistsGP(snowflake int64) bool {
	e, err := GroupExists(boil.GetDB(), snowflake)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// GroupExistsP checks if the Group row exists. Panics on error.
func GroupExistsP(exec boil.Executor, snowflake int64) bool {
	e, err := GroupExists(exec, snowflake)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
