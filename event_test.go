package event_test

import (
	. "github.com/puffinframework/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	evt := NewEvent("SomeEvent", 3, "someData")
	assert.Equal(t, "SomeEvent", evt.Header().Type)
	assert.Equal(t, 3, evt.Header().Version)
	assert.Equal(t, "someData", evt.Data())
}
