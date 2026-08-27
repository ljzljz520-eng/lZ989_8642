package api

import (
	"context"
	"customerfollowup/internal/model"
	"customerfollowup/internal/service"
	"encoding/json"
	"net/http"
	"strings"
)

func New(s *service.Service) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/records", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var v model.Record
			if json.NewDecoder(r.Body).Decode(&v) != nil {
				writeErr(w, 400)
				return
			}
			if s.Register(context.Background(), v) != nil {
				writeErr(w, 422)
				return
			}
			w.WriteHeader(201)
			return
		}
		id := r.URL.Query().Get("customer")
		v, e := s.Query(r.Context(), id)
		if e != nil {
			writeErr(w, 500)
			return
		}
		json.NewEncoder(w).Encode(v)
	})
	m.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/tasks/")
		var e error
		if r.Method == http.MethodPost {
			e = s.SubmitIndex(r.Context(), id)
		} else if r.Method == http.MethodDelete {
			e = s.StopIndex(r.Context(), id)
		} else {
			writeErr(w, 405)
			return
		}
		if e != nil {
			writeErr(w, 409)
			return
		}
		w.WriteHeader(204)
	})
	return m
}
func writeErr(w http.ResponseWriter, n int) { http.Error(w, http.StatusText(n), n) }
