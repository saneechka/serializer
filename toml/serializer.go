
package toml

import (
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/saneechka/serializer"
)


type TOMLSerializer struct{}


func New() serializer.Serializer {
	return &TOMLSerializer{}
}


func (s *TOMLSerializer) Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}


func (s *TOMLSerializer) Unmarshal(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}


func (s *TOMLSerializer) Format() string {
	return "TOML"
}
