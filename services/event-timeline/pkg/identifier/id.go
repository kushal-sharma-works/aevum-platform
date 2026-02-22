package identifier

import (
	"crypto/rand"
	"io"
	"time"

	"github.com/oklog/ulid/v2"
)

type Generator interface {
	New(t time.Time) (string, error)
}

type ULIDGenerator struct {
	entropy io.Reader
}

func NewULIDGenerator() *ULIDGenerator {
	return &ULIDGenerator{entropy: rand.Reader}
}

func (g *ULIDGenerator) New(t time.Time) (string, error) {
	id, err := ulid.New(ulid.Timestamp(t), g.entropy)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
