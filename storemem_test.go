package event_test

import (
	"testing"

	"github.com/puffinframework/event"
)

func TestMemEventStore(t *testing.T) {
	store := event.NewMemStore()
	testEventStore(t, store)
}
