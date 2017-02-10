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

// User is an object representing the database table.
type User struct {
	Snowflake int64       `boil:"snowflake" json:"snowflake" toml:"snowflake" yaml:"snowflake"`
	CreatedAt time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	DeletedAt null.Time   `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`
	Username  string      `boil:"username" json:"username" toml:"username" yaml:"username"`
	Email     null.String `boil:"email" json:"email,omitempty" toml:"email" yaml:"email,omitempty"`
	Avatar    []byte      `boil:"avatar" json:"avatar" toml:"avatar" yaml:"avatar"`

	R *userR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L userL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// userR is where relationships are stored.
type userR struct {
	RelUserGroups           RelUserGroupSlice
	Logins                  LoginSlice
	AuthorTopics            TopicSlice
	AuthorReplies           ReplySlice
	SenderPrivateMessages   PrivateMessageSlice
	ReceiverPrivateMessages PrivateMessageSlice
}

// userL is where Load methods for each relationship are stored.
type userL struct{}

var (
	userColumns               = []string{"snowflake", "created_at", "deleted_at", "username", "email", "avatar"}
	userColumnsWithoutDefault = []string{"snowflake", "deleted_at", "username", "email", "avatar"}
	userColumnsWithDefault    = []string{"created_at"}
	userPrimaryKeyColumns     = []string{"snowflake"}
)

type (
	// UserSlice is an alias for a slice of pointers to User.
	// This should generally be used opposed to []User.
	UserSlice []*User
	// UserHook is the signature for custom User hook methods
	UserHook func(boil.Executor, *User) error

	userQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	userType                 = reflect.TypeOf(&User{})
	userMapping              = queries.MakeStructMapping(userType)
	userPrimaryKeyMapping, _ = queries.BindMapping(userType, userMapping, userPrimaryKeyColumns)
	userInsertCacheMut       sync.RWMutex
	userInsertCache          = make(map[string]insertCache)
	userUpdateCacheMut       sync.RWMutex
	userUpdateCache          = make(map[string]updateCache)
	userUpsertCacheMut       sync.RWMutex
	userUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var userBeforeInsertHooks []UserHook
var userBeforeUpdateHooks []UserHook
var userBeforeDeleteHooks []UserHook
var userBeforeUpsertHooks []UserHook

var userAfterInsertHooks []UserHook
var userAfterSelectHooks []UserHook
var userAfterUpdateHooks []UserHook
var userAfterDeleteHooks []UserHook
var userAfterUpsertHooks []UserHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *User) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *User) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range userBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *User) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range userBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *User) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *User) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *User) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range userAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *User) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range userAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *User) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range userAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *User) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range userAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddUserHook registers your hook function for all future operations.
func AddUserHook(hookPoint boil.HookPoint, userHook UserHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		userBeforeInsertHooks = append(userBeforeInsertHooks, userHook)
	case boil.BeforeUpdateHook:
		userBeforeUpdateHooks = append(userBeforeUpdateHooks, userHook)
	case boil.BeforeDeleteHook:
		userBeforeDeleteHooks = append(userBeforeDeleteHooks, userHook)
	case boil.BeforeUpsertHook:
		userBeforeUpsertHooks = append(userBeforeUpsertHooks, userHook)
	case boil.AfterInsertHook:
		userAfterInsertHooks = append(userAfterInsertHooks, userHook)
	case boil.AfterSelectHook:
		userAfterSelectHooks = append(userAfterSelectHooks, userHook)
	case boil.AfterUpdateHook:
		userAfterUpdateHooks = append(userAfterUpdateHooks, userHook)
	case boil.AfterDeleteHook:
		userAfterDeleteHooks = append(userAfterDeleteHooks, userHook)
	case boil.AfterUpsertHook:
		userAfterUpsertHooks = append(userAfterUpsertHooks, userHook)
	}
}

// OneP returns a single user record from the query, and panics on error.
func (q userQuery) OneP() *User {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single user record from the query.
func (q userQuery) One() (*User, error) {
	o := &User{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for users")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all User records from the query, and panics on error.
func (q userQuery) AllP() UserSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all User records from the query.
func (q userQuery) All() (UserSlice, error) {
	var o UserSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to User slice")
	}

	if len(userAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all User records in the query, and panics on error.
func (q userQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all User records in the query.
func (q userQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count users rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q userQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q userQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if users exists")
	}

	return count > 0, nil
}

// RelUserGroupsG retrieves all the rel_user_group's rel user groups.
func (o *User) RelUserGroupsG(mods ...qm.QueryMod) relUserGroupQuery {
	return o.RelUserGroups(boil.GetDB(), mods...)
}

// RelUserGroups retrieves all the rel_user_group's rel user groups with an executor.
func (o *User) RelUserGroups(exec boil.Executor, mods ...qm.QueryMod) relUserGroupQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"user_id\"=?", o.Snowflake),
	)

	query := RelUserGroups(exec, queryMods...)
	queries.SetFrom(query.Query, "\"rel_user_groups\" as \"a\"")
	return query
}

// LoginsG retrieves all the login's logins.
func (o *User) LoginsG(mods ...qm.QueryMod) loginQuery {
	return o.Logins(boil.GetDB(), mods...)
}

// Logins retrieves all the login's logins with an executor.
func (o *User) Logins(exec boil.Executor, mods ...qm.QueryMod) loginQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"user_id\"=?", o.Snowflake),
	)

	query := Logins(exec, queryMods...)
	queries.SetFrom(query.Query, "\"logins\" as \"a\"")
	return query
}

// AuthorTopicsG retrieves all the topic's topics via author_id column.
func (o *User) AuthorTopicsG(mods ...qm.QueryMod) topicQuery {
	return o.AuthorTopics(boil.GetDB(), mods...)
}

// AuthorTopics retrieves all the topic's topics with an executor via author_id column.
func (o *User) AuthorTopics(exec boil.Executor, mods ...qm.QueryMod) topicQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"author_id\"=?", o.Snowflake),
	)

	query := Topics(exec, queryMods...)
	queries.SetFrom(query.Query, "\"topics\" as \"a\"")
	return query
}

// AuthorRepliesG retrieves all the reply's replies via author_id column.
func (o *User) AuthorRepliesG(mods ...qm.QueryMod) replyQuery {
	return o.AuthorReplies(boil.GetDB(), mods...)
}

// AuthorReplies retrieves all the reply's replies with an executor via author_id column.
func (o *User) AuthorReplies(exec boil.Executor, mods ...qm.QueryMod) replyQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"author_id\"=?", o.Snowflake),
	)

	query := Replies(exec, queryMods...)
	queries.SetFrom(query.Query, "\"replies\" as \"a\"")
	return query
}

// SenderPrivateMessagesG retrieves all the private_message's private messages via sender_id column.
func (o *User) SenderPrivateMessagesG(mods ...qm.QueryMod) privateMessageQuery {
	return o.SenderPrivateMessages(boil.GetDB(), mods...)
}

// SenderPrivateMessages retrieves all the private_message's private messages with an executor via sender_id column.
func (o *User) SenderPrivateMessages(exec boil.Executor, mods ...qm.QueryMod) privateMessageQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"sender_id\"=?", o.Snowflake),
	)

	query := PrivateMessages(exec, queryMods...)
	queries.SetFrom(query.Query, "\"private_messages\" as \"a\"")
	return query
}

// ReceiverPrivateMessagesG retrieves all the private_message's private messages via receiver_id column.
func (o *User) ReceiverPrivateMessagesG(mods ...qm.QueryMod) privateMessageQuery {
	return o.ReceiverPrivateMessages(boil.GetDB(), mods...)
}

// ReceiverPrivateMessages retrieves all the private_message's private messages with an executor via receiver_id column.
func (o *User) ReceiverPrivateMessages(exec boil.Executor, mods ...qm.QueryMod) privateMessageQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"receiver_id\"=?", o.Snowflake),
	)

	query := PrivateMessages(exec, queryMods...)
	queries.SetFrom(query.Query, "\"private_messages\" as \"a\"")
	return query
}

// LoadRelUserGroups allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (userL) LoadRelUserGroups(e boil.Executor, singular bool, maybeUser interface{}) error {
	var slice []*User
	var object *User

	count := 1
	if singular {
		object = maybeUser.(*User)
	} else {
		slice = *maybeUser.(*UserSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &userR{}
		}
		args[0] = object.Snowflake
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &userR{}
			}
			args[i] = obj.Snowflake
		}
	}

	query := fmt.Sprintf(
		"select * from \"rel_user_groups\" where \"user_id\" in (%s)",
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
			if local.Snowflake == foreign.UserID {
				local.R.RelUserGroups = append(local.R.RelUserGroups, foreign)
				break
			}
		}
	}

	return nil
}

// LoadLogins allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (userL) LoadLogins(e boil.Executor, singular bool, maybeUser interface{}) error {
	var slice []*User
	var object *User

	count := 1
	if singular {
		object = maybeUser.(*User)
	} else {
		slice = *maybeUser.(*UserSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &userR{}
		}
		args[0] = object.Snowflake
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &userR{}
			}
			args[i] = obj.Snowflake
		}
	}

	query := fmt.Sprintf(
		"select * from \"logins\" where \"user_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load logins")
	}
	defer results.Close()

	var resultSlice []*Login
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice logins")
	}

	if len(loginAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Logins = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.Snowflake == foreign.UserID {
				local.R.Logins = append(local.R.Logins, foreign)
				break
			}
		}
	}

	return nil
}

// LoadAuthorTopics allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (userL) LoadAuthorTopics(e boil.Executor, singular bool, maybeUser interface{}) error {
	var slice []*User
	var object *User

	count := 1
	if singular {
		object = maybeUser.(*User)
	} else {
		slice = *maybeUser.(*UserSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &userR{}
		}
		args[0] = object.Snowflake
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &userR{}
			}
			args[i] = obj.Snowflake
		}
	}

	query := fmt.Sprintf(
		"select * from \"topics\" where \"author_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load topics")
	}
	defer results.Close()

	var resultSlice []*Topic
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice topics")
	}

	if len(topicAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.AuthorTopics = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.Snowflake == foreign.AuthorID.Int64 {
				local.R.AuthorTopics = append(local.R.AuthorTopics, foreign)
				break
			}
		}
	}

	return nil
}

// LoadAuthorReplies allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (userL) LoadAuthorReplies(e boil.Executor, singular bool, maybeUser interface{}) error {
	var slice []*User
	var object *User

	count := 1
	if singular {
		object = maybeUser.(*User)
	} else {
		slice = *maybeUser.(*UserSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &userR{}
		}
		args[0] = object.Snowflake
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &userR{}
			}
			args[i] = obj.Snowflake
		}
	}

	query := fmt.Sprintf(
		"select * from \"replies\" where \"author_id\" in (%s)",
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
		object.R.AuthorReplies = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.Snowflake == foreign.AuthorID.Int64 {
				local.R.AuthorReplies = append(local.R.AuthorReplies, foreign)
				break
			}
		}
	}

	return nil
}

// LoadSenderPrivateMessages allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (userL) LoadSenderPrivateMessages(e boil.Executor, singular bool, maybeUser interface{}) error {
	var slice []*User
	var object *User

	count := 1
	if singular {
		object = maybeUser.(*User)
	} else {
		slice = *maybeUser.(*UserSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &userR{}
		}
		args[0] = object.Snowflake
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &userR{}
			}
			args[i] = obj.Snowflake
		}
	}

	query := fmt.Sprintf(
		"select * from \"private_messages\" where \"sender_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load private_messages")
	}
	defer results.Close()

	var resultSlice []*PrivateMessage
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice private_messages")
	}

	if len(privateMessageAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.SenderPrivateMessages = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.Snowflake == foreign.SenderID {
				local.R.SenderPrivateMessages = append(local.R.SenderPrivateMessages, foreign)
				break
			}
		}
	}

	return nil
}

// LoadReceiverPrivateMessages allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (userL) LoadReceiverPrivateMessages(e boil.Executor, singular bool, maybeUser interface{}) error {
	var slice []*User
	var object *User

	count := 1
	if singular {
		object = maybeUser.(*User)
	} else {
		slice = *maybeUser.(*UserSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &userR{}
		}
		args[0] = object.Snowflake
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &userR{}
			}
			args[i] = obj.Snowflake
		}
	}

	query := fmt.Sprintf(
		"select * from \"private_messages\" where \"receiver_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load private_messages")
	}
	defer results.Close()

	var resultSlice []*PrivateMessage
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice private_messages")
	}

	if len(privateMessageAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.ReceiverPrivateMessages = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.Snowflake == foreign.ReceiverID {
				local.R.ReceiverPrivateMessages = append(local.R.ReceiverPrivateMessages, foreign)
				break
			}
		}
	}

	return nil
}

// AddRelUserGroupsG adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.RelUserGroups.
// Sets related.R.User appropriately.
// Uses the global database handle.
func (o *User) AddRelUserGroupsG(insert bool, related ...*RelUserGroup) error {
	return o.AddRelUserGroups(boil.GetDB(), insert, related...)
}

// AddRelUserGroupsP adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.RelUserGroups.
// Sets related.R.User appropriately.
// Panics on error.
func (o *User) AddRelUserGroupsP(exec boil.Executor, insert bool, related ...*RelUserGroup) {
	if err := o.AddRelUserGroups(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddRelUserGroupsGP adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.RelUserGroups.
// Sets related.R.User appropriately.
// Uses the global database handle and panics on error.
func (o *User) AddRelUserGroupsGP(insert bool, related ...*RelUserGroup) {
	if err := o.AddRelUserGroups(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddRelUserGroups adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.RelUserGroups.
// Sets related.R.User appropriately.
func (o *User) AddRelUserGroups(exec boil.Executor, insert bool, related ...*RelUserGroup) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.UserID = o.Snowflake
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"rel_user_groups\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
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

			rel.UserID = o.Snowflake
		}
	}

	if o.R == nil {
		o.R = &userR{
			RelUserGroups: related,
		}
	} else {
		o.R.RelUserGroups = append(o.R.RelUserGroups, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &relUserGroupR{
				User: o,
			}
		} else {
			rel.R.User = o
		}
	}
	return nil
}

// AddLoginsG adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.Logins.
// Sets related.R.User appropriately.
// Uses the global database handle.
func (o *User) AddLoginsG(insert bool, related ...*Login) error {
	return o.AddLogins(boil.GetDB(), insert, related...)
}

// AddLoginsP adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.Logins.
// Sets related.R.User appropriately.
// Panics on error.
func (o *User) AddLoginsP(exec boil.Executor, insert bool, related ...*Login) {
	if err := o.AddLogins(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddLoginsGP adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.Logins.
// Sets related.R.User appropriately.
// Uses the global database handle and panics on error.
func (o *User) AddLoginsGP(insert bool, related ...*Login) {
	if err := o.AddLogins(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddLogins adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.Logins.
// Sets related.R.User appropriately.
func (o *User) AddLogins(exec boil.Executor, insert bool, related ...*Login) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.UserID = o.Snowflake
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"logins\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
				strmangle.WhereClause("\"", "\"", 2, loginPrimaryKeyColumns),
			)
			values := []interface{}{o.Snowflake, rel.Snowflake}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.UserID = o.Snowflake
		}
	}

	if o.R == nil {
		o.R = &userR{
			Logins: related,
		}
	} else {
		o.R.Logins = append(o.R.Logins, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &loginR{
				User: o,
			}
		} else {
			rel.R.User = o
		}
	}
	return nil
}

// AddAuthorTopicsG adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.AuthorTopics.
// Sets related.R.Author appropriately.
// Uses the global database handle.
func (o *User) AddAuthorTopicsG(insert bool, related ...*Topic) error {
	return o.AddAuthorTopics(boil.GetDB(), insert, related...)
}

// AddAuthorTopicsP adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.AuthorTopics.
// Sets related.R.Author appropriately.
// Panics on error.
func (o *User) AddAuthorTopicsP(exec boil.Executor, insert bool, related ...*Topic) {
	if err := o.AddAuthorTopics(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddAuthorTopicsGP adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.AuthorTopics.
// Sets related.R.Author appropriately.
// Uses the global database handle and panics on error.
func (o *User) AddAuthorTopicsGP(insert bool, related ...*Topic) {
	if err := o.AddAuthorTopics(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddAuthorTopics adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.AuthorTopics.
// Sets related.R.Author appropriately.
func (o *User) AddAuthorTopics(exec boil.Executor, insert bool, related ...*Topic) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.AuthorID.Int64 = o.Snowflake
			rel.AuthorID.Valid = true
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"topics\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"author_id"}),
				strmangle.WhereClause("\"", "\"", 2, topicPrimaryKeyColumns),
			)
			values := []interface{}{o.Snowflake, rel.Snowflake}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.AuthorID.Int64 = o.Snowflake
			rel.AuthorID.Valid = true
		}
	}

	if o.R == nil {
		o.R = &userR{
			AuthorTopics: related,
		}
	} else {
		o.R.AuthorTopics = append(o.R.AuthorTopics, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &topicR{
				Author: o,
			}
		} else {
			rel.R.Author = o
		}
	}
	return nil
}

// SetAuthorTopicsG removes all previously related items of the
// user replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Author's AuthorTopics accordingly.
// Replaces o.R.AuthorTopics with related.
// Sets related.R.Author's AuthorTopics accordingly.
// Uses the global database handle.
func (o *User) SetAuthorTopicsG(insert bool, related ...*Topic) error {
	return o.SetAuthorTopics(boil.GetDB(), insert, related...)
}

// SetAuthorTopicsP removes all previously related items of the
// user replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Author's AuthorTopics accordingly.
// Replaces o.R.AuthorTopics with related.
// Sets related.R.Author's AuthorTopics accordingly.
// Panics on error.
func (o *User) SetAuthorTopicsP(exec boil.Executor, insert bool, related ...*Topic) {
	if err := o.SetAuthorTopics(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetAuthorTopicsGP removes all previously related items of the
// user replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Author's AuthorTopics accordingly.
// Replaces o.R.AuthorTopics with related.
// Sets related.R.Author's AuthorTopics accordingly.
// Uses the global database handle and panics on error.
func (o *User) SetAuthorTopicsGP(insert bool, related ...*Topic) {
	if err := o.SetAuthorTopics(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetAuthorTopics removes all previously related items of the
// user replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Author's AuthorTopics accordingly.
// Replaces o.R.AuthorTopics with related.
// Sets related.R.Author's AuthorTopics accordingly.
func (o *User) SetAuthorTopics(exec boil.Executor, insert bool, related ...*Topic) error {
	query := "update \"topics\" set \"author_id\" = null where \"author_id\" = $1"
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
		for _, rel := range o.R.AuthorTopics {
			rel.AuthorID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Author = nil
		}

		o.R.AuthorTopics = nil
	}
	return o.AddAuthorTopics(exec, insert, related...)
}

// RemoveAuthorTopicsG relationships from objects passed in.
// Removes related items from R.AuthorTopics (uses pointer comparison, removal does not keep order)
// Sets related.R.Author.
// Uses the global database handle.
func (o *User) RemoveAuthorTopicsG(related ...*Topic) error {
	return o.RemoveAuthorTopics(boil.GetDB(), related...)
}

// RemoveAuthorTopicsP relationships from objects passed in.
// Removes related items from R.AuthorTopics (uses pointer comparison, removal does not keep order)
// Sets related.R.Author.
// Panics on error.
func (o *User) RemoveAuthorTopicsP(exec boil.Executor, related ...*Topic) {
	if err := o.RemoveAuthorTopics(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveAuthorTopicsGP relationships from objects passed in.
// Removes related items from R.AuthorTopics (uses pointer comparison, removal does not keep order)
// Sets related.R.Author.
// Uses the global database handle and panics on error.
func (o *User) RemoveAuthorTopicsGP(related ...*Topic) {
	if err := o.RemoveAuthorTopics(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveAuthorTopics relationships from objects passed in.
// Removes related items from R.AuthorTopics (uses pointer comparison, removal does not keep order)
// Sets related.R.Author.
func (o *User) RemoveAuthorTopics(exec boil.Executor, related ...*Topic) error {
	var err error
	for _, rel := range related {
		rel.AuthorID.Valid = false
		if rel.R != nil {
			rel.R.Author = nil
		}
		if err = rel.Update(exec, "author_id"); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.AuthorTopics {
			if rel != ri {
				continue
			}

			ln := len(o.R.AuthorTopics)
			if ln > 1 && i < ln-1 {
				o.R.AuthorTopics[i] = o.R.AuthorTopics[ln-1]
			}
			o.R.AuthorTopics = o.R.AuthorTopics[:ln-1]
			break
		}
	}

	return nil
}

// AddAuthorRepliesG adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.AuthorReplies.
// Sets related.R.Author appropriately.
// Uses the global database handle.
func (o *User) AddAuthorRepliesG(insert bool, related ...*Reply) error {
	return o.AddAuthorReplies(boil.GetDB(), insert, related...)
}

// AddAuthorRepliesP adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.AuthorReplies.
// Sets related.R.Author appropriately.
// Panics on error.
func (o *User) AddAuthorRepliesP(exec boil.Executor, insert bool, related ...*Reply) {
	if err := o.AddAuthorReplies(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddAuthorRepliesGP adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.AuthorReplies.
// Sets related.R.Author appropriately.
// Uses the global database handle and panics on error.
func (o *User) AddAuthorRepliesGP(insert bool, related ...*Reply) {
	if err := o.AddAuthorReplies(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddAuthorReplies adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.AuthorReplies.
// Sets related.R.Author appropriately.
func (o *User) AddAuthorReplies(exec boil.Executor, insert bool, related ...*Reply) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.AuthorID.Int64 = o.Snowflake
			rel.AuthorID.Valid = true
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"replies\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"author_id"}),
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

			rel.AuthorID.Int64 = o.Snowflake
			rel.AuthorID.Valid = true
		}
	}

	if o.R == nil {
		o.R = &userR{
			AuthorReplies: related,
		}
	} else {
		o.R.AuthorReplies = append(o.R.AuthorReplies, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &replyR{
				Author: o,
			}
		} else {
			rel.R.Author = o
		}
	}
	return nil
}

// SetAuthorRepliesG removes all previously related items of the
// user replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Author's AuthorReplies accordingly.
// Replaces o.R.AuthorReplies with related.
// Sets related.R.Author's AuthorReplies accordingly.
// Uses the global database handle.
func (o *User) SetAuthorRepliesG(insert bool, related ...*Reply) error {
	return o.SetAuthorReplies(boil.GetDB(), insert, related...)
}

// SetAuthorRepliesP removes all previously related items of the
// user replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Author's AuthorReplies accordingly.
// Replaces o.R.AuthorReplies with related.
// Sets related.R.Author's AuthorReplies accordingly.
// Panics on error.
func (o *User) SetAuthorRepliesP(exec boil.Executor, insert bool, related ...*Reply) {
	if err := o.SetAuthorReplies(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetAuthorRepliesGP removes all previously related items of the
// user replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Author's AuthorReplies accordingly.
// Replaces o.R.AuthorReplies with related.
// Sets related.R.Author's AuthorReplies accordingly.
// Uses the global database handle and panics on error.
func (o *User) SetAuthorRepliesGP(insert bool, related ...*Reply) {
	if err := o.SetAuthorReplies(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetAuthorReplies removes all previously related items of the
// user replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Author's AuthorReplies accordingly.
// Replaces o.R.AuthorReplies with related.
// Sets related.R.Author's AuthorReplies accordingly.
func (o *User) SetAuthorReplies(exec boil.Executor, insert bool, related ...*Reply) error {
	query := "update \"replies\" set \"author_id\" = null where \"author_id\" = $1"
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
		for _, rel := range o.R.AuthorReplies {
			rel.AuthorID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Author = nil
		}

		o.R.AuthorReplies = nil
	}
	return o.AddAuthorReplies(exec, insert, related...)
}

// RemoveAuthorRepliesG relationships from objects passed in.
// Removes related items from R.AuthorReplies (uses pointer comparison, removal does not keep order)
// Sets related.R.Author.
// Uses the global database handle.
func (o *User) RemoveAuthorRepliesG(related ...*Reply) error {
	return o.RemoveAuthorReplies(boil.GetDB(), related...)
}

// RemoveAuthorRepliesP relationships from objects passed in.
// Removes related items from R.AuthorReplies (uses pointer comparison, removal does not keep order)
// Sets related.R.Author.
// Panics on error.
func (o *User) RemoveAuthorRepliesP(exec boil.Executor, related ...*Reply) {
	if err := o.RemoveAuthorReplies(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveAuthorRepliesGP relationships from objects passed in.
// Removes related items from R.AuthorReplies (uses pointer comparison, removal does not keep order)
// Sets related.R.Author.
// Uses the global database handle and panics on error.
func (o *User) RemoveAuthorRepliesGP(related ...*Reply) {
	if err := o.RemoveAuthorReplies(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveAuthorReplies relationships from objects passed in.
// Removes related items from R.AuthorReplies (uses pointer comparison, removal does not keep order)
// Sets related.R.Author.
func (o *User) RemoveAuthorReplies(exec boil.Executor, related ...*Reply) error {
	var err error
	for _, rel := range related {
		rel.AuthorID.Valid = false
		if rel.R != nil {
			rel.R.Author = nil
		}
		if err = rel.Update(exec, "author_id"); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.AuthorReplies {
			if rel != ri {
				continue
			}

			ln := len(o.R.AuthorReplies)
			if ln > 1 && i < ln-1 {
				o.R.AuthorReplies[i] = o.R.AuthorReplies[ln-1]
			}
			o.R.AuthorReplies = o.R.AuthorReplies[:ln-1]
			break
		}
	}

	return nil
}

// AddSenderPrivateMessagesG adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.SenderPrivateMessages.
// Sets related.R.Sender appropriately.
// Uses the global database handle.
func (o *User) AddSenderPrivateMessagesG(insert bool, related ...*PrivateMessage) error {
	return o.AddSenderPrivateMessages(boil.GetDB(), insert, related...)
}

// AddSenderPrivateMessagesP adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.SenderPrivateMessages.
// Sets related.R.Sender appropriately.
// Panics on error.
func (o *User) AddSenderPrivateMessagesP(exec boil.Executor, insert bool, related ...*PrivateMessage) {
	if err := o.AddSenderPrivateMessages(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddSenderPrivateMessagesGP adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.SenderPrivateMessages.
// Sets related.R.Sender appropriately.
// Uses the global database handle and panics on error.
func (o *User) AddSenderPrivateMessagesGP(insert bool, related ...*PrivateMessage) {
	if err := o.AddSenderPrivateMessages(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddSenderPrivateMessages adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.SenderPrivateMessages.
// Sets related.R.Sender appropriately.
func (o *User) AddSenderPrivateMessages(exec boil.Executor, insert bool, related ...*PrivateMessage) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.SenderID = o.Snowflake
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"private_messages\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"sender_id"}),
				strmangle.WhereClause("\"", "\"", 2, privateMessagePrimaryKeyColumns),
			)
			values := []interface{}{o.Snowflake, rel.Snowflake}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.SenderID = o.Snowflake
		}
	}

	if o.R == nil {
		o.R = &userR{
			SenderPrivateMessages: related,
		}
	} else {
		o.R.SenderPrivateMessages = append(o.R.SenderPrivateMessages, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &privateMessageR{
				Sender: o,
			}
		} else {
			rel.R.Sender = o
		}
	}
	return nil
}

// AddReceiverPrivateMessagesG adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.ReceiverPrivateMessages.
// Sets related.R.Receiver appropriately.
// Uses the global database handle.
func (o *User) AddReceiverPrivateMessagesG(insert bool, related ...*PrivateMessage) error {
	return o.AddReceiverPrivateMessages(boil.GetDB(), insert, related...)
}

// AddReceiverPrivateMessagesP adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.ReceiverPrivateMessages.
// Sets related.R.Receiver appropriately.
// Panics on error.
func (o *User) AddReceiverPrivateMessagesP(exec boil.Executor, insert bool, related ...*PrivateMessage) {
	if err := o.AddReceiverPrivateMessages(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddReceiverPrivateMessagesGP adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.ReceiverPrivateMessages.
// Sets related.R.Receiver appropriately.
// Uses the global database handle and panics on error.
func (o *User) AddReceiverPrivateMessagesGP(insert bool, related ...*PrivateMessage) {
	if err := o.AddReceiverPrivateMessages(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddReceiverPrivateMessages adds the given related objects to the existing relationships
// of the user, optionally inserting them as new records.
// Appends related to o.R.ReceiverPrivateMessages.
// Sets related.R.Receiver appropriately.
func (o *User) AddReceiverPrivateMessages(exec boil.Executor, insert bool, related ...*PrivateMessage) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.ReceiverID = o.Snowflake
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"private_messages\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"receiver_id"}),
				strmangle.WhereClause("\"", "\"", 2, privateMessagePrimaryKeyColumns),
			)
			values := []interface{}{o.Snowflake, rel.Snowflake}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.ReceiverID = o.Snowflake
		}
	}

	if o.R == nil {
		o.R = &userR{
			ReceiverPrivateMessages: related,
		}
	} else {
		o.R.ReceiverPrivateMessages = append(o.R.ReceiverPrivateMessages, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &privateMessageR{
				Receiver: o,
			}
		} else {
			rel.R.Receiver = o
		}
	}
	return nil
}

// UsersG retrieves all records.
func UsersG(mods ...qm.QueryMod) userQuery {
	return Users(boil.GetDB(), mods...)
}

// Users retrieves all the records using an executor.
func Users(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	mods = append(mods, qm.From("\"users\""))
	return userQuery{NewQuery(exec, mods...)}
}

// FindUserG retrieves a single record by ID.
func FindUserG(snowflake int64, selectCols ...string) (*User, error) {
	return FindUser(boil.GetDB(), snowflake, selectCols...)
}

// FindUserGP retrieves a single record by ID, and panics on error.
func FindUserGP(snowflake int64, selectCols ...string) *User {
	retobj, err := FindUser(boil.GetDB(), snowflake, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindUser retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUser(exec boil.Executor, snowflake int64, selectCols ...string) (*User, error) {
	userObj := &User{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"users\" where \"snowflake\"=$1", sel,
	)

	q := queries.Raw(exec, query, snowflake)

	err := q.Bind(userObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from users")
	}

	return userObj, nil
}

// FindUserP retrieves a single record by ID with an executor, and panics on error.
func FindUserP(exec boil.Executor, snowflake int64, selectCols ...string) *User {
	retobj, err := FindUser(exec, snowflake, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *User) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *User) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *User) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *User) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no users provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	userInsertCacheMut.RLock()
	cache, cached := userInsertCache[key]
	userInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			userColumns,
			userColumnsWithDefault,
			userColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(userType, userMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(userType, userMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"users\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into users")
	}

	if !cached {
		userInsertCacheMut.Lock()
		userInsertCache[key] = cache
		userInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single User record. See Update for
// whitelist behavior description.
func (o *User) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single User record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *User) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the User, and panics on error.
// See Update for whitelist behavior description.
func (o *User) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the User.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *User) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	userUpdateCacheMut.RLock()
	cache, cached := userUpdateCache[key]
	userUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(userColumns, userPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update users, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"users\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, userPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userType, userMapping, append(wl, userPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update users row")
	}

	if !cached {
		userUpdateCacheMut.Lock()
		userUpdateCache[key] = cache
		userUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q userQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q userQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for users")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o UserSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o UserSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o UserSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"users\" SET %s WHERE (\"snowflake\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(userPrimaryKeyColumns), len(colNames)+1, len(userPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in user slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *User) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *User) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *User) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *User) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no users provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(userColumnsWithDefault, o)

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

	userUpsertCacheMut.RLock()
	cache, cached := userUpsertCache[key]
	userUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			userColumns,
			userColumnsWithDefault,
			userColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			userColumns,
			userPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert users, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(userPrimaryKeyColumns))
			copy(conflict, userPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"users\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(userType, userMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(userType, userMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert users")
	}

	if !cached {
		userUpsertCacheMut.Lock()
		userUpsertCache[key] = cache
		userUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single User record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *User) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single User record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *User) DeleteG() error {
	if o == nil {
		return errors.New("models: no User provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single User record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *User) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single User record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *User) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no User provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userPrimaryKeyMapping)
	sql := "DELETE FROM \"users\" WHERE \"snowflake\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from users")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q userQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q userQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no userQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from users")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o UserSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o UserSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no User slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o UserSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no User slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(userBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"users\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, userPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(userPrimaryKeyColumns), 1, len(userPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from user slice")
	}

	if len(userAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *User) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *User) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *User) ReloadG() error {
	if o == nil {
		return errors.New("models: no User provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *User) Reload(exec boil.Executor) error {
	ret, err := FindUser(exec, o.Snowflake)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *UserSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *UserSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty UserSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	users := UserSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"users\".* FROM \"users\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, userPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(userPrimaryKeyColumns), 1, len(userPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&users)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in UserSlice")
	}

	*o = users

	return nil
}

// UserExists checks if the User row exists.
func UserExists(exec boil.Executor, snowflake int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"users\" where \"snowflake\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, snowflake)
	}

	row := exec.QueryRow(sql, snowflake)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if users exists")
	}

	return exists, nil
}

// UserExistsG checks if the User row exists.
func UserExistsG(snowflake int64) (bool, error) {
	return UserExists(boil.GetDB(), snowflake)
}

// UserExistsGP checks if the User row exists. Panics on error.
func UserExistsGP(snowflake int64) bool {
	e, err := UserExists(boil.GetDB(), snowflake)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// UserExistsP checks if the User row exists. Panics on error.
func UserExistsP(exec boil.Executor, snowflake int64) bool {
	e, err := UserExists(exec, snowflake)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
