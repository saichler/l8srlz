package object

import (
	"github.com/saichler/shared/go/share/interfaces"
	"reflect"
)

type Map struct{}

func (this *Map) add(any interface{}) ([]byte, int) {
	if any == nil {
		sizeBytes, _ := sizeObjectType.add(int32(-1))
		return sizeBytes, 4
	}
	mapp := reflect.ValueOf(any)
	if mapp.Len() == 0 {
		sizeBytes, _ := sizeObjectType.add(int32(-1))
		return sizeBytes, 4
	}

	s, _ := sizeObjectType.add(int32(mapp.Len()))
	keys := mapp.MapKeys()

	for _, key := range keys {
		enc := New([]byte{}, 0, "", nil)
		enc.add(key.Interface())
		element := mapp.MapIndex(key).Interface()
		enc.add(element)
		s = append(s, enc.Data()...)
	}
	return s, len(s)
}

func (this *Map) get(data []byte, location int, typeName string, registry interfaces.IRegistry) (interface{}, int) {
	l, _ := sizeObjectType.get(data, location)
	size := l.(int32)
	if size == -1 || size == 0 {
		return nil, 4
	}
	location += 4

	enc := New(data, location, typeName, registry)
	mapp := make(map[interface{}]interface{})
	var mapKeyType reflect.Type
	var mapValueType reflect.Type
	found := false

	for i := 0; i < int(size); i++ {
		key, _ := enc.Get()
		value, _ := enc.Get()
		if !found && key != nil && value != nil {
			found = true
			mapKeyType = reflect.ValueOf(key).Type()
			mapValueType = reflect.ValueOf(value).Type()
		}
		mapp[key] = value
	}
	newMap := reflect.MakeMap(reflect.MapOf(mapKeyType, mapValueType))
	for k, v := range mapp {
		if v == nil {
			newValue := reflect.New(mapValueType)
			newValue.Elem().Set(reflect.Zero(newValue.Elem().Type()))
			newMap.SetMapIndex(reflect.ValueOf(k), newValue.Elem())
		} else {
			newMap.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
		}
	}

	return newMap.Interface(), enc.Location()
}
