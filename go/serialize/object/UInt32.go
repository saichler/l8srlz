package object

import "encoding/binary"

func addUInt32(any interface{}, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 4)
	i := any.(uint32)
	binary.BigEndian.PutUint32((*data)[*location:], i)
	*location += 4
}

func getUInt32(data *[]byte, location *int) interface{} {
	result := binary.BigEndian.Uint32((*data)[*location:])
	*location += 4
	return result
}
