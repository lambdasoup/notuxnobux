package app

import (
	"net/http"
	"net/url"

	"github.com/lambdasoup/notuxnobux/steam"
)

//go:generate elm-make App.elm

func init() {
	http.HandleFunc("/login", loginFunc)
}

func loginFunc(w http.ResponseWriter, r *http.Request) {

	switch steam.Mode(r) {
	case "":
		rtu, _ := url.Parse("http://localhost:8080/#login")
		http.Redirect(w, r, steam.AuthURLFor(r, rtu).String(), http.StatusSeeOther)
	case "cancel":
		w.Write([]byte("authorization cancelled"))
	default:
		steamID, err := steam.FetchID(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		//
		//  Do smth with steamId
		//
		w.Write([]byte(steamID))
	}
}
