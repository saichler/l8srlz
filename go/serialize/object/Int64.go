package object

import "encoding/binary"

func addInt64(any interface{}, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 8)
	i := any.(int64)
	binary.BigEndian.PutUint64((*data)[*location:], uint64(i))
	*location += 8
}

func getInt64(data *[]byte, location *int) interface{} {
	result := int64(binary.BigEndian.Uint64((*data)[*location:]))
	*location += 8
	return result
}
