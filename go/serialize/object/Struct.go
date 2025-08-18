package object

import (
	"errors"
	"github.com/saichler/l8types/go/ifs"
	"google.golang.org/protobuf/proto"
	"reflect"
)

func addStruct(any interface{}, data *[]byte, location *int) error {
	if any == nil {
		addInt32(int32(-1), data, location)
		return nil
	}

	val := reflect.ValueOf(any)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			addInt32(int32(-1), data, location)
			return nil
		}
		val = val.Elem()
	}

	typeName := val.Type().Name()

	pb := any.(proto.Message)
	pbData, err := proto.Marshal(pb)
	if err != nil {
		return errors.New("Failed To marshal proto " + typeName + " in protobuf object:" + err.Error())
	}

	size := len(pbData)
	checkAndEnlarge(data, location, 8+len(typeName)+size)
	if size == 0 {
		addInt32(int32(-2), data, location)
	} else {
		addInt32(int32(len(pbData)), data, location)
	}
	addString(typeName, data, location)
	if size > 0 {
		copy((*data)[*location:*location+len(pbData)], pbData)
		*location += len(pbData)
	}
	return nil
}

func getStruct(data *[]byte, location *int, registry ifs.IRegistry) (interface{}, error) {
	l := getInt32(data, location)
	size := int(l)

	if size == -1 || size == 0 {
		return nil, nil
	}

	typeName := getString(data, location)

	var info ifs.IInfo
	var err error
	var pb interface{}

	info, err = registry.Info(typeName)
	if err != nil {
		//panic("Unknown proto name " + typeName + " in registry, please register it.")
		return nil, errors.New("Unknown proto name " + typeName + " in registry, please register it.")
	}

	pb, err = info.NewInstance()
	if err != nil {
		return nil, errors.New("Error proto name " + typeName + " in registry, cannot instantiate.")
	}
	//if the size is -2 it is an empty interface
	if size == -2 {
		return pb, nil
	}

	protoData := make([]byte, size)
	copy(protoData, (*data)[*location:*location+size])

	err = proto.Unmarshal(protoData, pb.(proto.Message))
	if err != nil {
		return []byte{}, errors.New("Failed To unmarshal proto " + typeName + ":" + err.Error())
	}
	*location += size

	return pb, nil
}
