package object

import "encoding/binary"

func addUInt64(i uint64, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 8)
	binary.BigEndian.PutUint64((*data)[*location:], i)
	*location += 8
}

func getUInt64(data *[]byte, location *int) uint64 {
	result := binary.BigEndian.Uint64((*data)[*location:])
	*location += 8
	return result
}
