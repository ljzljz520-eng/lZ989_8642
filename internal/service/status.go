package service

import (
	"context"
	"customerfollowup/internal/model"
	"fmt"
)

func (x *Service) Status(ctx context.Context, id string) (model.IndexTask, error) {
	return x.s.Task(id)
}
func (x *Service) EnsureProfile(ctx context.Context, p model.Profile) error {
	if p.ID == "" {
		return errProfile()
	}
	return x.s.SaveProfile(p)
}
func errProfile() error { return fmt.Errorf("profile id required") }
