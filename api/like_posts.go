package api

import (
	"database/sql"
	"log"
	"net/http"
	"realtime/database"
	"strconv"
)

func (s *MyServer) LikePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(userIDKey).(int)
		if !ok {
			log.Println("User ID not found in context")
			http.Error(w, "User ID not found in context", http.StatusUnauthorized)
			return
		}

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			log.Println("failed to open database", err)
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		tx, err := DB.Begin()
		if err != nil {
			log.Println("Failed to begin transaction", err)
			http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
			return
		}

		err = ToggleLike(userID, postID, DB, tx)
		if err != nil {
			log.Println("Failed to toggle like", err)
			http.Error(w, "Failed to toggle like", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *MyServer) UnlikePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(userIDKey).(int)
		if !ok {
			log.Println("User ID not found in context")
			http.Error(w, "User ID not found in context", http.StatusUnauthorized)
			return
		}

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			log.Println("failed to open database", err)
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		tx, err := DB.Begin()
		if err != nil {
			log.Println("Failed to begin transaction", err)
			http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
			return
		}

		err = ToggleUnLike(userID, postID, DB, tx)
		if err != nil {
			log.Println("Failed to toggle unlike", err)
			http.Error(w, "Failed to toggle unlike", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

/*-------------------------------------------------------------------------*/
// ToggleLike gère l'ajout ou la suppression d'un like sur une publication
func ToggleLike(userID int, postID int, DB *sql.DB, tx *sql.Tx) error {
	_, err := DB.Exec(database.Likes_Table)
	if err != nil {
		log.Println("error create database")
		return err
	}

	liked, err := UserLikedPost(userID, postID, tx)
	if err != nil {
		log.Println("Error checking if user liked post:", err)
		return err
	}

	unliked, err := UserUnLikedPost(userID, postID, tx)
	if err != nil {
		log.Println("Error checking if user liked post:", err)
		return err
	}

	if liked {
		err := DeleteLike(userID, postID, tx)
		if err != nil {
			log.Println("Error delete like")
			return err
		}
	} else {
		if unliked {

			err := DeleteUnLike(userID, postID, tx)
			if err != nil {
				log.Println("Error create like")
				return err
			}

		}
		err := CreateLike(userID, postID, tx)
		if err != nil {
			log.Println("Error create like")
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction :", err)
		return err
	}

	return nil
}

func ToggleUnLike(userID int, postID int, DB *sql.DB, tx *sql.Tx) error {
	_, err := DB.Exec(database.Likes_Table)
	if err != nil {
		log.Println("error create database")
		return err
	}

	liked, err := UserLikedPost(userID, postID, tx)
	if err != nil {
		log.Println("Error checking if user liked post:", err)
		return err
	}

	unliked, err := UserUnLikedPost(userID, postID, tx)
	if err != nil {
		log.Println("Error checking if user liked post:", err)
		return err
	}

	if unliked {
		err := DeleteUnLike(userID, postID, tx)
		if err != nil {
			log.Println("Error create like")
			return err
		}
	} else {
		if liked {
			err := DeleteLike(userID, postID, tx)
			if err != nil {
				log.Println("Error create like")
				return err
			}
		}
		err := CreateUnLike(userID, postID, tx)
		if err != nil {
			log.Println("Error create like")
			return err
		}

	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction :", err)
		return err
	}

	return nil
}

/*--------------------------------------------------------------*/

func UserLikedPost(userID, postID int, tx *sql.Tx) (bool, error) {
	// Vérifier si l'utilisateur a aimé la publication en vérifiant s'il existe une entrée dans la table "likes" pour cet utilisateur et cette publication
	var count int
	err := tx.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil // return true, nil
}

func UserUnLikedPost(userID, postID int, tx *sql.Tx) (bool, error) {
	// Vérifier si l'utilisateur a aimé la publication en vérifiant s'il existe une entrée dans la table "likes" pour cet utilisateur et cette publication
	var count int
	err := tx.QueryRow("SELECT COUNT(*) FROM unlikes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil // return true, nil
}

func CreateLike(userID, postID int, tx *sql.Tx) error {
	// Insérer le nouveau like dans la base de données
	_, err := tx.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", userID, postID)
	if err != nil {
		return err
	}

	// Mettre à jour le compteur de likes dans la table posts
	_, err = tx.Exec("UPDATE posts SET total_likes = total_likes + 1 WHERE id = ?", postID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteLike(userID, postID int, tx *sql.Tx) error {
	// Supprimer le like de la base de données
	_, err := tx.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
	if err != nil {
		return err
	}

	// Mettre à jour le compteur de likes dans la table posts
	_, err = tx.Exec("UPDATE posts SET total_likes = total_likes - 1 WHERE id = ?", postID)
	if err != nil {
		return err
	}

	return nil
}

func CreateUnLike(userID, postID int, tx *sql.Tx) error {
	// Insérer le nouveau like dans la base de données
	_, err := tx.Exec("INSERT INTO unlikes (user_id, post_id) VALUES (?, ?)", userID, postID)
	if err != nil {
		return err
	}

	// Mettre à jour le compteur de likes dans la table posts
	_, err = tx.Exec("UPDATE posts SET total_unlikes = total_unlikes + 1 WHERE id = ?", postID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUnLike(userID, postID int, tx *sql.Tx) error {
	// Supprimer le like de la base de données
	_, err := tx.Exec("DELETE FROM unlikes WHERE user_id = ? AND post_id = ?", userID, postID)
	if err != nil {
		return err
	}

	// Mettre à jour le compteur de likes dans la table posts
	_, err = tx.Exec("UPDATE posts SET total_unlikes = total_unlikes - 1 WHERE id = ?", postID)
	if err != nil {
		return err
	}

	return nil
}

/*-------------------------------------------------------------------------------*/
