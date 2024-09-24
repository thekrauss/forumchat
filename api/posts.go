package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"realtime/database"
	"strings"
	"time"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Category  string    `json:"category"`
	Content   string    `json:"content"`
	Username  string    `json:"username"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ImagePath string    `json:"image_path"`
}

var CategoryList = []string{"All", "HTML/CSS", "JavaScript", "Java", "C#", "C++", "Python", "PHP", "Ruby"}

/*--------------------------------------------------------------------------------------------------------------------------------*/
func (s *MyServer) CreatePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var post Post
			if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
				log.Println("Failed to decode request payload:", err)
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}
			post.CreatedAt = time.Now()

			if len(post.Title) == 0 || len(post.Content) == 0 || len(post.Category) == 0 {
				log.Println("All fields are required")
				http.Error(w, "All fields are required", http.StatusBadRequest)
				return
			}

			isValidCategory := false
			for _, category := range CategoryList {
				if post.Category == category {
					isValidCategory = true
					break
				}
			}
			if !isValidCategory {
				log.Println("Invalid category")
				http.Error(w, "Invalid category", http.StatusBadRequest)
				return
			}

			var imagePath string
			if _, handler, err := r.FormFile("image"); err == nil {
				if !IsValidImageExtension(handler.Filename) {
					log.Println("Invalid image file extension")
					http.Error(w, "Invalid image file extension", http.StatusBadRequest)
					return
				}

				imagePath, err = UploadImages(w, r, "./frontEnd/images")
				if err != nil {
					log.Println("Failed to upload image:", err)
					http.Error(w, "Failed to upload image", http.StatusInternalServerError)
					return
				}
				post.ImagePath = imagePath
			}

			userID, ok := r.Context().Value(userIDKey).(int)
			if !ok {
				log.Println("User ID not found in context")
				http.Error(w, "User ID not found in context", http.StatusUnauthorized)
				return
			}

			post.UserID = userID

			postID, err := s.StorePost(post)
			if err != nil {
				log.Println("Failed to save post:", err)
				http.Error(w, "Failed to save post", http.StatusInternalServerError)
				return
			}

			// Set the post ID before responding
			post.ID = postID

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(post)
		} else {
			http.NotFound(w, r)
		}
	}
}
func (s *MyServer) StorePost(post Post) (int, error) {
	// Open the database
	DB, err := s.Store.OpenDatabase()
	if err != nil {
		return 0, fmt.Errorf("failed to open database: %v", err)
	}
	defer DB.Close()

	// Ensure the posts table exists
	_, err = DB.Exec(database.Posts_Table)
	if err != nil {
		log.Println("Error creating posts table:", err)
		return 0, fmt.Errorf("failed to ensure posts table exists: %v", err)
	}
	log.Println("Database and table ready")

	// Insert the post into the database
	query := `INSERT INTO posts (title, content, category, user_id, created_at, image_path) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := DB.Exec(query, post.Title, post.Content, post.Category, post.UserID, post.CreatedAt, post.ImagePath)
	if err != nil {
		return 0, fmt.Errorf("failed to insert post into database: %v", err)
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %v", err)
	}

	return int(lastInsertID), nil
}

/*--------------------------------------------------------------------------------------------------------------*/

func (s *MyServer) ListPostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DB, err := s.Store.OpenDatabase()
		if err != nil {
			log.Println("failed to open database for ListPost", err)
			http.Error(w, "failed to open database for ListPost", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		posts, err := GetPosts(DB)
		if err != nil {
			log.Println("Failed to retrieve posts from the database:", err)
			http.Error(w, "Failed to retrieve posts from the database", http.StatusInternalServerError)
			return
		}
		category := r.FormValue("category")
		if category != "" && category != "All" {
			var filterPosts []Post
			for _, post := range posts {
				if strings.EqualFold(post.Category, category) {
					filterPosts = append(filterPosts, post)
				}
			}
			posts = filterPosts
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(posts); err != nil {
			log.Println("Failed to encode posts to JSON:", err)
			http.Error(w, "Failed to encode posts to JSON", http.StatusInternalServerError)
		}
	}
}
