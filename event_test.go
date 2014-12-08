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

	ids := []string{}
	data := []string{}
	store.ForEachEventHeader(time.Unix(0, 0), func(header event.Header) bool {
		ids = append(ids, header.ID)
		d := &MyEventData{}
		store.MustLoadEvendData(header, d)
		data = append(data, d.Data)
		return true
	})
	assert.Equal(t, []string{"id1", "id2", "id3"}, ids)
	assert.Equal(t, []string{"data 1", "data 2", "data 3"}, data)

	ids = []string{}
	data = []string{}
	store.ForEachEventHeader(time.Unix(0, 10), func(header event.Header) bool {
		ids = append(ids, header.ID)
		d := &MyEventData{}
		store.MustLoadEvendData(header, d)
		data = append(data, d.Data)
		return true
	})
	assert.Equal(t, []string{"id2", "id3"}, ids)
	assert.Equal(t, []string{"data 2", "data 3"}, data)

	ids = []string{}
	data = []string{}
	store.ForEachEventHeader(time.Unix(0, 0), func(header event.Header) bool {
		ids = append(ids, header.ID)
		d := &MyEventData{}
		store.MustLoadEvendData(header, d)
		data = append(data, d.Data)
		return len(ids) < 2
	})
	assert.Equal(t, []string{"id1", "id2"}, ids)
	assert.Equal(t, []string{"data 1", "data 2"}, data)
}
