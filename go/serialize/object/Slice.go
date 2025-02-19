package object

import (
	"github.com/saichler/shared/go/share/interfaces"
	"reflect"
)

type Slice struct{}

func (this *Slice) add(any interface{}, log interfaces.ILogger) ([]byte, int) {
	if any == nil {
		sizeBytes, _ := sizeObjectType.add(int32(-1))
		return sizeBytes, 4
	}
	slice := reflect.ValueOf(any)
	if slice.Len() == 0 {
		sizeBytes, _ := sizeObjectType.add(int32(-1))
		return sizeBytes, 4
	}

	s, _ := sizeObjectType.add(int32(slice.Len()))

	data, ok := any.([]byte)
	if ok {
		s = append(s, byte(1))
		s = append(s, data...)
	} else {
		for i := 0; i < slice.Len(); i++ {
			enc := NewEncode([]byte{}, 0, log)
			element := slice.Index(i).Interface()
			enc.Add(element)
			s = append(s, enc.Data()...)
		}
	}
	return s, len(s)
}

func (this *Slice) get(data []byte, location int, typeName string, registry interfaces.IRegistry, log interfaces.ILogger) (interface{}, int) {
	l, _ := sizeObjectType.get(data, location)
	size := l.(int32)
	if size == -1 || size == 0 {
		return nil, 4
	}
	location += 4
	enc := NewDecode(data, location, typeName, registry, log)

	if data[location] == 1 {
		location++
		result := data[location : location+int(size)]
		location += int(size)
		return result, location
	}

	elems := make([]interface{}, 0)
	var sliceType reflect.Type

	for i := 0; i < int(size); i++ {
		element, _ := enc.Get()
		if i == 0 {
			sliceType = reflect.SliceOf(reflect.ValueOf(element).Type())
		}
		elems = append(elems, element)
	}

	newSlice := reflect.MakeSlice(sliceType, len(elems), len(elems))
	for i := 0; i < int(size); i++ {
		if elems[i] != nil {
			newSlice.Index(i).Set(reflect.ValueOf(elems[i]))
		}
	}

	return newSlice.Interface(), enc.Location()
}
