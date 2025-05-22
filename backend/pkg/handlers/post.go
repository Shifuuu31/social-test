package handlers

import (
	"fmt"
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

func (app *SocialApp) SetupRoutes(mux *http.ServeMux) {
	postMux := http.NewServeMux()

	postMux.HandleFunc("GET /new", app.NewPost)
	postMux.HandleFunc("POST /", app.GetFeedPosts)
	postMux.HandleFunc("POST /new/upload", UploadHandler)

	log.Println("Mounting post multiplexer at /post/")

	mux.Handle("/post/", http.StripPrefix("/post", postMux))
}

func (app *SocialApp) GetFeedPosts(w http.ResponseWriter, r *http.Request) {
	log.Printf("Post root path accessed: %s", r.URL.Path)

	if r.URL.Path != "/" {
		log.Printf("Not found within post handler: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	fmt.Fprintln(w, "Listing all posts")
}

func (app *SocialApp) NewPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("New post path accessed: %s", r.URL.Path)
	if r.URL.Path != "/new" {
		log.Printf("Not found within post handler: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	var post models.Post
	if err := utils.DecodeJson(r, &post); err != nil {
		log.Printf("internalServerERROR: %s", r.URL.Path)
		return
	}
	stmt, err := app.Posts.DB.Prepare(`
    INSERT INTO posts (user_id, group_id, content, image, privacy, created_at)
    VALUES (?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(post.OwnerId, post.GroupId, post.Content, post.Image, post.Privacy, time.Now()); err != nil {
		log.Fatal(err)
	}
	if post.Privacy == "private" {
		for _, id := range post.ChosenUsersIds {
			app.Posts.DB.Exec("INSERT INTO post_privacy (chosen_id , post_id) VALUES (?, ?)", id, post.Id)
		}
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// // Parse up to 10MB of form data
	// err := r.ParseMultipartForm(10 << 20) // 10MB
	// if err != nil {
	// 	http.Error(w, "Failed to parse form", http.StatusBadRequest)
	// 	return
	// }

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if it's an image
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		http.Error(w, "Only image files are allowed", http.StatusBadRequest)
		return
	}

	// Make sure the uploads directory exists
	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Save the file
	dst, err := os.Create(filepath.Join("images", handler.Filename)) // TODO might wann add user spicific folder assignment
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Image %s uploaded successfully", handler.Filename)
}
