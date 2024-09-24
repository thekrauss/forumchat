package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"realtime/database"
	"strconv"
)

func (s *MyServer) LikeCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Parse request form
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get comment ID from form
		commentID, err := strconv.Atoi(r.FormValue("comment_id"))
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		// Get user ID from session
		userID, err := s.getUserIDFromSession(r)
		if err != nil {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		// Open database connection
		DB, err := s.Store.OpenDatabase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		// Begin a transaction
		tx, err := DB.Begin()
		if err != nil {
			http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
			return
		}

		// Toggle like for the comment
		err = ToggleLikeComment(userID, commentID, DB, tx)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to toggle like", http.StatusInternalServerError)
			return
		}

		// Commit transaction
		err = tx.Commit()
		if err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Comment like toggled successfully"})
	}
}

func (s *MyServer) UnlikeCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Parse request form
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get comment ID from form
		commentID, err := strconv.Atoi(r.FormValue("comment_id"))
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		// Get user ID from session
		userID, err := s.getUserIDFromSession(r)
		if err != nil {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		// Open database connection
		DB, err := s.Store.OpenDatabase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		// Begin a transaction
		tx, err := DB.Begin()
		if err != nil {
			http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
			return
		}

		// Toggle unlike for the comment
		err = ToggleUnLikeComment(userID, commentID, DB, tx)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to toggle unlike", http.StatusInternalServerError)
			return
		}

		// Commit transaction
		err = tx.Commit()
		if err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Comment unlike toggled successfully"})
	}
}

func ToggleLikeComment(userID, commentID int, DB *sql.DB, tx *sql.Tx) error {
	liked, err := UserLikedComment(userID, commentID, tx)
	if err != nil {
		log.Println("Error checking if user liked comment:", err)
		return err
	}

	unliked, err := UserUnLikedComment(userID, commentID, tx)
	if err != nil {
		log.Println("Error checking if user unliked comment:", err)
		return err
	}

	if liked {
		err := DeleteLikeComment(userID, commentID, tx)
		if err != nil {
			log.Println("Error delete like comment:", err)
			return err
		}
	} else {
		if unliked {
			err := DeleteUnLikeComment(userID, commentID, tx)
			if err != nil {
				log.Println("Error delete unlike comment:", err)
				return err
			}
		}
		err := CreateLikeComment(userID, commentID, tx)
		if err != nil {
			log.Println("Error create like:", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction :", err)
		return err
	}

	return nil
}

func ToggleUnLikeComment(userID, commentID int, DB *sql.DB, tx *sql.Tx) error {
	_, err := DB.Exec(database.UnlikesComment_Table)
	if err != nil {
		log.Println("Error creating unlikesComments table:", err)
		return err
	}

	liked, err := UserLikedComment(userID, commentID, tx)
	if err != nil {
		log.Println("Error checking if user liked comment:", err)
		return err
	}

	unliked, err := UserUnLikedComment(userID, commentID, tx)
	if err != nil {
		log.Println("Error checking if user unliked comment:", err)
		return err
	}

	if unliked {
		err := DeleteUnLikeComment(userID, commentID, tx)
		if err != nil {
			log.Println("Error delete unlike comment:", err)
			return err
		}
	} else {
		if liked {
			err := DeleteLikeComment(userID, commentID, tx)
			if err != nil {
				log.Println("Error delete like comment:", err)
				return err
			}
		}
		err := CreateUnLikeComment(userID, commentID, DB, tx)
		if err != nil {
			log.Println("Error create unlike:", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction :", err)
		return err
	}

	return nil
}

/*--------------------------------------------------------------------*/

func DeleteLikeComment(userID, commentID int, tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM commentlikes WHERE user_id = ? AND comment_id = ?", userID, commentID)
	if err != nil {
		log.Println("Error delete commentlikes:", err)
		return err
	}

	_, err = tx.Exec("UPDATE comments SET total_likes = total_likes - 1 WHERE id = ?", commentID)
	if err != nil {
		log.Println("Error update total likes comments:", err)
		return err
	}

	return nil
}

func CreateLikeComment(userID, commentID int, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO commentlikes (user_id, comment_id) VALUES (?, ?)", userID, commentID)
	if err != nil {
		log.Println("Error create commentlikes:", err)
		return err
	}

	_, err = tx.Exec("UPDATE comments SET total_likes = total_likes + 1 WHERE id = ?", commentID)
	if err != nil {
		log.Println("Error update total likes comments:", err)
		return err
	}

	return nil
}

func DeleteUnLikeComment(userID, commentID int, tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM unlikescomment WHERE user_id = ? AND comment_id = ?", userID, commentID)
	if err != nil {
		log.Println("Error delete unlikescomment:", err)
		return err
	}

	_, err = tx.Exec("UPDATE comments SET total_unlikescomment = total_unlikescomment - 1 WHERE id = ?", commentID)
	if err != nil {
		log.Println("Error update total unlikes comments:", err)
		return err
	}

	return nil
}

func CreateUnLikeComment(userID, commentID int, DB *sql.DB, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO unlikescomment (comment_id, user_id) VALUES (?, ?)", commentID, userID)
	if err != nil {
		log.Println("Error create unlikescomment:", err)
		return err
	}

	_, err = tx.Exec("UPDATE comments SET total_unlikescomment = total_unlikescomment + 1 WHERE id = ?", commentID)
	if err != nil {
		log.Println("Error update total unlikes comments:", err)
		return err
	}

	return nil
}

func UserUnLikedComment(userID, commentID int, tx *sql.Tx) (bool, error) {
	var count int
	err := tx.QueryRow("SELECT COUNT(*) FROM unlikescomment WHERE user_id = ? AND comment_id = ?", userID, commentID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil // return true, nil
}

func UserLikedComment(userID, commentID int, tx *sql.Tx) (bool, error) {
	// verifie si l'utilisateur a aimé le commentaire
	var count int
	err := tx.QueryRow("SELECT COUNT(*) FROM commentlikes WHERE user_id = ? AND comment_id = ?", userID, commentID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
