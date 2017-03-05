package snowflakes

import (
	"sync"
)

type SequenceGenerator struct {
	Position int64
	mutex    *sync.Mutex
}

func (sg *SequenceGenerator) NewID() (int64, error) {
	if sg.mutex == nil {
		sg.mutex = new(sync.Mutex)
	}
	sg.mutex.Lock()
	defer sg.mutex.Unlock()

	sg.Position++

	return sg.Position, nil
}

var _ Fountain = (*SequenceGenerator)(nil)
