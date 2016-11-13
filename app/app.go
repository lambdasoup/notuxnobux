package app

import (
	"net/http"
	"net/url"

	"github.com/lambdasoup/notuxnobux/steam"
)

//go:generate elm-make App.elm

func init() {
	http.HandleFunc("/login", loginFunc)
	http.HandleFunc("/jwt", jwtFunc)
}

func loginFunc(w http.ResponseWriter, r *http.Request) {

	base := new(url.URL)
	if r.TLS != nil {
		base.Scheme = "https"
	} else {
		base.Scheme = "http"
	}
	base.Host = r.Host

	switch steam.Mode(r) {
	case "":
		rtu, _ := url.Parse(base.String() + "/login")
		http.Redirect(w, r, steam.AuthURLFor(r, rtu).String(), http.StatusSeeOther)
	case "cancel":
		w.Write([]byte("authorization cancelled"))
	default:
		steamID, err := steam.FetchID(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// TODO create user entity
		// TODO make 7-day JWT
		// TODO make auth token
		// TODO put auth token & jwt into memcache
		// TODO redirect with auth token
		w.Write([]byte(steamID))
		http.Redirect(w, r, base.String()+"/?token=testtoken", http.StatusSeeOther)
	}
}

func jwtFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{test:123}"))
}
