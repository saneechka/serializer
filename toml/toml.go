package toml

import (
	"github.com/pelletier/go-toml/v2"
)

type TOMLSerializer struct{}

func New() *TOMLSerializer {
	return &TOMLSerializer{}
}

func (s *TOMLSerializer) Marshal(v any) ([]byte, error) {
	return toml.Marshal(v)
}

func (s *TOMLSerializer) Unmarshal(data []byte, v any) error {
	return toml.Unmarshal(data, v)
}

func (s *TOMLSerializer) Format() string {
	return "TOML"
}
