package event_test

import (
	"errors"
	"testing"
	"time"

	"github.com/puffinframework/event"

	"github.com/stretchr/testify/assert"
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

func testEventStore(t *testing.T, store event.Store) {
	defer store.MustDestroy()
	assert.NotNil(t, store)

	time0 := time.Unix(0, 0)
	time1 := time.Unix(0, 1)
	time2 := time.Unix(0, 2)
	time3 := time.Unix(0, 3)

	header1 := event.Header{CreatedAt: time1, ID: "id1", Type: "TypeA", Version: 1}
	data1 := &MyEventData{Data: "data 1"}
	store.MustSaveEvent(header1, data1)

	header2 := event.Header{CreatedAt: time2, ID: "id2", Type: "TypeA", Version: 1}
	data2 := &MyEventData{Data: "data 2"}
	store.MustSaveEvent(header2, data2)

	header3 := event.Header{CreatedAt: time3, ID: "id3", Type: "TypeA", Version: 1}
	data3 := &MyEventData{Data: "data 3"}
	store.MustSaveEvent(header3, data3)

	data11 := &MyEventData{}
	store.MustLoadEvent(header1, data11)
	assert.Equal(t, data11, data1)

	data22 := &MyEventData{}
	store.MustLoadEvent(header2, data22)
	assert.Equal(t, data22, data2)

	data33 := &MyEventData{}
	store.MustLoadEvent(header3, data33)
	assert.Equal(t, data33, data3)

	ids := []string{}
	data := []string{}
	store.ForEachEventHeader(time0, func(header event.Header) (bool, error) {
		ids = append(ids, header.ID)
		d := &MyEventData{}
		store.MustLoadEvent(header, d)
		data = append(data, d.Data)
		return true, nil
	})
	assert.Equal(t, []string{"id1", "id2", "id3"}, ids)
	assert.Equal(t, []string{"data 1", "data 2", "data 3"}, data)

	ids = []string{}
	data = []string{}
	store.ForEachEventHeader(time1, func(header event.Header) (bool, error) {
		ids = append(ids, header.ID)
		d := &MyEventData{}
		store.MustLoadEvent(header, d)
		data = append(data, d.Data)
		return true, nil
	})
	assert.Equal(t, []string{"id2", "id3"}, ids)
	assert.Equal(t, []string{"data 2", "data 3"}, data)

	ids = []string{}
	data = []string{}
	store.ForEachEventHeader(time0, func(header event.Header) (bool, error) {
		ids = append(ids, header.ID)
		d := &MyEventData{}
		store.MustLoadEvent(header, d)
		data = append(data, d.Data)
		return len(ids) < 2, nil
	})
	assert.Equal(t, []string{"id1", "id2"}, ids)
	assert.Equal(t, []string{"data 1", "data 2"}, data)

	err := store.ForEachEventHeader(time0, func(header event.Header) (bool, error) {
		return true, errors.New("callback error")
	})
	assert.Equal(t, "callback error", err.Error())

}
