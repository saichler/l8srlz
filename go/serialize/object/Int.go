package object

import "encoding/binary"

func addInt(i int, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 8)
	binary.BigEndian.PutUint64((*data)[*location:], uint64(i))
	*location += 8
}

func getInt(data *[]byte, location *int) int {
	result := int(binary.BigEndian.Uint64((*data)[*location:]))
	*location += 8
	return result
}
