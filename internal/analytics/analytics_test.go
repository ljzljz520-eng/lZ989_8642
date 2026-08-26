package analytics

import (
	"customerfollowup/internal/store"
	"testing"
)

func TestSummary(t *testing.T) {
	s, e := store.Open(t.TempDir() + "/x")
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	v, e := Build(s)
	if e != nil || v.Total != 0 {
		t.Fatal(e, v)
	}
}
