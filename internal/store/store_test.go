package store

import (
	"customerfollowup/internal/model"
	"testing"
)

func TestStoreRoundTrip(t *testing.T) {
	p := t.TempDir() + "/a.db"
	s, e := Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	r := model.Record{ID: "1", CustomerID: "c", Summary: "s"}
	if e = s.SaveRecord(r); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Record("1"); e != nil {
		t.Fatal(e)
	}
}
