package service

import (
	"context"
	"customerfollowup/internal/model"
	"errors"
	"time"
)

func (x *Service) TouchProfile(ctx context.Context, id string) error {
	p, e := x.s.Profile(id)
	if e != nil {
		return e
	}
	p.LastContact = now()
	return x.s.SaveProfile(p)
}
func (x *Service) AddEvent(ctx context.Context, e model.Event) error {
	if e.ID == "" || e.RecordID == "" {
		return fmtEvent()
	}
	return x.s.SaveEvent(e)
}
func now() time.Time  { return time.Now() }
func fmtEvent() error { return errors.New("event identity required") }
