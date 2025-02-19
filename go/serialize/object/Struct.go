package object

import (
	"github.com/saichler/shared/go/share/interfaces"
	"google.golang.org/protobuf/proto"
	"reflect"
)

type Struct struct {
}

func (this *Struct) add(any interface{}, log interfaces.ILogger) ([]byte, int) {

	if any == nil || reflect.ValueOf(any).IsNil() {
		sizeBytes, _ := sizeObjectType.add(int32(-1))
		return sizeBytes, 4
	}

	typ := reflect.ValueOf(any).Elem().Type()
	typeName := typ.Name()

	pb := any.(proto.Message)
	pbData, err := proto.Marshal(pb)
	if err != nil {
		log.Error("Failed To marshal proto ", typeName, " in protobuf object:", err)
		return []byte{}, 0
	}

	data, _ := stringObjectType.add(typeName)
	sizeData, _ := sizeObjectType.add(int32(len(pbData)))
	data = append(data, sizeData...)
	data = append(data, pbData...)

	return data, len(data)
}

func (this *Struct) get(data []byte, location int, typeName string, registry interfaces.IRegistry, log interfaces.ILogger) (interface{}, int) {
	l, _ := sizeObjectType.get(data, location)
	size := l.(int32)
	if size == -1 || size == 0 {
		return nil, 4
	}

	typeN, typeSize := stringObjectType.get(data, location)
	typeName = typeN.(string)
	info, err := registry.Info(typeName)
	if err != nil {
		log.Error("Unknown proto name ", typeName, " in registry, please register it.")
		return []byte{}, 0
	}

	pb, err := info.NewInstance()
	if err != nil {
		log.Error("Unknown proto name ", typeName, " in registry, please register it.")
		return []byte{}, 0
	}

	location += typeSize
	s, _ := sizeObjectType.get(data, location)
	size = s.(int32)
	location += 4
	protoData := data[location : location+int(size)]

	err = proto.Unmarshal(protoData, pb.(proto.Message))
	if err != nil {
		log.Error("Failed To unmarshal proto ", typeName, ":", err)
		return []byte{}, 0
	}
	return pb, typeSize + 4 + int(size)
}
