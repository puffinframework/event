package event_test

import (
	"testing"

	"github.com/puffinframework/event"
)

func TestLeveldbEventStore(t *testing.T) {
	store := event.NewLeveldbStore("test-event-store")
	testEventStore(t, store)
}
