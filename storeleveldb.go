package event

import (
	"log"
	"os"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/util"
	"labix.org/v2/mgo/bson"
)

type storeLeveldb struct {
	dir string
	db  *leveldb.DB
}

func NewLeveldbStore(dir string) Store {
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		log.Print(err)
		log.Panic(ErrOpenStore)
	}

	return &storeLeveldb{dir: dir, db: db}
}

func (self *storeLeveldb) ForEachEventHeader(since time.Time, callback func(header Header) (bool, error)) (callbackErr error) {
	startHeader := Header{CreatedAt: since.Add(1 * time.Nanosecond), ID: "", Type: "", Version: 0}
	startKey := MustEncodeEventHeader(startHeader)

	iter := self.db.NewIterator(&util.Range{Start: startKey}, nil)
	for iter.Next() {
		key := iter.Key()
		header := MustDecodeEventHeader(key)
		cont, err := callback(header)
		if err != nil {
			log.Print(err)
			callbackErr = err
			break
		}
		if cont == false {
			break
		}
	}
	iter.Release()

	if err := iter.Error(); err != nil {
		log.Print(err)
		log.Panic(ErrForEachEventHeader)
	}

	return callbackErr
}

func (self *storeLeveldb) MustLoadEventData(header Header, data interface{}) {
	key := MustEncodeEventHeader(header)

	value, err := self.db.Get(key, nil)
	if err != nil {
		if err == errors.ErrNotFound {
			return
		} else {
			log.Print(err)
			log.Panic(ErrGetEventData)
		}
	}

	if err = bson.Unmarshal(value, data); err != nil {
		log.Print(err)
		log.Panic(ErrUnmarshalEventData)
	}
}

func (self *storeLeveldb) MustSaveEventData(header Header, data interface{}) {
	key := MustEncodeEventHeader(header)

	value, err := bson.Marshal(data)
	if err != nil {
		log.Print(err)
		log.Panic(ErrMarshalEventData)
	}

	if err = self.db.Put(key, value, nil); err != nil {
		log.Print(err)
		log.Panic(ErrPutEventData)
	}
}

func (self *storeLeveldb) MustClose() {
	if err := self.db.Close(); err != nil {
		log.Print(err)
		log.Panic(ErrCloseStore)
	}
}

func (self *storeLeveldb) MustDestroy() {
	self.MustClose()
	if err := os.RemoveAll(self.dir); err != nil {
		log.Print(err)
		log.Panic(ErrDestroyStore)
	}
}
