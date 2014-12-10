package event

import (
	"errors"
	"github.com/puffinframework/config"
	"github.com/satori/go.uuid"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbErrors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/util"
	"labix.org/v2/mgo/bson"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ErrDecodeEventHeader  error = errors.New("event: couldn't decode the event header")
	ErrEncodeEventHeader  error = errors.New("event: couldn't encode the event header")
	ErrOpenStore          error = errors.New("event: couldn't open the store")
	ErrCloseStore         error = errors.New("event: couldn't close the store")
	ErrDestroyStore       error = errors.New("event: couldn't destroy the store")
	ErrGetEventData       error = errors.New("event: couldn't get event data from the store")
	ErrPutEventData       error = errors.New("event: couldn't put event data into the store")
	ErrMarshalEventData   error = errors.New("event: couldn't marshal the event data")
	ErrUnmarshalEventData error = errors.New("event: couldn't unmarshal the event data")
	ErrForEachEventHeader error = errors.New("event: there was a problem during the iteration of event headers")
)

type Header struct {
	CreatedAt time.Time
	ID        string
	Type      string
	Version   int
}

func NewHeader(eventType string, version int) Header {
	return Header{
		CreatedAt: time.Now(),
		ID:        uuid.NewV1().String(),
		Type:      eventType,
		Version:   version,
	}
}

func MustEncodeEventHeader(header Header) []byte {
	createdAt, err := header.CreatedAt.MarshalBinary()
	if err != nil {
		log.Print(err)
		log.Panic(ErrEncodeEventHeader)
	}

	tokens := []string{
		string(createdAt),
		header.ID,
		header.Type,
		strconv.Itoa(header.Version),
	}
	return []byte(strings.Join(tokens, "::"))
}

func MustDecodeEventHeader(encoded []byte) Header {
	tokens := strings.Split(string(encoded), "::")

	createdAt := time.Unix(0, 0)
	err := createdAt.UnmarshalBinary([]byte(tokens[0]))
	if err != nil {
		log.Print(err)
		log.Panic(ErrDecodeEventHeader)
	}

	version, err := strconv.Atoi(tokens[3])
	if err != nil {
		log.Print(err)
		log.Panic(ErrDecodeEventHeader)
	}

	return Header{
		CreatedAt: createdAt,
		ID:        tokens[1],
		Type:      tokens[2],
		Version:   version,
	}
}

type Store interface {
	ForEachEventHeader(since time.Time, callback func(header Header) bool)
	MustLoadEventData(header Header, data interface{})
	MustSaveEventData(header Header, data interface{})
	MustClose()
	MustDestroy()
}

type leveldbStoreConfig struct {
	EventStore struct {
		LeveldbDir string
	}
}

type leveldbStore struct {
	dir string
	db  *leveldb.DB
}

func NewLeveldbStore() Store {
	cfg := &leveldbStoreConfig{}
	config.MustReadConfig(cfg)

	dir := cfg.EventStore.LeveldbDir

	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		log.Print(err)
		log.Panic(ErrOpenStore)
	}

	return &leveldbStore{dir: dir, db: db}
}

func (self *leveldbStore) ForEachEventHeader(since time.Time, callback func(header Header) bool) {
	startHeader := Header{CreatedAt: since.Add(1 * time.Nanosecond), ID: "", Type: "", Version: 0}
	startKey := MustEncodeEventHeader(startHeader)

	iter := self.db.NewIterator(&util.Range{Start: startKey}, nil)
	for iter.Next() {
		key := iter.Key()
		header := MustDecodeEventHeader(key)
		if callback(header) == false {
			break
		}
	}
	iter.Release()

	if err := iter.Error(); err != nil {
		log.Print(err)
		log.Panic(ErrForEachEventHeader)
	}
}

func (self *leveldbStore) MustLoadEventData(header Header, data interface{}) {
	key := MustEncodeEventHeader(header)

	value, err := self.db.Get(key, nil)
	if err != nil {
		if err == leveldbErrors.ErrNotFound {
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

func (self *leveldbStore) MustSaveEventData(header Header, data interface{}) {
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

func (self *leveldbStore) MustClose() {
	if err := self.db.Close(); err != nil {
		log.Print(err)
		log.Panic(ErrCloseStore)
	}
}

func (self *leveldbStore) MustDestroy() {
	self.MustClose()
	if err := os.RemoveAll(self.dir); err != nil {
		log.Print(err)
		log.Panic(ErrDestroyStore)
	}
}
