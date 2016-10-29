package app

import "net/http"

//go:generate elm-make App.elm

func init() {
	http.HandleFunc("/", handleFunc)
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}
