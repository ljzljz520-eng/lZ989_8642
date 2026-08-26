package model

import "time"

type Record struct {
	ID, CustomerID, Summary, Status string
	CreatedAt, UpdatedAt            time.Time
	Tags                            []string
}
type Profile struct {
	ID, Name, Email, Segment string
	Active                   bool
	LastContact              time.Time
}
type Event struct {
	ID, RecordID, Kind, Actor, Payload string
	At                                 time.Time
}
type Audit struct {
	ID, Action, EntityID, Actor, Detail string
	At                                  time.Time
}
type IndexTask struct {
	ID, Name, State        string
	SubmittedAt, StoppedAt time.Time
	Processed              int
}

const (
	StateRunning  = "running"
	StateStopped  = "stopped"
	StateArchived = "archived"
)
