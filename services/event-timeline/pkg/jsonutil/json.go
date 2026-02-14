package jsonutil

import "encoding/json"

func MustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func Unmarshal(data []byte, out any) error {
	return json.Unmarshal(data, out)
}
