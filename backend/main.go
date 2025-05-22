package main

import (
	"fmt"
	"log"
	"net/http"

	"social-network/pkg/db/sqlite"
	"social-network/pkg/handlers"
	"social-network/pkg/models"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	db := sqlite.OpenDB()
	defer db.Close()
	log.Println("Database connection established")

	mainMux := http.NewServeMux()

	// authHandler := handlers.NewAuthHandler(database)
	// profileHandler := handlers.NewProfileHandler(database)
	// postHandler := handlers.NewPostHandler(db)

	log.Println("Handlers created")

	// authHandler.SetupRoutes(mainMux)
	// profileHandler.SetupRoutes(mainMux)
	// postHandler.SetupRoutes(mainMux)
	app := &handlers.SocialApp{
		Posts: &models.PostModel{
			DB: db,
		},
	}
	app.SetupPostRoutes(mainMux) // Setup routes for posts
	app.SetupCommentRoutes(mainMux) // Setup routes for comments

	log.Println("Routes set up")

	mainMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			log.Printf("404 Not Found: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Welcome to the Social Network")
	})

	handler := loggingMiddleware(mainMux)

	serverAddr := ":8080"
	log.Printf("Server starting at http://localhost%s", serverAddr)

	if err := http.ListenAndServe(serverAddr, handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
