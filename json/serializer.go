
package json

import (
	"encoding/json"

	"github.com/saneechka/serializer"
)


type JSONSerializer struct{}


func New() serializer.Serializer {
	return &JSONSerializer{}
}


func (s *JSONSerializer) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}


func (s *JSONSerializer) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}


func (s *JSONSerializer) Format() string {
	return "JSON"
}
