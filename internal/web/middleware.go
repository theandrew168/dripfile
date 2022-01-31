package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/theandrew168/dripfile/internal/core"
)

// TODO: mware for ParseForm + badRequestResponse?

type contextKey string

const (
	contextKeySession = contextKey("session")
)

// helper for pulling a session value of a request
func (app *Application) requestSession(r *http.Request) (core.Session, error) {
	session, ok := r.Context().Value(contextKeySession).(core.Session)
	if !ok {
		return core.Session{}, fmt.Errorf("failed context value cast to core.Session")
	}

	return session, nil
}

// check for valid session, redirect to /login if not found
func (app *Application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check for session cookie
		sessionID, err := r.Cookie(sessionIDCookieName)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// check for session in database
		session, err := app.storage.Session.Read(sessionID.Value)
		if err != nil {
			// user has an expired session cookie, delete it
			if errors.Is(err, core.ErrNotExist) {
				cookie := NewExpiredCookie(sessionIDCookieName)
				http.SetCookie(w, &cookie)
			}

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// attach session to the request context
		ctx := context.WithValue(r.Context(), contextKeySession, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
