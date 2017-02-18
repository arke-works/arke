package http

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

type Ping struct {
	snowflake int64
}

func (p *Ping) Snowflake() int64 {
	return p.snowflake
}

func (p *Ping) Type() string {
	return "ping"
}

func (p *Ping) StripReadOnly() error {
	return nil
}

func (p *Ping) MarshalJSON() ([]byte, error) {
	return []byte("\"pong\""), nil
}

func (p *Ping) UnmarshalJSON(e []byte) error {
	if string(e) != "\"ping\"" {
		return errors.New("Invalid Ping Serialization")
	}
	p.snowflake = -1
	return nil
}

var _ Resource = (*Ping)(nil)

type PingEndpoint struct{}

func (pe *PingEndpoint) Name() string {
	return "ping"
}

func (pe *PingEndpoint) Find(snowflake int64) (Resource, error) {
	return &Ping{snowflake: snowflake}, nil
}

func (pe *PingEndpoint) FindAll(page, size int64) ([]Resource, error) {
	if size < 1 {
		return []Resource{}, nil
	}
	var ret = []Resource{}
	for i := int64(0); i < size; i++ {
		ret = append(ret, &Ping{snowflake: -1})
	}
	return ret, nil
}
