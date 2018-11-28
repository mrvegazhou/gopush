package pubsub

type Timeout struct {
	lastping  int64
	alive     int32
	report    chan bool
	onTimeout func()
}

func NewTimeout(timeout time.Duration, onTimeout func()) (*Timeout, error) {
	if timeout <= 0 {
		return nil, errors.New("positive timeout value expected")
	}
	tor := &Timeout{
		alive:     yes,
		report:    make(chan bool, 1),
		onTimeout: onTimeout,
	}

	tor.Ping()
	go tor.handle(int64(timeout))
	return tor, nil
}

func (tor *Timeout) Ping() {
	if tor.IsAlive() {
		atomic.StoreInt64(&tor.lastping, tor.now())
	}
}

func (tor *Timeout) elapsed() int64 {
	return tor.now() - atomic.LoadInt64(&tor.lastping)
}

func (tor *Timeout) now() int64 {
	return time.Now().UnixNano()
}

func (tor *Timeout) handle(timeout int64) {
	hundredth := timeout / 100
	for tor.elapsed() < timeout && tor.IsAlive() {
		time.Sleep(time.Duration(hundredth))
	}
	if tor.IsAlive() {
		atomic.StoreInt32(&tor.alive, no)
		if tor.onTimeout != nil {
			go tor.onTimeout()
		}
	}
	// tor.report <- true
}

func (tor *Timeout) Drop() {
	atomic.StoreInt32(&tor.alive, no)
}
