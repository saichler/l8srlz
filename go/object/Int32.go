package object

import (
	"reflect"
)

type Int32 struct{}

func (this *Int32) add(any interface{}) ([]byte, int) {
	i, ok := any.(int32)
	//When it is an int32 derived type
	if !ok {
		i = int32(reflect.ValueOf(any).Int())
	}
	bytes := make([]byte, 4)
	bytes[3] = byte((i >> 24) & 0xff)
	bytes[2] = byte((i >> 16) & 0xff)
	bytes[1] = byte((i >> 8) & 0xff)
	bytes[0] = byte(i & 0xff)

	return bytes, 4
}

func (this *Int32) get(data []byte, location int) (interface{}, int) {
	var result int32
	result = (0xff&int32(data[location+3])<<24 |
		0xff&int32(data[location+2])<<16 |
		0xff&int32(data[location+1])<<8 |
		0xff&int32(data[location]))
	return result, 4
}
