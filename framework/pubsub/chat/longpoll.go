package longpolling

import (
	"chatroom/chatroom"
)

func JoinChat(uname string) {
	// Join chat room.
	Join(uname, nil)
	return
}

func PostMsg(uname, content string) {
	if len(uname) == 0 || len(content) == 0 {
		return
	}

	publish <- newEvent(chatroom.EVENT_MESSAGE, uname, content)
}

func FetchMsgs(lastReceived int) []chatroom.Event {

	events := chatroom.GetEvents(lastReceived)
	if len(events) > 0 {
		return events
	}

	// Wait for new message(s).
	ch := make(chan bool)
	waitingList.PushBack(ch)
	<-ch

	return models.GetEvents(lastReceived)
}
