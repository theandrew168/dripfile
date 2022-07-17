package web

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/postgresql"
)

type contextKey string

const (
	contextKeySession  = contextKey("session")
	requestBodyMaxSize = 4096
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

		// check for session in postgresql
		sessionHash := fmt.Sprintf("%x", sha256.Sum256([]byte(sessionID.Value)))
		session, err := app.store.Session.Read(sessionHash)
		if err != nil {
			// user has an invalid session cookie, delete it
			if errors.Is(err, postgresql.ErrNotExist) {
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

func (app *Application) setSecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *Application) limitRequestSize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, requestBodyMaxSize)

		next.ServeHTTP(w, r)
	})
}
