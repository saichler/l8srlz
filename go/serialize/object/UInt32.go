package object

import "encoding/binary"

type UInt32 struct{}

func (this *UInt32) add(any interface{}, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 4)
	i := any.(uint32)
	binary.BigEndian.PutUint32((*data)[*location:], i)
	*location += 4
	/*
		checkAndEnlarge(data, location, 4)
		i := any.(uint32)
		(*data)[*location+3] = byte(i)
		(*data)[*location+2] = byte(i >> 8)
		(*data)[*location+1] = byte(i >> 16)
		(*data)[*location] = byte(i >> 24)
		*location += 4
	*/
}

func (this *UInt32) get(data *[]byte, location *int) interface{} {
	result := binary.BigEndian.Uint32((*data)[*location:])
	*location += 4
	return result
	/*
		var result uint32
		v1 := uint32((*data)[*location]) << 24
		v2 := uint32((*data)[*location+1]) << 16
		v3 := uint32((*data)[*location+2]) << 8
		v4 := uint32((*data)[*location+3])
		result = v1 + v2 + v3 + v4
		*location += 4
		return result
	*/
}
