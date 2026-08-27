package store

import (
	"customerfollowup/internal/model"
	"time"
)

func (s *Store) SeedProfile(id, name string) error {
	return s.SaveProfile(model.Profile{ID: id, Name: name, Active: true, LastContact: time.Now()})
}
func (s *Store) RemoveRecord(id string) error { return s.Delete("records", id) }
func (s *Store) Healthy() bool                { return s.Count("records") >= 0 }
