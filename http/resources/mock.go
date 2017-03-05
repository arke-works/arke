package resources

import (
	"encoding/json"
	"errors"
	"iris.arke.works/forum/snowflakes"
	"time"
)

func SetupMockDatatype() {
	var err error
	err = RegisterResource("mock", func(fountain snowflakes.Fountain) (Resource, error) {
		if fountain == nil {
			return &MockResource{SnowflakeField: -1}, nil
		}
		id, err := fountain.NewID()
		return &MockResource{SnowflakeField: id}, err
	})
	if err != nil {
		panic(err)
	}
	err = RegisterResource("mock_nofactory", func(fountain snowflakes.Fountain) (Resource, error) {
		if fountain == nil {
			return &MockResource{SnowflakeField: -1}, nil
		}
		id, err := fountain.NewID()
		return &MockResource{SnowflakeField: id}, err
	})
	if err != nil {
		panic(err)
	}
	err = RegisterResourceEndpoint("mock", &MockRegistry{})
	if err != nil {
		panic(err)
	}
	err = RegisterResourceEndpoint("mock_noresource", &MockRegistry{})
	if err != nil {
		panic(err)
	}
}

type MockResource struct {
	SnowflakeField int64      `json:"id_field"`
	TextField      string     `json:"text"`
	OtherTextField string     `json:"other_text"`
	IntField       int        `json:"int"`
	SliceField     []byte     `json:"bytes"`
	TimeField      time.Time  `json:"time"`
	OptTimeField   *time.Time `json:"time"`
}

type prvMockResource struct {
	SnowflakeField int64      `json:"id_field"`
	TextField      string     `json:"text"`
	OtherTextField string     `json:"other_text"`
	IntField       int        `json:"int"`
	SliceField     []byte     `json:"bytes"`
	TimeField      time.Time  `json:"time"`
	OptTimeField   *time.Time `json:"time"`
}

func (m *MockResource) MarshalJSON() ([]byte, error) {
	return json.Marshal((*prvMockResource)(m))
}

func (m *MockResource) Merge(i Resource) error {
	if m2, ok := i.(*MockResource); ok {
		if m2.SnowflakeField != 0 {
			m.SnowflakeField = m2.SnowflakeField
		}
		if m2.TextField != "" {
			m.TextField = m2.TextField
		}
		if m2.IntField != 0 {
			m.IntField = m2.IntField
		}
		if m2.OptTimeField != nil {
			m.OptTimeField = m2.OptTimeField
		}
		if m2.TimeField.Unix() != 0 {
			m.TimeField = m2.TimeField
		}
		if m2.SliceField != nil {
			m.SliceField = m2.SliceField
		}
		if m2.OtherTextField != "" {
			m.OtherTextField = m2.OtherTextField
		}
		return nil
	}
	return errors.New("Merging with incompatible type")
}

func (m *MockResource) Snowflake() int64 {
	return m.SnowflakeField
}

func (m *MockResource) StripReadOnly() error {
	m.SnowflakeField = 0
	m.TextField = ""
	return nil
}

func (m *MockResource) Type() string {
	return "mock"
}

func (m *MockResource) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*prvMockResource)(m))
}

var _ Resource = (*MockResource)(nil)

type MockRegistry struct {
	Database map[int64]*MockResource
}

func (mr *MockRegistry) cai() {
	if mr.Database == nil {
		mr.Database = map[int64]*MockResource{}
	}
}

func (mr *MockRegistry) New(r Resource) error {
	mr.cai()

	if m, ok := r.(*MockResource); ok {
		if m != nil {
			mr.Database[m.Snowflake()] = m
			return nil
		}
		return errors.New("Resource was nil")
	}
	return errors.New("Wrong resource endpoint or incompatible resource")
}

func (mr *MockRegistry) Name() string {
	return "mock"
}
