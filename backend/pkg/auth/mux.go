package auth

import "net/http"

func AuthMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", RegisterHandler)
	mux.HandleFunc("POST /login", LoginHandler)

	return mux
}
