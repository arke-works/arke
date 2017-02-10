package models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/nullbio/null.v6"
	"testing"
)

func TestUser_Insert(t *testing.T) {
	assert := assert.New(t)

	test := makeUser()
	assert.NoError(test.InsertG())
}

func TestUser_Delete(t *testing.T) {
	assert := assert.New(t)

	test := makeUser()
	assert.NoError(test.InsertG())

	assert.NoError(test.DeleteG())
}

func TestUser_Read(t *testing.T) {
	assert := assert.New(t)

	test := makeUser()

	assert.NoError(test.InsertG())

	cat, err := FindUserG(test.Snowflake)
	assert.NoError(err)

	// Fix monotonic clock in compare *and* also fix postgre rounding my time
	assert.EqualValues(test.CreatedAt.Unix(), cat.CreatedAt.Unix())
	test.CreatedAt = cat.CreatedAt

	assert.EqualValues(cat, test)
}

func TestUser_Update(t *testing.T) {
	assert := assert.New(t)

	test := makeUser()

	assert.NoError(test.InsertG())

	test.Username = fmt.Sprintf("new-user-%d", test.Snowflake)

	assert.NoError(test.UpdateG())
}

func TestUser_Group(t *testing.T) {
	assert := assert.New(t)

	user := makeUser()
	group := makeGroup()

	assert.NoError(user.InsertG())
	assert.NoError(group.InsertG())

	rel := &RelUserGroup{
		GroupID: group.Snowflake,
		UserID:  user.Snowflake,
	}

	assert.NoError(rel.InsertG())

	rel.ReloadGP()

	relRet := user.RelUserGroupsG().OneP()

	// Fix monotonic clock in compare *and* also fix postgre rounding my time
	assert.EqualValues(rel.CreatedAt.Unix(), relRet.CreatedAt.Unix())
	relRet.CreatedAt = rel.CreatedAt

	assert.EqualValues(rel, relRet)

	relRet.DeleteGP()

	assert.EqualValues(0, user.RelUserGroupsG().CountP())
}

func makeUser() *User {
	flake, err := snow.NewID()
	if err != nil {
		panic(err)
	}
	return &User{
		Snowflake: int64(flake),
		Avatar:    []byte{0x00},
		Email:     null.StringFrom(fmt.Sprintf("%d@exampleorg", flake)),
		Username:  fmt.Sprintf("user-%d", flake),
	}
}
