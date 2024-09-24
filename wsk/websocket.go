package wsk

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upGradeWebsocket = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     CheckOrigin,
}

func CheckOrigin(r *http.Request) bool {
	log.Printf("%s %s %s %v", r.Method, r.Host, r.RequestURI, r.Proto)
	return r.Method == http.MethodGet
}

type WebsocketChat struct {
	Users          map[string]*UserChat
	JoinChannel    userChannel
	LeaveChannel   userChannel
	MessageChannel messageChannel
	MessageHistory map[string][]*Message
	Mu             sync.Mutex
}

func NewWebsocketChat() *WebsocketChat {
	w := &WebsocketChat{
		Users:          make(map[string]*UserChat),
		JoinChannel:    make(userChannel),
		LeaveChannel:   make(userChannel),
		MessageChannel: make(messageChannel),
		MessageHistory: make(map[string][]*Message),
	}
	go w.UsersChatManager()
	return w
}

func (w *WebsocketChat) UsersChatManager() {
	for {
		select {
		case user := <-w.JoinChannel:
			w.Mu.Lock()
			w.Users[user.Username] = user
			w.sendHistory(user)
			w.Mu.Unlock()

		case user := <-w.LeaveChannel:
			w.Mu.Lock()
			delete(w.Users, user.Username)
			w.Mu.Unlock()

		case msg := <-w.MessageChannel:
			w.Mu.Lock()

			if msg.Type == "typing" {
				// Gestion des événements "typing"
				if targetUser, ok := w.Users[msg.TargetUsername]; ok {
					log.Printf("%s is typing...", msg.SenderUsername)
					err := targetUser.Connection.WriteJSON(msg)
					if err != nil {
						log.Printf("Error sending typing status to %s: %v", targetUser.Username, err)
					}
				}
			} else if msg.Type == "message" && msg.Content != "" {
				// Gestion des messages réels
				log.Printf("Received message: %+v", msg)

				// Envoi du message à tous si c'est un message global
				if msg.TargetUsername == "all" {
					for _, user := range w.Users {
						err := user.Connection.WriteJSON(msg)
						if err != nil {
							log.Printf("Error sending message to %s: %v", user.Username, err)
							user.Connection.Close()
							delete(w.Users, user.Username)
						}
					}
				} else {
					// Envoi du message au destinataire
					if targetUser, ok := w.Users[msg.TargetUsername]; ok {
						log.Printf("Sending message from %s to %s", msg.SenderUsername, msg.TargetUsername)
						err := targetUser.Connection.WriteJSON(msg)
						if err != nil {
							log.Printf("Error sending message to %s: %v", targetUser.Username, err)
						}
					}

					// Envoi de la confirmation à l'expéditeur
					if senderUser, ok := w.Users[msg.SenderUsername]; ok {
						err := senderUser.Connection.WriteJSON(msg)
						if err != nil {
							log.Printf("Error sending message to %s: %v", senderUser.Username, err)
						}
					}
				}

				// Sauvegarde du message dans l'historique des deux utilisateurs (expéditeur et destinataire).
				w.MessageHistory[msg.SenderUsername] = append(w.MessageHistory[msg.SenderUsername], msg)
				w.MessageHistory[msg.TargetUsername] = append(w.MessageHistory[msg.TargetUsername], msg)
			}

			w.Mu.Unlock()
		}
	}
}

// Fonction pour envoyer l'historique des messages à l'utilisateur nouvellement connecté.
func (w *WebsocketChat) sendHistory(user *UserChat) {
	if messages, ok := w.MessageHistory[user.Username]; ok {
		for _, msg := range messages {
			user.Connection.WriteJSON(msg) // Envoi de chaque message de l'historique.
		}
	}
}
