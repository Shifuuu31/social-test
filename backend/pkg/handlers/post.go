package handlers

import (
	"fmt"
	"log"
	"net/http"
)

// Posts
// After a user is logged in he/she can create posts and comments on already created posts. While creating a post or a comment, the user can include an image or GIF.

// The user must be able to specify the privacy of the post:

// public (all users in the social network will be able to see the post) [no condition to fetch]
// almost private (only followers of the creator of the post will be able to see the post)
// private (only the followers chosen by the creator of the post will be able to see it)

func (app *socialApp) SetupRoutes(mux *http.ServeMux) {
	postMux := http.NewServeMux()

	postMux.HandleFunc("POST /new", app.NewPost)
	postMux.HandleFunc("POST /", app.GetFeedPosts)

	log.Println("Mounting post multiplexer at /post/")

	mux.Handle("/post/", http.StripPrefix("/post", postMux))
}

func (app *socialApp) DebugRoutes() string {
	return "Available post routes: /, /new, /view/{id}"
}

func (app *socialApp) GetFeedPosts(w http.ResponseWriter, r *http.Request) {
	log.Printf("Post root path accessed: %s", r.URL.Path)

	if r.URL.Path != "/" {
		log.Printf("Not found within post handler: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	fmt.Fprintln(w, "Listing all posts")
}
func (app *socialApp) NewPost(w http.ResponseWriter, r *http.Request) {}
