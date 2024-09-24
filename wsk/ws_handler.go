package wsk

import (
	"log"
	"net/http"
	"time"
)

func (w *WebsocketChat) HanderUsersConnection(wr http.ResponseWriter, r *http.Request) {
	// Établissement de la connexion WebSocket avec le client.
	conn, err := upGradeWebsocket.Upgrade(wr, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}

	// Récupération du nom d'utilisateur à partir du cookie.
	var username string
	cookie, err := r.Cookie("username")
	if err != nil {
		log.Printf("Failed to get username cookie: %v", err)
		return
	}
	username = cookie.Value

	// Log pour indiquer que l'utilisateur est connecté.
	log.Printf("User %s connected", username)

	// Création d'une nouvelle instance de chat pour l'utilisateur connecté.
	userChat := NewUserChat(&Channel{
		messageChannel: w.MessageChannel,
		leaveChannel:   w.LeaveChannel,
	}, username, conn)

	// Ajout de l'utilisateur au canal de discussion.
	w.JoinChannel <- userChat

	// Lancement d'une goroutine pour écouter les messages de l'utilisateur.
	go userChat.listenForMessages()
}

func (u *UserChat) listenForMessages() {
	defer func() {
		// Gestion de la déconnexion de l'utilisateur et fermeture de la connexion WebSocket.
		u.channels.leaveChannel <- u
		u.Connection.Close()
		log.Printf("User %s disconnected", u.Username)
	}()

	// Boucle pour lire les messages envoyés par l'utilisateur.
	for {
		var msg Message
		err := u.Connection.ReadJSON(&msg) // Lecture d'un message JSON.
		if err != nil {
			log.Printf("Error reading json from user %s: %v", u.Username, err)
			break
		}

		// Ajout du nom d'utilisateur et de l'horodatage au message.
		msg.SenderUsername = u.Username
		msg.Timestamp = time.Now()

		// Gestion des différents types de messages.
		switch msg.Type {
		case "typing":
			// Message de type 'typing' (indicateur de saisie).
			log.Printf("%s is typing...", msg.SenderUsername)
			u.channels.messageChannel <- &msg
		case "message":
			// Vérification si le message est vide, sinon, l'envoyer.
			if msg.Content == "" {
				log.Printf("Empty message from %s ignored", msg.SenderUsername)
				continue
			}
			log.Printf("Message to send: %+v", msg)
			u.channels.messageChannel <- &msg
		default:
			// Message de type inconnu.
			log.Printf("Unknown message type: %s", msg.Type)
		}
	}
}
