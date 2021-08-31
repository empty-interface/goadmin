package routes

import (
	"html/template"
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("home.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	data := Page{
		Title: "GoAdminer v1",
		Data: map[string]string{
			"postgres": "PostgreSQL",
			"mysql":    "mySQL",
		},
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

// func HandleConnect(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	r.
// }

type Page struct {
	Title string
	Data  interface{}
}
