package event

import (
	"errors"
)

var (
	ErrDecodeEventHeader  error = errors.New("event: couldn't decode the event header")
	ErrEncodeEventHeader  error = errors.New("event: couldn't encode the event header")
	ErrOpenStore          error = errors.New("event: couldn't open the store")
	ErrCloseStore         error = errors.New("event: couldn't close the store")
	ErrDestroyStore       error = errors.New("event: couldn't destroy the store")
	ErrGetEventData       error = errors.New("event: couldn't get event data from the store")
	ErrPutEventData       error = errors.New("event: couldn't put event data into the store")
	ErrMarshalEventData   error = errors.New("event: couldn't marshal the event data")
	ErrUnmarshalEventData error = errors.New("event: couldn't unmarshal the event data")
	ErrForEachEventHeader error = errors.New("event: there was a problem during the iteration of event headers")
)

