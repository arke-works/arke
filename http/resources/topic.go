package resources

import (
	"encoding/json"
	"errors"
	"iris.arke.works/forum/db/models"
)

type TopicResource struct {
	int *models.Topic
}

func (t *TopicResource) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.int)
}

func (t *TopicResource) Merge(res Resource) error {
	if t2, ok := res.(*TopicResource); ok {
		if t2.int.AuthorID.Valid &&
			t2.int.AuthorID.Int64 != 0 {
			t.int.AuthorID = t2.int.AuthorID
		}
		if t2.int.Snowflake != 0 {
			t.int.Snowflake = t2.int.Snowflake
		}
		if t2.int.CreatedAt != nil {
			t.int.CreatedAt = t2.int.CreatedAt
		}
		if t2.int.Revision >= 0 {
			t.int.Revision = t2.int.Revision
		}
		return nil
	}
	return errors.New("Merging with incompatible type")
}

func (t *TopicResource) Snowflake() int64 {
	return t.int.Snowflake
}

func (t *TopicResource) StripReadOnly() error {
	t.int.AuthorID.Valid = false
	t.int.AuthorID.Int64 = 0
	t.int.Snowflake = 0
	t.int.CreatedAt = nil
	t.int.Revision = -1
}

func (t *TopicResource) Type() string {
	return "topic"
}

func (t *TopicResource) UnmarshalJSON(data []byte) error {
	if t.int == nil {
		t.int = &models.Topic{}
	}
	return json.Unmarshal(data, t.int)
}
