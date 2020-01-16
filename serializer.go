package redis

// Serializer 序列化
type Serializer interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type serializer struct {
}

func (s *serializer) Marshal(interface{}) ([]byte, error) {
	return nil, nil
}

func (s *serializer) Unmarshal([]byte, interface{}) error {
	return nil
}
