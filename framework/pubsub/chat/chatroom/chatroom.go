package chatroom

import (
	"container/list"
	"encoding/json"
	"fmt"
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
	Publish = make(chan Event, 10)
	// Long polling waiting list.
	WaitingList = list.New()
	subscribers = list.New()
	// 群消息
	subscribersIngroup = list.New()
	// P2P消息
	subscribersInFriends = list.New()
)

// type Subscription struct {
// 	Archive []Event      // All the events from the archive.
// 	New     <-chan Event // New events coming in.
// }

func NewEvent(ep EventType, user, msg string) Event {
	return Event{ep, user, int(time.Now().Unix()), msg}
}

func Join(user string, ws *websocket.Conn) {
	subscribe <- Subscriber{Name: user, Conn: ws}
}

func Leave(user string) {
	unsubscribe <- user
}

func isUserExist(subscribers *list.List, user string) bool {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).Name == user {
			return true
		}
	}
	return false
}

func broadcastWebSocket(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Println("Fail to marshal event:", err)
		return
	}
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				unsubscribe <- sub.Value.(Subscriber).Name
			}
		}
	}
}

func chatroom() {
	for {
		select {
		case sub := <-subscribe:
			if !isUserExist(subscribers, sub.Name) {
				subscribers.PushBack(sub) // Add user to the end of list.
				// Publish a JOIN event.
				Publish <- NewEvent(EVENT_JOIN, sub.Name, "")
				log.Println("New user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			} else {
				log.Println("Old user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			}
		case event := <-Publish:
			fmt.Print("event:", event)
			// Notify waiting list.
			for ch := WaitingList.Back(); ch != nil; ch = ch.Prev() {
				ch.Value.(chan bool) <- true
				WaitingList.Remove(ch)
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
					Publish <- NewEvent(EVENT_LEAVE, unsub, "") // Publish a LEAVE event.
					break
				}
			}
		}
	}
}

func init() {
	go chatroom()
}
