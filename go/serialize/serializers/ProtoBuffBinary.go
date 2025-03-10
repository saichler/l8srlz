package serializers

import (
	"errors"
	"github.com/saichler/types/go/common"
	"google.golang.org/protobuf/proto"
)

type ProtoBuffBinary struct{}

func (s *ProtoBuffBinary) Mode() common.SerializerMode {
	return common.BINARY
}

func (s *ProtoBuffBinary) Marshal(any interface{}, registry common.IRegistry) ([]byte, error) {
	if any == nil {
		return nil, errors.New("attempting to marshal nil interface")
	}

	pb, ok := any.(proto.Message)
	if !ok {
		return nil, errors.New("interface is no a protobuf message")
	}
	return proto.Marshal(pb)
}

func (s *ProtoBuffBinary) Unmarshal(data []byte, typeName string, registry common.IRegistry) (interface{}, error) {
	info, err := registry.Info(typeName)
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
