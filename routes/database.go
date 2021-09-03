package routes

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/empty-interface/goadmin/dbms"
	"github.com/empty-interface/goadmin/session"
)

func HandleDatabase(w http.ResponseWriter, r *http.Request, currentSession *session.Session) (int, error) {
	// we get session
	tmpl := template.Must(template.ParseFiles("html/database.html"))
	tables := currentSession.Conn.GetTables()
	page := databasePage{
		Title:                 fmt.Sprintf("GoAdmin - %s", currentSession.DBname),
		DisconnectPath:        fmt.Sprintf("%s?a=%v", DisconnectPath, time.Now().UnixMilli()),
		Tables:                tables,
		Driver:                currentSession.Driver,
		DBname:                currentSession.DBname,
		Username:              currentSession.Username,
		Password:              currentSession.Password,
		Uuid:                  currentSession.Uuid,
		SaveConnectionLocally: currentSession.SavedLocally,
		ItemName:              session.GetGlobalSessionManager().ItemName,
		Junk:                  time.Now().UnixMilli(),
	}
	buffer := bytes.NewBufferString("")
	err := tmpl.Execute(buffer, page)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	buffer.WriteTo(w)
	return -1, nil
}
func getCurrentSession(r *http.Request) *session.Session {
	sessionManager := session.GetGlobalSessionManager()
	cookie, err := r.Cookie(sessionManager.Name)
	if err != nil {
		return nil
	}
	return sessionManager.Get(cookie.Value)
}

type databasePage struct {
	Title                 string
	DisconnectPath        string
	Tables                []dbms.Table
	Driver                string
	Username              string
	DBname                string
	Password              string
	Uuid                  string
	SaveConnectionLocally bool
	ItemName              string
	Junk                  int64
}
