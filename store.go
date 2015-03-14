package event

import (
	"time"
)

type Store interface {
	ForEachEventHeader(since time.Time, callback func(header Header) (bool, error)) error
	MustLoadEventData(header Header, data interface{})
	MustSaveEventData(header Header, data interface{})
	MustClose()
	MustDestroy()
}

