package resources

import (
	"iris.arke.works/forum/db/models"
	"encoding/json"
	"errors"
)

type ReplyResource struct {
	int *models.Reply
}

func (r *ReplyResource) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.int)
}

func (r *ReplyResource) Merge(res Resource) error {
	if r2, ok := res.(*ReplyResource); ok {
		if r2.int.Snowflake != 0 {
			r.int.Snowflake = r2.int.Snowflake
		}
		if r2.int.CreatedAt != nil {
			r.int.CreatedAt = r2.int.CreatedAt
		}
		if r2.int.AuthorID.Valid && r2.int.AuthorID.Int64 != 0 {
			r.int.AuthorID = r2.int.AuthorID
		}
		if r2.int.ParentID.Valid && r2.int.ParentID.Int64 != 0 {
			r.int.ParentID = r2.int.ParentID
		}
		if r2.int.TopicID != 0 {
			r.int.TopicID = r2.int.TopicID
		}
		return nil
	}
	return errors.New("Merging with incompatible type")
}

func (r *ReplyResource) Snowflake() int64 {
	return r.int.Snowflake
}

func (r *ReplyResource) StripReadOnly() error {
	r.int.Snowflake = 0
	r.int.CreatedAt = nil
	r.int.AuthorID.Valid = false
	r.int.AuthorID.Int64 = 0
	r.int.ParentID.Valid = false
	r.int.ParentID.Int64 = 0
	r.int.TopicID = 0
}

func (r *ReplyResource) Type() string {
	return "reply"
}

func (r *ReplyResource) UnmarshalJSON(data []byte) error {
	if r.int == nil {
		r.int = &models.Reply{}
	}
	return json.Unmarshal(data, r.int)
}


