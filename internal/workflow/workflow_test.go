package workflow

import (
	"context"
	"customerfollowup/internal/audit"
	"customerfollowup/internal/model"
	"customerfollowup/internal/service"
	"customerfollowup/internal/store"
	"testing"
)

func TestRegisterAndShow(t *testing.T) {
	s, e := store.Open(t.TempDir() + "/x")
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	x := service.New(s, audit.New(s))
	if _, e = RegisterAndShow(context.Background(), x, model.Record{ID: "1", CustomerID: "c", Summary: "s"}); e != nil {
		t.Fatal(e)
	}
}
