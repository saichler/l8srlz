package object

import "encoding/binary"

func addInt64(i int64, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 8)
	binary.BigEndian.PutUint64((*data)[*location:], uint64(i))
	*location += 8
}

func getInt64(data *[]byte, location *int) int64 {
	result := int64(binary.BigEndian.Uint64((*data)[*location:]))
	*location += 8
	return result
}
