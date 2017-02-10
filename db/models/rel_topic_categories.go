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

// RelTopicCategory is an object representing the database table.
type RelTopicCategory struct {
	TopicID    int64     `boil:"topic_id" json:"topic_id" toml:"topic_id" yaml:"topic_id"`
	CategoryID int64     `boil:"category_id" json:"category_id" toml:"category_id" yaml:"category_id"`
	CreatedAt  time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	DeletedAt  null.Time `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`

	R *relTopicCategoryR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L relTopicCategoryL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// relTopicCategoryR is where relationships are stored.
type relTopicCategoryR struct {
	Topic    *Topic
	Category *Category
}

// relTopicCategoryL is where Load methods for each relationship are stored.
type relTopicCategoryL struct{}

var (
	relTopicCategoryColumns               = []string{"topic_id", "category_id", "created_at", "deleted_at"}
	relTopicCategoryColumnsWithoutDefault = []string{"topic_id", "category_id", "deleted_at"}
	relTopicCategoryColumnsWithDefault    = []string{"created_at"}
	relTopicCategoryPrimaryKeyColumns     = []string{"topic_id", "category_id"}
)

type (
	// RelTopicCategorySlice is an alias for a slice of pointers to RelTopicCategory.
	// This should generally be used opposed to []RelTopicCategory.
	RelTopicCategorySlice []*RelTopicCategory
	// RelTopicCategoryHook is the signature for custom RelTopicCategory hook methods
	RelTopicCategoryHook func(boil.Executor, *RelTopicCategory) error

	relTopicCategoryQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	relTopicCategoryType                 = reflect.TypeOf(&RelTopicCategory{})
	relTopicCategoryMapping              = queries.MakeStructMapping(relTopicCategoryType)
	relTopicCategoryPrimaryKeyMapping, _ = queries.BindMapping(relTopicCategoryType, relTopicCategoryMapping, relTopicCategoryPrimaryKeyColumns)
	relTopicCategoryInsertCacheMut       sync.RWMutex
	relTopicCategoryInsertCache          = make(map[string]insertCache)
	relTopicCategoryUpdateCacheMut       sync.RWMutex
	relTopicCategoryUpdateCache          = make(map[string]updateCache)
	relTopicCategoryUpsertCacheMut       sync.RWMutex
	relTopicCategoryUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var relTopicCategoryBeforeInsertHooks []RelTopicCategoryHook
var relTopicCategoryBeforeUpdateHooks []RelTopicCategoryHook
var relTopicCategoryBeforeDeleteHooks []RelTopicCategoryHook
var relTopicCategoryBeforeUpsertHooks []RelTopicCategoryHook

var relTopicCategoryAfterInsertHooks []RelTopicCategoryHook
var relTopicCategoryAfterSelectHooks []RelTopicCategoryHook
var relTopicCategoryAfterUpdateHooks []RelTopicCategoryHook
var relTopicCategoryAfterDeleteHooks []RelTopicCategoryHook
var relTopicCategoryAfterUpsertHooks []RelTopicCategoryHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *RelTopicCategory) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range relTopicCategoryBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *RelTopicCategory) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range relTopicCategoryBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *RelTopicCategory) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range relTopicCategoryBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *RelTopicCategory) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range relTopicCategoryBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *RelTopicCategory) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range relTopicCategoryAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *RelTopicCategory) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range relTopicCategoryAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *RelTopicCategory) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range relTopicCategoryAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *RelTopicCategory) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range relTopicCategoryAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *RelTopicCategory) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range relTopicCategoryAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddRelTopicCategoryHook registers your hook function for all future operations.
func AddRelTopicCategoryHook(hookPoint boil.HookPoint, relTopicCategoryHook RelTopicCategoryHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		relTopicCategoryBeforeInsertHooks = append(relTopicCategoryBeforeInsertHooks, relTopicCategoryHook)
	case boil.BeforeUpdateHook:
		relTopicCategoryBeforeUpdateHooks = append(relTopicCategoryBeforeUpdateHooks, relTopicCategoryHook)
	case boil.BeforeDeleteHook:
		relTopicCategoryBeforeDeleteHooks = append(relTopicCategoryBeforeDeleteHooks, relTopicCategoryHook)
	case boil.BeforeUpsertHook:
		relTopicCategoryBeforeUpsertHooks = append(relTopicCategoryBeforeUpsertHooks, relTopicCategoryHook)
	case boil.AfterInsertHook:
		relTopicCategoryAfterInsertHooks = append(relTopicCategoryAfterInsertHooks, relTopicCategoryHook)
	case boil.AfterSelectHook:
		relTopicCategoryAfterSelectHooks = append(relTopicCategoryAfterSelectHooks, relTopicCategoryHook)
	case boil.AfterUpdateHook:
		relTopicCategoryAfterUpdateHooks = append(relTopicCategoryAfterUpdateHooks, relTopicCategoryHook)
	case boil.AfterDeleteHook:
		relTopicCategoryAfterDeleteHooks = append(relTopicCategoryAfterDeleteHooks, relTopicCategoryHook)
	case boil.AfterUpsertHook:
		relTopicCategoryAfterUpsertHooks = append(relTopicCategoryAfterUpsertHooks, relTopicCategoryHook)
	}
}

// OneP returns a single relTopicCategory record from the query, and panics on error.
func (q relTopicCategoryQuery) OneP() *RelTopicCategory {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single relTopicCategory record from the query.
func (q relTopicCategoryQuery) One() (*RelTopicCategory, error) {
	o := &RelTopicCategory{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for rel_topic_categories")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all RelTopicCategory records from the query, and panics on error.
func (q relTopicCategoryQuery) AllP() RelTopicCategorySlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all RelTopicCategory records from the query.
func (q relTopicCategoryQuery) All() (RelTopicCategorySlice, error) {
	var o RelTopicCategorySlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to RelTopicCategory slice")
	}

	if len(relTopicCategoryAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all RelTopicCategory records in the query, and panics on error.
func (q relTopicCategoryQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all RelTopicCategory records in the query.
func (q relTopicCategoryQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count rel_topic_categories rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q relTopicCategoryQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q relTopicCategoryQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if rel_topic_categories exists")
	}

	return count > 0, nil
}

// TopicG pointed to by the foreign key.
func (o *RelTopicCategory) TopicG(mods ...qm.QueryMod) topicQuery {
	return o.Topic(boil.GetDB(), mods...)
}

// Topic pointed to by the foreign key.
func (o *RelTopicCategory) Topic(exec boil.Executor, mods ...qm.QueryMod) topicQuery {
	queryMods := []qm.QueryMod{
		qm.Where("snowflake=?", o.TopicID),
	}

	queryMods = append(queryMods, mods...)

	query := Topics(exec, queryMods...)
	queries.SetFrom(query.Query, "\"topics\"")

	return query
}

// CategoryG pointed to by the foreign key.
func (o *RelTopicCategory) CategoryG(mods ...qm.QueryMod) categoryQuery {
	return o.Category(boil.GetDB(), mods...)
}

// Category pointed to by the foreign key.
func (o *RelTopicCategory) Category(exec boil.Executor, mods ...qm.QueryMod) categoryQuery {
	queryMods := []qm.QueryMod{
		qm.Where("snowflake=?", o.CategoryID),
	}

	queryMods = append(queryMods, mods...)

	query := Categories(exec, queryMods...)
	queries.SetFrom(query.Query, "\"categories\"")

	return query
}

// LoadTopic allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (relTopicCategoryL) LoadTopic(e boil.Executor, singular bool, maybeRelTopicCategory interface{}) error {
	var slice []*RelTopicCategory
	var object *RelTopicCategory

	count := 1
	if singular {
		object = maybeRelTopicCategory.(*RelTopicCategory)
	} else {
		slice = *maybeRelTopicCategory.(*RelTopicCategorySlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &relTopicCategoryR{}
		}
		args[0] = object.TopicID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &relTopicCategoryR{}
			}
			args[i] = obj.TopicID
		}
	}

	query := fmt.Sprintf(
		"select * from \"topics\" where \"snowflake\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Topic")
	}
	defer results.Close()

	var resultSlice []*Topic
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Topic")
	}

	if len(relTopicCategoryAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		object.R.Topic = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.TopicID == foreign.Snowflake {
				local.R.Topic = foreign
				break
			}
		}
	}

	return nil
}

// LoadCategory allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (relTopicCategoryL) LoadCategory(e boil.Executor, singular bool, maybeRelTopicCategory interface{}) error {
	var slice []*RelTopicCategory
	var object *RelTopicCategory

	count := 1
	if singular {
		object = maybeRelTopicCategory.(*RelTopicCategory)
	} else {
		slice = *maybeRelTopicCategory.(*RelTopicCategorySlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &relTopicCategoryR{}
		}
		args[0] = object.CategoryID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &relTopicCategoryR{}
			}
			args[i] = obj.CategoryID
		}
	}

	query := fmt.Sprintf(
		"select * from \"categories\" where \"snowflake\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Category")
	}
	defer results.Close()

	var resultSlice []*Category
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Category")
	}

	if len(relTopicCategoryAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		object.R.Category = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.CategoryID == foreign.Snowflake {
				local.R.Category = foreign
				break
			}
		}
	}

	return nil
}

// SetTopicG of the rel_topic_category to the related item.
// Sets o.R.Topic to related.
// Adds o to related.R.RelTopicCategories.
// Uses the global database handle.
func (o *RelTopicCategory) SetTopicG(insert bool, related *Topic) error {
	return o.SetTopic(boil.GetDB(), insert, related)
}

// SetTopicP of the rel_topic_category to the related item.
// Sets o.R.Topic to related.
// Adds o to related.R.RelTopicCategories.
// Panics on error.
func (o *RelTopicCategory) SetTopicP(exec boil.Executor, insert bool, related *Topic) {
	if err := o.SetTopic(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetTopicGP of the rel_topic_category to the related item.
// Sets o.R.Topic to related.
// Adds o to related.R.RelTopicCategories.
// Uses the global database handle and panics on error.
func (o *RelTopicCategory) SetTopicGP(insert bool, related *Topic) {
	if err := o.SetTopic(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetTopic of the rel_topic_category to the related item.
// Sets o.R.Topic to related.
// Adds o to related.R.RelTopicCategories.
func (o *RelTopicCategory) SetTopic(exec boil.Executor, insert bool, related *Topic) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"rel_topic_categories\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"topic_id"}),
		strmangle.WhereClause("\"", "\"", 2, relTopicCategoryPrimaryKeyColumns),
	)
	values := []interface{}{related.Snowflake, o.TopicID, o.CategoryID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.TopicID = related.Snowflake

	if o.R == nil {
		o.R = &relTopicCategoryR{
			Topic: related,
		}
	} else {
		o.R.Topic = related
	}

	if related.R == nil {
		related.R = &topicR{
			RelTopicCategories: RelTopicCategorySlice{o},
		}
	} else {
		related.R.RelTopicCategories = append(related.R.RelTopicCategories, o)
	}

	return nil
}

// SetCategoryG of the rel_topic_category to the related item.
// Sets o.R.Category to related.
// Adds o to related.R.RelTopicCategories.
// Uses the global database handle.
func (o *RelTopicCategory) SetCategoryG(insert bool, related *Category) error {
	return o.SetCategory(boil.GetDB(), insert, related)
}

// SetCategoryP of the rel_topic_category to the related item.
// Sets o.R.Category to related.
// Adds o to related.R.RelTopicCategories.
// Panics on error.
func (o *RelTopicCategory) SetCategoryP(exec boil.Executor, insert bool, related *Category) {
	if err := o.SetCategory(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetCategoryGP of the rel_topic_category to the related item.
// Sets o.R.Category to related.
// Adds o to related.R.RelTopicCategories.
// Uses the global database handle and panics on error.
func (o *RelTopicCategory) SetCategoryGP(insert bool, related *Category) {
	if err := o.SetCategory(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetCategory of the rel_topic_category to the related item.
// Sets o.R.Category to related.
// Adds o to related.R.RelTopicCategories.
func (o *RelTopicCategory) SetCategory(exec boil.Executor, insert bool, related *Category) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"rel_topic_categories\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"category_id"}),
		strmangle.WhereClause("\"", "\"", 2, relTopicCategoryPrimaryKeyColumns),
	)
	values := []interface{}{related.Snowflake, o.TopicID, o.CategoryID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.CategoryID = related.Snowflake

	if o.R == nil {
		o.R = &relTopicCategoryR{
			Category: related,
		}
	} else {
		o.R.Category = related
	}

	if related.R == nil {
		related.R = &categoryR{
			RelTopicCategories: RelTopicCategorySlice{o},
		}
	} else {
		related.R.RelTopicCategories = append(related.R.RelTopicCategories, o)
	}

	return nil
}

// RelTopicCategoriesG retrieves all records.
func RelTopicCategoriesG(mods ...qm.QueryMod) relTopicCategoryQuery {
	return RelTopicCategories(boil.GetDB(), mods...)
}

// RelTopicCategories retrieves all the records using an executor.
func RelTopicCategories(exec boil.Executor, mods ...qm.QueryMod) relTopicCategoryQuery {
	mods = append(mods, qm.From("\"rel_topic_categories\""))
	return relTopicCategoryQuery{NewQuery(exec, mods...)}
}

// FindRelTopicCategoryG retrieves a single record by ID.
func FindRelTopicCategoryG(topicID int64, categoryID int64, selectCols ...string) (*RelTopicCategory, error) {
	return FindRelTopicCategory(boil.GetDB(), topicID, categoryID, selectCols...)
}

// FindRelTopicCategoryGP retrieves a single record by ID, and panics on error.
func FindRelTopicCategoryGP(topicID int64, categoryID int64, selectCols ...string) *RelTopicCategory {
	retobj, err := FindRelTopicCategory(boil.GetDB(), topicID, categoryID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindRelTopicCategory retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindRelTopicCategory(exec boil.Executor, topicID int64, categoryID int64, selectCols ...string) (*RelTopicCategory, error) {
	relTopicCategoryObj := &RelTopicCategory{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"rel_topic_categories\" where \"topic_id\"=$1 AND \"category_id\"=$2", sel,
	)

	q := queries.Raw(exec, query, topicID, categoryID)

	err := q.Bind(relTopicCategoryObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from rel_topic_categories")
	}

	return relTopicCategoryObj, nil
}

// FindRelTopicCategoryP retrieves a single record by ID with an executor, and panics on error.
func FindRelTopicCategoryP(exec boil.Executor, topicID int64, categoryID int64, selectCols ...string) *RelTopicCategory {
	retobj, err := FindRelTopicCategory(exec, topicID, categoryID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *RelTopicCategory) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *RelTopicCategory) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *RelTopicCategory) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *RelTopicCategory) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no rel_topic_categories provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(relTopicCategoryColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	relTopicCategoryInsertCacheMut.RLock()
	cache, cached := relTopicCategoryInsertCache[key]
	relTopicCategoryInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			relTopicCategoryColumns,
			relTopicCategoryColumnsWithDefault,
			relTopicCategoryColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(relTopicCategoryType, relTopicCategoryMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(relTopicCategoryType, relTopicCategoryMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"rel_topic_categories\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into rel_topic_categories")
	}

	if !cached {
		relTopicCategoryInsertCacheMut.Lock()
		relTopicCategoryInsertCache[key] = cache
		relTopicCategoryInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single RelTopicCategory record. See Update for
// whitelist behavior description.
func (o *RelTopicCategory) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single RelTopicCategory record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *RelTopicCategory) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the RelTopicCategory, and panics on error.
// See Update for whitelist behavior description.
func (o *RelTopicCategory) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the RelTopicCategory.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *RelTopicCategory) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	relTopicCategoryUpdateCacheMut.RLock()
	cache, cached := relTopicCategoryUpdateCache[key]
	relTopicCategoryUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(relTopicCategoryColumns, relTopicCategoryPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update rel_topic_categories, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"rel_topic_categories\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, relTopicCategoryPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(relTopicCategoryType, relTopicCategoryMapping, append(wl, relTopicCategoryPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update rel_topic_categories row")
	}

	if !cached {
		relTopicCategoryUpdateCacheMut.Lock()
		relTopicCategoryUpdateCache[key] = cache
		relTopicCategoryUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q relTopicCategoryQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q relTopicCategoryQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for rel_topic_categories")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o RelTopicCategorySlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o RelTopicCategorySlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o RelTopicCategorySlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o RelTopicCategorySlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), relTopicCategoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"rel_topic_categories\" SET %s WHERE (\"topic_id\",\"category_id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(relTopicCategoryPrimaryKeyColumns), len(colNames)+1, len(relTopicCategoryPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in relTopicCategory slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *RelTopicCategory) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *RelTopicCategory) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *RelTopicCategory) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *RelTopicCategory) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no rel_topic_categories provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(relTopicCategoryColumnsWithDefault, o)

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

	relTopicCategoryUpsertCacheMut.RLock()
	cache, cached := relTopicCategoryUpsertCache[key]
	relTopicCategoryUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			relTopicCategoryColumns,
			relTopicCategoryColumnsWithDefault,
			relTopicCategoryColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			relTopicCategoryColumns,
			relTopicCategoryPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert rel_topic_categories, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(relTopicCategoryPrimaryKeyColumns))
			copy(conflict, relTopicCategoryPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"rel_topic_categories\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(relTopicCategoryType, relTopicCategoryMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(relTopicCategoryType, relTopicCategoryMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert rel_topic_categories")
	}

	if !cached {
		relTopicCategoryUpsertCacheMut.Lock()
		relTopicCategoryUpsertCache[key] = cache
		relTopicCategoryUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single RelTopicCategory record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *RelTopicCategory) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single RelTopicCategory record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *RelTopicCategory) DeleteG() error {
	if o == nil {
		return errors.New("models: no RelTopicCategory provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single RelTopicCategory record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *RelTopicCategory) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single RelTopicCategory record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *RelTopicCategory) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no RelTopicCategory provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), relTopicCategoryPrimaryKeyMapping)
	sql := "DELETE FROM \"rel_topic_categories\" WHERE \"topic_id\"=$1 AND \"category_id\"=$2"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from rel_topic_categories")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q relTopicCategoryQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q relTopicCategoryQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no relTopicCategoryQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from rel_topic_categories")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o RelTopicCategorySlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o RelTopicCategorySlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no RelTopicCategory slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o RelTopicCategorySlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o RelTopicCategorySlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no RelTopicCategory slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(relTopicCategoryBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), relTopicCategoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"rel_topic_categories\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, relTopicCategoryPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(relTopicCategoryPrimaryKeyColumns), 1, len(relTopicCategoryPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from relTopicCategory slice")
	}

	if len(relTopicCategoryAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *RelTopicCategory) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *RelTopicCategory) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *RelTopicCategory) ReloadG() error {
	if o == nil {
		return errors.New("models: no RelTopicCategory provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *RelTopicCategory) Reload(exec boil.Executor) error {
	ret, err := FindRelTopicCategory(exec, o.TopicID, o.CategoryID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *RelTopicCategorySlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *RelTopicCategorySlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *RelTopicCategorySlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty RelTopicCategorySlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *RelTopicCategorySlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	relTopicCategories := RelTopicCategorySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), relTopicCategoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"rel_topic_categories\".* FROM \"rel_topic_categories\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, relTopicCategoryPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(relTopicCategoryPrimaryKeyColumns), 1, len(relTopicCategoryPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&relTopicCategories)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in RelTopicCategorySlice")
	}

	*o = relTopicCategories

	return nil
}

// RelTopicCategoryExists checks if the RelTopicCategory row exists.
func RelTopicCategoryExists(exec boil.Executor, topicID int64, categoryID int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"rel_topic_categories\" where \"topic_id\"=$1 AND \"category_id\"=$2 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, topicID, categoryID)
	}

	row := exec.QueryRow(sql, topicID, categoryID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if rel_topic_categories exists")
	}

	return exists, nil
}

// RelTopicCategoryExistsG checks if the RelTopicCategory row exists.
func RelTopicCategoryExistsG(topicID int64, categoryID int64) (bool, error) {
	return RelTopicCategoryExists(boil.GetDB(), topicID, categoryID)
}

// RelTopicCategoryExistsGP checks if the RelTopicCategory row exists. Panics on error.
func RelTopicCategoryExistsGP(topicID int64, categoryID int64) bool {
	e, err := RelTopicCategoryExists(boil.GetDB(), topicID, categoryID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// RelTopicCategoryExistsP checks if the RelTopicCategory row exists. Panics on error.
func RelTopicCategoryExistsP(exec boil.Executor, topicID int64, categoryID int64) bool {
	e, err := RelTopicCategoryExists(exec, topicID, categoryID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
