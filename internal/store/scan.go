package store

import (
	"customerfollowup/internal/model"
	"encoding/json"
)

func (s *Store) Records() ([]model.Record, error) {
	raw, e := s.List("records")
	if e != nil {
		return nil, e
	}
	out := []model.Record{}
	for _, b := range raw {
		var r model.Record
		if json.Unmarshal(b, &r) == nil {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Store) Profiles() ([]model.Profile, error) {
	raw, e := s.List("profiles")
	if e != nil {
		return nil, e
	}
	out := []model.Profile{}
	for _, b := range raw {
		var p model.Profile
		if json.Unmarshal(b, &p) == nil {
			out = append(out, p)
		}
	}
	return out, nil
}
