package models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogin_Insert(t *testing.T) {
	assert := assert.New(t)

	testingCat := makeLogin()
	assert.NoError(testingCat.InsertG())
}

func TestLogin_Delete(t *testing.T) {
	assert := assert.New(t)

	testingCat := makeLogin()
	assert.NoError(testingCat.InsertG())

	assert.NoError(testingCat.DeleteG())
}

func TestLogin_Read(t *testing.T) {
	assert := assert.New(t)

	testingCat := makeLogin()

	assert.NoError(testingCat.InsertG())

	cat, err := FindLoginG(testingCat.Snowflake)
	assert.NoError(err)

	// Fix monotonic clock in compare *and* also fix postgre rounding my time
	assert.EqualValues(testingCat.CreatedAt.Unix(), cat.CreatedAt.Unix())
	testingCat.CreatedAt = cat.CreatedAt

	assert.EqualValues(cat, testingCat)
}

func TestLogin_Update(t *testing.T) {
	assert := assert.New(t)

	testLogin := makeLogin()

	assert.NoError(testLogin.InsertG())

	testLogin.Data = []byte("New data")

	assert.NoError(testLogin.UpdateG())
}

func makeLogin() *Login {
	flake, err := snow.NewID()
	if err != nil {
		panic(err)
	}
	user := makeUser()
	user.InsertGP()

	return &Login{
		Snowflake:  int64(flake),
		Data:       []byte("Some login data"),
		Identifier: fmt.Sprintf("login-%d", flake),
		Type:       1,
		UserID:     user.Snowflake,
	}
}
