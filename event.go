package event

import (
	"time"
)

type Header struct {
	UUID      string
	Timestamp time.Time
	Type      string
	Version   int
}

type Data interface{}
