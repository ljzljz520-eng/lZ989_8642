package service

import (
	"context"
	"customerfollowup/internal/audit"
	"customerfollowup/internal/model"
	"customerfollowup/internal/store"
	"testing"
)

func TestServiceValidation(t *testing.T) {
	s, e := store.Open(t.TempDir() + "/x")
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	x := New(s, audit.New(s))
	if e = x.Register(context.Background(), model.Record{}); e == nil {
		t.Fatal("expected validation")
	}
}
