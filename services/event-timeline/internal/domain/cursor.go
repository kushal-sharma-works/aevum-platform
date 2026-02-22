package domain

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

const (
	DirectionForward  = "forward"
	DirectionBackward = "backward"
)

type Cursor struct {
	StreamID  string
	Sequence  int64
	Direction string
}

func (c Cursor) Encode() string {
	raw := fmt.Sprintf("%s:%d:%s", c.StreamID, c.Sequence, c.Direction)
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func DecodeCursor(v string) (Cursor, error) {
	b, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return Cursor{}, fmt.Errorf("decode cursor base64: %w", err)
	}
	parts := strings.Split(string(b), ":")
	if len(parts) != 3 {
		return Cursor{}, fmt.Errorf("invalid cursor format: %w", ErrValidation)
	}
	seq, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return Cursor{}, fmt.Errorf("parse cursor sequence: %w", err)
	}
	if parts[2] != DirectionForward && parts[2] != DirectionBackward {
		return Cursor{}, fmt.Errorf("invalid cursor direction: %w", ErrValidation)
	}
	return Cursor{StreamID: parts[0], Sequence: seq, Direction: parts[2]}, nil
}
