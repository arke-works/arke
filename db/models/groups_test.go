package models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/nullbio/null.v6"
	"testing"
)

func TestGroup_Insert(t *testing.T) {
	assert := assert.New(t)

	testingGroup := makeGroup()

	assert.NoError(testingGroup.InsertG())
}

func TestGroup_Delete(t *testing.T) {
	assert := assert.New(t)

	testingGroup := makeGroup()

	assert.NoError(testingGroup.InsertG())

	assert.NoError(testingGroup.DeleteG())
}

func TestGroup_Read(t *testing.T) {
	assert := assert.New(t)

	testingGroup := makeGroup()

	assert.NoError(testingGroup.InsertG())

	grp, err := FindGroupG(testingGroup.Snowflake)
	assert.NoError(err)

	// Fix monotonic clock in compare *and* also fix postgre rounding my time
	assert.EqualValues(testingGroup.CreatedAt.Unix(), grp.CreatedAt.Unix())
	grp.CreatedAt = testingGroup.CreatedAt

	assert.EqualValues(grp, testingGroup)
}

func TestGroup_Update(t *testing.T) {
	assert := assert.New(t)

	flake, err := snow.NewID()
	assert.NoError(err)
	testingGroup := makeGroup()

	assert.NoError(testingGroup.InsertG())

	testingGroup.Name = fmt.Sprintf("A new name - %d", flake)

	assert.NoError(testingGroup.UpdateG())
}

func TestGroup_Parent(t *testing.T) {
	assert := assert.New(t)

	group1 := makeGroup()
	assert.NoError(group1.InsertG())

	group2 := makeGroup()
	group2.ParentID = null.Int64From(group1.Snowflake)
	assert.NoError(group2.InsertG())

	parent := group2.ParentG().OneP()

	assert.EqualValues(group1.CreatedAt.Unix(), parent.CreatedAt.Unix())

	parent.CreatedAt = group1.CreatedAt

	assert.EqualValues(group1, parent)
}

func makeGroup() *Group {
	flake, err := snow.NewID()
	if err != nil {
		panic(err)
	}
	return &Group{
		Snowflake:  int64(flake),
		Permission: null.BytesFrom([]byte{0x00, 0xFF}),
		Name:       fmt.Sprintf("%d-test", flake),
		ParentID:   null.Int64From(int64(flake)),
	}
}
