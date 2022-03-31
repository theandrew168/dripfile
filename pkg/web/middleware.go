package web

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"

	"github.com/theandrew168/dripfile/pkg/core"
)

type contextKey string

const (
	contextKeySession = contextKey("session")
)

// helper for pulling a session value of a request
func (app *Application) requestSession(r *http.Request) (core.Session, error) {
	session, ok := r.Context().Value(contextKeySession).(core.Session)
	if !ok {
		return core.Session{}, fmt.Errorf("invalid or missing session")
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
		sessionHash := fmt.Sprintf("%x", sha256.Sum256([]byte(sessionID.Value)))
		session, err := app.storage.Session.Read(sessionHash)
		if err != nil {
			// user has an invalid session cookie, delete it
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

// check for verified billing, request payment info if not
func (app *Application) requireBilling(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check for session
		session, err := app.requestSession(r)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// check if project has verified billing
		if !session.Account.Project.BillingVerified {
			// TODO: create setup intent
			// TODO: redirect to payment element
		}

		next.ServeHTTP(w, r)
	})
}

// set size limit and attempt to parse POSTed form data
func (app *Application) parseForm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 4096)

		err := r.ParseForm()
		if err != nil {
			app.badRequestResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// helper for wrapping HandlerFuncs
func (app *Application) parseFormFunc(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return app.parseForm(http.HandlerFunc(f))
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
