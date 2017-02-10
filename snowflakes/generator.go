package snowflakes // import "iris.arke.works/forum/snowflakes"

import (
	"errors"
	"github.com/osamingo/indigo/base58"
	"sync"
	"time"
)

const (
	counterLen  = 13
	instanceLen = 7
	counterMask = -1 ^ (-1 << counterLen)
)

var (
	errNoFuture    = errors.New("Start Time cannot be set in the future")
	errBadInstance = errors.New("Instance ID must be smaller than 129")
)

// Generator is a fountain for new snowflakes. StartTime must be
// initialized to a past point in time and Instance ID can be any
// positive value or 0.
//
// If any value is not correctly set, new IDs cannot be produced.
type Generator struct {
	StartTime  int64
	InstanceID int8
	mutex      *sync.Mutex
	sequence   int32
	now        int64
}

// NewID generates a new, unique snowflake value
//
// Up to 8192 snowflakes per second can be requested
// If exhausted, it blocks and sleeps until a new second
// of unix time starts.
//
// The return value is signed but always positive.
//
// Additionally, the return value is monotonic for a single
// instance and weakly monotonic for many instances.
func (g *Generator) NewID() (int64, error) {
	if g.mutex == nil {
		g.mutex = new(sync.Mutex)
	}
	if g.StartTime > time.Now().Unix() {
		return 0, errNoFuture
	}
	if g.InstanceID < 0 {
		return 0, errBadInstance
	}
	g.mutex.Lock()
	defer g.mutex.Unlock()

	var (
		now   int64
		flake int64
	)
	now = int64(time.Now().Unix())

	if now == g.now {
		g.sequence = (g.sequence + 1) & counterMask
		if g.sequence == 0 {
			for now <= g.now {
				now = int64(time.Now().Unix())
				time.Sleep(time.Microsecond * 100)
			}
		}
	} else {
		g.sequence = 0
	}

	g.now = now

	flake = int64(
		((now - g.StartTime) << (instanceLen + counterLen)) |
			(int64(g.sequence) << (instanceLen)) |
			(int64(g.InstanceID)))

	return flake, nil
}

// IDToEncoded encodes an incoming ID to a Base58 string
func IDToEncoded(id int64) string {
	return base58.StdEncoding.Encode(uint64(id))
}

// EncodedToID will attempt to encode a string using Base58 and return
// the ID
func EncodedToID(idStr string) (int64, error) {
	id, err := base58.StdEncoding.Decode(idStr)
	return int64(id), err
}
