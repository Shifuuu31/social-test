package handlers

import (
	"log"
	"net/http"
	"time"

	"social-network/pkg/models"
	"social-network/pkg/utils"
)

func (app *SocialApp) SetupCommentRoutes(mux *http.ServeMux) {
	commentMux := http.NewServeMux()

	commentMux.HandleFunc("GET /new", app.NewPost)
	commentMux.HandleFunc("POST /", app.GetFeedPosts)
	// commentMux.HandleFunc("POST /new/upload", UploadHandler) // TODO need to handel image

	log.Println("Mounting post multiplexer at /comment/")

	mux.Handle("/comment/", http.StripPrefix("/comment", commentMux))
}

func (app *SocialApp) GetFeedComments(w http.ResponseWriter, r *http.Request) {
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
	var comment models.Comment
	if err := utils.DecodeJson(r, &comment); err != nil {
		utils.EncodeJson(w, 500, nil)
		return
	}
	stmt, err := app.Posts.DB.Prepare(`
    INSERT INTO comments (user_id, post_id, content, image, created_at) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		utils.EncodeJson(w, 500, nil)
		return
	}
	defer stmt.Close()

	file, handler, err := r.FormFile("image")
	if err != nil {
		utils.EncodeJson(w, 500, nil)
		return
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

	if _, err = stmt.Exec(comment.OwnerId, comment.Post_id, comment.Content, temp, time.Now()); err != nil {
		utils.EncodeJson(w, 200, nil)
		return
	}

	utils.EncodeJson(w, 200, nil)
}
