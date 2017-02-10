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

// PrivateMessage is an object representing the database table.
type PrivateMessage struct {
	Snowflake  int64      `boil:"snowflake" json:"snowflake" toml:"snowflake" yaml:"snowflake"`
	CreatedAt  time.Time  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	DeletedAt  null.Time  `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`
	Title      string     `boil:"title" json:"title" toml:"title" yaml:"title"`
	Body       string     `boil:"body" json:"body" toml:"body" yaml:"body"`
	SenderID   int64      `boil:"sender_id" json:"sender_id" toml:"sender_id" yaml:"sender_id"`
	ReceiverID int64      `boil:"receiver_id" json:"receiver_id" toml:"receiver_id" yaml:"receiver_id"`
	ParentID   null.Int64 `boil:"parent_id" json:"parent_id,omitempty" toml:"parent_id" yaml:"parent_id,omitempty"`

	R *privateMessageR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L privateMessageL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// privateMessageR is where relationships are stored.
type privateMessageR struct {
	Sender                *User
	Receiver              *User
	Parent                *PrivateMessage
	ParentPrivateMessages PrivateMessageSlice
}

// privateMessageL is where Load methods for each relationship are stored.
type privateMessageL struct{}

var (
	privateMessageColumns               = []string{"snowflake", "created_at", "deleted_at", "title", "body", "sender_id", "receiver_id", "parent_id"}
	privateMessageColumnsWithoutDefault = []string{"snowflake", "deleted_at", "title", "body", "sender_id", "receiver_id", "parent_id"}
	privateMessageColumnsWithDefault    = []string{"created_at"}
	privateMessagePrimaryKeyColumns     = []string{"snowflake"}
)

type (
	// PrivateMessageSlice is an alias for a slice of pointers to PrivateMessage.
	// This should generally be used opposed to []PrivateMessage.
	PrivateMessageSlice []*PrivateMessage
	// PrivateMessageHook is the signature for custom PrivateMessage hook methods
	PrivateMessageHook func(boil.Executor, *PrivateMessage) error

	privateMessageQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	privateMessageType                 = reflect.TypeOf(&PrivateMessage{})
	privateMessageMapping              = queries.MakeStructMapping(privateMessageType)
	privateMessagePrimaryKeyMapping, _ = queries.BindMapping(privateMessageType, privateMessageMapping, privateMessagePrimaryKeyColumns)
	privateMessageInsertCacheMut       sync.RWMutex
	privateMessageInsertCache          = make(map[string]insertCache)
	privateMessageUpdateCacheMut       sync.RWMutex
	privateMessageUpdateCache          = make(map[string]updateCache)
	privateMessageUpsertCacheMut       sync.RWMutex
	privateMessageUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var privateMessageBeforeInsertHooks []PrivateMessageHook
var privateMessageBeforeUpdateHooks []PrivateMessageHook
var privateMessageBeforeDeleteHooks []PrivateMessageHook
var privateMessageBeforeUpsertHooks []PrivateMessageHook

var privateMessageAfterInsertHooks []PrivateMessageHook
var privateMessageAfterSelectHooks []PrivateMessageHook
var privateMessageAfterUpdateHooks []PrivateMessageHook
var privateMessageAfterDeleteHooks []PrivateMessageHook
var privateMessageAfterUpsertHooks []PrivateMessageHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *PrivateMessage) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range privateMessageBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *PrivateMessage) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range privateMessageBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *PrivateMessage) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range privateMessageBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *PrivateMessage) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range privateMessageBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *PrivateMessage) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range privateMessageAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *PrivateMessage) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range privateMessageAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *PrivateMessage) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range privateMessageAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *PrivateMessage) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range privateMessageAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *PrivateMessage) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range privateMessageAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddPrivateMessageHook registers your hook function for all future operations.
func AddPrivateMessageHook(hookPoint boil.HookPoint, privateMessageHook PrivateMessageHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		privateMessageBeforeInsertHooks = append(privateMessageBeforeInsertHooks, privateMessageHook)
	case boil.BeforeUpdateHook:
		privateMessageBeforeUpdateHooks = append(privateMessageBeforeUpdateHooks, privateMessageHook)
	case boil.BeforeDeleteHook:
		privateMessageBeforeDeleteHooks = append(privateMessageBeforeDeleteHooks, privateMessageHook)
	case boil.BeforeUpsertHook:
		privateMessageBeforeUpsertHooks = append(privateMessageBeforeUpsertHooks, privateMessageHook)
	case boil.AfterInsertHook:
		privateMessageAfterInsertHooks = append(privateMessageAfterInsertHooks, privateMessageHook)
	case boil.AfterSelectHook:
		privateMessageAfterSelectHooks = append(privateMessageAfterSelectHooks, privateMessageHook)
	case boil.AfterUpdateHook:
		privateMessageAfterUpdateHooks = append(privateMessageAfterUpdateHooks, privateMessageHook)
	case boil.AfterDeleteHook:
		privateMessageAfterDeleteHooks = append(privateMessageAfterDeleteHooks, privateMessageHook)
	case boil.AfterUpsertHook:
		privateMessageAfterUpsertHooks = append(privateMessageAfterUpsertHooks, privateMessageHook)
	}
}

// OneP returns a single privateMessage record from the query, and panics on error.
func (q privateMessageQuery) OneP() *PrivateMessage {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single privateMessage record from the query.
func (q privateMessageQuery) One() (*PrivateMessage, error) {
	o := &PrivateMessage{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for private_messages")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all PrivateMessage records from the query, and panics on error.
func (q privateMessageQuery) AllP() PrivateMessageSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all PrivateMessage records from the query.
func (q privateMessageQuery) All() (PrivateMessageSlice, error) {
	var o PrivateMessageSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to PrivateMessage slice")
	}

	if len(privateMessageAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all PrivateMessage records in the query, and panics on error.
func (q privateMessageQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all PrivateMessage records in the query.
func (q privateMessageQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count private_messages rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q privateMessageQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q privateMessageQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if private_messages exists")
	}

	return count > 0, nil
}

// SenderG pointed to by the foreign key.
func (o *PrivateMessage) SenderG(mods ...qm.QueryMod) userQuery {
	return o.Sender(boil.GetDB(), mods...)
}

// Sender pointed to by the foreign key.
func (o *PrivateMessage) Sender(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("snowflake=?", o.SenderID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(exec, queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// ReceiverG pointed to by the foreign key.
func (o *PrivateMessage) ReceiverG(mods ...qm.QueryMod) userQuery {
	return o.Receiver(boil.GetDB(), mods...)
}

// Receiver pointed to by the foreign key.
func (o *PrivateMessage) Receiver(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("snowflake=?", o.ReceiverID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(exec, queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// ParentG pointed to by the foreign key.
func (o *PrivateMessage) ParentG(mods ...qm.QueryMod) privateMessageQuery {
	return o.Parent(boil.GetDB(), mods...)
}

// Parent pointed to by the foreign key.
func (o *PrivateMessage) Parent(exec boil.Executor, mods ...qm.QueryMod) privateMessageQuery {
	queryMods := []qm.QueryMod{
		qm.Where("snowflake=?", o.ParentID),
	}

	queryMods = append(queryMods, mods...)

	query := PrivateMessages(exec, queryMods...)
	queries.SetFrom(query.Query, "\"private_messages\"")

	return query
}

// ParentPrivateMessagesG retrieves all the private_message's private messages via parent_id column.
func (o *PrivateMessage) ParentPrivateMessagesG(mods ...qm.QueryMod) privateMessageQuery {
	return o.ParentPrivateMessages(boil.GetDB(), mods...)
}

// ParentPrivateMessages retrieves all the private_message's private messages with an executor via parent_id column.
func (o *PrivateMessage) ParentPrivateMessages(exec boil.Executor, mods ...qm.QueryMod) privateMessageQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"parent_id\"=?", o.Snowflake),
	)

	query := PrivateMessages(exec, queryMods...)
	queries.SetFrom(query.Query, "\"private_messages\" as \"a\"")
	return query
}

// LoadSender allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (privateMessageL) LoadSender(e boil.Executor, singular bool, maybePrivateMessage interface{}) error {
	var slice []*PrivateMessage
	var object *PrivateMessage

	count := 1
	if singular {
		object = maybePrivateMessage.(*PrivateMessage)
	} else {
		slice = *maybePrivateMessage.(*PrivateMessageSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &privateMessageR{}
		}
		args[0] = object.SenderID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &privateMessageR{}
			}
			args[i] = obj.SenderID
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

	if len(privateMessageAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		object.R.Sender = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.SenderID == foreign.Snowflake {
				local.R.Sender = foreign
				break
			}
		}
	}

	return nil
}

// LoadReceiver allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (privateMessageL) LoadReceiver(e boil.Executor, singular bool, maybePrivateMessage interface{}) error {
	var slice []*PrivateMessage
	var object *PrivateMessage

	count := 1
	if singular {
		object = maybePrivateMessage.(*PrivateMessage)
	} else {
		slice = *maybePrivateMessage.(*PrivateMessageSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &privateMessageR{}
		}
		args[0] = object.ReceiverID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &privateMessageR{}
			}
			args[i] = obj.ReceiverID
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

	if len(privateMessageAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		object.R.Receiver = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ReceiverID == foreign.Snowflake {
				local.R.Receiver = foreign
				break
			}
		}
	}

	return nil
}

// LoadParent allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (privateMessageL) LoadParent(e boil.Executor, singular bool, maybePrivateMessage interface{}) error {
	var slice []*PrivateMessage
	var object *PrivateMessage

	count := 1
	if singular {
		object = maybePrivateMessage.(*PrivateMessage)
	} else {
		slice = *maybePrivateMessage.(*PrivateMessageSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &privateMessageR{}
		}
		args[0] = object.ParentID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &privateMessageR{}
			}
			args[i] = obj.ParentID
		}
	}

	query := fmt.Sprintf(
		"select * from \"private_messages\" where \"snowflake\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load PrivateMessage")
	}
	defer results.Close()

	var resultSlice []*PrivateMessage
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice PrivateMessage")
	}

	if len(privateMessageAfterSelectHooks) != 0 {
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

// LoadParentPrivateMessages allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (privateMessageL) LoadParentPrivateMessages(e boil.Executor, singular bool, maybePrivateMessage interface{}) error {
	var slice []*PrivateMessage
	var object *PrivateMessage

	count := 1
	if singular {
		object = maybePrivateMessage.(*PrivateMessage)
	} else {
		slice = *maybePrivateMessage.(*PrivateMessageSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &privateMessageR{}
		}
		args[0] = object.Snowflake
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &privateMessageR{}
			}
			args[i] = obj.Snowflake
		}
	}

	query := fmt.Sprintf(
		"select * from \"private_messages\" where \"parent_id\" in (%s)",
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
		object.R.ParentPrivateMessages = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.Snowflake == foreign.ParentID.Int64 {
				local.R.ParentPrivateMessages = append(local.R.ParentPrivateMessages, foreign)
				break
			}
		}
	}

	return nil
}

// SetSenderG of the private_message to the related item.
// Sets o.R.Sender to related.
// Adds o to related.R.SenderPrivateMessages.
// Uses the global database handle.
func (o *PrivateMessage) SetSenderG(insert bool, related *User) error {
	return o.SetSender(boil.GetDB(), insert, related)
}

// SetSenderP of the private_message to the related item.
// Sets o.R.Sender to related.
// Adds o to related.R.SenderPrivateMessages.
// Panics on error.
func (o *PrivateMessage) SetSenderP(exec boil.Executor, insert bool, related *User) {
	if err := o.SetSender(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetSenderGP of the private_message to the related item.
// Sets o.R.Sender to related.
// Adds o to related.R.SenderPrivateMessages.
// Uses the global database handle and panics on error.
func (o *PrivateMessage) SetSenderGP(insert bool, related *User) {
	if err := o.SetSender(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetSender of the private_message to the related item.
// Sets o.R.Sender to related.
// Adds o to related.R.SenderPrivateMessages.
func (o *PrivateMessage) SetSender(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"private_messages\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"sender_id"}),
		strmangle.WhereClause("\"", "\"", 2, privateMessagePrimaryKeyColumns),
	)
	values := []interface{}{related.Snowflake, o.Snowflake}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.SenderID = related.Snowflake

	if o.R == nil {
		o.R = &privateMessageR{
			Sender: related,
		}
	} else {
		o.R.Sender = related
	}

	if related.R == nil {
		related.R = &userR{
			SenderPrivateMessages: PrivateMessageSlice{o},
		}
	} else {
		related.R.SenderPrivateMessages = append(related.R.SenderPrivateMessages, o)
	}

	return nil
}

// SetReceiverG of the private_message to the related item.
// Sets o.R.Receiver to related.
// Adds o to related.R.ReceiverPrivateMessages.
// Uses the global database handle.
func (o *PrivateMessage) SetReceiverG(insert bool, related *User) error {
	return o.SetReceiver(boil.GetDB(), insert, related)
}

// SetReceiverP of the private_message to the related item.
// Sets o.R.Receiver to related.
// Adds o to related.R.ReceiverPrivateMessages.
// Panics on error.
func (o *PrivateMessage) SetReceiverP(exec boil.Executor, insert bool, related *User) {
	if err := o.SetReceiver(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetReceiverGP of the private_message to the related item.
// Sets o.R.Receiver to related.
// Adds o to related.R.ReceiverPrivateMessages.
// Uses the global database handle and panics on error.
func (o *PrivateMessage) SetReceiverGP(insert bool, related *User) {
	if err := o.SetReceiver(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetReceiver of the private_message to the related item.
// Sets o.R.Receiver to related.
// Adds o to related.R.ReceiverPrivateMessages.
func (o *PrivateMessage) SetReceiver(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"private_messages\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"receiver_id"}),
		strmangle.WhereClause("\"", "\"", 2, privateMessagePrimaryKeyColumns),
	)
	values := []interface{}{related.Snowflake, o.Snowflake}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ReceiverID = related.Snowflake

	if o.R == nil {
		o.R = &privateMessageR{
			Receiver: related,
		}
	} else {
		o.R.Receiver = related
	}

	if related.R == nil {
		related.R = &userR{
			ReceiverPrivateMessages: PrivateMessageSlice{o},
		}
	} else {
		related.R.ReceiverPrivateMessages = append(related.R.ReceiverPrivateMessages, o)
	}

	return nil
}

// SetParentG of the private_message to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentPrivateMessages.
// Uses the global database handle.
func (o *PrivateMessage) SetParentG(insert bool, related *PrivateMessage) error {
	return o.SetParent(boil.GetDB(), insert, related)
}

// SetParentP of the private_message to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentPrivateMessages.
// Panics on error.
func (o *PrivateMessage) SetParentP(exec boil.Executor, insert bool, related *PrivateMessage) {
	if err := o.SetParent(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParentGP of the private_message to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentPrivateMessages.
// Uses the global database handle and panics on error.
func (o *PrivateMessage) SetParentGP(insert bool, related *PrivateMessage) {
	if err := o.SetParent(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParent of the private_message to the related item.
// Sets o.R.Parent to related.
// Adds o to related.R.ParentPrivateMessages.
func (o *PrivateMessage) SetParent(exec boil.Executor, insert bool, related *PrivateMessage) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"private_messages\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"parent_id"}),
		strmangle.WhereClause("\"", "\"", 2, privateMessagePrimaryKeyColumns),
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
		o.R = &privateMessageR{
			Parent: related,
		}
	} else {
		o.R.Parent = related
	}

	if related.R == nil {
		related.R = &privateMessageR{
			ParentPrivateMessages: PrivateMessageSlice{o},
		}
	} else {
		related.R.ParentPrivateMessages = append(related.R.ParentPrivateMessages, o)
	}

	return nil
}

// RemoveParentG relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle.
func (o *PrivateMessage) RemoveParentG(related *PrivateMessage) error {
	return o.RemoveParent(boil.GetDB(), related)
}

// RemoveParentP relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Panics on error.
func (o *PrivateMessage) RemoveParentP(exec boil.Executor, related *PrivateMessage) {
	if err := o.RemoveParent(exec, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParentGP relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle and panics on error.
func (o *PrivateMessage) RemoveParentGP(related *PrivateMessage) {
	if err := o.RemoveParent(boil.GetDB(), related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParent relationship.
// Sets o.R.Parent to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *PrivateMessage) RemoveParent(exec boil.Executor, related *PrivateMessage) error {
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

	for i, ri := range related.R.ParentPrivateMessages {
		if o.ParentID.Int64 != ri.ParentID.Int64 {
			continue
		}

		ln := len(related.R.ParentPrivateMessages)
		if ln > 1 && i < ln-1 {
			related.R.ParentPrivateMessages[i] = related.R.ParentPrivateMessages[ln-1]
		}
		related.R.ParentPrivateMessages = related.R.ParentPrivateMessages[:ln-1]
		break
	}
	return nil
}

// AddParentPrivateMessagesG adds the given related objects to the existing relationships
// of the private_message, optionally inserting them as new records.
// Appends related to o.R.ParentPrivateMessages.
// Sets related.R.Parent appropriately.
// Uses the global database handle.
func (o *PrivateMessage) AddParentPrivateMessagesG(insert bool, related ...*PrivateMessage) error {
	return o.AddParentPrivateMessages(boil.GetDB(), insert, related...)
}

// AddParentPrivateMessagesP adds the given related objects to the existing relationships
// of the private_message, optionally inserting them as new records.
// Appends related to o.R.ParentPrivateMessages.
// Sets related.R.Parent appropriately.
// Panics on error.
func (o *PrivateMessage) AddParentPrivateMessagesP(exec boil.Executor, insert bool, related ...*PrivateMessage) {
	if err := o.AddParentPrivateMessages(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddParentPrivateMessagesGP adds the given related objects to the existing relationships
// of the private_message, optionally inserting them as new records.
// Appends related to o.R.ParentPrivateMessages.
// Sets related.R.Parent appropriately.
// Uses the global database handle and panics on error.
func (o *PrivateMessage) AddParentPrivateMessagesGP(insert bool, related ...*PrivateMessage) {
	if err := o.AddParentPrivateMessages(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddParentPrivateMessages adds the given related objects to the existing relationships
// of the private_message, optionally inserting them as new records.
// Appends related to o.R.ParentPrivateMessages.
// Sets related.R.Parent appropriately.
func (o *PrivateMessage) AddParentPrivateMessages(exec boil.Executor, insert bool, related ...*PrivateMessage) error {
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
				"UPDATE \"private_messages\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"parent_id"}),
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

			rel.ParentID.Int64 = o.Snowflake
			rel.ParentID.Valid = true
		}
	}

	if o.R == nil {
		o.R = &privateMessageR{
			ParentPrivateMessages: related,
		}
	} else {
		o.R.ParentPrivateMessages = append(o.R.ParentPrivateMessages, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &privateMessageR{
				Parent: o,
			}
		} else {
			rel.R.Parent = o
		}
	}
	return nil
}

// SetParentPrivateMessagesG removes all previously related items of the
// private_message replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentPrivateMessages accordingly.
// Replaces o.R.ParentPrivateMessages with related.
// Sets related.R.Parent's ParentPrivateMessages accordingly.
// Uses the global database handle.
func (o *PrivateMessage) SetParentPrivateMessagesG(insert bool, related ...*PrivateMessage) error {
	return o.SetParentPrivateMessages(boil.GetDB(), insert, related...)
}

// SetParentPrivateMessagesP removes all previously related items of the
// private_message replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentPrivateMessages accordingly.
// Replaces o.R.ParentPrivateMessages with related.
// Sets related.R.Parent's ParentPrivateMessages accordingly.
// Panics on error.
func (o *PrivateMessage) SetParentPrivateMessagesP(exec boil.Executor, insert bool, related ...*PrivateMessage) {
	if err := o.SetParentPrivateMessages(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParentPrivateMessagesGP removes all previously related items of the
// private_message replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentPrivateMessages accordingly.
// Replaces o.R.ParentPrivateMessages with related.
// Sets related.R.Parent's ParentPrivateMessages accordingly.
// Uses the global database handle and panics on error.
func (o *PrivateMessage) SetParentPrivateMessagesGP(insert bool, related ...*PrivateMessage) {
	if err := o.SetParentPrivateMessages(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetParentPrivateMessages removes all previously related items of the
// private_message replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Parent's ParentPrivateMessages accordingly.
// Replaces o.R.ParentPrivateMessages with related.
// Sets related.R.Parent's ParentPrivateMessages accordingly.
func (o *PrivateMessage) SetParentPrivateMessages(exec boil.Executor, insert bool, related ...*PrivateMessage) error {
	query := "update \"private_messages\" set \"parent_id\" = null where \"parent_id\" = $1"
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
		for _, rel := range o.R.ParentPrivateMessages {
			rel.ParentID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Parent = nil
		}

		o.R.ParentPrivateMessages = nil
	}
	return o.AddParentPrivateMessages(exec, insert, related...)
}

// RemoveParentPrivateMessagesG relationships from objects passed in.
// Removes related items from R.ParentPrivateMessages (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
// Uses the global database handle.
func (o *PrivateMessage) RemoveParentPrivateMessagesG(related ...*PrivateMessage) error {
	return o.RemoveParentPrivateMessages(boil.GetDB(), related...)
}

// RemoveParentPrivateMessagesP relationships from objects passed in.
// Removes related items from R.ParentPrivateMessages (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
// Panics on error.
func (o *PrivateMessage) RemoveParentPrivateMessagesP(exec boil.Executor, related ...*PrivateMessage) {
	if err := o.RemoveParentPrivateMessages(exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParentPrivateMessagesGP relationships from objects passed in.
// Removes related items from R.ParentPrivateMessages (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
// Uses the global database handle and panics on error.
func (o *PrivateMessage) RemoveParentPrivateMessagesGP(related ...*PrivateMessage) {
	if err := o.RemoveParentPrivateMessages(boil.GetDB(), related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveParentPrivateMessages relationships from objects passed in.
// Removes related items from R.ParentPrivateMessages (uses pointer comparison, removal does not keep order)
// Sets related.R.Parent.
func (o *PrivateMessage) RemoveParentPrivateMessages(exec boil.Executor, related ...*PrivateMessage) error {
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
		for i, ri := range o.R.ParentPrivateMessages {
			if rel != ri {
				continue
			}

			ln := len(o.R.ParentPrivateMessages)
			if ln > 1 && i < ln-1 {
				o.R.ParentPrivateMessages[i] = o.R.ParentPrivateMessages[ln-1]
			}
			o.R.ParentPrivateMessages = o.R.ParentPrivateMessages[:ln-1]
			break
		}
	}

	return nil
}

// PrivateMessagesG retrieves all records.
func PrivateMessagesG(mods ...qm.QueryMod) privateMessageQuery {
	return PrivateMessages(boil.GetDB(), mods...)
}

// PrivateMessages retrieves all the records using an executor.
func PrivateMessages(exec boil.Executor, mods ...qm.QueryMod) privateMessageQuery {
	mods = append(mods, qm.From("\"private_messages\""))
	return privateMessageQuery{NewQuery(exec, mods...)}
}

// FindPrivateMessageG retrieves a single record by ID.
func FindPrivateMessageG(snowflake int64, selectCols ...string) (*PrivateMessage, error) {
	return FindPrivateMessage(boil.GetDB(), snowflake, selectCols...)
}

// FindPrivateMessageGP retrieves a single record by ID, and panics on error.
func FindPrivateMessageGP(snowflake int64, selectCols ...string) *PrivateMessage {
	retobj, err := FindPrivateMessage(boil.GetDB(), snowflake, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindPrivateMessage retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPrivateMessage(exec boil.Executor, snowflake int64, selectCols ...string) (*PrivateMessage, error) {
	privateMessageObj := &PrivateMessage{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"private_messages\" where \"snowflake\"=$1", sel,
	)

	q := queries.Raw(exec, query, snowflake)

	err := q.Bind(privateMessageObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from private_messages")
	}

	return privateMessageObj, nil
}

// FindPrivateMessageP retrieves a single record by ID with an executor, and panics on error.
func FindPrivateMessageP(exec boil.Executor, snowflake int64, selectCols ...string) *PrivateMessage {
	retobj, err := FindPrivateMessage(exec, snowflake, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *PrivateMessage) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *PrivateMessage) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *PrivateMessage) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *PrivateMessage) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no private_messages provided for insertion")
	}

	var err error
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(privateMessageColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	privateMessageInsertCacheMut.RLock()
	cache, cached := privateMessageInsertCache[key]
	privateMessageInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			privateMessageColumns,
			privateMessageColumnsWithDefault,
			privateMessageColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(privateMessageType, privateMessageMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(privateMessageType, privateMessageMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"private_messages\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into private_messages")
	}

	if !cached {
		privateMessageInsertCacheMut.Lock()
		privateMessageInsertCache[key] = cache
		privateMessageInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single PrivateMessage record. See Update for
// whitelist behavior description.
func (o *PrivateMessage) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single PrivateMessage record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *PrivateMessage) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the PrivateMessage, and panics on error.
// See Update for whitelist behavior description.
func (o *PrivateMessage) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the PrivateMessage.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *PrivateMessage) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	privateMessageUpdateCacheMut.RLock()
	cache, cached := privateMessageUpdateCache[key]
	privateMessageUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(privateMessageColumns, privateMessagePrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update private_messages, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"private_messages\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, privateMessagePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(privateMessageType, privateMessageMapping, append(wl, privateMessagePrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update private_messages row")
	}

	if !cached {
		privateMessageUpdateCacheMut.Lock()
		privateMessageUpdateCache[key] = cache
		privateMessageUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q privateMessageQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q privateMessageQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for private_messages")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o PrivateMessageSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o PrivateMessageSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o PrivateMessageSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PrivateMessageSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), privateMessagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"private_messages\" SET %s WHERE (\"snowflake\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(privateMessagePrimaryKeyColumns), len(colNames)+1, len(privateMessagePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in privateMessage slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *PrivateMessage) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *PrivateMessage) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *PrivateMessage) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *PrivateMessage) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no private_messages provided for upsert")
	}
	currTime := time.Now().In(boil.GetLocation())

	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(privateMessageColumnsWithDefault, o)

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

	privateMessageUpsertCacheMut.RLock()
	cache, cached := privateMessageUpsertCache[key]
	privateMessageUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			privateMessageColumns,
			privateMessageColumnsWithDefault,
			privateMessageColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			privateMessageColumns,
			privateMessagePrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert private_messages, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(privateMessagePrimaryKeyColumns))
			copy(conflict, privateMessagePrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"private_messages\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(privateMessageType, privateMessageMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(privateMessageType, privateMessageMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert private_messages")
	}

	if !cached {
		privateMessageUpsertCacheMut.Lock()
		privateMessageUpsertCache[key] = cache
		privateMessageUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single PrivateMessage record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *PrivateMessage) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single PrivateMessage record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *PrivateMessage) DeleteG() error {
	if o == nil {
		return errors.New("models: no PrivateMessage provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single PrivateMessage record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *PrivateMessage) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single PrivateMessage record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *PrivateMessage) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no PrivateMessage provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), privateMessagePrimaryKeyMapping)
	sql := "DELETE FROM \"private_messages\" WHERE \"snowflake\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from private_messages")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q privateMessageQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q privateMessageQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no privateMessageQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from private_messages")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o PrivateMessageSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o PrivateMessageSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no PrivateMessage slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o PrivateMessageSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PrivateMessageSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no PrivateMessage slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(privateMessageBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), privateMessagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"private_messages\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, privateMessagePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(privateMessagePrimaryKeyColumns), 1, len(privateMessagePrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from privateMessage slice")
	}

	if len(privateMessageAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *PrivateMessage) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *PrivateMessage) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *PrivateMessage) ReloadG() error {
	if o == nil {
		return errors.New("models: no PrivateMessage provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *PrivateMessage) Reload(exec boil.Executor) error {
	ret, err := FindPrivateMessage(exec, o.Snowflake)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *PrivateMessageSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *PrivateMessageSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PrivateMessageSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty PrivateMessageSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PrivateMessageSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	privateMessages := PrivateMessageSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), privateMessagePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"private_messages\".* FROM \"private_messages\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, privateMessagePrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(privateMessagePrimaryKeyColumns), 1, len(privateMessagePrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&privateMessages)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in PrivateMessageSlice")
	}

	*o = privateMessages

	return nil
}

// PrivateMessageExists checks if the PrivateMessage row exists.
func PrivateMessageExists(exec boil.Executor, snowflake int64) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"private_messages\" where \"snowflake\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, snowflake)
	}

	row := exec.QueryRow(sql, snowflake)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if private_messages exists")
	}

	return exists, nil
}

// PrivateMessageExistsG checks if the PrivateMessage row exists.
func PrivateMessageExistsG(snowflake int64) (bool, error) {
	return PrivateMessageExists(boil.GetDB(), snowflake)
}

// PrivateMessageExistsGP checks if the PrivateMessage row exists. Panics on error.
func PrivateMessageExistsGP(snowflake int64) bool {
	e, err := PrivateMessageExists(boil.GetDB(), snowflake)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// PrivateMessageExistsP checks if the PrivateMessage row exists. Panics on error.
func PrivateMessageExistsP(exec boil.Executor, snowflake int64) bool {
	e, err := PrivateMessageExists(exec, snowflake)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
