package analytics

import (
	"customerfollowup/internal/store"
	"encoding/json"
)

func JSON(s *store.Store) ([]byte, error) {
	v, e := Build(s)
	if e != nil {
		return nil, e
	}
	return json.Marshal(v)
}
