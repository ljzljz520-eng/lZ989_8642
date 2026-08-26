package main

import (
	"context"
	"customerfollowup/internal/audit"
	"customerfollowup/internal/model"
	"customerfollowup/internal/service"
	"customerfollowup/internal/store"
	"customerfollowup/internal/workflow"
	"testing"
)

func setup(t *testing.T) (*service.Service, func()) {
	p := t.TempDir() + "/x.db"
	s, e := store.Open(p)
	if e != nil {
		t.Fatal(e)
	}
	return service.New(s, audit.New(s)), func() { s.Close() }
}
func TestWorkflowOne(t *testing.T) {
	s, done := setup(t)
	defer done()
	r, e := workflow.RegisterAndShow(context.Background(), s, model.Record{ID: "r1", CustomerID: "c1", Summary: "call"})
	if e != nil || r.ID != "r1" {
		t.Fatal(e, r)
	}
}
func TestWorkflowTwo(t *testing.T) {
	s, done := setup(t)
	defer done()
	s.Register(context.Background(), model.Record{ID: "r2", CustomerID: "c2", Summary: "call"})
	if e := s.Process(context.Background(), "r2", "u"); e != nil {
		t.Fatal(e)
	}
	if e := s.Archive(context.Background(), "r2", "u"); e != nil {
		t.Fatal(e)
	}
}
func TestWorkflowThree(t *testing.T) {
	s, done := setup(t)
	defer done()
	s.CreateTask(context.Background(), model.IndexTask{ID: "t1", Name: "search"})
	if e := s.SubmitIndex(context.Background(), "t1"); e != nil {
		t.Fatal(e)
	}
}
func TestBusinessChain34(t *testing.T) {
	s, done := setup(t)
	defer done()
	s.CreateTask(context.Background(), model.IndexTask{ID: "t2", Name: "search"})
	if e := s.SubmitIndex(context.Background(), "t2"); e != nil {
		t.Fatal(e)
	}
	if e := s.StopIndex(context.Background(), "t2"); e != nil {
		t.Fatal(e)
	}
	before, e := s.Status(context.Background(), "t2")
	if e != nil {
		t.Fatal(e)
	}
	if e := s.SubmitIndex(context.Background(), "t2"); e == nil {
		t.Fatal("stopped task must reject submission")
	}
	after, e := s.Status(context.Background(), "t2")
	if e != nil {
		t.Fatal(e)
	}
	if after.State != model.StateStopped || after.Processed != before.Processed {
		t.Fatalf("state=%s processed=%d", after.State, after.Processed)
	}
}
