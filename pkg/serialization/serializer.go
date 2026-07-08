package serialization

type Serializer interface {

	Marshal(v any) ([]byte, error)

	Unmarshal(data []byte,v any) error

	ContentType() string
}