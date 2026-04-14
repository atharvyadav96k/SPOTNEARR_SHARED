package common

import (
	"encoding/json"
	"fmt"
)

func ToJSON(v any) ([]byte, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}
	return bytes, nil
}

func FromJSON[T any](data []byte) (T, error) {
	var v T
	err := json.Unmarshal([]byte(data), &v)
	if err != nil {
		return v, fmt.Errorf("unmarshal error: %w", err)
	}
	return v, nil
}
