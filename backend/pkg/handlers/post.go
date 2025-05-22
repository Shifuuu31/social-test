package handlers

import (
	"log"
	"net/http"

	"social-network/pkg/models"
	"social-network/pkg/utils"
)

// Posts
// After a user is logged in he/she can create posts and comments on already created posts. While creating a post or a comment, the user can include an image or GIF.

// The user must be able to specify the privacy of the post:

// public (all users in the social network will be able to see the post) [no condition to fetch]
// almost private (only followers of the creator of the post will be able to see the post)
// private (only the followers chosen by the creator of the post will be able to see it)

func (app *SocialApp) SetupRoutes(mux *http.ServeMux) {
	postMux := http.NewServeMux()

	postMux.HandleFunc("POST /new", app.NewPost)
	postMux.HandleFunc("POST /", app.GetFeedPosts)

	log.Println("Mounting post multiplexer at /post/")

	mux.Handle("/post/", http.StripPrefix("/post", postMux))
}

func (app *SocialApp) DebugRoutes() string {
	return "Available post routes: /, /new, /view/{id}"
}

func (app *SocialApp) GetFeedPosts(w http.ResponseWriter, r *http.Request) {
	var filter *models.PostFilter

	if err := utils.DecodeJson(r, &filter); err != nil {
		log.Println(err)
		return
	}

	posts, err := app.Posts.GetPosts(filter)
	if err != nil {
		utils.EncodeJson(w, http.StatusInternalServerError, nil)
		return
	}
	if err := utils.EncodeJson(w, http.StatusOK, posts); err != nil {
		log.Println(err)
	}
	
}
func (app *SocialApp) NewPost(w http.ResponseWriter, r *http.Request) {}
