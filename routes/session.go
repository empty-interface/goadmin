package routes

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"time"

	"github.com/empty-interface/goadmin/dbms"
	"github.com/google/uuid"
)

var _sessionManager *sessionManager

const sessionTimeToLive time.Duration = time.Minute * 10

type Session struct {
	Uuid         string
	Driver       string
	Username     string
	Password     string
	DBname       string
	createdAt    time.Time
	Conn         *dbms.GormConnection
	SavedLocally bool
}
type infileSession struct {
	Uuid      string `json:"uuid"`
	Driver    string `json:"driver"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	DBname    string `json:"dbname"`
	CreatedAt int64  `json:"createdAt"`
}

func sessionToInfileSession(sess *Session) infileSession {
	return infileSession{
		Uuid:      sess.Uuid,
		Driver:    sess.Driver,
		Username:  sess.Username,
		Password:  sess.Password,
		DBname:    sess.DBname,
		CreatedAt: sess.createdAt.UnixMilli(),
	}
}

func infileSessiontoSession(sess *infileSession) *Session {
	return &Session{
		Uuid:      sess.Uuid,
		Driver:    sess.Driver,
		Username:  sess.Username,
		Password:  sess.Password,
		DBname:    sess.DBname,
		createdAt: time.UnixMilli(sess.CreatedAt),
		Conn:      nil,
	}
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
func NewSession(driver, username, password, dbname string, saved bool) (*Session, error) {
	_uuid := ""
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	_uuid = id.String()
	return &Session{
		Uuid:   _uuid,
		Driver: driver, Username: username, Password: password, DBname: dbname,
		createdAt:    time.Now(),
		SavedLocally: saved,
	}, nil
}

type sessionManager struct {
	ItemName       string
	outputFilename string
	Name           string
	aliveSessions  map[string]*Session
}

func NewSessionManager(name string) {
	if name == "" {
		name = "goadminv1"
	}
	_sessionManager = &sessionManager{
		ItemName:       "sessions",
		Name:           name,
		aliveSessions:  make(map[string]*Session),
		outputFilename: "./sessions.json",
	}
}
func GetGlobalSessionManager() *sessionManager {
	if _sessionManager == nil {
		NewSessionManager("")
		_sessionManager.loadSessionsFromFile()
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
func (manager *sessionManager) Close() {
	manager.saveSessionsToFile()
}
func (manager *sessionManager) loadSessionsFromFile() {
	b, err := ioutil.ReadFile(manager.outputFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file: %s\n", err.Error())
		return
	}
	sessions := make(map[string]infileSession)
	err = json.Unmarshal(b, &sessions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not unmarshal sessions: %s\n", err.Error())
		return
	}
	for _, sess := range sessions {
		_s := infileSessiontoSession(&sess)
		if _s.alive() {
			manager.aliveSessions[sess.Uuid] = _s
		}
	}
	fmt.Println("Loaded sessions")
}

func (manager *sessionManager) saveSessionsToFile() {
	sessions := make(map[string]infileSession)
	for _, sess := range manager.aliveSessions {
		sessions[sess.Uuid] = sessionToInfileSession(sess)
	}
	b, err := json.Marshal(sessions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not marshal sessions: %s", err.Error())
		return
	}
	err = ioutil.WriteFile(manager.outputFilename, b, fs.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not write to file: %s", err.Error())
		return
	}
	fmt.Println("Sessions saved")
}
func (manager *sessionManager) set(sess *Session) {
	//we should maybe return an error if session is already there
	manager.aliveSessions[sess.Uuid] = sess
}
func (manager *sessionManager) delete(id string) {
	delete(manager.aliveSessions, id)
}
