package model

import "testing"

func TestNormalizeTags(t *testing.T) {
	v := NormalizeTags([]string{" A ", "a", ""})
	if len(v) != 1 || v[0] != "a" {
		t.Fatal(v)
	}
}
