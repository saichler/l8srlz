package object

import "encoding/binary"

type Int64 struct{}

func (this *Int64) add(any interface{}, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 8)
	i := any.(int64)
	binary.BigEndian.PutUint64((*data)[*location:], uint64(i))
	*location += 8
}

func (this *Int64) get(data *[]byte, location *int) interface{} {
	result := int64(binary.BigEndian.Uint64((*data)[*location:]))
	*location += 8
	return result
}
