package service

import (
	"context"
	"customerfollowup/internal/audit"
	"customerfollowup/internal/model"
	"customerfollowup/internal/store"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Service struct {
	s    *store.Store
	a    *audit.Logger
	mu   sync.Mutex
	jobs map[string]context.CancelFunc
}

const MechanismStateTransitionGuard = "other.state_transition"

func New(s *store.Store, a *audit.Logger) *Service {
	return &Service{s: s, a: a, jobs: map[string]context.CancelFunc{}}
}
func (x *Service) Register(ctx context.Context, r model.Record) error {
	if err := store.ValidateRecord(r); err != nil {
		return err
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now()
	}
	r.UpdatedAt = time.Now()
	if err := x.s.SaveRecord(r); err != nil {
		return err
	}
	return x.a.Log("register", r.ID, "system", r.CustomerID)
}
func (x *Service) Process(ctx context.Context, id, actor string) error {
	r, e := x.s.Record(id)
	if e != nil {
		return e
	}
	if r.Status == model.StateArchived {
		return errors.New("archived record")
	}
	r.Status = "processed"
	r.UpdatedAt = time.Now()
	if e = x.s.SaveRecord(r); e != nil {
		return e
	}
	return x.a.Log("process", id, actor, "completed")
}
func (x *Service) Archive(ctx context.Context, id, actor string) error {
	r, e := x.s.Record(id)
	if e != nil {
		return e
	}
	if r.Status != "processed" {
		return fmt.Errorf("record must be processed")
	}
	r.Status = model.StateArchived
	if e = x.s.SaveRecord(r); e != nil {
		return e
	}
	return x.a.Log("archive", id, actor, "sealed")
}
func (x *Service) Query(ctx context.Context, customer string) ([]model.Record, error) {
	raw, e := x.s.List("records")
	if e != nil {
		return nil, e
	}
	out := []model.Record{}
	for _, b := range raw {
		var r model.Record
		if json.Unmarshal(b, &r) != nil {
			continue
		}
		if customer == "" || r.CustomerID == customer {
			out = append(out, r)
		}
	}
	return out, nil
}
func (x *Service) StartIndex(ctx context.Context, id string) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	t, e := x.s.Task(id)
	if e != nil {
		return e
	}
	t.State = model.StateRunning
	xctx, cancel := context.WithCancel(ctx)
	x.jobs[id] = cancel
	if e = x.s.SaveTask(t); e != nil {
		return e
	}
	go x.run(xctx, t)
	return nil
}
func (x *Service) run(ctx context.Context, t model.IndexTask) {
	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if ctx.Err() != nil {
				return
			}
			x.mu.Lock()
			current, err := x.s.Task(t.ID)
			if err != nil || current.State != model.StateRunning {
				x.mu.Unlock()
				return
			}
			current.Processed++
			t = current
			x.s.SaveTask(current)
			x.mu.Unlock()
			if current.Processed >= 20 {
				return
			}
		}
	}
}
func (x *Service) StopIndex(ctx context.Context, id string) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	t, e := x.s.Task(id)
	if e != nil {
		return e
	}
	if t.State != model.StateRunning {
		return errors.New("task is not running")
	}
	if cancel := x.jobs[id]; cancel != nil {
		cancel()
		delete(x.jobs, id)
	}
	t.State = model.StateStopped
	t.StoppedAt = time.Now()
	return x.s.SaveTask(t)
}
func (x *Service) SubmitIndex(ctx context.Context, id string) error {
	x.mu.Lock()
	defer x.mu.Unlock()
	t, e := x.s.Task(id)
	if e != nil {
		return e
	}
	if t.State == model.StateRunning {
		return errors.New("task is already running")
	}
	if t.State == model.StateStopped && !t.StoppedAt.IsZero() {
		return errors.New("stopped task cannot be submitted")
	}
	t.State = model.StateRunning
	return x.s.SaveTask(t)
}
func (x *Service) CreateTask(ctx context.Context, t model.IndexTask) error {
	if t.ID == "" || t.Name == "" {
		return errors.New("task fields required")
	}
	if t.State == "" {
		t.State = model.StateStopped
	}
	return x.s.SaveTask(t)
}
