package serializers

import (
	"github.com/saichler/serializer/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
)

type ProtoBuffBinary struct{}

func (s *ProtoBuffBinary) Mode() ifs.SerializerMode {
	return ifs.BINARY
}

func (s *ProtoBuffBinary) Marshal(any interface{}, registry ifs.IRegistry) ([]byte, error) {
	obj := object.NewEncode()
	obj.Add(any)
	return obj.Data(), nil
}

func (s *ProtoBuffBinary) Unmarshal(data []byte, registry ifs.IRegistry) (interface{}, error) {
	obj := object.NewDecode(data, 0, registry)
	return obj.Get()
}
