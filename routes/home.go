package routes

import (
	"html/template"
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request) (int, error) {
	tmpl, err := template.ParseFiles("html/home.html")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	page := struct {
		Title  string
		Action string
		Select interface{}
	}{
		Title: "GoAdminer v1",
		Select: map[string]string{
			"mysql":    "mySQL",
			"postgres": "PostgreSQL",
		},
		Action: ConnectPath,
	}
	w.Header().Set("Content-type", "text/html")
	err = tmpl.Execute(w, page)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return -1, nil
}
