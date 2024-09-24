package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"realtime/database"
	"realtime/wsk"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

const (
	ColorGreen = "\033[32m"
	ColorBlue  = "\033[34m"
	ColorReset = "\033[0m"
	port       = ":8079"
)

type MyServer struct {
	Store             database.Store
	Router            *http.ServeMux
	Server            *http.Server
	WebSocketChat     *wsk.WebsocketChat
	GoogleOAuthConfig *oauth2.Config
	GitHubOAuthConfig *oauth2.Config
}

func NewServer(store database.Store, wsChat *wsk.WebsocketChat) *MyServer {
	router := http.NewServeMux()
	server := &MyServer{
		Store:         store,
		Router:        router,
		WebSocketChat: wsChat,
		GoogleOAuthConfig: &oauth2.Config{
			ClientID:     "your-google-client-id",
			ClientSecret: "your-google-client-secret",
			RedirectURL:  "http://localhost:8079/auth/google/callback",
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		}, GitHubOAuthConfig: &oauth2.Config{
			ClientID:     "your-github-client-id",
			ClientSecret: "your-github-client-secret",
			RedirectURL:  "http://localhost:8079/auth/github/callback",
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
	}

	server.routes()

	fmt.Println(ColorBlue, "(http://localhost:8079) - Server started on port", port, ColorReset)
	fmt.Println(ColorGreen, "[SERVER_INFO] : To stop the server : Ctrl + c", ColorReset)

	srv := &http.Server{
		Addr:              "localhost:8079",
		Handler:           router,
		ReadHeaderTimeout: 15 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	server.Server = srv

	return server
}

func (s *MyServer) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

func LogRequestMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%v], %v", r.Method, r.RequestURI)
		next(w, r)
	}
}

// function Chain pour empiler les middlewares
func Chain(f http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middlewares := range middlewares {
		f = middlewares(f)
	}
	return f
}
