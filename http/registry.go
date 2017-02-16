package http

import "errors"

// ResourceFactory produces a new instance of a particular resource, if withID is true, it should also
// generate a new ID for the resource, otherwise leave the ID field on default or 0.
type ResourceFactory func(withID bool) (Resource, error)

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
	if _, ok := resources[name]; ok {
		return errors.New("Resource Endpoint already registered")
	}
	resourceEndpoints[name] = r
	return nil
}

func makeResource(name string, withID bool) (Resource, error) {
	if _, ok := resources[name]; !ok {
		return nil, errors.New("Resource not registered")
	}
	return resources[name](withID)
}
