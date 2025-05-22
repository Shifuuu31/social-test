package handlers

import (
	"log"
	"net/http"
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

func (app *SocialApp) SetupPostRoutes(mux *http.ServeMux) {
	postMux := http.NewServeMux()

	postMux.HandleFunc("GET /new", app.NewPost)
	postMux.HandleFunc("POST /", app.GetFeedPosts)
	// postMux.HandleFunc("POST /new/upload", UploadHandler)

	log.Println("Mounting post multiplexer at /post/")

	mux.Handle("/post/", http.StripPrefix("/post", postMux))
}

func (app *SocialApp) NewPost(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	defer stmt.Close()

	file, handler, err := r.FormFile("image")
	if err != nil {
		utils.EncodeJson(w, 500, nil)
	}

	defer file.Close()
	temp, status := "", 0
	if handler.Filename != "" {
		temp, status = utils.UploadHandler(file, handler)
		if status != 200 {
			utils.EncodeJson(w, status, nil)
			return
		}
	}

	if _, err = stmt.Exec(post.OwnerId, post.GroupId, post.Content, temp, post.Privacy, time.Now()); err != nil {
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
	utils.EncodeJson(w, 200, nil)
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
