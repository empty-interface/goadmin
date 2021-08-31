package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

const (
	HomePath    = "/"
	ConnectPath = "/connect"
	DBPath      = "/database"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("home.html")
	if err != nil {
		handleError(w, r, err, http.StatusInternalServerError)
		return
	}
	page := struct {
		Title  string
		Action string
		Select interface{}
	}{
		Title: "GoAdminer v1",
		Select: map[string]string{
			"postgres": "PostgreSQL",
			"mysql":    "mySQL",
		},
		Action: ConnectPath,
	}
	w.Header().Set("Content-type", "text/html")
	err = tmpl.Execute(w, page)
	if err != nil {
		handleError(w, r, err, http.StatusInternalServerError)
		return
	}
}
func HandleConnect(connect func(string, string, string, string) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		form, err := parseConnectForm(r.PostForm)
		if err != nil {
			handleError(w, r, err, http.StatusBadRequest)
			return
		}

		err = connect(form["driver"], form["username"], form["password"], form["dbname"])
		if err != nil {
			handleError(w, r, err, http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, DBPath, http.StatusMovedPermanently)
	})
}

func HandleDatabase(w http.ResponseWriter, r *http.Request) {
}
func handleError(w http.ResponseWriter, r *http.Request, _err error, status int) {
	tmpl, err := template.ParseFiles("error.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	page := struct{ Title, ErrorMsg, HomePath string }{
		Title:    "GoAdminer v1",
		ErrorMsg: _err.Error(),
		HomePath: HomePath,
	}
	err = tmpl.Execute(w, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func parseConnectForm(map_ url.Values) (map[string]string, error) {
	expected := map[string]string{
		"driver":   "System",
		"dbname":   "Database name",
		"password": "Password",
		"username": "Username",
	}
	ret := make(map[string]string)
	for name, field := range expected {
		v, ok := map_[name]
		if !ok || v[0] == "" {
			return nil, fmt.Errorf("Missing form field:%s", field)
		}
		ret[name] = v[0]
	}
	return ret, nil
}
