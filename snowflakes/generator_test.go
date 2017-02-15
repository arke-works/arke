package snowflakes

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerator_NewID(t *testing.T) {
	assert := assert.New(t)
	generator := Generator{
		StartTime:  time.Date(1998, time.November, 19, 0, 0, 0, 0, time.UTC).Unix(),
		InstanceID: 18,
	}

	id, err := generator.NewID()
	assert.NoError(err)
	assert.True(id > 0)

	lastID := id
	for i := 0; i < 30000; i++ {
		id, err := generator.NewID()
		assert.NoError(err)
		assert.True(id > lastID)
		lastID = id
	}

	generator.InstanceID = -2

	_, err = generator.NewID()
	assert.Equal(err, errBadInstance)

	generator.StartTime = time.Now().Unix() + 10000

	_, err = generator.NewID()
	assert.Equal(err, errNoFuture)
}

func BenchmarkGenerator_NewID(b *testing.B) {
	generator := Generator{
		StartTime:  time.Date(1998, time.November, 19, 0, 0, 0, 0, time.UTC).Unix(),
		InstanceID: 18,
	}
	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		generator.NewID()
	}
}

func TestEncodedToID(t *testing.T) {
	assert := assert.New(t)

	id, err := EncodedToID("JVzh")
	assert.NoError(err)
	assert.EqualValues(3414442, id)
}

func TestIDToEncoded(t *testing.T) {
	assert := assert.New(t)

	assert.EqualValues("JVzh", IDToEncoded(3414442))
}
