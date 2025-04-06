package object

import (
	"errors"
	"github.com/saichler/types/go/common"
	"google.golang.org/protobuf/proto"
	"reflect"
)

type Struct struct {
}

func (this *Struct) add(any interface{}, data *[]byte, location *int) error {
	if any == nil {
		sizeObjectType.add(int32(-1), data, location)
		return nil
	}

	val := reflect.ValueOf(any)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			sizeObjectType.add(int32(-1), data, location)
			return nil
		}
		val = val.Elem()
	}

	typeName := val.Type().Name()
	var pbData []byte

	if typeName == "Transaction" {
		pbData, _ = TransactionSerializer.Marshal(any, nil)
	} else {
		pb := any.(proto.Message)
		pbd, err := proto.Marshal(pb)
		if err != nil {
			return errors.New("Failed To marshal proto " + typeName + " in protobuf object:" + err.Error())
		}
		pbData = pbd
	}

	checkAndEnlarge(data, location, 8+len(typeName)+len(pbData))
	sizeObjectType.add(int32(len(pbData)), data, location)
	stringObjectType.add(typeName, data, location)
	copy((*data)[*location:*location+len(pbData)], pbData)
	*location += len(pbData)
	return nil
}

func (this *Struct) get(data *[]byte, location *int, registry common.IRegistry) (interface{}, error) {
	l := sizeObjectType.get(data, location)
	size := int(l.(int32))

	if size == -1 || size == 0 {
		return nil, nil
	}

	typeN := stringObjectType.get(data, location)
	typeName := typeN.(string)

	var info common.IInfo
	var err error
	var pb interface{}

	isTransaction := typeName == "Transaction"
	if !isTransaction {
		info, err = registry.Info(typeName)
		if err != nil {
			return nil, errors.New("Unknown proto name " + typeName + " in registry, please register it.")
		}

		pb, err = info.NewInstance()
		if err != nil {
			return nil, errors.New("Error proto name " + typeName + " in registry, cannot instantiate.")
		}
	}

	protoData := make([]byte, size)
	copy(protoData, (*data)[*location:*location+size])

	if isTransaction {
		pb, _ = TransactionSerializer.Unmarshal(protoData, nil)
	} else {
		err = proto.Unmarshal(protoData, pb.(proto.Message))
		if err != nil {
			return []byte{}, errors.New("Failed To unmarshal proto " + typeName + ":" + err.Error())
		}
	}
	*location += size

	return pb, nil
}
