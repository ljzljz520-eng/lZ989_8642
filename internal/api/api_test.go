package api

import (
	"customerfollowup/internal/audit"
	"customerfollowup/internal/service"
	"customerfollowup/internal/store"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	s, e := store.Open(t.TempDir() + "/x")
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	h := Health(s)
	r := httptest.NewRecorder()
	h(r, httptest.NewRequest("GET", "/health", nil))
	if r.Code != 200 {
		t.Fatal(r.Code)
	}
	_ = New(service.New(s, audit.New(s)))
}
