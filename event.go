package event

import (
	"github.com/satori/go.uuid"
	"time"
)

type EventHeader struct {
	CreatedAt time.Time
	ID        string
	Type      string
	Version   int
}

func NewEventHeader(eventType string, version int) *EventHeader {
	return &EventHeader{
		CreatedAt: time.Now(),
		ID:        uuid.NewV1().String(),
		Type:      eventType,
		Version:   version,
	}
}
