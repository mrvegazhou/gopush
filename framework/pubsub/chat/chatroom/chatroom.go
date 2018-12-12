package chatroom

import (
	"container/list"
	"time"
)

type Event struct {
	Type      string // "join", "leave", or "message"
	User      string
	Timestamp int    // Unix timestamp (secs)
	Text      string // What the user said (if Type == "message")
}

func newEvent(typ, user, msg string) Event {
	return Event{typ, user, int(time.Now().Unix()), msg}
}

type Subscription struct {
	Archive []Event      // All the events from the archive.
	New     <-chan Event // New events coming in.
}

// Owner of a subscription must cancel it when they stop listening to events.
func (s Subscription) Cancel() {
	unsubscribe <- s.New // Unsubscribe the channel.
	drain(s.New)         // Drain it, just in case there was a pending publish.
}

// Drains a given channel of any messages.
func drain(ch <-chan Event) {
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		default:
			return
		}
	}
}

const archiveSize = 10

var (
	// 当有新用户加入时，初始化的一些订阅信息
	subscribe = make(chan (chan<- Subscription), 10) //只能用来接收
	// 在信道channel中发送退订
	unsubscribe = make(chan (<-chan Event), 10) //发送消息
	// 在这里发送事件
	publish = make(chan Event, 10)
)

func Subscribe() Subscription {
	resp := make(chan Subscription)
	subscribe <- resp
	return <-resp
}

func Join(user string) {
	publish <- newEvent("join", user, "")
}

func Say(user, message string) {
	publish <- newEvent("message", user, message)
}

func Leave(user string) {
	publish <- newEvent("leave", user, "")
}

// This function loops forever, handling the chat room pubsub
func chatroom() {
	// 最近的几条聊天记录
	archive := list.New()
	// 订阅者列表
	subscribers := list.New()
	for {
		select {
		case ch := <-subscribe:
			var events []Event
			// 获取list的第一个元素
			for e := archive.Front(); e != nil; e = e.Next() {
				events = append(events, e.Value.(Event))
			}
			// 在list l的末尾插入值为v的元素，并返回该元素

		// 发送 Event 给所有订阅者，并且增加到聊天记录中
		case event := <-publish:
			for ch := subscribers.Front(); ch != nil; ch = ch.Next() {
				ch.Value.(chan Event) <- event
			}
			if archive.Len() >= archiveSize {
				archive.Remove(archive.Front())
			}
			archive.PushBack(event)

		// 订阅者离开后，从用户列表中删除离开者
		case unsub := <-unsubscribe:
			for ch := subscribers.Front(); ch != nil; ch = ch.Next() {
				if ch.Value.(chan Event) == unsub {
					subscribers.Remove(ch)
					break
				}
			}
		}
	}
}

func init() {
	go chatroom()
}
