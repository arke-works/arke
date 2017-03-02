package resources

import (
	"errors"
	"iris.arke.works/forum/snowflakes"
)

func init() {
	var err error
	err = RegisterResource("ping", func(fountain snowflakes.Fountain) (Resource, error) {
		if fountain == nil {
			return &Ping{-1}, nil
		}
		id, err := fountain.NewID()
		return &Ping{id}, err
	})
	if err != nil {
		panic(err)
	}
	err = RegisterResourceEndpoint("ping", &PingEndpoint{})
	if err != nil {
		panic(err)
	}
}

// Ping is a mostly inactive resource that is used to test if the API is reactive and online
type Ping struct {
	snowflake int64
}

// Snowflake returns the unique ID of the Ping
func (p *Ping) Snowflake() int64 {
	return p.snowflake
}

// Type returns "ping"
func (p *Ping) Type() string {
	return "ping"
}

// StripReadOnly removes the snowflake value from a Ping
func (p *Ping) StripReadOnly() error {
	p.snowflake = -1
	return nil
}

// MarshalJSON parses the Ping into a "pong" value as JSON
func (p *Ping) MarshalJSON() ([]byte, error) {
	return []byte("\"pong\""), nil
}

// UnmarshalJSON tests if the incoming string is JSON value of "ping" and returns an error
// otherwise
func (p *Ping) UnmarshalJSON(e []byte) error {
	if string(e) != "\"ping\"" {
		return errors.New("Invalid Ping Serialization")
	}
	p.snowflake = -1
	return nil
}

// Merge checks if the given resource is of type Ping and sets the internal snowflake
// to that of the given Ping.
func (p *Ping) Merge(r Resource) error {
	var (
		pr *Ping
		ok bool
	)
	if pr, ok = r.(*Ping); !ok {
		return ErrMergeTypeMismatch
	}
	if pr.snowflake > 0 {
		p.snowflake = pr.snowflake
	}
	return nil
}

var _ Resource = (*Ping)(nil)

// PingEndpoint is used to retrieve pings, a special resource used to test the API connectivity
type PingEndpoint struct{}

// Name returns the value "ping"
func (pe *PingEndpoint) Name() string {
	return "ping"
}

// Find returns a ping resource with a set snowflake.
func (pe *PingEndpoint) Find(snowflake int64) (Resource, error) {
	return &Ping{snowflake: snowflake}, nil
}

// FindAll returns a list of ping resources of the given size. The page value is ignored.
func (pe *PingEndpoint) FindAll(_, size int64) ([]Resource, error) {
	if size < 1 {
		return []Resource{}, nil
	}
	var ret = []Resource{}
	for i := int64(0); i < size; i++ {
		ret = append(ret, &Ping{snowflake: -1})
	}
	return ret, nil
}
