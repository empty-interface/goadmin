package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/empty-interface/goadmin/session"
)

const ConnectPath = "/connect"

type HandleError func(w http.ResponseWriter, r *http.Request) (int, error)

func (handler HandleError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if code, originalErr := handler(w, r); originalErr != nil {
		tmpl, err := template.ParseFiles("html/error.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		page := struct{ Title, ErrorMsg, HomePath string }{
			Title:    "GoAdminer v1",
			ErrorMsg: originalErr.Error(),
			HomePath: HomePath,
		}
		err = tmpl.Execute(w, page)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(code)
	}
}

func HandleConnect(connect func(*session.Session) error) HandleError {
	return HandleError(func(w http.ResponseWriter, r *http.Request) (int, error) {
		r.ParseForm()
		form, err := parseConnectForm(r.PostForm)
		if err != nil {
			return http.StatusBadRequest, err
		}
		saved := false
		if form["rememberme"] == "on" {
			saved = true
		}
		sess, err := session.NewSession(form["driver"], form["username"], form["password"], form["dbname"], saved)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		err = connect(sess)
		if err != nil {
			return http.StatusBadRequest, err
		}
		addOrRefreshLoginSession(w, r, sess)
		http.Redirect(w, r, HomePath, http.StatusPermanentRedirect)
		return -1, nil
	})
}
func addOrRefreshLoginSession(w http.ResponseWriter, r *http.Request, sess *session.Session) {
	sessionManager := session.GetGlobalSessionManager()
	if _sess := sessionManager.Get(sess.Uuid); _sess != nil {
		_sess.Refresh()
		sess = _sess
	} else {
		sessionManager.Set(sess)
	}
	fmt.Println("Creating a new session ,expires at ", sess.ExpiresAt().String())
	http.SetCookie(w, &http.Cookie{
		Name:     sessionManager.Name,
		Value:    sess.Uuid,
		Expires:  sess.ExpiresAt(),
		HttpOnly: true,
	})
}
func handleExpiredSession(w http.ResponseWriter, r *http.Request, sess *session.Session) {
	sessionManager := session.GetGlobalSessionManager()
	if sess != nil {
		sessionManager.Delete(sess.Uuid)
		fmt.Println("Deleted session")
	}
	http.SetCookie(w, &http.Cookie{
		Name:   sessionManager.Name,
		Value:  "",
		MaxAge: -1,
		// Expires: time.Now(),
	})
}

type handlerWithSession func(http.ResponseWriter, *http.Request, *session.Session) (int, error)

func HandleSession(connect func(*session.Session) error, next handlerWithSession) HandleError {
	return HandleError(func(w http.ResponseWriter, r *http.Request) (int, error) {
		sessionManager := session.GetGlobalSessionManager()
		cookie, err := r.Cookie(sessionManager.Name)
		if err != nil {
			fmt.Println("No session", err.Error())
			handleExpiredSession(w, r, nil)
			HandleHome(w, r)
			return -1, nil
		}
		sess := session.GetGlobalSessionManager().Get(cookie.Value)
		if sess == nil || sess.Expired() {
			fmt.Println("Session expired")
			handleExpiredSession(w, r, sess)
			HandleHome(w, r)
			return -1, nil
		}
		if sess.Conn == nil {
			err := connect(sess)
			if err != nil {
				handleExpiredSession(w, r, sess)
				HandleHome(w, r)
				return -1, nil
			}
		}
		fmt.Println("Session is alive", sess.Uuid)
		return next(w, r, sess)
	})
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
	ret["rememberme"] = "off"
	v, ok := map_["rememberme"]
	if ok && v[0] == "on" {
		ret["rememberme"] = "on"
	}
	return ret, nil
}
