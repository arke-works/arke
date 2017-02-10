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

// Reply is an object representing the database table.
type Reply struct {
	Snowflake int64      `boil:"snowflake" json:"snowflake" toml:"snowflake" yaml:"snowflake"`
	CreatedAt time.Time  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	DeletedAt null.Time  `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`
	AuthorID  null.Int64 `boil:"author_id" json:"author_id,omitempty" toml:"author_id" yaml:"author_id,omitempty"`
	Body      string     `boil:"body" json:"body" toml:"body" yaml:"body"`
	ParentID  null.Int64 `boil:"parent_id" json:"parent_id,omitempty" toml:"parent_id" yaml:"parent_id,omitempty"`
	TopicID   int64      `boil:"topic_id" json:"topic_id" toml:"topic_id" yaml:"topic_id"`

	R *replyR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L replyL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// replyR is where relationships are stored.
type replyR struct {
	Author        *User
	Parent        *Reply
	Topic         *Topic
	ParentReplies ReplySlice
}

// replyL is where Load methods for each relationship are stored.
type replyL struct{}

var (
	replyColumns               = []string{"snowflake", "created_at", "deleted_at", "author_id", "body", "parent_id", "topic_id"}
	replyColumnsWithoutDefault = []string{"snowflake", "deleted_at", "author_id", "body", "parent_id", "topic_id"}
	replyColumnsWithDefault    = []string{"created_at"}
	replyPrimaryKeyColumns     = []string{"snowflake"}
)

type (
	// ReplySlice is an alias for a slice of pointers to Reply.
	// This should generally be used opposed to []Reply.
	ReplySlice []*Reply
	// ReplyHook is the signature for custom Reply hook methods
	ReplyHook func(boil.Executor, *Reply) error

	replyQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	replyType                 = reflect.TypeOf(&Reply{})
	replyMapping              = queries.MakeStructMapping(replyType)
	replyPrimaryKeyMapping, _ = queries.BindMapping(replyType, replyMapping, replyPrimaryKeyColumns)
	replyInsertCacheMut       sync.RWMutex
	replyInsertCache          = make(map[string]insertCache)
	replyUpdateCacheMut       sync.RWMutex
	replyUpdateCache          = make(map[string]updateCache)
	replyUpsertCacheMut       sync.RWMutex
	replyUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var replyBeforeInsertHooks []ReplyHook
var replyBeforeUpdateHooks []ReplyHook
var replyBeforeDeleteHooks []ReplyHook
var replyBeforeUpsertHooks []ReplyHook

var replyAfterInsertHooks []ReplyHook
var replyAfterSelectHooks []ReplyHook
var replyAfterUpdateHooks []ReplyHook
var replyAfterDeleteHooks []ReplyHook
var replyAfterUpsertHooks []ReplyHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Reply) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range replyBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Reply) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range replyBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Reply) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range replyBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Reply) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range replyBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Reply) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range replyAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Reply) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range replyAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Reply) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range replyAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Reply) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range replyAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Reply) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range replyAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddReplyHook registers your hook function for all future operations.
func AddReplyHook(hookPoint boil.HookPoint, replyHook ReplyHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		replyBeforeInsertHooks = append(replyBeforeInsertHooks, replyHook)
	case boil.BeforeUpdateHook:
		replyBeforeUpdateHooks = append(replyBeforeUpdateHooks, replyHook)
	case boil.BeforeDeleteHook:
		replyBeforeDeleteHooks = append(replyBeforeDeleteHooks, replyHook)
	case boil.BeforeUpsertHook:
		replyBeforeUpsertHooks = append(replyBeforeUpsertHooks, replyHook)
	case boil.AfterInsertHook:
		replyAfterInsertHooks = append(replyAfterInsertHooks, replyHook)
	case boil.AfterSelectHook:
		replyAfterSelectHooks = append(replyAfterSelectHooks, replyHook)
	case boil.AfterUpdateHook:
		replyAfterUpdateHooks = append(replyAfterUpdateHooks, replyHook)
	case boil.AfterDeleteHook:
		replyAfterDeleteHooks = append(replyAfterDeleteHooks, replyHook)
	case boil.AfterUpsertHook:
		replyAfterUpsertHooks = append(replyAfterUpsertHooks, replyHook)
	}
}

// OneP returns a single reply record from the query, and panics on error.
func (q replyQuery) OneP() *Reply {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single reply record from the query.
func (q replyQuery) One() (*Reply, error) {
	o := &Reply{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for replies")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Reply records from the query, and panics on error.
func (q replyQuery) AllP() ReplySlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Reply records from the query.
func (q replyQuery) All() (ReplySlice, error) {
	var o ReplySlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Reply slice")
	}

	if len(replyAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Reply records in the query, and panics on error.
func (q replyQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Reply records in the query.
func (q replyQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count replies rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q replyQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q replyQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if replies exists")
	}

	return count > 0, nil
}

// AuthorG pointed to by the foreign key.
func (o *Reply) AuthorG(mods ...qm.QueryMod) userQuery {
	return o.Author(boil.GetDB(), mods...)
}

// Author pointed to by the foreign key.
func (o *Reply) Author(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("snowflake=?", o.AuthorID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(exec, queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// ParentG pointed to by the foreign key.
func (o *Reply) ParentG(mods ...qm.QueryMod) replyQuery {
	return o.Parent(boil.GetDB(), mods...)
}

// Parent pointed to by the foreign key.
func (o *Reply) Parent(exec boil.Executor, mods ...qm.QueryMod) replyQuery {
	queryMods := []qm.QueryMod{
		qm.Where("snowflake=?", o.ParentID),
	}

	queryMods = append(queryMods, mods...)

	query := Replies(exec, queryMods...)
	queries.SetFrom(query.Query, "\"replies\"")

	return query
}

// TopicG pointed to by the foreign key.
func (o *Reply) TopicG(mods ...qm.QueryMod) topicQuery {
	return o.Topic(boil.GetDB(), mods...)
}

// Topic pointed to by the foreign key.
func (o *Reply) Topic(exec boil.Executor, mods ...qm.QueryMod) topicQuery {
	queryMods := []qm.QueryMod{
		qm.Where("snowflake=?", o.TopicID),
	}

	queryMods = append(queryMods, mods...)

	query := Topics(exec, queryMods...)
	queries.SetFrom(query.Query, "\"topics\"")

	return query
}

// ParentRepliesG retrieves all the reply's replies via parent_id column.
func (o *Reply) ParentRepliesG(mods ...qm.QueryMod) replyQuery {
	return o.ParentReplies(boil.GetDB(), mods...)
}

// ParentReplies retrieves all the reply's replies with an executor via parent_id column.
func (o *Reply) ParentReplies(exec boil.Executor, mods ...qm.QueryMod) replyQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"parent_id\"=?", o.Snowflake),
	)

	query := Replies(exec, queryMods...)
	queries.SetFrom(query.Query, "\"replies\" as \"a\"")
	return query
}

// LoadAuthor allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (replyL) LoadAuthor(e boil.Executor, singular bool, maybeReply interface{}) error {
	var slice []*Reply
	var object *Reply

	count := 1
	if singular {
		object = maybeReply.(*Reply)
	} else {
		slice = *maybeReply.(*ReplySlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &replyR{}
		}
		args[0] = object.AuthorID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &replyR{}
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

	if len(replyAfterSelectHooks) != 0 {
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

// LoadParent allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (replyL) LoadParent(e boil.Executor, singular bool, maybeReply interface{}) error {
	var slice []*Reply
	var object *Reply

	count := 1
	if singular {
		object = maybeReply.(*Reply)
	} else {
		slice = *maybeReply.(*ReplySlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &replyR{}
		}
		args[0] = object.ParentID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &replyR{}
			}
			args[i] = obj.ParentID
		}
	}

	query := fmt.Sprintf(
		"select * from \"replies\" where \"snowflake\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Reply")
	}
	defer results.Close()

	var resultSlice []*Reply
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Reply")
	}

	if len(replyAfterSelectHooks) != 0 {
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

// LoadTopic allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (replyL) LoadTopic(e boil.Executor, singular bool, maybeReply interface{}) error {
	var slice []*Reply
	var object *Reply

	count := 1
	if singular {
		object = maybeReply.(*Reply)
	} else {
		slice = *maybeReply.(*ReplySlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &replyR{}
		}
		args[0] = object.TopicID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &replyR{}
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

	if len(replyAfterSelectHooks) != 0 {
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

// LoadParentReplies allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (replyL) LoadParentReplies(e boil.Executor, singular bool, maybeReply interface{}) error {
	var slice []*Reply
	var object *Reply

	count := 1
	if singular {
		object = maybeReply.(*Reply)
	} else {
		slice = *maybeReply.(*ReplySlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &replyR{}
		}
		args[0] = object.Snowflake
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &replyR{}
			}
			args[i] = obj.Snowflake
		}
	}

	query := fmt.Sprintf(
		"select * from \"replies\" where \"parent_id\" in (%s)",
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
		object.R.ParentReplies = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.Snowflake == foreign.ParentID.Int64 {
				local.R.ParentReplies = append(local.R.ParentReplies, foreign)
				break
			}
		}
	}

	return nil
}

// SetAuthorG of the reply to the related item.
// Sets o.R.Author to related.
// Adds o to related.R.AuthorReplies.
// Uses the global database handle.
func (o *Reply) SetAuthorG(insert bool, related *User) error {
	return o.SetAuthor(boil.GetDB(), insert, related)
}

// SetAuthorP of the reply to the related item.
// Sets o.R.Author to related.
// Adds o to related.R.AuthorReplies.
// Panics on error.
func (o *Reply) SetAuthorP(exec boil.Executor, insert bool, related *User) {
	if err := o.SetAuthor(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetAuthorGP of the reply to the related item.
// Sets o.R.Author to related.
// Adds o to related.R.AuthorReplies.
// Uses the global database handle and panics on error.
func (o *Reply) SetAuthorGP(insert bool, related *User) {
	if err := o.SetAuthor(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetAuthor of the reply to the related item.
// Sets o.R.Author to related.
// Adds o to related.R.AuthorReplies.
func (o *Reply) SetAuthor(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"replies\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"author_id"}),
		strmangle.WhereClause("\"", "\"", 2, replyPrimaryKeyColumns),
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
		o.R = &replyR{
			Author: related,
		}
	} else {
		o.R.Author = related
	}

	if related.R == nil {
		related.R = &userR{
			AuthorReplies: ReplySlice{o},
		}
	} else {
		related.R.AuthorReplies = append(related.R.AuthorReplies, o)
	}

	return nil
}

// RemoveAuthorG relationship.
// Sets o.R.Author to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle.
func (o *Reply) RemoveAuthorG(related *User) error {
	return o.RemoveAuthor(boil.GetDB(), related)
}

// RemoveAuthorP relationship.
// Sets o.R.Author to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Panics on error.
func (o *Reply) RemoveAuthorP(exec boil.Executor, related *User) {
	if err := o.RemoveAuthor(exec, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveAuthorGP relationship.
// Sets o.R.Author to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle and panics on error.
func (o *Reply) RemoveAuthorGP(related *User) {
	if err := o.RemoveAuthor(boil.GetDB(), related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveAuthor relationship.
// Sets o.R.Author to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *Reply) RemoveAuthor(exec boil.Executor, related *User) error {
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

	for i, ri := range related.R.AuthorReplies {
		if o.AuthorID.Int64 != ri.AuthorID.Int64 {
			continue
		}

		ln := len(related.R.AuthorReplies)
		if ln > 1 && i < ln-1 {
			related.R.AuthorReplies[i] = related.R.AuthorReplies[ln-1]
		}
		related.R.AuthorReplies = related.R.AuthorReplies[:ln-1]
		break
	}
	return nil
}

// SetParentG of the reply to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentReplies.
// Uses the global database handle.
func (o *Reply) SetParentG(insert bool, related *Reply) error {
	return o.SetParent(boil.GetDB(), insert, related)
}

// SetParentP of the reply to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentReplies.
// Panics on error.
func (o *Reply) SetParentP(exec boil.Executor, insert bool, related *Reply) {
	if err := o.SetParent(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParentGP of the reply to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentReplies.
// Uses the global database handle and panics on error.
func (o *Reply) SetParentGP(insert bool, related *Reply) {
	if err := o.SetParent(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParent of the reply to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentReplies.
func (o *Reply) SetParent(exec boil.Executor, insert bool, related *Reply) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"replies\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"parent_id"}),
		strmangle.WhereClause("\"", "\"", 2, replyPrimaryKeyColumns),
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
		o.R = &replyR{
			Parent: related,
		}
	} else {
		o.R.Parent = related
	}

	if related.R == nil {
		related.R = &replyR{
			ParentReplies: ReplySlice{o},
		}
	} else {
		related.R.ParentReplies = append(related.R.ParentReplies, o)
	}

	return nil
}

// RemoveParentG relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle.
func (o *Reply) RemoveParentG(related *Reply) error {
	return o.RemoveParent(boil.GetDB(), related)
}

// RemoveParentP relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Panics on error.
func (o *Reply) RemoveParentP(exec boil.Executor, related *Reply) {
	if err := o.RemoveParent(exec, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParentGP relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle and panics on error.
func (o *Reply) RemoveParentGP(related *Reply) {
	if err := o.RemoveParent(boil.GetDB(), related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParent relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *Reply) RemoveParent(exec boil.Executor, related *Reply) error {
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

	for i, ri := range related.R.ParentReplies {
		if o.ParentID.Int64 != ri.ParentID.Int64 {
			continue
		}

		ln := len(related.R.ParentReplies)
		if ln > 1 && i < ln-1 {
			related.R.ParentReplies[i] = related.R.ParentReplies[ln-1]
		}
		related.R.ParentReplies = related.R.ParentReplies[:ln-1]
		break
	}
	return nil
}

// SetTopicG of the reply to the related item.
// Sets o.R.Topic to related.
// Adds o to related.R.Replies.
// Uses the global database handle.
func (o *Reply) SetTopicG(insert bool, related *Topic) error {
	return o.SetTopic(boil.GetDB(), insert, related)
}

// SetTopicP of the reply to the related item.
// Sets o.R.Topic to related.
// Adds o to related.R.Replies.
// Panics on error.
func (o *Reply) SetTopicP(exec boil.Executor, insert bool, related *Topic) {
	if err := o.SetTopic(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetTopicGP of the reply to the related item.
// Sets o.R.Topic to related.
// Adds o to related.R.Replies.
// Uses the global database handle and panics on error.
func (o *Reply) SetTopicGP(insert bool, related *Topic) {
	if err := o.SetTopic(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetTopic of the reply to the related item.
// Sets o.R.Topic to related.
// Adds o to related.R.Replies.
func (o *Reply) SetTopic(exec boil.Executor, insert bool, related *Topic) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"replies\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"topic_id"}),
		strmangle.WhereClause("\"", "\"", 2, replyPrimaryKeyColumns),
	)
	values := []interface{}{related.Snowflake, o.Snowflake}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.TopicID = related.Snowflake

	if o.R == nil {
		o.R = &replyR{
			Topic: related,
		}
	} else {
		o.R.Topic = related
	}

	if related.R == nil {
		related.R = &topicR{
			Replies: ReplySlice{o},
		}
	} else {
		related.R.Replies = append(related.R.Replies, o)
	}

	return nil
}

// AddParentRepliesG adds the given related objects to the existing relationships
// of the reply, optionally inserting them as new records.
// Appends related to o.R.ParentReplies.
// Sets related.R.Parent appropriately.
// Uses the global database handle.
func (o *Reply) AddParentRepliesG(insert bool, related ...*Reply) error {
	return o.AddParentReplies(boil.GetDB(), insert, related...)
}

// AddParentRepliesP adds the given related objects to the existing relationships
// of the reply, optionally inserting them as new records.
// Appends related to o.R.ParentReplies.
// Sets related.R.Parent appropriately.
// Panics on error.
func (o *Reply) AddParentRepliesP(exec boil.Executor, insert bool, related ...*Reply) {
	if err := o.AddParentReplies(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddParentRepliesGP adds the given related objects to the existing relationships
// of the reply, optionally inserting them as new records.
// Appends related to o.R.ParentReplies.
// Sets related.R.Parent appropriately.
// Uses the global database handle and panics on error.
func (o *Reply) AddParentRepliesGP(insert bool, related ...*Reply) {
	if err := o.AddParentReplies(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddParentReplies adds the given related objects to the existing relationships
// of the reply, optionally inserting them as new records.
// Appends related to o.R.ParentReplies.
// Sets related.R.Parent appropriately.
func (o *Reply) AddParentReplies(exec boil.Executor, insert bool, related ...*Reply) error {
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
				"UPDATE \"replies\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"parent_id"}),
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

			rel.ParentID.Int64 = o.Snowflake
			rel.ParentID.Valid = true
		}
	}

	if o.R == nil {
		o.R = &replyR{
			ParentReplies: related,
		}
	} else {
		o.R.ParentReplies = append(o.R.ParentReplies, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &replyR{
				Parent: o,
			}
		} else {
			rel.R.Parent = o
		}
	}
	return nil
}

// SetParentRepliesG removes all previously related items of the
// reply replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentReplies accordingly.
// Replaces o.R.ParentReplies with related.
// Sets related.R.Parent's ParentReplies accordingly.
// Uses the global database handle.
func (o *Reply) SetParentRepliesG(insert bool, related ...*Reply) error {
	return o.SetParentReplies(boil.GetDB(), insert, related...)
}

// SetParentRepliesP removes all previously related items of the
// reply replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentReplies accordingly.
// Replaces o.R.ParentReplies with related.
// Sets related.R.Parent's ParentReplies accordingly.
// Panics on error.
func (o *Reply) SetParentRepliesP(exec boil.Executor, insert bool, related ...*Reply) {
	if err := o.SetParentReplies(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParentRepliesGP removes all previously related items of the
// reply replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentReplies accordingly.
// Replaces o.R.ParentReplies with related.
// Sets related.R.Parent's ParentReplies accordingly.
// Uses the global database handle and panics on error.
func (o *Reply) SetParentRepliesGP(insert bool, related ...*Reply) {
	if err := o.SetParentReplies(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParentReplies removes all previously related items of the
// reply replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentReplies accordingly.
// Replaces o.R.ParentReplies with related.
// Sets related.R.Parent's ParentReplies accordingly.
func (o *Reply) SetParentReplies(exec boil.Executor, insert bool, related ...*Reply) error {
	query := "update \"replies\" set \"parent_id\" = null where \"parent_id\" = $1"
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
		for _, rel := range o.R.ParentReplies {
			rel.ParentID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Parent = nil
		}

		o.R.ParentReplies = nil
	}
	return o.AddParentReplies(exec, insert, related...)
}

// RemoveParentRepliesG relationships from objects passed in.
// Removes related items from R.ParentReplies (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
// Uses the global database handle.
func (o *Reply) RemoveParentRepliesG(related ...*Reply) error {
	return o.RemoveParentReplies(boil.GetDB(), related...)
}

// RemoveParentRepliesP relationships from objects passed in.
// Removes related items from R.ParentReplies (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
// Panics on error.
func (o *Reply) RemoveParentRepliesP(exec boil.Executor, related ...*Reply) {
	if err := o.RemoveParentReplies(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParentRepliesGP relationships from objects passed in.
// Removes related items from R.ParentReplies (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
// Uses the global database handle and panics on error.
func (o *Reply) RemoveParentRepliesGP(related ...*Reply) {
	if err := o.RemoveParentReplies(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParentReplies relationships from objects passed in.
// Removes related items from R.ParentReplies (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
func (o *Reply) RemoveParentReplies(exec boil.Executor, related ...*Reply) error {
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
		for i, ri := range o.R.ParentReplies {
			if rel != ri {
				continue
			}

			ln := len(o.R.ParentReplies)
			if ln > 1 && i < ln-1 {
				o.R.ParentReplies[i] = o.R.ParentReplies[ln-1]
			}
			o.R.ParentReplies = o.R.ParentReplies[:ln-1]
			break
		}
	}

	return nil
}

// RepliesG retrieves all records.
func RepliesG(mods ...qm.QueryMod) replyQuery {
	return Replies(boil.GetDB(), mods...)
}

// Replies retrieves all the records using an executor.
func Replies(exec boil.Executor, mods ...qm.QueryMod) replyQuery {
	mods = append(mods, qm.From("\"replies\""))
	return replyQuery{NewQuery(exec, mods...)}
}

// FindReplyG retrieves a single record by ID.
func FindReplyG(snowflake int64, selectCols ...string) (*Reply, error) {
	return FindReply(boil.GetDB(), snowflake, selectCols...)
}

// FindReplyGP retrieves a single record by ID, and panics on error.
func FindReplyGP(snowflake int64, selectCols ...string) *Reply {
	retobj, err := FindReply(boil.GetDB(), snowflake, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindReply retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindReply(exec boil.Executor, snowflake int64, selectCols ...string) (*Reply, error) {
	replyObj := &Reply{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"replies\" where \"snowflake\"=$1", sel,
	)

	q := queries.Raw(exec, query, snowflake)

	err := q.Bind(replyObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from replies")
	}

	return replyObj, nil
}

// FindReplyP retrieves a single record by ID with an executor, and panics on error.
func FindReplyP(exec boil.Executor, snowflake int64, selectCols ...string) *Reply {
	retobj, err := FindReply(exec, snowflake, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Reply) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Reply) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Reply) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Reply) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no replies provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(replyColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	replyInsertCacheMut.RLock()
	cache, cached := replyInsertCache[key]
	replyInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			replyColumns,
			replyColumnsWithDefault,
			replyColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(replyType, replyMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(replyType, replyMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"replies\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into replies")
	}

	if !cached {
		replyInsertCacheMut.Lock()
		replyInsertCache[key] = cache
		replyInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Reply record. See Update for
// whitelist behavior description.
func (o *Reply) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Reply record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Reply) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Reply, and panics on error.
// See Update for whitelist behavior description.
func (o *Reply) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Reply.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Reply) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	replyUpdateCacheMut.RLock()
	cache, cached := replyUpdateCache[key]
	replyUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(replyColumns, replyPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update replies, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"replies\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, replyPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(replyType, replyMapping, append(wl, replyPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update replies row")
	}

	if !cached {
		replyUpdateCacheMut.Lock()
		replyUpdateCache[key] = cache
		replyUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q replyQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q replyQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for replies")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o ReplySlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o ReplySlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o ReplySlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ReplySlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), replyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"replies\" SET %s WHERE (\"snowflake\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(replyPrimaryKeyColumns), len(colNames)+1, len(replyPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in reply slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Reply) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Reply) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Reply) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Reply) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no replies provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(replyColumnsWithDefault, o)

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

	replyUpsertCacheMut.RLock()
	cache, cached := replyUpsertCache[key]
	replyUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			replyColumns,
			replyColumnsWithDefault,
			replyColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			replyColumns,
			replyPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert replies, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(replyPrimaryKeyColumns))
			copy(conflict, replyPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"replies\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(replyType, replyMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(replyType, replyMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert replies")
	}

	if !cached {
		replyUpsertCacheMut.Lock()
		replyUpsertCache[key] = cache
		replyUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Reply record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Reply) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Reply record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Reply) DeleteG() error {
	if o == nil {
		return errors.New("models: no Reply provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Reply record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Reply) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Reply record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Reply) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Reply provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), replyPrimaryKeyMapping)
	sql := "DELETE FROM \"replies\" WHERE \"snowflake\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from replies")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q replyQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q replyQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no replyQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from replies")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o ReplySlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o ReplySlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Reply slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o ReplySlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ReplySlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Reply slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(replyBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), replyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"replies\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, replyPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(replyPrimaryKeyColumns), 1, len(replyPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from reply slice")
	}

	if len(replyAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Reply) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Reply) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Reply) ReloadG() error {
	if o == nil {
		return errors.New("models: no Reply provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Reply) Reload(exec boil.Executor) error {
	ret, err := FindReply(exec, o.Snowflake)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ReplySlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ReplySlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ReplySlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty ReplySlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ReplySlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	replies := ReplySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), replyPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"replies\".* FROM \"replies\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, replyPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(replyPrimaryKeyColumns), 1, len(replyPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&replies)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ReplySlice")
	}

	*o = replies

	return nil
}

// ReplyExists checks if the Reply row exists.
func ReplyExists(exec boil.Executor, snowflake int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"replies\" where \"snowflake\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, snowflake)
	}

	row := exec.QueryRow(sql, snowflake)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if replies exists")
	}

	return exists, nil
}

// ReplyExistsG checks if the Reply row exists.
func ReplyExistsG(snowflake int64) (bool, error) {
	return ReplyExists(boil.GetDB(), snowflake)
}

// ReplyExistsGP checks if the Reply row exists. Panics on error.
func ReplyExistsGP(snowflake int64) bool {
	e, err := ReplyExists(boil.GetDB(), snowflake)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// ReplyExistsP checks if the Reply row exists. Panics on error.
func ReplyExistsP(exec boil.Executor, snowflake int64) bool {
	e, err := ReplyExists(exec, snowflake)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
