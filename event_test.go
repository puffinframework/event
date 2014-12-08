package event_test

import (
	"github.com/puffinframework/config"
	"github.com/puffinframework/event"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
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

type MyEventData struct {
	Data string
}

func TestEventStore(t *testing.T) {
	os.Setenv(config.ENV_VAR_NAME, config.MODE_TEST)
	store := event.NewLeveldbStore()
	assert.NotNil(t, store)
	defer store.MustDestroy()

	header1 := event.Header{CreatedAt: time.Unix(0, 10), ID: "id1", Type: "TypeA", Version: 1}
	data1 := &MyEventData{Data: "data 1"}
	store.MustSaveEventData(header1, data1)

	header2 := event.Header{CreatedAt: time.Unix(0, 20), ID: "id2", Type: "TypeA", Version: 1}
	data2 := &MyEventData{Data: "data 2"}
	store.MustSaveEventData(header2, data2)

	header3 := event.Header{CreatedAt: time.Unix(0, 30), ID: "id3", Type: "TypeA", Version: 1}
	data3 := &MyEventData{Data: "data 3"}
	store.MustSaveEventData(header3, data3)

	data11 := &MyEventData{}
	store.MustLoadEvendData(header1, data11)
	assert.Equal(t, data11, data1)

	data22 := &MyEventData{}
	store.MustLoadEvendData(header2, data22)
	assert.Equal(t, data22, data2)

	data33 := &MyEventData{}
	store.MustLoadEvendData(header3, data33)
	assert.Equal(t, data33, data3)
}
