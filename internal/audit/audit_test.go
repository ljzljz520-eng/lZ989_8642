package audit

import (
	"customerfollowup/internal/store"
	"testing"
)

func TestLogger(t *testing.T) {
	s, e := store.Open(t.TempDir() + "/x")
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	if e = New(s).Log("x", "1", "u", "d"); e != nil {
		t.Fatal(e)
	}
}
