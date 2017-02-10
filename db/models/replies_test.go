package models

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/nullbio/null.v6"
	"testing"
)

func TestReply_Insert(t *testing.T) {
	assert := assert.New(t)

	test := makeReply()
	assert.NoError(test.InsertG())
}

func TestReply_Delete(t *testing.T) {
	assert := assert.New(t)

	test := makeReply()
	assert.NoError(test.InsertG())

	assert.NoError(test.DeleteG())
}

func TestReply_Read(t *testing.T) {
	assert := assert.New(t)

	test := makeReply()

	assert.NoError(test.InsertG())

	cat, err := FindReplyG(test.Snowflake)
	assert.NoError(err)

	// Fix monotonic clock in compare *and* also fix postgre rounding my time
	assert.EqualValues(test.CreatedAt.Unix(), cat.CreatedAt.Unix())
	test.CreatedAt = cat.CreatedAt

	assert.EqualValues(cat, test)
}

func TestReply_Update(t *testing.T) {
	assert := assert.New(t)

	test := makeReply()

	assert.NoError(test.InsertG())

	test.Body = "A new Description"

	assert.NoError(test.UpdateG())
}

func makeReply() *Reply {
	flake, err := snow.NewID()
	if err != nil {
		panic(err)
	}
	author := makeUser()
	author.InsertGP()
	topic := makeTopic()
	topic.InsertGP()
	return &Reply{
		Snowflake: int64(flake),
		AuthorID:  null.Int64From(author.Snowflake),
		Body:      "Some body **test**",
		TopicID:   topic.Snowflake,
	}
}
