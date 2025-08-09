package object

import "encoding/binary"

type Int64 struct{}

func (this *Int64) add(any interface{}, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 8)
	i := any.(int64)
	binary.BigEndian.PutUint64((*data)[*location:], uint64(i))
	*location += 8
	/*
		checkAndEnlarge(data, location, 8)
		i := any.(int64)
		(*data)[*location] = byte((i >> 56) & 0xff)
		(*data)[*location+1] = byte((i >> 48) & 0xff)
		(*data)[*location+2] = byte((i >> 40) & 0xff)
		(*data)[*location+3] = byte((i >> 32) & 0xff)
		(*data)[*location+4] = byte((i >> 24) & 0xff)
		(*data)[*location+5] = byte((i >> 16) & 0xff)
		(*data)[*location+6] = byte((i >> 8) & 0xff)
		(*data)[*location+7] = byte((i) & 0xff)
		*location += 8
	*/
}

func (this *Int64) get(data *[]byte, location *int) interface{} {
	result := int64(binary.BigEndian.Uint64((*data)[*location:]))
	*location += 8
	return result
	/*
		var result int64
		result = int64(0xff&(*data)[*location])<<56 |
			int64(0xff&(*data)[*location+1])<<48 |
			int64(0xff&(*data)[*location+2])<<40 |
			int64(0xff&(*data)[*location+3])<<32 |
			int64(0xff&(*data)[*location+4])<<24 |
			int64(0xff&(*data)[*location+5])<<16 |
			int64(0xff&(*data)[*location+6])<<8 |
			int64(0xff&(*data)[*location+7])
		*location += 8
		return result
	*/
}
