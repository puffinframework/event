package event

import (
	"errors"
	"github.com/puffinframework/config"
	"github.com/satori/go.uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
	"time"
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

var (
	ErrOpenStore    error = errors.New("snapshot: couldn't open store")
	ErrCloseStore   error = errors.New("snapshot: couldn't close store")
	ErrDestroyStore error = errors.New("snapshot: couldn't destroy store")
)

type Store interface {
	ForEach(since time.Time, handler func(header Header))
	MustLoad(header Header, data interface{})
	MustSave(header Header, data interface{})
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
		log.Panic(ErrOpenStore)
	}

	return &leveldbStore{dir: dir, db: db}
}

func (self *leveldbStore) ForEach(since time.Time, handler func(header Header)) {
}

func (self *leveldbStore) MustLoad(header Header, data interface{}) {
}

func (self *leveldbStore) MustSave(header Header, data interface{}) {
}

func (self *leveldbStore) MustClose() {
	if err := self.db.Close(); err != nil {
		log.Panic(ErrCloseStore)
	}
}

func (self *leveldbStore) MustDestroy() {
	self.MustClose()
	if err := os.RemoveAll(self.dir); err != nil {
		log.Panic(ErrDestroyStore)
	}
}
