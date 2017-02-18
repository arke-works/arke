package http

import (
	"errors"
	"iris.arke.works/forum/snowflakes"
)

// ResourceFactory produces a new instance of a particular resource
//
// The fountain can either be a value, in which case a new ID is to be generated
// or it is nil, in which case no ID value is desired (used to generate new resources
// for unmarshalling)
type ResourceFactory func(fountain snowflakes.Fountain) (Resource, error)

var resources = map[string]ResourceFactory{}

var resourceEndpoints = map[string]ResourceEndpoint{}

// RegisterResource registers a resource on the global map. If the name already exists
// the function errors.
func RegisterResource(name string, r ResourceFactory) error {
	if _, ok := resources[name]; ok {
		return errors.New("Resource already registered")
	}
	resources[name] = r
	return nil
}

// RegisterResourceEndpoint saves a named resource endpoint into the global registry.
// It returns an error if the resource already exists.
func RegisterResourceEndpoint(name string, r ResourceEndpoint) error {
	if _, ok := resourceEndpoints[name]; ok {
		return errors.New("Resource Endpoint already registered")
	}
	resourceEndpoints[name] = r
	return nil
}

func makeResource(name string, fountain snowflakes.Fountain) (Resource, error) {
	if _, ok := resources[name]; !ok {
		return nil, errors.New("Resource not registered")
	}
	return resources[name](fountain)
}
