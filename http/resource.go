package http

import "encoding/json"

// Resource is a generic interface for objects to be represented in a REST API
type Resource interface {
	Snowflake() int64
	Type() string
	// StripReadOnly defaults or zeroes all fields that should not be accessible to userinput
	StripReadOnly() error

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
// (minus any fields cleared by StripReadOnly)
type ResourceEndpointUpdate interface {
	ResourceEndpoint

	Merge(Resource, Resource) (Resource, error)
	Replace(Resource, Resource) (Resource, error)
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

	SoftDelete(Resource) error
	HardDelete(Resource) error
}
