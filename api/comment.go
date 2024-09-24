package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"realtime/database"
	"strconv"
	"time"
)

type Comment struct {
	ID                 int         `json:"id"`
	PostID             int         `json:"post_id"`
	Comments           []Comment   `json:"comments,omitempty"`
	Content            string      `json:"content"`
	UserID             int         `json:"user_id,omitempty"`
	User               User        `json:"user"`
	UsernamePost       string      `json:"username_post"`
	TotalLikeComment   map[int]int `json:"total_like_comment,omitempty"`
	TotalUnLikeComment map[int]int `json:"total_unlike_comment,omitempty"`
	CreatedAt          string      `json:"created_at"`
	CreatePosts        Post        `json:"create_posts,omitempty"`
}

func (s *MyServer) CreateCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var comment Comment

			if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
				log.Printf("Failed to decode comment request payload: %v. Payload: %v", err, r.Body)
				http.Error(w, "Invalid request comment payload", http.StatusBadRequest)
				return
			}

			fmt.Printf("data: %v\n", comment.Content)
			log.Printf("Received PostID: %v for comment creation", comment.PostID)

			if comment.PostID <= 0 {
				log.Println("Invalid post ID")
				http.Error(w, "Invalid post ID", http.StatusBadRequest)
				return
			}

			// Vérifier si le post_id existe dans la base de données
			if !s.PostExists(comment.PostID) {
				log.Println("Post ID does not exist")
				http.Error(w, "Post ID does not exist", http.StatusBadRequest)
				return
			}

			comment.CreatedAt = time.Now().Format(time.RFC3339)

			userID, ok := r.Context().Value(userIDKey).(int)
			if !ok {
				log.Println("User ID not found in context")
				http.Error(w, "User ID not found in context", http.StatusUnauthorized)
				return
			}
			comment.UserID = userID

			fmt.Printf("userID: %v\n", comment.UserID)

			if err := s.StoreComment(comment); err != nil {
				log.Println("Failed to store comment:", err)
				http.Error(w, "Failed to store comment", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(comment)
		} else {
			http.NotFound(w, r)
		}
	}
}

func (s *MyServer) StoreComment(comment Comment) error {
	DB, err := s.Store.OpenDatabase()
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer DB.Close()

	username, err := GetUsernameByID(DB, comment.UserID)
	if err != nil {
		log.Println("failed to get username", err)
	}
	comment.UsernamePost = username

	_, err = DB.Exec(database.Comments_table)
	if err != nil {
		log.Println("Error creating comment table:", err)
		return fmt.Errorf("failed to open database: %v", err)
	}
	log.Println("Database and table comment ready")

	query := `INSERT INTO comments (post_id, content, user_id, username, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err = DB.Exec(query, comment.PostID, comment.Content, comment.UserID, comment.UsernamePost, comment.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert comment into database: %v", err)
	}

	return nil
}

func (s *MyServer) ListCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		postIDStr := query.Get("post_id")
		if postIDStr == "" {
			http.Error(w, "Post ID not provided", http.StatusBadRequest)
			return
		}

		postID, err := strconv.Atoi(postIDStr)
		if err != nil || postID <= 0 {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		comments, err := GetCommentsByPost(DB, postID)
		if err != nil {
			http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"comments": comments,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Failed to encode comments to JSON:", err)
			http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
		}
	}
}

func GetCommentsByPost(DB *sql.DB, id int) ([]Comment, error) {
	if DB == nil {
		return nil, errors.New("database connection is nil")
	}
	rows, err := DB.Query("SELECT id, content, post_id, user_id, created_at FROM comments WHERE post_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var commentRow Comment
		err := rows.Scan(&commentRow.ID, &commentRow.Content, &commentRow.PostID, &commentRow.UserID, &commentRow.CreatedAt)
		if err != nil {
			return nil, err
		}
		username, err := GetUsernamePost(DB, commentRow.PostID)
		if err != nil {
			return nil, err
		}
		commentRow.UsernamePost = username
		comments = append(comments, commentRow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func GetUsernamePost(DB *sql.DB, postID int) (string, error) {
	var username string
	err := DB.QueryRow("SELECT username FROM users INNER JOIN posts ON users.id = posts.user_id WHERE posts.id = ?", postID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (s *MyServer) PostExists(postID int) bool {
	DB, err := s.Store.OpenDatabase()
	if err != nil {
		log.Println("Failed to open database:", err)
		return false
	}
	defer DB.Close()

	var exists bool
	err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
	if err != nil {
		log.Println("Failed to check if post exists:", err)
		return false
	}
	return exists
}
