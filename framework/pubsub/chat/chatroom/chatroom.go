package chatroom

import (
	"container/list"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

/////////////////////////////////begin archive/////////////////////////////////////////
type EventType int

const (
	EVENT_JOIN = iota
	EVENT_LEAVE
	EVENT_MESSAGE
)
const archiveSize = 20

type Event struct {
	Type      EventType // JOIN, LEAVE, MESSAGE
	User      string
	Timestamp int // Unix timestamp (secs)
	Content   string
}

var archive = list.New()

func newArchive(event Event) {
	if archive.Len() >= archiveSize {
		archive.Remove(archive.Front())
	}
	archive.PushBack(event)
}

func GetEvents(lastReceived int) []Event {
	events := make([]Event, 0, archive.Len())
	for event := archive.Front(); event != nil; event = event.Next() {
		e := event.Value.(Event)
		if e.Timestamp > int(lastReceived) {
			events = append(events, e)
		}
	}
	return events
}

///////////////////////////////end archive//////////////////////////////////////////////

type Subscriber struct {
	Name string
	Conn *websocket.Conn // Only for WebSocket users; otherwise nil.
}

var (
	// Channel for new join users.
	subscribe = make(chan Subscriber, 10)
	// Channel for exit users.
	unsubscribe = make(chan string, 10)
	// Send events here to publish them.
	publish = make(chan models.Event, 10)
	// Long polling waiting list.
	waitingList = list.New()
	subscribers = list.New()
)

type Subscription struct {
	Archive []models.Event      // All the events from the archive.
	New     <-chan models.Event // New events coming in.
}

func newEvent(ep models.EventType, user, msg string) models.Event {
	return models.Event{ep, user, int(time.Now().Unix()), msg}
}

func Join(user string, ws *websocket.Conn) {
	subscribe <- Subscriber{Name: user, Conn: ws}
}

func isUserExist(subscribers *list.List, user string) bool {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).Name == user {
			return true
		}
	}
	return false
}

func chatroom() {
	for {
		select {
		case sub := <-subscribe:
			if !isUserExist(subscribers, sub.Name) {
				subscribers.PushBack(sub) // Add user to the end of list.
				// Publish a JOIN event.
				publish <- newEvent(EVENT_JOIN, sub.Name, "")
				log.Println("New user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			} else {
				log.Println("Old user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			}
		case event := <-publish:
			// Notify waiting list.
			for ch := waitingList.Back(); ch != nil; ch = ch.Prev() {
				ch.Value.(chan bool) <- true
				waitingList.Remove(ch)
			}

			broadcastWebSocket(event)
			newArchive(event)

			if event.Type == EVENT_MESSAGE {
				log.Println("Message from", event.User, ";Content:", event.Content)
			}
		case unsub := <-unsubscribe:
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).Name == unsub {
					subscribers.Remove(sub)
					// Clone connection.
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
						log.Println("WebSocket closed:", unsub)
					}
					publish <- newEvent(EVENT_LEAVE, unsub, "") // Publish a LEAVE event.
					break
				}
			}
		}
	}
}
