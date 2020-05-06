package redis

import "encoding/json"

// Serializer 序列化
type Serializer interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

// NewDefaultSerializer NewDefaultSerializer
func NewDefaultSerializer() Serializer {
	return &serializer{}
}

type serializer struct {
}

func (s *serializer) Marshal(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func (s *serializer) Unmarshal(data []byte, obj interface{}) error {
	return json.Unmarshal(data, obj)
}
