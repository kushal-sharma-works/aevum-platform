package identifier

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type failingReader struct{}

func (failingReader) Read(_ []byte) (int, error) {
	return 0, errors.New("entropy read failed")
}

func TestNewULIDGeneratorGeneratesID(t *testing.T) {
	generator := NewULIDGenerator()
	id, err := generator.New(time.Now().UTC())
	require.NoError(t, err)
	require.NotEmpty(t, id)
}

func TestULIDGeneratorReturnsErrorFromEntropy(t *testing.T) {
	generator := &ULIDGenerator{entropy: failingReader{}}
	_, err := generator.New(time.Now().UTC())
	require.Error(t, err)
}
