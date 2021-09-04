package routes

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/empty-interface/goadmin/dbms"
	"github.com/empty-interface/goadmin/session"
)

func HandleDatabase(w http.ResponseWriter, r *http.Request, currentSession *session.Session) (int, error) {
	// we get session
	tmpl := template.Must(template.ParseFiles("html/database.html"))
	schemas := currentSession.Conn.GetDBSchemas()
	schema := getParam(r, "schema", "public")
	tables := currentSession.Conn.GetSchemaTables(schema)
	page := databasePage{
		Title:                 fmt.Sprintf("GoAdmin - %s", currentSession.DBname),
		DisconnectPath:        DisconnectPath,
		Tables:                tables,
		SaveConnectionLocally: currentSession.SavedLocally,
		ItemName:              session.GetGlobalSessionManager().ItemName,
		Schemas:               schemas,
		CurrentSchema:         schema,

		Driver:   currentSession.Driver,
		DBname:   currentSession.DBname,
		Username: currentSession.Username,
		Password: currentSession.Password,
		Uuid:     currentSession.Uuid,
		Host:     currentSession.Host,
		Port:     currentSession.Port,
	}
	buffer := bytes.NewBufferString("")
	err := tmpl.Execute(buffer, page)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	buffer.WriteTo(w)
	return -1, nil
}
func getParam(r *http.Request, name string, default_ string) string {
	params, exist := r.URL.Query()[name]
	if !exist || len(params) < 0 {
		return default_
	}
	return params[0]
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
	SaveConnectionLocally bool
	ItemName              string
	Schemas               []dbms.Schema
	CurrentSchema         string

	Driver   string
	Username string
	DBname   string
	Password string
	Uuid     string
	Host     string
	Port     string
}
