package routes

import (
	"time"

	"github.com/google/uuid"
)

var sessionManager SessionManager
var maxTimeToLife time.Duration = time.Hour

type Session struct {
	uuid                               string
	driver, username, password, dbname string
	createdAt                          time.Time
}

func (session *Session) expired() bool {
	return time.Now().Sub(session.createdAt) < maxTimeToLife
}
func newSession(driver, username, password, dbname string) (*Session, error) {
	_uuid := ""
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	_uuid = id.String()
	return &Session{
		uuid:   _uuid,
		driver: driver, username: username, password: password, dbname: dbname,
		createdAt: time.Now(),
	}, nil
}

type SessionManager struct {
	aliveSessions map[string]*Session
}

func GetSessionManager() *SessionManager {
	return &sessionManager
}
func (manager *SessionManager) get(id string) *Session {
	sess, ok := manager.aliveSessions[id]
	if !ok {
		return nil
	}
	return sess
}
func (manager *SessionManager) set(id string, sess Session) {
	//we should maybe return an error if session is already there
	manager.aliveSessions[id] = &sess
}
func (manager *SessionManager) delete(id string) {
	delete(manager.aliveSessions, id)
}
