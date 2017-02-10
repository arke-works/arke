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

// Topic is an object representing the database table.
type Topic struct {
	Snowflake int64      `boil:"snowflake" json:"snowflake" toml:"snowflake" yaml:"snowflake"`
	CreatedAt time.Time  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	DeletedAt null.Time  `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`
	AuthorID  null.Int64 `boil:"author_id" json:"author_id,omitempty" toml:"author_id" yaml:"author_id,omitempty"`
	Title     string     `boil:"title" json:"title" toml:"title" yaml:"title"`
	Body      string     `boil:"body" json:"body" toml:"body" yaml:"body"`
	Revision  int64      `boil:"revision" json:"revision" toml:"revision" yaml:"revision"`

	R *topicR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L topicL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// topicR is where relationships are stored.
type topicR struct {
	Author             *User
	RelTopicCategories RelTopicCategorySlice
	Replies            ReplySlice
}

// topicL is where Load methods for each relationship are stored.
type topicL struct{}

var (
	topicColumns               = []string{"snowflake", "created_at", "deleted_at", "author_id", "title", "body", "revision"}
	topicColumnsWithoutDefault = []string{"snowflake", "deleted_at", "author_id", "title", "body", "revision"}
	topicColumnsWithDefault    = []string{"created_at"}
	topicPrimaryKeyColumns     = []string{"snowflake"}
)

type (
	// TopicSlice is an alias for a slice of pointers to Topic.
	// This should generally be used opposed to []Topic.
	TopicSlice []*Topic
	// TopicHook is the signature for custom Topic hook methods
	TopicHook func(boil.Executor, *Topic) error

	topicQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	topicType                 = reflect.TypeOf(&Topic{})
	topicMapping              = queries.MakeStructMapping(topicType)
	topicPrimaryKeyMapping, _ = queries.BindMapping(topicType, topicMapping, topicPrimaryKeyColumns)
	topicInsertCacheMut       sync.RWMutex
	topicInsertCache          = make(map[string]insertCache)
	topicUpdateCacheMut       sync.RWMutex
	topicUpdateCache          = make(map[string]updateCache)
	topicUpsertCacheMut       sync.RWMutex
	topicUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var topicBeforeInsertHooks []TopicHook
var topicBeforeUpdateHooks []TopicHook
var topicBeforeDeleteHooks []TopicHook
var topicBeforeUpsertHooks []TopicHook

var topicAfterInsertHooks []TopicHook
var topicAfterSelectHooks []TopicHook
var topicAfterUpdateHooks []TopicHook
var topicAfterDeleteHooks []TopicHook
var topicAfterUpsertHooks []TopicHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Topic) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range topicBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Topic) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range topicBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Topic) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range topicBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Topic) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range topicBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Topic) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range topicAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Topic) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range topicAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Topic) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range topicAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Topic) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range topicAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Topic) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range topicAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddTopicHook registers your hook function for all future operations.
func AddTopicHook(hookPoint boil.HookPoint, topicHook TopicHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		topicBeforeInsertHooks = append(topicBeforeInsertHooks, topicHook)
	case boil.BeforeUpdateHook:
		topicBeforeUpdateHooks = append(topicBeforeUpdateHooks, topicHook)
	case boil.BeforeDeleteHook:
		topicBeforeDeleteHooks = append(topicBeforeDeleteHooks, topicHook)
	case boil.BeforeUpsertHook:
		topicBeforeUpsertHooks = append(topicBeforeUpsertHooks, topicHook)
	case boil.AfterInsertHook:
		topicAfterInsertHooks = append(topicAfterInsertHooks, topicHook)
	case boil.AfterSelectHook:
		topicAfterSelectHooks = append(topicAfterSelectHooks, topicHook)
	case boil.AfterUpdateHook:
		topicAfterUpdateHooks = append(topicAfterUpdateHooks, topicHook)
	case boil.AfterDeleteHook:
		topicAfterDeleteHooks = append(topicAfterDeleteHooks, topicHook)
	case boil.AfterUpsertHook:
		topicAfterUpsertHooks = append(topicAfterUpsertHooks, topicHook)
	}
}

// OneP returns a single topic record from the query, and panics on error.
func (q topicQuery) OneP() *Topic {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single topic record from the query.
func (q topicQuery) One() (*Topic, error) {
	o := &Topic{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for topics")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Topic records from the query, and panics on error.
func (q topicQuery) AllP() TopicSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Topic records from the query.
func (q topicQuery) All() (TopicSlice, error) {
	var o TopicSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Topic slice")
	}

	if len(topicAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Topic records in the query, and panics on error.
func (q topicQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Topic records in the query.
func (q topicQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count topics rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q topicQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q topicQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if topics exists")
	}

	return count > 0, nil
}

// AuthorG pointed to by the foreign key.
func (o *Topic) AuthorG(mods ...qm.QueryMod) userQuery {
	return o.Author(boil.GetDB(), mods...)
}

// Author pointed to by the foreign key.
func (o *Topic) Author(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("snowflake=?", o.AuthorID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(exec, queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// RelTopicCategoriesG retrieves all the rel_topic_category's rel topic categories.
func (o *Topic) RelTopicCategoriesG(mods ...qm.QueryMod) relTopicCategoryQuery {
	return o.RelTopicCategories(boil.GetDB(), mods...)
}

// RelTopicCategories retrieves all the rel_topic_category's rel topic categories with an executor.
func (o *Topic) RelTopicCategories(exec boil.Executor, mods ...qm.QueryMod) relTopicCategoryQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"topic_id\"=?", o.Snowflake),
	)

	query := RelTopicCategories(exec, queryMods...)
	queries.SetFrom(query.Query, "\"rel_topic_categories\" as \"a\"")
	return query
}

// RepliesG retrieves all the reply's replies.
func (o *Topic) RepliesG(mods ...qm.QueryMod) replyQuery {
	return o.Replies(boil.GetDB(), mods...)
}

// Replies retrieves all the reply's replies with an executor.
func (o *Topic) Replies(exec boil.Executor, mods ...qm.QueryMod) replyQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"topic_id\"=?", o.Snowflake),
	)

	query := Replies(exec, queryMods...)
	queries.SetFrom(query.Query, "\"replies\" as \"a\"")
	return query
}

// LoadAuthor allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (topicL) LoadAuthor(e boil.Executor, singular bool, maybeTopic interface{}) error {
	var slice []*Topic
	var object *Topic

	count := 1
	if singular {
		object = maybeTopic.(*Topic)
	} else {
		slice = *maybeTopic.(*TopicSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &topicR{}
		}
		args[0] = object.AuthorID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &topicR{}
			}
			args[i] = obj.AuthorID
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

	if len(topicAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		object.R.Author = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.AuthorID.Int64 == foreign.Snowflake {
				local.R.Author = foreign
				break
			}
		}
	}

	return nil
}

// LoadRelTopicCategories allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (topicL) LoadRelTopicCategories(e boil.Executor, singular bool, maybeTopic interface{}) error {
	var slice []*Topic
	var object *Topic

	count := 1
	if singular {
		object = maybeTopic.(*Topic)
	} else {
		slice = *maybeTopic.(*TopicSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &topicR{}
		}
		args[0] = object.Snowflake
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &topicR{}
			}
			args[i] = obj.Snowflake
		}
	}

	query := fmt.Sprintf(
		"select * from \"rel_topic_categories\" where \"topic_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load rel_topic_categories")
	}
	defer results.Close()

	var resultSlice []*RelTopicCategory
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice rel_topic_categories")
	}

	if len(relTopicCategoryAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.RelTopicCategories = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.Snowflake == foreign.TopicID {
				local.R.RelTopicCategories = append(local.R.RelTopicCategories, foreign)
				break
			}
		}
	}

	return nil
}

// LoadReplies allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (topicL) LoadReplies(e boil.Executor, singular bool, maybeTopic interface{}) error {
	var slice []*Topic
	var object *Topic

	count := 1
	if singular {
		object = maybeTopic.(*Topic)
	} else {
		slice = *maybeTopic.(*TopicSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &topicR{}
		}
		args[0] = object.Snowflake
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &topicR{}
			}
			args[i] = obj.Snowflake
		}
	}

	query := fmt.Sprintf(
		"select * from \"replies\" where \"topic_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load replies")
	}
	defer results.Close()

	var resultSlice []*Reply
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice replies")
	}

	if len(replyAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Replies = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.Snowflake == foreign.TopicID {
				local.R.Replies = append(local.R.Replies, foreign)
				break
			}
		}
	}

	return nil
}

// SetAuthorG of the topic to the related item.
// Sets o.R.Author to related.
// Adds o to related.R.AuthorTopics.
// Uses the global database handle.
func (o *Topic) SetAuthorG(insert bool, related *User) error {
	return o.SetAuthor(boil.GetDB(), insert, related)
}

// SetAuthorP of the topic to the related item.
// Sets o.R.Author to related.
// Adds o to related.R.AuthorTopics.
// Panics on error.
func (o *Topic) SetAuthorP(exec boil.Executor, insert bool, related *User) {
	if err := o.SetAuthor(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetAuthorGP of the topic to the related item.
// Sets o.R.Author to related.
// Adds o to related.R.AuthorTopics.
// Uses the global database handle and panics on error.
func (o *Topic) SetAuthorGP(insert bool, related *User) {
	if err := o.SetAuthor(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetAuthor of the topic to the related item.
// Sets o.R.Author to related.
// Adds o to related.R.AuthorTopics.
func (o *Topic) SetAuthor(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"topics\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"author_id"}),
		strmangle.WhereClause("\"", "\"", 2, topicPrimaryKeyColumns),
	)
	values := []interface{}{related.Snowflake, o.Snowflake}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.AuthorID.Int64 = related.Snowflake
	o.AuthorID.Valid = true

	if o.R == nil {
		o.R = &topicR{
			Author: related,
		}
	} else {
		o.R.Author = related
	}

	if related.R == nil {
		related.R = &userR{
			AuthorTopics: TopicSlice{o},
		}
	} else {
		related.R.AuthorTopics = append(related.R.AuthorTopics, o)
	}

	return nil
}

// RemoveAuthorG relationship.
// Sets o.R.Author to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle.
func (o *Topic) RemoveAuthorG(related *User) error {
	return o.RemoveAuthor(boil.GetDB(), related)
}

// RemoveAuthorP relationship.
// Sets o.R.Author to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Panics on error.
func (o *Topic) RemoveAuthorP(exec boil.Executor, related *User) {
	if err := o.RemoveAuthor(exec, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveAuthorGP relationship.
// Sets o.R.Author to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle and panics on error.
func (o *Topic) RemoveAuthorGP(related *User) {
	if err := o.RemoveAuthor(boil.GetDB(), related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveAuthor relationship.
// Sets o.R.Author to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *Topic) RemoveAuthor(exec boil.Executor, related *User) error {
	var err error

	o.AuthorID.Valid = false
	if err = o.Update(exec, "author_id"); err != nil {
		o.AuthorID.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.Author = nil
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.AuthorTopics {
		if o.AuthorID.Int64 != ri.AuthorID.Int64 {
			continue
		}

		ln := len(related.R.AuthorTopics)
		if ln > 1 && i < ln-1 {
			related.R.AuthorTopics[i] = related.R.AuthorTopics[ln-1]
		}
		related.R.AuthorTopics = related.R.AuthorTopics[:ln-1]
		break
	}
	return nil
}

// AddRelTopicCategoriesG adds the given related objects to the existing relationships
// of the topic, optionally inserting them as new records.
// Appends related to o.R.RelTopicCategories.
// Sets related.R.Topic appropriately.
// Uses the global database handle.
func (o *Topic) AddRelTopicCategoriesG(insert bool, related ...*RelTopicCategory) error {
	return o.AddRelTopicCategories(boil.GetDB(), insert, related...)
}

// AddRelTopicCategoriesP adds the given related objects to the existing relationships
// of the topic, optionally inserting them as new records.
// Appends related to o.R.RelTopicCategories.
// Sets related.R.Topic appropriately.
// Panics on error.
func (o *Topic) AddRelTopicCategoriesP(exec boil.Executor, insert bool, related ...*RelTopicCategory) {
	if err := o.AddRelTopicCategories(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddRelTopicCategoriesGP adds the given related objects to the existing relationships
// of the topic, optionally inserting them as new records.
// Appends related to o.R.RelTopicCategories.
// Sets related.R.Topic appropriately.
// Uses the global database handle and panics on error.
func (o *Topic) AddRelTopicCategoriesGP(insert bool, related ...*RelTopicCategory) {
	if err := o.AddRelTopicCategories(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddRelTopicCategories adds the given related objects to the existing relationships
// of the topic, optionally inserting them as new records.
// Appends related to o.R.RelTopicCategories.
// Sets related.R.Topic appropriately.
func (o *Topic) AddRelTopicCategories(exec boil.Executor, insert bool, related ...*RelTopicCategory) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.TopicID = o.Snowflake
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"rel_topic_categories\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"topic_id"}),
				strmangle.WhereClause("\"", "\"", 2, relTopicCategoryPrimaryKeyColumns),
			)
			values := []interface{}{o.Snowflake, rel.TopicID, rel.CategoryID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.TopicID = o.Snowflake
		}
	}

	if o.R == nil {
		o.R = &topicR{
			RelTopicCategories: related,
		}
	} else {
		o.R.RelTopicCategories = append(o.R.RelTopicCategories, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &relTopicCategoryR{
				Topic: o,
			}
		} else {
			rel.R.Topic = o
		}
	}
	return nil
}

// AddRepliesG adds the given related objects to the existing relationships
// of the topic, optionally inserting them as new records.
// Appends related to o.R.Replies.
// Sets related.R.Topic appropriately.
// Uses the global database handle.
func (o *Topic) AddRepliesG(insert bool, related ...*Reply) error {
	return o.AddReplies(boil.GetDB(), insert, related...)
}

// AddRepliesP adds the given related objects to the existing relationships
// of the topic, optionally inserting them as new records.
// Appends related to o.R.Replies.
// Sets related.R.Topic appropriately.
// Panics on error.
func (o *Topic) AddRepliesP(exec boil.Executor, insert bool, related ...*Reply) {
	if err := o.AddReplies(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddRepliesGP adds the given related objects to the existing relationships
// of the topic, optionally inserting them as new records.
// Appends related to o.R.Replies.
// Sets related.R.Topic appropriately.
// Uses the global database handle and panics on error.
func (o *Topic) AddRepliesGP(insert bool, related ...*Reply) {
	if err := o.AddReplies(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddReplies adds the given related objects to the existing relationships
// of the topic, optionally inserting them as new records.
// Appends related to o.R.Replies.
// Sets related.R.Topic appropriately.
func (o *Topic) AddReplies(exec boil.Executor, insert bool, related ...*Reply) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.TopicID = o.Snowflake
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"replies\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"topic_id"}),
				strmangle.WhereClause("\"", "\"", 2, replyPrimaryKeyColumns),
			)
			values := []interface{}{o.Snowflake, rel.Snowflake}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.TopicID = o.Snowflake
		}
	}

	if o.R == nil {
		o.R = &topicR{
			Replies: related,
		}
	} else {
		o.R.Replies = append(o.R.Replies, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &replyR{
				Topic: o,
			}
		} else {
			rel.R.Topic = o
		}
	}
	return nil
}

// TopicsG retrieves all records.
func TopicsG(mods ...qm.QueryMod) topicQuery {
	return Topics(boil.GetDB(), mods...)
}

// Topics retrieves all the records using an executor.
func Topics(exec boil.Executor, mods ...qm.QueryMod) topicQuery {
	mods = append(mods, qm.From("\"topics\""))
	return topicQuery{NewQuery(exec, mods...)}
}

// FindTopicG retrieves a single record by ID.
func FindTopicG(snowflake int64, selectCols ...string) (*Topic, error) {
	return FindTopic(boil.GetDB(), snowflake, selectCols...)
}

// FindTopicGP retrieves a single record by ID, and panics on error.
func FindTopicGP(snowflake int64, selectCols ...string) *Topic {
	retobj, err := FindTopic(boil.GetDB(), snowflake, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindTopic retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindTopic(exec boil.Executor, snowflake int64, selectCols ...string) (*Topic, error) {
	topicObj := &Topic{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"topics\" where \"snowflake\"=$1", sel,
	)

	q := queries.Raw(exec, query, snowflake)

	err := q.Bind(topicObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from topics")
	}

	return topicObj, nil
}

// FindTopicP retrieves a single record by ID with an executor, and panics on error.
func FindTopicP(exec boil.Executor, snowflake int64, selectCols ...string) *Topic {
	retobj, err := FindTopic(exec, snowflake, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Topic) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Topic) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Topic) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Topic) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no topics provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(topicColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	topicInsertCacheMut.RLock()
	cache, cached := topicInsertCache[key]
	topicInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			topicColumns,
			topicColumnsWithDefault,
			topicColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(topicType, topicMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(topicType, topicMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"topics\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into topics")
	}

	if !cached {
		topicInsertCacheMut.Lock()
		topicInsertCache[key] = cache
		topicInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Topic record. See Update for
// whitelist behavior description.
func (o *Topic) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Topic record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Topic) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Topic, and panics on error.
// See Update for whitelist behavior description.
func (o *Topic) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Topic.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Topic) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	topicUpdateCacheMut.RLock()
	cache, cached := topicUpdateCache[key]
	topicUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(topicColumns, topicPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update topics, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"topics\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, topicPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(topicType, topicMapping, append(wl, topicPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update topics row")
	}

	if !cached {
		topicUpdateCacheMut.Lock()
		topicUpdateCache[key] = cache
		topicUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q topicQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q topicQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for topics")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o TopicSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o TopicSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o TopicSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o TopicSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), topicPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"topics\" SET %s WHERE (\"snowflake\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(topicPrimaryKeyColumns), len(colNames)+1, len(topicPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in topic slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Topic) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Topic) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Topic) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Topic) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no topics provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(topicColumnsWithDefault, o)

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

	topicUpsertCacheMut.RLock()
	cache, cached := topicUpsertCache[key]
	topicUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			topicColumns,
			topicColumnsWithDefault,
			topicColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			topicColumns,
			topicPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert topics, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(topicPrimaryKeyColumns))
			copy(conflict, topicPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"topics\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(topicType, topicMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(topicType, topicMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert topics")
	}

	if !cached {
		topicUpsertCacheMut.Lock()
		topicUpsertCache[key] = cache
		topicUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Topic record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Topic) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Topic record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Topic) DeleteG() error {
	if o == nil {
		return errors.New("models: no Topic provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Topic record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Topic) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Topic record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Topic) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Topic provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), topicPrimaryKeyMapping)
	sql := "DELETE FROM \"topics\" WHERE \"snowflake\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from topics")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q topicQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q topicQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no topicQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from topics")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o TopicSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o TopicSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Topic slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o TopicSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o TopicSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Topic slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(topicBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), topicPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"topics\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, topicPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(topicPrimaryKeyColumns), 1, len(topicPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from topic slice")
	}

	if len(topicAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Topic) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Topic) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Topic) ReloadG() error {
	if o == nil {
		return errors.New("models: no Topic provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Topic) Reload(exec boil.Executor) error {
	ret, err := FindTopic(exec, o.Snowflake)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *TopicSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *TopicSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *TopicSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty TopicSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *TopicSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	topics := TopicSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), topicPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"topics\".* FROM \"topics\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, topicPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(topicPrimaryKeyColumns), 1, len(topicPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&topics)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in TopicSlice")
	}

	*o = topics

	return nil
}

// TopicExists checks if the Topic row exists.
func TopicExists(exec boil.Executor, snowflake int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"topics\" where \"snowflake\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, snowflake)
	}

	row := exec.QueryRow(sql, snowflake)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if topics exists")
	}

	return exists, nil
}

// TopicExistsG checks if the Topic row exists.
func TopicExistsG(snowflake int64) (bool, error) {
	return TopicExists(boil.GetDB(), snowflake)
}

// TopicExistsGP checks if the Topic row exists. Panics on error.
func TopicExistsGP(snowflake int64) bool {
	e, err := TopicExists(boil.GetDB(), snowflake)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// TopicExistsP checks if the Topic row exists. Panics on error.
func TopicExistsP(exec boil.Executor, snowflake int64) bool {
	e, err := TopicExists(exec, snowflake)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
