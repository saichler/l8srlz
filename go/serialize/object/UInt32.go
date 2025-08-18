package object

import "encoding/binary"

func addUInt32(i uint32, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 4)
	binary.BigEndian.PutUint32((*data)[*location:], i)
	*location += 4
}

func getUInt32(data *[]byte, location *int) uint32 {
	result := binary.BigEndian.Uint32((*data)[*location:])
	*location += 4
	return result
}
