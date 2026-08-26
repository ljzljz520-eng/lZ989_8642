package notify

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Message struct {
	To, Body string
	SentAt   time.Time
}
type Sender struct {
	mu  sync.Mutex
	out []Message
}

func New() *Sender { return &Sender{} }
func (s *Sender) Send(ctx context.Context, to, body string) error {
	if to == "" || body == "" {
		return fmt.Errorf("recipient and body required")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	s.out = append(s.out, Message{to, body, time.Now()})
	s.mu.Unlock()
	return nil
}
func (s *Sender) Messages() []Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Message(nil), s.out...)
}
