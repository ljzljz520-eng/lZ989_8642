package store

import (
	"customerfollowup/internal/model"
	"encoding/json"
	"errors"
	"go.etcd.io/bbolt"
	"strings"
)

func decodeRecords(raw [][]byte) []model.Record {
	out := []model.Record{}
	for _, b := range raw {
		var r model.Record
		if json.Unmarshal(b, &r) == nil {
			out = append(out, r)
		}
	}
	return out
}
func (s *Store) SearchSummary(term string) ([]model.Record, error) {
	raw, e := s.List("records")
	if e != nil {
		return nil, e
	}
	term = strings.ToLower(strings.TrimSpace(term))
	out := []model.Record{}
	for _, r := range decodeRecords(raw) {
		if term == "" || strings.Contains(strings.ToLower(r.Summary), term) || strings.Contains(strings.ToLower(r.CustomerID), term) {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Store) RecordsByStatus(status string) ([]model.Record, error) {
	raw, e := s.List("records")
	if e != nil {
		return nil, e
	}
	out := []model.Record{}
	for _, r := range decodeRecords(raw) {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Store) TagCounts() map[string]int {
	m := map[string]int{}
	raw, e := s.List("records")
	if e != nil {
		return m
	}
	for _, r := range decodeRecords(raw) {
		for _, t := range model.NormalizeTags(r.Tags) {
			m[t]++
		}
	}
	return m
}
func (s *Store) TaskStates() map[string]int {
	m := map[string]int{}
	raw, e := s.List("tasks")
	if e != nil {
		return m
	}
	for _, b := range raw {
		var t model.IndexTask
		if json.Unmarshal(b, &t) == nil {
			m[t.State]++
		}
	}
	return m
}
func (s *Store) Has(bucket, key string) bool { _, e := s.getBytes(bucket, key); return e == nil }
func (s *Store) getBytes(bucket, key string) ([]byte, error) {
	var v []byte
	e := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket)).Get([]byte(key))
		if b == nil {
			return errors.New("not found")
		}
		v = append([]byte(nil), b...)
		return nil
	})
	return v, e
}
