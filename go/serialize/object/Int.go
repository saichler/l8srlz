package object

import "encoding/binary"

func addInt(any interface{}, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 8)
	i := any.(int)
	binary.BigEndian.PutUint64((*data)[*location:], uint64(i))
	*location += 8
}

func getInt(data *[]byte, location *int) interface{} {
	result := int(binary.BigEndian.Uint64((*data)[*location:]))
	*location += 8
	return result
}
