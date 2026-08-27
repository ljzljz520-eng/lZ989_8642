package internal_test

import (
	"customerfollowup/internal/model"
	"customerfollowup/internal/store"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	p := t.TempDir() + "/persist.db"
	s, e := store.Open(p)
	if e != nil {
		t.Fatal(e)
	}
	if e = s.SaveRecord(model.Record{ID: "persist", CustomerID: "c", Summary: "saved"}); e != nil {
		t.Fatal(e)
	}
	s.Close()
	s, e = store.Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	if _, e = s.Record("persist"); e != nil {
		t.Fatal(e)
	}
}
