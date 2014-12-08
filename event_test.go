package event_test

import (
	"github.com/puffinframework/config"
	"github.com/puffinframework/event"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestHeader(t *testing.T) {
	header := event.NewHeader("EventType1", 3)
	assert.NotNil(t, header.CreatedAt)
	assert.NotNil(t, header.ID)
	assert.NotEqual(t, "", header.ID)
	assert.Equal(t, "EventType1", header.Type)
	assert.Equal(t, 3, header.Version)

	encoded := event.MustEncodeEventHeader(header)
	decoded := event.MustDecodeEventHeader(encoded)
	assert.Equal(t, header, decoded)
}

func TestEventStore(t *testing.T) {
	os.Setenv(config.ENV_VAR_NAME, config.MODE_TEST)
	store := event.NewLeveldbStore()
	assert.NotNil(t, store)
	defer store.MustDestroy()
}
