package serializer

type Serializer interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
	Format() string
}

type Error struct {
	message string
}

func NewError(message string) *Error {
	return &Error{message: message}
}

func (e *Error) Error() string {
	return e.message
}

var (
	ErrUnsupportedFormat = NewError("unsupported serialization format")
)
