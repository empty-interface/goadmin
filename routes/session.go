package routes

import (
	"time"

	"github.com/google/uuid"
	"github.com/kalitheniks/goadmin/dbms"
)

var _sessionManager *sessionManager

const sessionTimeToLive time.Duration = time.Minute * 10

type Session struct {
	uuid                               string
	Driver, Username, Password, DBname string
	createdAt                          time.Time
	Conn                               *dbms.GormConnection
}

func (sess *Session) expired() bool {
	return time.Now().Sub(sess.createdAt) > sessionTimeToLive
}

func (sess *Session) alive() bool {
	return !sess.expired()
}
func (sess *Session) expiresAt() time.Time {
	return sess.createdAt.Add(sessionTimeToLive)
}
func (sess *Session) refresh() {
	sess.createdAt = time.Now()
}
func NewSession(driver, username, password, dbname string) (*Session, error) {
	_uuid := ""
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	_uuid = id.String()
	return &Session{
		uuid:   _uuid,
		Driver: driver, Username: username, Password: password, DBname: dbname,
		createdAt: time.Now(),
	}, nil
}

type sessionManager struct {
	Name          string
	aliveSessions map[string]*Session
}

func NewSessionManager(name string) {
	if name == "" {
		name = "goadminv1"
	}
	_sessionManager = &sessionManager{
		Name:          name,
		aliveSessions: make(map[string]*Session),
	}
}
func GetGlobalSessionManager() *sessionManager {
	if _sessionManager == nil {
		NewSessionManager("")
	}
	return _sessionManager
}
func (manager *sessionManager) get(id string) *Session {
	sess, ok := manager.aliveSessions[id]
	if !ok {
		return nil
	}
	return sess
}
func (manager *sessionManager) set(sess *Session) {
	//we should maybe return an error if session is already there
	manager.aliveSessions[sess.uuid] = sess
}
func (manager *sessionManager) delete(id string) {
	delete(manager.aliveSessions, id)
}
