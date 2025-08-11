package object

import (
	"encoding/binary"
)

type Int32 struct{}

func (this *Int32) add(any interface{}, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 4)
	i, _ := any.(int32)
	binary.BigEndian.PutUint32((*data)[*location:], uint32(i))
	*location += 4
}

func (this *Int32) get(data *[]byte, location *int) interface{} {
	result := int32(binary.BigEndian.Uint32((*data)[*location:]))
	*location += 4
	return result
}
