package object

import (
	"encoding/binary"
)

func addInt32(i int32, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 4)
	binary.BigEndian.PutUint32((*data)[*location:], uint32(i))
	*location += 4
}

func getInt32(data *[]byte, location *int) int32 {
	result := int32(binary.BigEndian.Uint32((*data)[*location:]))
	*location += 4
	return result
}
