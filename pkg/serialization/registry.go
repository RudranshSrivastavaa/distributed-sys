package serialization

type Registry struct {

	serializers map[string]Serializer
}

func NewRegistry() *Registry {

	r := &Registry{

		serializers: make(map[string]Serializer),
	}

	r.Register(
		NewJSONSerializer(),
	)

	return r

}

func (r *Registry) Register(
	s Serializer,
) {

	r.serializers[
		s.ContentType(),
	] = s

}

func (r *Registry) Get(contentType string) (Serializer, bool) {

	s, ok := r.serializers[
		contentType,
	]

	return s, ok

}