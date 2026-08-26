package workflow

import (
	"context"
	"customerfollowup/internal/model"
	"customerfollowup/internal/service"
	"fmt"
)

func ValidateChain(ctx context.Context, s *service.Service, id string) error {
	r, e := s.Query(ctx, "")
	if e != nil {
		return e
	}
	for _, v := range r {
		if v.ID == id && model.IsTerminal(v.Status) {
			return fmt.Errorf("terminal record")
		}
	}
	return nil
}
func Describe(status string) string {
	if status == model.StateArchived {
		return "sealed"
	}
	if status == "processed" {
		return "handled"
	}
	return "open"
}
