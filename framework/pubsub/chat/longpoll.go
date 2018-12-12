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
	if err == nil {
		lp.mx.Lock()
		lp.chcache = nil
		lp.chmap[ch.id] = ch
		lp.mx.Unlock()
		return ch.id, nil
	}
	return "", err
}

func (lp *LongPoll) Get(id string, polltime time.Duration) (chan []interface{}, error) {
	if !lp.IsAlive() {
		return nil, errors.New("pubsub is down")
	}
	if ch, ok := lp.Channel(id); ok {
		return ch.Get(polltime)
	}
	return nil, fmt.Errorf("no channel for Id %v", id)
}

func (lp *LongPoll) Publish(data interface{}, topics ...string) error {
	if !lp.IsAlive() {
		return errors.New("pubsub is down")
	}
	if len(topics) == 0 {
		return errors.New("expected at least one topic")
	}
	for _, ch := range lp.Channels() {
		for _, topic := range topics {
			ch.Publish(data, topic) // errors ignored
		}
	}
	return nil
}

func (lp *LongPoll) Channel(id string) (*Channel, bool) {
	if !lp.IsAlive() {
		return nil, false
	}
	lp.mx.Lock()
	defer lp.mx.Unlock()
	res, ok := lp.chmap[id]
	return res, ok && res.IsAlive()
}

func (lp *LongPoll) Channels() []*Channel {
	if !lp.IsAlive() {
		return nil
	}
	lp.mx.Lock()
	defer lp.mx.Unlock()
	if len(lp.chcache) == 0 {
		for _, ch := range lp.chmap {
			if ch.IsAlive() {
				lp.chcache = append(lp.chcache, ch)
			}
		}
	}
	return lp.chcache
}

func (lp *LongPoll) IsAlive() bool {
	return atomic.LoadInt32(&lp.alive) == yes
}

func (lp *LongPoll) drop(id string) {
	lp.mx.Lock()
	defer lp.mx.Unlock()
	lp.chcache = nil
	delete(lp.chmap, id)
}
