package workflow

import (
	"context"
	"customerfollowup/internal/model"
	"customerfollowup/internal/service"
	"fmt"
)

type Step struct {
	Code, Label string
	Required    bool
}
type Catalog struct{ Steps []Step }

func DefaultCatalog() Catalog {
	return Catalog{Steps: []Step{{"receive", "Receive", true}, {"validate", "Validate", true}, {"save", "Save", true}, {"show", "Show", false}}}
}
func (c Catalog) RequiredCount() int {
	n := 0
	for _, s := range c.Steps {
		if s.Required {
			n++
		}
	}
	return n
}
func (c Catalog) Has(code string) bool {
	for _, s := range c.Steps {
		if s.Code == code {
			return true
		}
	}
	return false
}
func (c Catalog) Labels() []string {
	out := []string{}
	for _, s := range c.Steps {
		out = append(out, s.Label)
	}
	return out
}
func Run(ctx context.Context, s *service.Service, r model.Record, c Catalog) (model.Record, error) {
	if !c.Has("receive") || !c.Has("save") {
		return model.Record{}, fmt.Errorf("incomplete catalog")
	}
	return RegisterAndShow(ctx, s, r)
}
func CanArchive(r model.Record) bool { return r.Status == "processed" }
func CanProcess(r model.Record) bool { return r.Status != "archived" }
func NextStatus(s string) string {
	switch s {
	case "":
		return "processed"
	case "processed":
		return model.StateArchived
	default:
		return s
	}
}
func ChainLength(c Catalog) int { return len(c.Steps) }
func RequiredLabels(c Catalog) []string {
	out := []string{}
	for _, s := range c.Steps {
		if s.Required {
			out = append(out, s.Label)
		}
	}
	return out
}
