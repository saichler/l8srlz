package object

import (
	"reflect"
)

type Int32 struct{}

func (this *Int32) add(any interface{}, data *[]byte, location *int) {
	i, ok := any.(int32)
	//When it is an int32 derived type
	if !ok {
		i = int32(reflect.ValueOf(any).Int())
	}
	(*data)[*location+3] = byte((i >> 24) & 0xff)
	(*data)[*location+2] = byte((i >> 16) & 0xff)
	(*data)[*location+1] = byte((i >> 8) & 0xff)
	(*data)[*location] = byte(i & 0xff)
	*location += 4
}

func (this *Int32) get(data *[]byte, location *int) interface{} {
	var result int32
	result = (0xff&int32((*data)[*location+3])<<24 |
		0xff&int32((*data)[*location+2])<<16 |
		0xff&int32((*data)[*location+1])<<8 |
		0xff&int32((*data)[*location]))
	*location += 4
	return result
}
