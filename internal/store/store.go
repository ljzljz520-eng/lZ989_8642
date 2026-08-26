package store

import (
	"customerfollowup/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"sync"
	"time"
)

var buckets = [][]byte{[]byte("records"), []byte("profiles"), []byte("events"), []byte("audits"), []byte("tasks")}

type Store struct {
	db *bbolt.DB
	mu sync.RWMutex
}

func Open(path string) (*Store, error) {
	db, e := bbolt.Open(path, 0600, &bbolt.Options{Timeout: time.Second})
	if e != nil {
		return nil, e
	}
	s := &Store{db: db}
	e = db.Update(func(tx *bbolt.Tx) error {
		for _, b := range buckets {
			if _, x := tx.CreateBucketIfNotExists(b); x != nil {
				return x
			}
		}
		return nil
	})
	if e != nil {
		db.Close()
		return nil, e
	}
	return s, nil
}
func (s *Store) Close() error      { return s.db.Close() }
func encode(v any) ([]byte, error) { return json.Marshal(v) }
func (s *Store) put(bucket, key string, v any) error {
	raw, e := encode(v)
	if e != nil {
		return e
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(bucket)).Put([]byte(key), raw) })
}
func (s *Store) get(bucket, key string, v any) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		val := tx.Bucket([]byte(bucket)).Get([]byte(key))
		if val == nil {
			return errors.New("not found")
		}
		return json.Unmarshal(val, v)
	})
}
func (s *Store) Delete(bucket, key string) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(bucket)).Delete([]byte(key)) })
}
func (s *Store) SaveRecord(v model.Record) error { return s.put("records", v.ID, v) }
func (s *Store) Record(id string) (model.Record, error) {
	var v model.Record
	e := s.get("records", id, &v)
	return v, e
}
func (s *Store) SaveProfile(v model.Profile) error { return s.put("profiles", v.ID, v) }
func (s *Store) Profile(id string) (model.Profile, error) {
	var v model.Profile
	e := s.get("profiles", id, &v)
	return v, e
}
func (s *Store) SaveEvent(v model.Event) error    { return s.put("events", v.ID, v) }
func (s *Store) SaveAudit(v model.Audit) error    { return s.put("audits", v.ID, v) }
func (s *Store) SaveTask(v model.IndexTask) error { return s.put("tasks", v.ID, v) }
func (s *Store) Task(id string) (model.IndexTask, error) {
	var v model.IndexTask
	e := s.get("tasks", id, &v)
	return v, e
}
func (s *Store) List(bucket string) ([][]byte, error) {
	var out [][]byte
	e := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(bucket)).ForEach(func(k, v []byte) error {
			if v != nil {
				out = append(out, append([]byte(nil), v...))
			}
			return nil
		})
	})
	return out, e
}
func (s *Store) Count(bucket string) int {
	n := 0
	s.db.View(func(tx *bbolt.Tx) error { n = tx.Bucket([]byte(bucket)).Stats().KeyN; return nil })
	return n
}
func ValidateRecord(r model.Record) error {
	if r.ID == "" || r.CustomerID == "" {
		return fmt.Errorf("record identity required")
	}
	if r.Summary == "" {
		return fmt.Errorf("summary required")
	}
	return nil
}
