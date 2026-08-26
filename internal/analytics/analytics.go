package analytics

import (
	"customerfollowup/internal/model"
	"customerfollowup/internal/store"
	"encoding/json"
)

type Summary struct{ Total, Processed, Archived int }

func Build(s *store.Store) (Summary, error) {
	raw, e := s.List("records")
	if e != nil {
		return Summary{}, e
	}
	out := Summary{Total: len(raw)}
	for _, b := range raw {
		var r model.Record
		if json.Unmarshal(b, &r) != nil {
			continue
		}
		switch r.Status {
		case "processed":
			out.Processed++
		case model.StateArchived:
			out.Archived++
		}
	}
	return out, nil
}
func CompletionRate(v Summary) float64 {
	if v.Total == 0 {
		return 0
	}
	return float64(v.Archived) / float64(v.Total)
}
