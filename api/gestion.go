package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func GetUserIDbyEmail(db *sql.DB, email string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func GetUsernameByEmail(db *sql.DB, email string) (string, error) {
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE email = ?", email).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func GetPasswordByEmail(db *sql.DB, email string) (string, error) {
	var passWordId string
	err := db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&passWordId)
	if err != nil {
		return "", err
	}
	return passWordId, nil
}

func GetUsernameByID(db *sql.DB, userID int) (string, error) {
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func GetUserIDbyUsername(db *sql.DB, username string) (int, error) {
	var userID int
	query := "SELECT id FROM users WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No user found with username:", username)
		}
		return 0, fmt.Errorf("failed to get user ID by username: %w", err)
	}
	return userID, nil
}

func GetPasswordByUsername(db *sql.DB, username string) (string, error) {
	var password string
	query := "SELECT password FROM users WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No password found for username:", username)
		}
		return "", fmt.Errorf("failed to get password by username: %w", err)
	}
	return password, nil
}

func GetPosts(DB *sql.DB) ([]Post, error) {
	if DB == nil {
		return nil, errors.New("database connection is nil")
	}

	rows, err := DB.Query("SELECT p.id, p.title, p.content, p.image_path, p.user_id, p.created_at, p.category, u.username FROM posts p JOIN users u ON p.user_id = u.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.ImagePath, &post.UserID, &post.CreatedAt, &post.Category, &post.Username)
		if err != nil {
			return nil, err
		}
		log.Println(post)
		posts = append(posts, post)
	}
	return posts, nil
}

func IsValidImageExtension(filename string) bool {
	validExtensios := []string{".jpg", ".jpeg", ".png", ".gif"}
	fileExt := strings.ToLower(filepath.Ext(filename))
	for _, ext := range validExtensios {
		if ext == fileExt {
			return true
		}
	}
	return false
}

func (s *MyServer) getUserIDFromSession(r *http.Request) (int, error) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return 0, fmt.Errorf("no user ID cookie found")
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID cookie value")
	}

	return userID, nil
}
