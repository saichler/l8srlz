package object

import (
	"github.com/saichler/types/go/common"
	"reflect"
)

type Map struct{}

func (this *Map) add(any interface{}) ([]byte, int, error) {
	if any == nil {
		sizeBytes, _ := sizeObjectType.add(int32(-1))
		return sizeBytes, 4, nil
	}
	mapp := reflect.ValueOf(any)
	if mapp.Len() == 0 {
		sizeBytes, _ := sizeObjectType.add(int32(-1))
		return sizeBytes, 4, nil
	}

	obj := NewEncode()
	obj.appendBytes(sizeObjectType.add(int32(mapp.Len())))

	keys := mapp.MapKeys()

	for _, key := range keys {
		obj.Add(key.Interface())
		element := mapp.MapIndex(key).Interface()
		obj.Add(element)
	}
	return obj.Data(), obj.location, nil
}

func (this *Map) get(data []byte, location int, registry common.IRegistry) (interface{}, int, error) {
	l, _ := sizeObjectType.get(data, location)
	size := l.(int32)
	if size == -1 || size == 0 {
		return nil, 4, nil
	}
	location += 4

	enc := NewDecode(data, location, registry)
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

	return newMap.Interface(), enc.Location(), nil
}
