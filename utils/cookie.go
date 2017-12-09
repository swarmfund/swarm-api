package utils

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"gitlab.com/swarmfund/api/internal/lorem"
	"gitlab.com/swarmfund/go/hash"
)

const DUIDPrefix = "dev_unique_id"

func DeviceUIDCookie(id, domain string) *http.Cookie {
	cookieName := DeviceUIDCookieName(id)
	return &http.Cookie{
		Name:     cookieName,
		Value:    lorem.Token(),
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().AddDate(0, 0, 14),
	}
}

func UpdateCookieExpires(cookie *http.Cookie) {
	cookie.Expires = time.Now().AddDate(0, 0, 14)
}

func DeviceUIDCookieName(id string) string {
	cookieNameHash := hash.Hash([]byte(id + DUIDPrefix))
	name := base64.URLEncoding.EncodeToString(cookieNameHash[:])
	return strings.TrimRight(name, "=")
}
