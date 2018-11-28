package longpoll

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

type LongPoll struct {
	mx    sync.Mutex
	chmap map[string]*Channel
	alive int32
	// performance optimisation: channel list cache between updates to avoid reconstructing it
	// from chmap values and unlocking the thread ASAP. Reset to nil on any alterations to chmap
	chcache []*Channel
}

func (lp *LongPoll) Subscribe(timeout time.Duration, topics ...string) (string, error) {
	if !lp.IsAlive() {
		return "", errors.New("pubsub is down")
	}
}

func (lp *LongPoll) IsAlive() bool {
	return atomic.LoadInt32(&lp.alive) == yes
}
