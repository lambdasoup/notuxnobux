package app

import "net/http"

func init() {
	http.HandleFunc("/", handleFunc)
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}
