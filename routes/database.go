package routes

import (
	"fmt"
	"html/template"
	"net/http"
)

func HandleDatabase(w http.ResponseWriter, r *http.Request, currentSession *Session) (int, error) {
	// we get session
	tmpl := template.Must(template.ParseFiles("html/database.html"))
	currentSession.Conn.GetTables()
	page := struct {
		Title                 string
		DisconnectPath        string
		Tables                map[string]string
		Driver                string
		Username              string
		DBname                string
		Password              string
		Uuid                  string
		SaveConnectionLocally bool
	}{
		Title:                 fmt.Sprintf("GoAdmin - %s", currentSession.DBname),
		DisconnectPath:        DisconnectPath,
		Tables:                make(map[string]string),
		Driver:                currentSession.Driver,
		DBname:                currentSession.DBname,
		Username:              currentSession.Username,
		Password:              currentSession.Password,
		Uuid:                  currentSession.Uuid,
		SaveConnectionLocally: currentSession.SavedLocally,
	}
	err := tmpl.Execute(w, page)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return -1, nil
}
