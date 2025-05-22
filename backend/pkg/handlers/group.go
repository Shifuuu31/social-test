package handlers

import "net/http"

func GroupMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", CreateGroup)
	mux.HandleFunc("GET /", FetchGroup)
	mux.HandleFunc("POST /{id}/invite/{user_id}", HandleInvite)
	mux.HandleFunc("POST /{id}/join", HandleJoin)
	return mux
}



func CreateGroup(w http.ResponseWriter, r *http.Request) {

 }

 func FetchGroup(w http.ResponseWriter, r *http.Request) {

 }

  func HandleInvite(w http.ResponseWriter, r *http.Request) {

 }

  func HandleJoin(w http.ResponseWriter, r *http.Request) {

 }