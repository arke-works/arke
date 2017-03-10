package resources

import (
	"encoding/json"
	"errors"
	"iris.arke.works/forum/db/models"
	"iris.arke.works/forum/snowflakes"
	"time"
)

func init() {
	if err := RegisterResource("user", func(fountain snowflakes.Fountain) (Resource, error) {
		if fountain == nil {
			return &UserResource{intUser: &models.User{}}, nil
		}
		id, err := fountain.NewID()
		return &UserResource{intUser: &models.User{Snowflake: id}}, err
	}); err != nil {
		panic(err)
	}
}

func RegisterUserEndpoint(db models.XODB) error {
	if err := RegisterResourceEndpoint("user", &UserEndpoint{db: db}); err != nil {
		return err
	}
	return nil
}

type UserResource struct {
	intUser *models.User
}

func (u *UserResource) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.intUser)
}

func (u *UserResource) Merge(r Resource) error {
	if u2, ok := r.(*UserResource); ok {
		if u2.intUser.Snowflake != 0 {
			u.intUser.Snowflake = u2.intUser.Snowflake
		}
		if len(u2.intUser.Avatar) != 0 {
			u.intUser.Avatar = u2.intUser.Avatar
		}
		if u2.intUser.CreatedAt != nil {
			u.intUser.CreatedAt = u2.intUser.CreatedAt
		}
		if u2.intUser.DeletedAt.Valid && u2.intUser.DeletedAt.Time.Unix() != 0 {
			u.intUser.DeletedAt = u2.intUser.DeletedAt
		}
		if u2.intUser.Username != "" {
			u.intUser.Username = u2.intUser.Username
		}
		if u2.intUser.Email.Valid && u2.intUser.Email.String != "" {
			u.intUser.Email = u2.intUser.Email
		}
		return nil
	}
	return errors.New("Merging with incompatible type")
}

func (u *UserResource) Snowflake() int64 {
	return u.intUser.Snowflake
}

func (u *UserResource) StripReadOnly() error {
	u.intUser.Snowflake = 0
	u.intUser.CreatedAt = nil
	return nil
}

func (u *UserResource) Type() string {
	return "user"
}

func (u *UserResource) UnmarshalJSON(data []byte) error {
	if u.intUser == nil {
		u.intUser = &models.User{}
	}
	return json.Unmarshal(data, u.intUser)
}

type UserEndpoint struct {
	db models.XODB
}

func (ue *UserEndpoint) Find(snowflake int64) (Resource, error) {
	intUser, err := models.UserBySnowflake(ue.db, snowflake)
	if err != nil {
		return nil, err
	}
	return &UserResource{intUser}, nil
}

func (ue *UserEndpoint) FindAll(pivot, size int64) ([]Resource, error) {
	intUsers, err := models.UsersByPivot(ue.db, pivot, size)
	if err != nil {
		return nil, err
	}
	var users []Resource
	for _, v := range intUsers {
		users = append(users, &UserResource{v})
	}
	return users, nil
}

func (ue *UserEndpoint) HardDelete(s int64) error {
	user, err := models.UserBySnowflake(ue.db, s)
	if err != nil {
		return err
	}
	return user.Delete(ue.db)
}

func (ue *UserEndpoint) SoftDelete(s int64) error {
	user, err := models.UserBySnowflake(ue.db, s)
	if err != nil {
		return err
	}

	user.DeletedAt.Valid = true
	user.DeletedAt.Time = time.Now().UTC()

	return user.Save(ue.db)
}

func (ue *UserEndpoint) New(r Resource) error {
	if u, ok := r.(*UserResource); ok {
		if u != nil && u.intUser != nil {
			u.intUser.Insert(ue.db)
		}
		return errors.New("Resource was nil")
	}
	return errors.New("Wrong resource endpoint or incompatible resource")
}

func (ue *UserEndpoint) Name() string {
	return "user"
}
