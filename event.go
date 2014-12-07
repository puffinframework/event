package event

import (
	"github.com/satori/go.uuid"
	"time"
)

type Header struct {
	CreatedAt time.Time
	ID        string
	Type      string
	Version   int
}

func NewHeader(eventType string, version int) *Header {
	return &Header{
		CreatedAt: time.Now(),
		ID:        uuid.NewV1().String(),
		Type:      eventType,
		Version:   version,
	}
}

type Store interface {
	ForEach(since time.Time, handler func(header *Header))
	MustLoad(header Header, data interface{})
	MustSave(header Header, data interface{})
}
