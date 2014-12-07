package event_test

import (
	. "github.com/puffinframework/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	header := NewEventHeader("SomeEvent", 3)
	assert.NotNil(t, header.CreatedAt)
	assert.NotNil(t, header.ID)
	assert.Equal(t, "SomeEvent", header.Type)
	assert.Equal(t, 3, header.Version)
}
