package session

import (
	"encoding/json"
	"fmt"
	"sync"

	"nostr-relay/pkg/models"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var id = 0

// .. //
type Session struct {
	id    int
	conn  *websocket.Conn
	mutex sync.Mutex
	subID *string
}

func NewSession(conn *websocket.Conn) *Session {
	id++
	return &Session{
		id:   id,
		conn: conn,
	}
}

func (t *Session) OnEvent(fromID int, data []byte) error {

	if t.ID() != fromID {
		if !t.isReq() {
			return nil
		}
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.conn.WriteMessage(websocket.TextMessage, data)
}

// func (t *Session) WriteMessage(messageType int, data []byte) error {
// 	t.mutex.Lock()
// 	defer t.mutex.Unlock()
// 	return t.conn.WriteMessage(messageType, data)
// }

func (t *Session) WriteMessage(data []byte) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.conn.WriteMessage(websocket.TextMessage, data)
}

func (t *Session) WriteJson(v interface{}) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.conn.WriteJSON(v)
}

func (t *Session) ID() int {
	return t.id
}
func (t *Session) Start() {

	trackSession(t, true)
	for {
		_, data, err := t.conn.ReadMessage()
		if err != nil {
			log.Infof(" %s | read err: %v", t.basicInfo(), err)
			break
		}
		log.Infof("ReadMessage %s", string(data))
		//log.Infof("ReadMessage %v", data)

		err = t.msgHandle(data)
		if err != nil {
			log.Infof(" %s | msgHandle err: %v", t.basicInfo(), err)
			break
		}

	}

	log.Infof(" %s | closed", t.basicInfo())
}

func (t *Session) Close() {
	t.conn.Close()
}

func (t *Session) basicInfo() string {
	return fmt.Sprintf("%15d", id)
}

func (t *Session) msgHandle(message []byte) error {

	// Parse the message as a JSON array
	var msg []interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		e := fmt.Errorf("Session msgHandle: json unmarshal error:%s", err.Error())
		return e
	}
	// Handle each message type
	switch msg[0] {
	case "EVENT":
		// Parse the event JSON
		var event models.Msg
		jsonData, _ := json.Marshal(msg[1])
		if err := json.Unmarshal(jsonData, &event); err != nil {
			e := fmt.Errorf("Session msgHandle: json unmarshal error:%s", err.Error())
			return e
		}
		// Print the event data
		fmt.Printf("Received event: %+v\n", event)

		ForEachSession(func(s SessionF) {
			s.OnEvent(t.ID(), message)
		})

	case "REQ":
		// Subscription has been closed
		fmt.Printf("Subscription %s req\n", msg[1])

		tmp, ok := msg[1].(string)
		if ok {
			t.setReq(&tmp)
		}

	case "CLOSE":
		// Subscription has been closed
		fmt.Printf("Subscription %s closed\n", msg[1])
		t.setReq(nil)
	default:
		log.Printf("Unknown message type: %s\n", msg[0])
	}

	return nil
}

func (t *Session) setReq(subID *string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.subID = subID
}

func (t *Session) isReq() bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.subID != nil
}
