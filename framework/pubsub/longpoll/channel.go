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

func (ch *Channel) Publish(data interface{}, topic string) error {
	if !ch.IsAlive() {
		return errors.New("subscription channel is down")
	}
	if _, ok := ch.topics[topic]; !ok {
		return nil
	}
	go func() {
		ch.mx.Lock()
		defer ch.mx.Unlock()

		if ch.IsAlive() {
			ch.data = append(ch.data, data)
			// if ch.notif != nil && !ch.notif.pinged {
			// 	ch.notif.pinged = true
			// 	ch.notif.ping <- true
			// }
		}
	}()

	defer runtime.Gosched()
	return nil
}

func (ch *Channel) Get(polltime time.Duration) (chan []interface{}, error) {
	if !ch.IsAlive() {
		return nil, errors.New("subscription channel is down")
	}
	if polltime <= 0 {
		return nil, errors.New("positive polltime value expected")
	}
	resp := make(chan []interface{}, 1)
	go func() {
		ch.tor.Ping()
		ch.mx.Lock()
		defer ch.mx.Unlock()
		if !ch.IsAlive() {
			resp <- nil
			ch.mx.Unlock()
			return
		}

		if ch.onDataWaiting(resp) {
			ch.mx.Unlock()
			return
		}

		notif := &getnotifier{ping: make(chan bool, 1), pinged: false}
		ch.notif = notif
		ch.mx.Unlock()

		gotdata := no
		pollend := make(chan bool, 1)
		go ch.startLongpollTimer(polltime, pollend, &gotdata)

	}()
}

func (ch *Channel) startLongpollTimer(polltime time.Duration, pollend chan bool, gotdata *int32) {
	
}


func (ch *Channel) onDataWaiting(resp chan []interface{}) bool {
	if len(ch.data) > 0 {
		resp <- ch.data
		ch.data = nil
		ch.notif = nil
		return true
	}
	return false
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

func (ch *Channel) IsAlive() bool {
	return atomic.LoadInt32(&ch.alive) == yes
}