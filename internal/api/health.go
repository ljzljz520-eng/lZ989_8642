package api

import (
	"customerfollowup/internal/store"
	"net/http"
)

func Health(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.Healthy() {
			w.WriteHeader(200)
			w.Write([]byte("ok"))
			return
		}
		w.WriteHeader(503)
	}
}
