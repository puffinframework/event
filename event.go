package event

import (
	"github.com/satori/go.uuid"
	"time"
)

type Event interface {
	Header() Header
	Data() Data
}

type Header struct {
	UUID      string
	Timestamp time.Time
	Type      Type
	Version   int
}

type Type string

type Data interface{}

type eventImpl struct {
	header Header
	data   Data
}

func NewEvent(eventType Type, version int, data interface{}) Event {
	return eventImpl{header: Header{UUID: uuid.NewV1().String(), Timestamp: time.Now(), Type: eventType, Version: version}, data: data}
}

func (self eventImpl) Header() Header {
	return self.header
}

func (self eventImpl) Data() Data {
	return self.data
}
