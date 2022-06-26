package web

import (
	"crypto/sha256"
	"fmt"
	"net/http"
)

func (app *Application) handleLogoutForm(w http.ResponseWriter, r *http.Request) {
	// check for session cookie
	sessionID, err := r.Cookie(sessionIDCookieName)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// check for session in database
	sessionHash := fmt.Sprintf("%x", sha256.Sum256([]byte(sessionID.Value)))
	session, err := app.store.Session.Read(sessionHash)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// delete session from database
	err = app.store.Session.Delete(session)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// expire the existing session cookie
	cookie := NewExpiredCookie(sessionIDCookieName)
	http.SetCookie(w, &cookie)

	app.logger.Info("account logout", map[string]string{
		"project_id": session.Account.Project.ID,
		"account_id": session.Account.ID,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
