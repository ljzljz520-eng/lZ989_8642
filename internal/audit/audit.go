package audit

import (
	"customerfollowup/internal/model"
	"customerfollowup/internal/store"
	"encoding/json"
	"fmt"
	"time"
)

type Logger struct{ s *store.Store }

func New(s *store.Store) *Logger { return &Logger{s: s} }
func (l *Logger) Log(action, id, actor, detail string) error {
	return l.s.SaveAudit(model.Audit{ID: fmt.Sprintf("%d", time.Now().UnixNano()), Action: action, EntityID: id, Actor: actor, Detail: detail, At: time.Now()})
}
func (l *Logger) History() ([]model.Audit, error) {
	raw, e := l.s.List("audits")
	if e != nil {
		return nil, e
	}
	out := make([]model.Audit, 0, len(raw))
	for _, b := range raw {
		var a model.Audit
		if jsonErr := unmarshal(b, &a); jsonErr != nil {
			return nil, jsonErr
		}
		out = append(out, a)
	}
	return out, nil
}
func unmarshal(b []byte, v any) error { return json.Unmarshal(b, v) }
