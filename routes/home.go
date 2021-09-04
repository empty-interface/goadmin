package routes

import (
	"html/template"
	"net/http"

	"github.com/empty-interface/goadmin/dbms"
	"github.com/empty-interface/goadmin/session"
)

const HomePath = "/"

func HandleHome(w http.ResponseWriter, r *http.Request) (int, error) {
	tmpl, err := template.ParseFiles("html/home.html")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	sessionManager := session.GetGlobalSessionManager()
	supportedDrivers := dbms.GetSupportedDrivers()
	page := homePage{
		Title:    "GoAdminer v1",
		Select:   supportedDrivers,
		ItemName: sessionManager.ItemName,
		Action:   ConnectPath,
	}
	w.Header().Set("Content-type", "text/html")
	err = tmpl.Execute(w, page)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return -1, nil
}

type homePage struct {
	Title    string
	Action   string
	ItemName string
	Select   interface{}
}
