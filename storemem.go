package event

import (
	"encoding/json"
	"log"
	"sort"
	"time"
)

type storeMem struct {
	db map[string]interface{}
}

func NewMemStore() Store {
	return &storeMem{db: make(map[string]interface{})}
}

func (self *storeMem) ForEachEventHeader(since time.Time, callback func(header Header) (bool, error)) (callbackErr error) {
	startHeader := Header{CreatedAt: since.Add(1 * time.Nanosecond), ID: "", Type: "", Version: 0}
	startHeaderStr := string(MustEncodeEventHeader(startHeader))

	headerStrArr := []string{}
	for headerStr, _ := range self.db {
		if headerStr >= startHeaderStr {
			headerStrArr = append(headerStrArr, headerStr)
		}
	}

	sort.Strings(headerStrArr)
	for _, headerStr := range headerStrArr {
		header := MustDecodeEventHeader([]byte(headerStr))
		cont, err := callback(header)
		if err != nil {
			callbackErr = err
			break
		}
		if cont == false {
			break
		}
	}

	return callbackErr
}

func (self *storeMem) MustLoadEvent(header Header, data interface{}) {
	key := string(MustEncodeEventHeader(header))

	// The following is absurd, all we want is to assign self.db[key] to data
	// but I haven't figured out how to do this with poiters
	value, err := json.Marshal(self.db[key])
	if err != nil {
		log.Print(err)
		log.Panic(ErrMarshalEventData)
	}
	if err = json.Unmarshal(value, data); err != nil {
		log.Print(err)
		log.Panic(ErrUnmarshalEventData)
	}
}

func (self *storeMem) MustSaveEvent(header Header, data interface{}) {
	key := string(MustEncodeEventHeader(header))
	self.db[key] = data
}

func (self *storeMem) MustClose() {
}

func (self *storeMem) MustDestroy() {
}
