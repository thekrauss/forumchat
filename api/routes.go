package api

import (
	"fmt"
	"net/http"
)

func (s *MyServer) routes() {
	s.Router.Handle("/frontEnd/", http.StripPrefix("/frontEnd/", http.FileServer(http.Dir("frontEnd"))))

	s.Router.HandleFunc("/", Chain(s.AccueilHandler(), LogRequestMiddleware))
	s.Router.HandleFunc("/login-form", Chain(s.LoginHandler(), LogRequestMiddleware))
	s.Router.HandleFunc("/register-form", Chain(s.RegisterHandler(), LogRequestMiddleware))
	s.Router.HandleFunc("/nav-logout", Chain(s.LogoutHandler(), LogRequestMiddleware))

	s.Router.HandleFunc("/protected", Chain(s.ProtectedHandler(), LogRequestMiddleware, s.Authenticate))

	s.Router.HandleFunc("/auth/google/login-form", Chain(s.GoogleLogin(), LogRequestMiddleware))
	s.Router.HandleFunc("/auth/google/callback", Chain(s.GoogleCallBackHandler(), LogRequestMiddleware))
	s.Router.HandleFunc("/auth/google/register-form", Chain(s.GoogleRegisterHandler(), LogRequestMiddleware))
	s.Router.HandleFunc("/auth/google/callback/register", Chain(s.GoogleCallbackRegisterHandler(), LogRequestMiddleware))

	s.Router.HandleFunc("/auth/github/login-form", Chain(s.GitHubLoginHandler(), LogRequestMiddleware))
	s.Router.HandleFunc("/auth/github/callback", Chain(s.GitHubCallbackHandler(), LogRequestMiddleware))
	s.Router.HandleFunc("/auth/github/register-form", Chain(s.GitHubRegisterHandler(), LogRequestMiddleware))
	s.Router.HandleFunc("/auth/github/callback/register", Chain(s.GitHubCallbackRegisterHandler(), LogRequestMiddleware))

	s.Router.HandleFunc("/create-post-form", Chain(s.CreatePostHandler(), LogRequestMiddleware, s.Authenticate))
	s.Router.HandleFunc("/list-posts-form", Chain(s.ListPostsHandler(), LogRequestMiddleware))
	s.Router.HandleFunc("/create-comment-form", Chain(s.CreateCommentHandler(), LogRequestMiddleware, s.Authenticate))
	s.Router.HandleFunc("/list-comment-form", Chain(s.ListCommentHandler(), LogRequestMiddleware))

	s.Router.HandleFunc("/user-list", Chain(s.OnlineUsersHandler(), LogRequestMiddleware))
	s.Router.HandleFunc("/message", Chain(s.MessagesHandler(), LogRequestMiddleware, s.Authenticate))
	s.Router.HandleFunc("/ws", s.WebSocketChat.HanderUsersConnection)
}

func (s *MyServer) ProtectedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)
		w.Write([]byte(fmt.Sprintf("Hello, user %d", userID)))
	}
}
