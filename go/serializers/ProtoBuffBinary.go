package serializers

import (
	"errors"
	"github.com/saichler/shared/go/share/interfaces"
	"google.golang.org/protobuf/proto"
)

type ProtoBuffBinary struct{}

func (s *ProtoBuffBinary) Mode() interfaces.SerializerMode {
	return interfaces.BINARY
}

func (s *ProtoBuffBinary) Marshal(any interface{}, registry interfaces.ITypeRegistry) ([]byte, error) {
	if any == nil {
		return nil, errors.New("attempting to marshal nil interface")
	}

	pb, ok := any.(proto.Message)
	if !ok {
		return nil, errors.New("interface is no a protobuf message")
	}
	return proto.Marshal(pb)
}

func (s *ProtoBuffBinary) Unmarshal(data []byte, typeName string, registry interfaces.ITypeRegistry) (interface{}, error) {
	info, err := registry.TypeInfo(typeName)
	if err != nil {
		return nil, errors.New("No type info found for type " + typeName)
	}

	ins, err := info.NewInstance()
	if err != nil {
		return nil, err
	}

	pb := ins.(proto.Message)
	err = proto.Unmarshal(data, pb)
	return pb, err
}
