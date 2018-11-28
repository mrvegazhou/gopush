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
	ch, err := NewChannel(timeout, lp.drop, topics...)
}

func (lp *LongPoll) IsAlive() bool {
	return atomic.LoadInt32(&lp.alive) == yes
}

func (lp *LongPoll) drop(id string) {
	lp.mx.Lock()
	lp.chcache = nil
	delete(lp.chmap, id)
	lp.mx.Unlock()
}
