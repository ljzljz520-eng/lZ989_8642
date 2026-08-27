package workflow

import (
	"context"
	"customerfollowup/internal/model"
	"customerfollowup/internal/service"
	"fmt"
)

func RegisterAndShow(ctx context.Context, s *service.Service, r model.Record) (model.Record, error) {
	if err := s.Register(ctx, r); err != nil {
		return model.Record{}, err
	}
	rows, err := s.Query(ctx, r.CustomerID)
	if err != nil {
		return model.Record{}, err
	}
	if len(rows) == 0 {
		return model.Record{}, fmt.Errorf("registration not visible")
	}
	return rows[len(rows)-1], nil
}
func ProcessAndArchive(ctx context.Context, s *service.Service, id string) error {
	if err := s.Process(ctx, id, "workflow"); err != nil {
		return err
	}
	return s.Archive(ctx, id, "workflow")
}
func SubmitAndTrack(ctx context.Context, s *service.Service, id string) error {
	if err := s.SubmitIndex(ctx, id); err != nil {
		return err
	}
	t, err := s.Status(ctx, id)
	if err != nil {
		return err
	}
	if t.State != "running" {
		return fmt.Errorf("task not running")
	}
	return nil
}
