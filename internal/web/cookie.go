package web

import (
	"net/http"
	"time"
)

var (
	sessionIDCookieName = "session_id"
)

func NewSessionCookie(name, value string) http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/", // applies to the whole site
		Domain:   "",  // will default to the server's base domain
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return cookie
}

func NewPermanentCookie(name, value string, expiry time.Time) http.Cookie {
	// round expiration up to nearest second
	cookie := NewSessionCookie(name, value)
	cookie.Expires = time.Unix(expiry.Unix()+1, 0)
	cookie.MaxAge = int(time.Until(expiry).Seconds() + 1)
	return cookie
}

func NewExpiredCookie(name string) http.Cookie {
	// expires now
	cookie := NewSessionCookie(name, "")
	cookie.Expires = time.Unix(1, 0)
	cookie.MaxAge = -1
	return cookie
}
