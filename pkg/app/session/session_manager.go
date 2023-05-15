package session

import (
	"sync"
)

type SessionF interface {
	ID() int
	Start()
	Close()
	OnEvent(fromID int, data []byte) error
}

var allSession map[SessionF]struct{}
var allSessionMu sync.RWMutex

func init() {
	allSession = make(map[SessionF]struct{})
}

func trackSession(s SessionF, add bool) {
	allSessionMu.Lock()
	defer allSessionMu.Unlock()

	if add {
		allSession[s] = struct{}{}
	} else {
		delete(allSession, s)
	}
}

func ForEachSession(fn func(SessionF)) {
	allSessionMu.RLock()
	for s := range allSession {
		fn(s)
	}
	allSessionMu.RUnlock()
}

func CountSession() int {
	allSessionMu.RLock()
	defer allSessionMu.RUnlock()
	return len(allSession)
}

func DeleteSession(s SessionF) error {
	trackSession(s, false)
	return nil
}
