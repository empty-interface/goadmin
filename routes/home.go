package routes

import (
	"html/template"
	"net/http"
)

const HomePath = "/"

func HandleHome(w http.ResponseWriter, r *http.Request) (int, error) {
	tmpl, err := template.ParseFiles("html/home.html")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	sessionManager := GetGlobalSessionManager()
	page := homePage{
		Title: "GoAdminer v1",
		Select: map[string]string{
			"postgres": "PostgreSQL",
			// "mysql":    "MySQL",
		},
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
