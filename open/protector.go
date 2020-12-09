package open

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

type Protector interface {
}

func RedirectToLastLocation(w http.ResponseWriter, r *http.Request) {
	location, err := r.Cookie("location")

	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, location.Value, http.StatusFound)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		MaxAge:   0,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	return state
}

func setLastLocationCookie(w http.ResponseWriter, url string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "location",
		Value:    url,
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
		Path:     "/",
	})
}
