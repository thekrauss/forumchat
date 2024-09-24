package api

import (
	"encoding/json"
	"log"
	"net/http"
	"realtime/wsk"
	"strconv"
	"time"
)

func (s *MyServer) OnlineUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.WebSocketChat.Mu.Lock()
		defer s.WebSocketChat.Mu.Unlock()

		users := make([]string, 0, len(s.WebSocketChat.Users))
		for username := range s.WebSocketChat.Users {
			users = append(users, username)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		}
	}
}

func (s *MyServer) MessagesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request to", r.Method, "method")

		switch r.Method {
		case http.MethodGet:
			username := r.URL.Query().Get("user")
			offsetStr := r.URL.Query().Get("offset")
			offset, err := strconv.Atoi(offsetStr)
			if err != nil {
				offset = 0
			}

			log.Println("Fetching messages for user:", username, "with offset:", offset)

			DB, err := s.Store.OpenDatabase()
			if err != nil {
				log.Println("Failed to open database:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			defer DB.Close()

			query := `
                SELECT sender_username, target_username, content, timestamp 
                FROM chatHistory 
                WHERE sender_username = ? OR target_username = ? 
                ORDER BY timestamp DESC 
                LIMIT 10 OFFSET ?
            `
			rows, err := DB.Query(query, username, username, offset)
			if err != nil {
				log.Println("Failed to fetch chatHistory:", err)
				http.Error(w, "Failed to fetch chatHistory", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var messages []wsk.Message
			for rows.Next() {
				var msg wsk.Message
				if err := rows.Scan(&msg.SenderUsername, &msg.TargetUsername, &msg.Content, &msg.Timestamp); err != nil {
					log.Println("Failed to scan message:", err)
					http.Error(w, "Failed to scan message", http.StatusInternalServerError)
					return
				}
				messages = append(messages, msg)
			}

			log.Println("Fetched", len(messages), "messages for user:", username)

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(messages); err != nil {
				log.Println("Failed to encode messages:", err)
				http.Error(w, "Failed to encode messages", http.StatusInternalServerError)
			}

		case http.MethodPost:
			log.Println("Received POST request to send a message")
			var msg wsk.Message
			if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
				log.Println("Failed to decode message:", err)
				http.Error(w, "Failed to decode message", http.StatusBadRequest)
				return
			}
			log.Printf("Message content: %v", msg)

			msg.Timestamp = time.Now()
			sender, ok := r.Context().Value("username").(string)
			if !ok {
				log.Println("Username not found in context")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			msg.SenderUsername = sender

			DB, err := s.Store.OpenDatabase()
			if err != nil {
				log.Println("Failed to open database:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			defer DB.Close()

			log.Printf("Saving message: %+v", msg)

			query := `INSERT INTO chatHistory (sender_username, target_username, content, timestamp) VALUES (?, ?, ?, ?)`
			log.Printf("Inserting message with sender: %s, target: %s, content: %s", msg.SenderUsername, msg.TargetUsername, msg.Content)

			if msg.SenderUsername == "" {
				log.Println("Warning: SenderUsername is empty before saving to the database")
			}
			result, err := DB.Exec(query, msg.SenderUsername, msg.TargetUsername, msg.Content, msg.Timestamp)

			if err != nil {
				log.Printf("Failed to save message: %v", err)
				http.Error(w, "Failed to save message", http.StatusInternalServerError)
				return
			}

			log.Printf("Saving message from %s to %s: %v", msg.SenderUsername, msg.TargetUsername, msg)

			id, err := result.LastInsertId()
			if err != nil {
				log.Printf("Failed to retrieve last inserted ID: %v", err)
			} else {
				log.Printf("Saved message with ID: %d from %s to %s", id, msg.SenderUsername, msg.TargetUsername)
			}

			log.Println("Saved message from", msg.SenderUsername, "to", msg.TargetUsername)
			log.Printf("Saving message: %v", msg)

			s.WebSocketChat.MessageChannel <- &msg
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
