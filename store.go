package event

import (
	"time"
)

type Store interface {
	ForEachEventHeader(since time.Time, callback func(header Header) (bool, error)) error
	MustLoadEvent(header Header, data interface{})
	MustSaveEvent(header Header, data interface{})
	MustClose()
	MustDestroy()
}
