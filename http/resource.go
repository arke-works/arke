package http

import (
	"encoding/json"
	"errors"
)

var (
	// ErrMergeTypeMismatch indicates that a merge failed because the incoming type
	// was not the same as the expected resource type
	ErrMergeTypeMismatch = errors.New("Type mismatch during merge operation")
)

// Resource is a generic interface for objects to be represented in a REST API
type Resource interface {
	Snowflake() int64
	Type() string
	// StripReadOnly defaults or zeroes all fields that should not be accessible to userinput
	// This will usually be called before writing data into the database
	StripReadOnly() error

	// Merge will overwrite all fields of the current resource with those of the specified
	// resource that are not of a null value
	//
	// If nil is provided, it should take no action and return no error
	//
	// The incoming resource **must** be type checked, if the types mismatch, no merging
	// is to take place. If this is the case, the error ErrMergeTypeMismatch must be returned
	Merge(Resource) error

	// A Resource MUST implement JSON Marshaler and Unmarshaler
	json.Marshaler
	json.Unmarshaler
}

// ResourceEndpoint is the minimal interface for a REST Endpoint that allows no HTTP Methods but OPTION which
// at minimum returns only OPTION as accepted Method.
type ResourceEndpoint interface {
	Name() string
}

// ResourceEndpointNew defines an interface for REST Endpoints that allows creating new instances of a resource
type ResourceEndpointNew interface {
	ResourceEndpoint

	New(Resource) error
}

// ResourceEndpointFind defines an interface for REST Endpoints that allows GET for a single or multiple objects
type ResourceEndpointFind interface {
	ResourceEndpoint

	Find(snowflake int64) (Resource, error)
	FindAll(page, size int64) ([]Resource, error)
}

// ResourceEndpointUpdate defines an interface for merging or replacing an object
//
// Merging refers to supplying a Resource where only some fields will be set, only fields that are not empty or null
// should overwrite existing fields
//
// Replace refers to supplying a full Resource where the given object should completely replace the existing one
type ResourceEndpointUpdate interface {
	ResourceEndpoint

	Merge(int64, Resource) (Resource, error)
	Replace(int64, Resource) (Resource, error)
}

// ResourceEndpointDelete defines an interface for deleting an object
//
// A soft delete refers to removing an object only virtually by marking the object deleted without actually deleting it
// This preserves existing relationship and should be prefered over a hard delete. Soft deletes are reversible (but this
// might not be exposed)
//
// A hard delete issues a delete into the underlying database and removes the object permanently.
type ResourceEndpointDelete interface {
	ResourceEndpoint

	SoftDelete(int64) error
	HardDelete(int64) error
}
