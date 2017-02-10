package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrivateMessage_Insert(t *testing.T) {
	assert := assert.New(t)

	test := makePrivateMessage()
	assert.NoError(test.InsertG())
}

func TestPrivateMessage_Delete(t *testing.T) {
	assert := assert.New(t)

	test := makePrivateMessage()
	assert.NoError(test.InsertG())

	assert.NoError(test.DeleteG())
}

func TestPrivateMessage_Read(t *testing.T) {
	assert := assert.New(t)

	test := makePrivateMessage()

	assert.NoError(test.InsertG())

	cat, err := FindPrivateMessageG(test.Snowflake)
	assert.NoError(err)

	// Fix monotonic clock in compare *and* also fix postgre rounding my time
	assert.EqualValues(test.CreatedAt.Unix(), cat.CreatedAt.Unix())
	test.CreatedAt = cat.CreatedAt

	assert.EqualValues(cat, test)
}

func TestPrivateMessage_Update(t *testing.T) {
	assert := assert.New(t)

	test := makePrivateMessage()

	assert.NoError(test.InsertG())

	test.Body = "A new Description"

	assert.NoError(test.UpdateG())
}

func makePrivateMessage() *PrivateMessage {
	flake, err := snow.NewID()
	if err != nil {
		panic(err)
	}
	sender := makeUser()
	sender.InsertGP()
	recv := makeUser()
	recv.InsertGP()
	return &PrivateMessage{
		Snowflake:  int64(flake),
		Body:       "Some text",
		SenderID:   sender.Snowflake,
		ReceiverID: recv.Snowflake,
		Title:      "Some title",
	}
}
