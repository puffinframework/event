package event

import (
	"errors"
	"github.com/puffinframework/config"
	"github.com/satori/go.uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ErrDecodeHeader error = errors.New("event: couldn't decode header")
	ErrEncodeHeader error = errors.New("event: couldn't encode header")
	ErrOpenStore    error = errors.New("event: couldn't open store")
	ErrCloseStore   error = errors.New("event: couldn't close store")
	ErrDestroyStore error = errors.New("event: couldn't destroy store")
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

func EncodeHeader(header Header) []byte {
	createdAt, err := header.CreatedAt.MarshalBinary()
	if err != nil {
		log.Panic(ErrEncodeHeader)
	}

	tokens := []string{
		string(createdAt),
		header.ID,
		header.Type,
		strconv.Itoa(header.Version),
	}
	return []byte(strings.Join(tokens, "::"))
}

func DecodeHeader(encoded []byte) Header {
	tokens := strings.Split(string(encoded), "::")

	createdAt := time.Unix(0, 0)
	err := createdAt.UnmarshalBinary([]byte(tokens[0]))
	if err != nil {
		log.Panic(ErrDecodeHeader)
	}

	version, err := strconv.Atoi(tokens[3])
	if err != nil {
		log.Panic(ErrDecodeHeader)
	}

	return Header{
		CreatedAt: createdAt,
		ID:        tokens[1],
		Type:      tokens[2],
		Version:   version,
	}
}

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
