package event

import (
	"time"
)

type Event interface {
	Header() Header
	Data() Data
}

type Header struct {
	UUID      string
	Timestamp time.Time
	Type      string
	Version   int
}

type Data interface{}
