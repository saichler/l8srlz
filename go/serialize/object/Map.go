package object

import (
	"github.com/saichler/l8types/go/ifs"
	"reflect"
)

func addMap(any interface{}, data *[]byte, location *int) error {
	if any == nil {
		addInt32(int32(-1), data, location)
		return nil
	}
	mapp := reflect.ValueOf(any)
	if mapp.Len() == 0 {
		addInt32(int32(-1), data, location)
		return nil
	}

	addInt32(int32(mapp.Len()), data, location)

	obj := newDecode(data, location, nil)
	keys := mapp.MapKeys()

	for _, key := range keys {
		obj.Add(key.Interface())
		element := mapp.MapIndex(key).Interface()
		obj.Add(element)
	}

	return nil
}

func getMap(data *[]byte, location *int, registry ifs.IRegistry) (interface{}, error) {
	l := getInt32(data, location)
	size := int(l)
	if size == -1 || size == 0 {
		return nil, nil
	}

	enc := newDecode(data, location, registry)
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
	return newMap.Interface(), nil
}
