package models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/nullbio/null.v6"
	"testing"
)

func TestCategory_Insert(t *testing.T) {
	assert := assert.New(t)

	test := makeCategory()
	assert.NoError(test.InsertG())
}

func TestCategory_Delete(t *testing.T) {
	assert := assert.New(t)

	test := makeCategory()
	assert.NoError(test.InsertG())

	assert.NoError(test.DeleteG())
}

func TestCategory_Read(t *testing.T) {
	assert := assert.New(t)

	test := makeCategory()

	assert.NoError(test.InsertG())

	cat, err := FindCategoryG(test.Snowflake)
	assert.NoError(err)

	// Fix monotonic clock in compare *and* also fix postgre rounding my time
	assert.EqualValues(test.CreatedAt.Unix(), cat.CreatedAt.Unix())
	test.CreatedAt = cat.CreatedAt

	assert.EqualValues(cat, test)
}

func TestCategory_Update(t *testing.T) {
	assert := assert.New(t)

	test := makeCategory()

	assert.NoError(test.InsertG())

	test.Description = null.StringFrom("A new Description")

	assert.NoError(test.UpdateG())
}

func makeCategory() *Category {
	flake, err := snow.NewID()
	if err != nil {
		panic(err)
	}
	return &Category{
		Snowflake:   int64(flake),
		Description: null.StringFrom("Testing category"),
		Color:       null.IntFrom(0xFFFFFF),
		Title:       fmt.Sprintf("%d-test", flake),
	}
}
