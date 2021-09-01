package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

const (
	HomePath       = "/"
	ConnectPath    = "/connect"
	DisconnectPath = "/disconnect"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("html/home.html")
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
			"mysql":    "mySQL",
			"postgres": "PostgreSQL",
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
func HandleConnect(connect func(*Session) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		form, err := parseConnectForm(r.PostForm)
		if err != nil {
			handleError(w, r, err, http.StatusBadRequest)
			return
		}
		sess, err := NewSession(form["driver"], form["username"], form["password"], form["dbname"])
		if err != nil {
			handleError(w, r, err, http.StatusInternalServerError)
			return
		}
		err = connect(sess)
		if err != nil {
			handleError(w, r, err, http.StatusBadRequest)
			return
		}
		addOrRefreshLoginSession(w, r, sess)
		http.Redirect(w, r, HomePath, http.StatusPermanentRedirect)
	})
}
func addOrRefreshLoginSession(w http.ResponseWriter, r *http.Request, sess *Session) {
	sessionManager := *GetGlobalSessionManager()
	if _sess := sessionManager.get(sess.uuid); _sess != nil {
		_sess.refresh()
		sess = _sess
	} else {
		sessionManager.set(sess)
	}
	fmt.Println("Creating a new session ,expires at ", sess.expiresAt().String())
	http.SetCookie(w, &http.Cookie{
		Name:    sessionManager.Name,
		Value:   sess.uuid,
		Expires: sess.expiresAt(),
	})
}
func handleExpiredSession(w http.ResponseWriter, r *http.Request, sess *Session) {
	sessionManager := *GetGlobalSessionManager()
	if sess != nil {
		sessionManager.delete(sess.uuid)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   sessionManager.Name,
		MaxAge: -1,
		// Value:  "",
	})
	fmt.Println("Deleted session")
}
func HandleDisconnect(w http.ResponseWriter, r *http.Request) {
	sessionManager := *GetGlobalSessionManager()
	var sess *Session = nil
	cookie, err := r.Cookie(sessionManager.Name)
	if err == nil {
		sess = sessionManager.get(cookie.Value)
	}
	handleExpiredSession(w, r, sess)
	http.Redirect(w, r, HomePath, http.StatusPermanentRedirect)
}
func HandleSession(connect func(*Session) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionManager := *GetGlobalSessionManager()
		cookie, err := r.Cookie(sessionManager.Name)
		if err != nil {
			fmt.Println("No session", err.Error())
			handleExpiredSession(w, r, nil)
			HandleHome(w, r)
			return
		}
		sess := GetGlobalSessionManager().get(cookie.Value)
		if sess == nil || sess.expired() {
			fmt.Println("Session expired")
			handleExpiredSession(w, r, sess)
			HandleHome(w, r)
			return
		}
		if sess.Conn == nil {
			err := connect(sess)
			if err != nil {
				handleExpiredSession(w, r, sess)
				HandleHome(w, r)
				return
			}
		}
		fmt.Println("Session is alive", sess.uuid)
		HandleDatabase(w, r, sess)
	})
}
func HandleDatabase(w http.ResponseWriter, r *http.Request, currentSession *Session) {
	// we get session
	tmpl, err := template.ParseFiles("html/database.html")
	if err != nil {
		handleError(w, r, err, http.StatusInternalServerError)
		return
	}
	currentSession.Conn.GetTables()
	page := struct {
		Title          string
		DisconnectPath string
		Tables         map[string]string
	}{
		Title:          fmt.Sprintf("GoAdmin - %s", currentSession.DBname),
		DisconnectPath: DisconnectPath,
		Tables:         make(map[string]string),
	}
	err = tmpl.Execute(w, page)
	if err != nil {
		handleError(w, r, err, http.StatusInternalServerError)
		return
	}
}
func handleError(w http.ResponseWriter, r *http.Request, _err error, status int) {
	tmpl, err := template.ParseFiles("html/error.html")
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
