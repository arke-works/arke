package models

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/nullbio/null.v6"
	"testing"
)

func TestTopic_Insert(t *testing.T) {
	assert := assert.New(t)

	test := makeTopic()
	assert.NoError(test.InsertG())
}

func TestTopic_Delete(t *testing.T) {
	assert := assert.New(t)

	test := makeTopic()
	assert.NoError(test.InsertG())

	assert.NoError(test.DeleteG())
}

func TestTopic_Read(t *testing.T) {
	assert := assert.New(t)

	test := makeTopic()

	assert.NoError(test.InsertG())

	cat, err := FindTopicG(test.Snowflake)
	assert.NoError(err)

	// Fix monotonic clock in compare *and* also fix postgre rounding my time
	assert.EqualValues(test.CreatedAt.Unix(), cat.CreatedAt.Unix())
	test.CreatedAt = cat.CreatedAt

	assert.EqualValues(cat, test)
}

func TestTopic_Update(t *testing.T) {
	assert := assert.New(t)

	test := makeTopic()

	assert.NoError(test.InsertG())

	test.Body = "A new Description"

	assert.NoError(test.UpdateG())
}

func TestTopicCategory(t *testing.T) {
	assert := assert.New(t)

	topic := makeTopic()
	category := makeCategory()

	topCat := &RelTopicCategory{
		TopicID:    topic.Snowflake,
		CategoryID: category.Snowflake,
	}

	topic.InsertGP()
	topic.ReloadGP()
	category.InsertGP()
	category.ReloadGP()
	topCat.InsertGP()
	topCat.ReloadGP()

	retCat := topic.RelTopicCategoriesG().OneP()

	// Fix monotonic clock in compare *and* also fix postgre rounding my time
	assert.EqualValues(topCat.CreatedAt.Unix(), retCat.CreatedAt.Unix())
	topCat.CreatedAt = retCat.CreatedAt

	assert.EqualValues(topCat, retCat)

	topCat.DeleteGP()

	assert.EqualValues(0, category.RelTopicCategoriesG().CountP())
}

func makeTopic() *Topic {
	flake, err := snow.NewID()
	if err != nil {
		panic(err)
	}
	user := makeUser()
	user.InsertGP()
	return &Topic{
		Snowflake: int64(flake),
		AuthorID:  null.Int64From(user.Snowflake),
		Body:      "Some text *in markdown*",
		Revision:  0,
		Title:     "Test",
	}
}
