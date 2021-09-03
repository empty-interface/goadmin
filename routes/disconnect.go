package routes

import (
	"fmt"
	"net/http"

	"github.com/empty-interface/goadmin/session"
)

const DisconnectPath = "/logout"

func HandleDisconnect(w http.ResponseWriter, r *http.Request) {
	sessionManager := session.GetGlobalSessionManager()
	var sess *session.Session = nil
	cookie, err := r.Cookie(sessionManager.Name)
	if err == nil {
		sess = sessionManager.Get(cookie.Value)
	}
	fmt.Println("Session to delete: ", sess)
	handleExpiredSession(w, r, sess)

	http.Redirect(w, r, HomePath, http.StatusPermanentRedirect)
}
