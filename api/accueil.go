package api

import "net/http"

func (s *MyServer) AccueilHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontEnd/index.html")
	}
}
