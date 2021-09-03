package routes

import (
	"fmt"
	"net/http"
)

const DisconnectPath = "/logout"

func HandleDisconnect(w http.ResponseWriter, r *http.Request) {
	sessionManager := GetGlobalSessionManager()
	var sess *Session = nil
	cookie, err := r.Cookie(sessionManager.Name)
	if err == nil {
		sess = sessionManager.get(cookie.Value)
	}
	fmt.Println("Session to delete: ", sess)
	handleExpiredSession(w, r, sess)

	http.Redirect(w, r, HomePath, http.StatusPermanentRedirect)
}
