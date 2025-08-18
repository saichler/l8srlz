package object

import (
	"github.com/saichler/l8types/go/ifs"
	"reflect"
)

func addSlice(any interface{}, data *[]byte, location *int) error {
	if any == nil {
		addInt32(int32(-1), data, location)
		return nil
	}
	slice := reflect.ValueOf(any)
	if slice.Len() == 0 {
		addInt32(int32(-1), data, location)
		return nil
	}

	addInt32(int32(slice.Len()), data, location)
	dataByte, ok := any.([]byte)
	if ok {
		(*data)[*location] = 1
		*location += 1
		checkAndEnlarge(data, location, len(dataByte))
		copy((*data)[*location:*location+len(dataByte)], dataByte)
		*location += len(dataByte)
	} else {
		(*data)[*location] = 0
		*location += 1
		obj := newDecode(data, location, nil)
		for i := 0; i < slice.Len(); i++ {
			element := slice.Index(i).Interface()
			obj.Add(element)
		}
	}
	return nil
}

func getSlice(data *[]byte, location *int, registry ifs.IRegistry) (interface{}, error) {
	l := getInt32(data, location)
	size := int(l.(int32))
	if size == -1 || size == 0 {
		return nil, nil
	}

	if (*data)[*location] == 1 {
		*location += 1
		result := make([]byte, size)
		copy(result, (*data)[*location:*location+size])
		*location += size
		return result, nil
	} else {
		*location += 1
	}

	elems := make([]interface{}, 0)
	var sliceType reflect.Type

	obj := newDecode(data, location, registry)

	for i := 0; i < size; i++ {
		element, _ := obj.Get()
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

	return newSlice.Interface(), nil
}
