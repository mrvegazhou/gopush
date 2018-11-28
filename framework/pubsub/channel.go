package pubsub

import (
	"errors"
	"github.com/teris-io/shortid"
	"sync/atomic"
	"time"
)

type Channel struct {
	mx      sync.Mutex
	id      string
	onClose func(id string)
	topics  map[string]bool
	data    []interface{}
	alive   int32
	notif   *getnotifier
	tor     *Timeout
}

type getnotifier struct {
	ping   chan bool
	pinged bool
}

func NewChannel(timeout time.Duration, onClose func(id string), topics ...string) (*Channel, error) {
	if len(topics) == 0 {
		return nil, errors.New("at least one topic expected")
	}
	id, err := shortid.Generate()
	if err != nil {
		return nil, err
	}
	ch := Channel{
		id:      id,
		onClose: onClose,
		topics:  make(map[string]bool),
		alive:   yes,
	}
	for _, topic := range topics {
		ch.topics[topic] = true
	}
	if tor, err := NewTimeout(timeout, ch.Drop); err == nil {
		ch.tor = tor
	} else {
		return nil, err
	}

}

func (ch *Channel) Drop() {
	if !ch.IsAlive() {
		return
	}
	atomic.StoreInt32(&ch.alive, no)
	go func() {
		ch.mx.Lock()
		defer ch.mx.Unlock()
		ch.tor.Drop()
		ch.data = nil
		// if ch.notif != nil && !ch.notif.pinged {
		// 	ch.notif.ping <- true
		// }
		ch.notif = nil
		if ch.onClose != nil {
			ch.onClose(ch.id)
		}
	}
}
