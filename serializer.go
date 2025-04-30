package serializer

import (
	"github.com/saneechka/serializer/json"
	"github.com/saneechka/serializer/toml"
)

func New(format string) (Serializer, error) {
	switch format {
	case "json", "JSON":
		return json.New(), nil
	case "toml", "TOML":
		return toml.New(), nil
	default:
		return nil, ErrUnsupportedFormat
	}
}

type GinSerializer struct {
	serializer Serializer
}

func NewGin(format string) (*GinSerializer, error) {
	s, err := New(format)
	if err != nil {
		return nil, err
	}
	return &GinSerializer{serializer: s}, nil
}

func (g *GinSerializer) Bind(data []byte, obj any) error {
	return g.serializer.Unmarshal(data, obj)
}

func (g *GinSerializer) Render(obj any) ([]byte, error) {
	return g.serializer.Marshal(obj)
}

func (g *GinSerializer) Format() string {
	return g.serializer.Format()
}
