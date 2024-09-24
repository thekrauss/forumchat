package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"realtime/database"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username  string `json:"username"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
type Response struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}

func (s *MyServer) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var user User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				log.Println("Failed to decode request payload:", err)
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}
			log.Println("Decoded user:", user)
			if len(user.Username) > 16 {
				http.Error(w, "Username must not be between 6 and 16 caracteres", http.StatusBadRequest)
				return
			}
			if len(user.FirstName) > 16 || len(user.LastName) > 16 {
				http.Error(w, "firstname must be between 6 and 16 caractères", http.StatusBadRequest)
				return
			}
			if !isValidGender(user.Gender) {
				http.Error(w, "Gender must be 'male' or 'female'.", http.StatusBadRequest)
				return
			}
			if len(user.Password) < 6 || len(user.Password) > 16 {
				http.Error(w, "Password must be between 6 and 16 characters long", http.StatusBadRequest)
				return
			}
			DB, err := s.Store.OpenDatabase()
			if err != nil {
				log.Println("Failed to open database:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			_, err = DB.Exec(database.Users_Table)
			if err != nil {
				log.Println("Error creating users table:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			defer DB.Close()

			log.Println("Database and table ready")
			if err := RegisterUser(w, r, DB, user); err != nil {
				log.Println("Failed to create user:", err)
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}
			response := Response{
				Message: "User registered successfully",
				User:    user,
			}
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Println("Failed to encode response:", err)
				http.Error(w, "Failed to send response", http.StatusInternalServerError)
				return
			}
		} else {
			http.NotFound(w, r)
		}
	}
}
func RegisterUser(w http.ResponseWriter, r *http.Request, DB *sql.DB, user User) error {
	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}
	err := CreateUser(DB, user)
	if err != nil {
		return err
	}
	return nil
}
func CreateUser(db *sql.DB, user User) error {
	var countEmail, countUsername int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", user.Email).Scan(&countEmail)
	if err != nil {
		log.Println("Failed to check email existence:", err)
		return err
	}
	if countEmail > 0 {
		log.Println("Email already exists")
		return fmt.Errorf("email already exists")
	}
	log.Println("Email does not exist, proceeding")
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", user.Username).Scan(&countUsername)
	if err != nil {
		log.Println("Failed to check username existence:", err)
		return err
	}
	if countUsername > 0 {
		log.Println("Username already exists")
		return fmt.Errorf("username already exists")
	}
	log.Println("Username does not exist, proceeding")
	var hashedPassword []byte
	if user.Password != "" {
		// If user provides a password, hash it
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Failed to hash password:", err)
			return fmt.Errorf("failed to hash password: %w", err)
		}
	} else {
		// If no password provided (OAuth case), set hashedPassword to empty slice
		hashedPassword = []byte("")
	}
	query := `INSERT INTO users (username, age, gender, firstname, lastname, email, password)
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(query, user.Username, user.Age, user.Gender, user.FirstName, user.LastName, user.Email, string(hashedPassword))
	if err != nil {
		log.Println("Failed to execute insert query:", err)
		return fmt.Errorf("failed to create user: %w", err)
	}
	log.Println("User successfully created")
	return nil
}
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	at := strings.Index(email, "@")
	dot := strings.LastIndex(email, ".")
	return at > 0 && dot > at+1 && dot < len(email)-1
}
func isValidGender(gender string) bool {
	return gender == "Homme" || gender == "Femme"
}
