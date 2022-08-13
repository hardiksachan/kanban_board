package json

import (
	"encoding/json"
	"io"
)

func Parse[T interface{}](r io.Reader) (*T, error) {
	var result T
	err := json.NewDecoder(r).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
