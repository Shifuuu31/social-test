package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"social-network/pkg/models"
	"social-network/pkg/utils"
)

// Posts
// After a user is logged in he/she can create posts and comments on already created posts. While creating a post or a comment, the user can include an image or GIF.

// The user must be able to specify the privacy of the post:

// public (all users in the social network will be able to see the post) [no condition to fetch]
// almost private (only followers of the creator of the post will be able to see the post)
// private (only the followers chosen by the creator of the post will be able to see it)

func (app *SocialApp) SetupCommentRoutes(mux *http.ServeMux) {
	commentMux := http.NewServeMux()

	commentMux.HandleFunc("GET /new", app.NewPost)
	commentMux.HandleFunc("POST /", app.GetFeedPosts)
	// commentMux.HandleFunc("POST /new/upload", UploadHandler) // TODO need to handel image 

	log.Println("Mounting post multiplexer at /comment/")

	mux.Handle("/comment/", http.StripPrefix("/comment", commentMux))
}

func (app *SocialApp) GetFeedComments(w http.ResponseWriter, r *http.Request) {
	// log.Printf("Post root path accessed: %s", r.URL.Path)
	// TODO need to specify the methode
	if r.URL.Path != "/" {
		utils.EncodeJson(w, 404, nil)

		return
	}
	// fmt.Fprintln(w, "Listing all posts")
}

func (app *SocialApp) NewComment(w http.ResponseWriter, r *http.Request) {
	log.Printf("New post path accessed: %s", r.URL.Path)
	if r.URL.Path != "/new" {
		utils.EncodeJson(w, 500, nil)

		return
	}
	var post models.Post
	if err := utils.DecodeJson(r, &post); err != nil {
		utils.EncodeJson(w, 500, nil)
		return
	}
	stmt, err := app.Posts.DB.Prepare(`
    INSERT INTO posts (user_id, group_id, content, image, privacy, created_at)
    VALUES (?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		utils.EncodeJson(w, 500, nil)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(post.OwnerId, post.GroupId, post.Content, post.Image, post.Privacy, time.Now()); err != nil {
		log.Fatal(err)
	}
	if post.Privacy == "private" {
		for _, id := range post.ChosenUsersIds {
			if _, err := app.Posts.DB.Exec("INSERT INTO post_privacy (chosen_id , post_id) VALUES (?, ?)", id, post.Id); err != nil {
				utils.EncodeJson(w, 500, nil)
				return
			}
		}
	}
	utils.EncodeJson(w, 200, "done")
}


