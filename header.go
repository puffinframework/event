package event

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/satori/go.uuid"
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

