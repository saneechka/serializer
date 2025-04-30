package json

import (
	"encoding/json"
)

type JSONSerializer struct{}

func New() *JSONSerializer {
	return &JSONSerializer{}
}

func (s *JSONSerializer) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (s *JSONSerializer) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (s *JSONSerializer) Format() string {
	return "JSON"
}
