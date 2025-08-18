package object

import "encoding/binary"

type Int struct{}

func (this *Int) add(any interface{}, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 8)
	i := any.(int)
	binary.BigEndian.PutUint64((*data)[*location:], uint64(i))
	*location += 8
}

func (this *Int) get(data *[]byte, location *int) interface{} {
	result := int(binary.BigEndian.Uint64((*data)[*location:]))
	*location += 8
	return result
}
