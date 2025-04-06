package object

import (
	"math"
)

type Float64 struct{}

func (this *Float64) add(any interface{}, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 8)
	f := any.(float64)
	i := math.Float64bits(f)
	(*data)[*location] = byte((i >> 56) & 0xff)
	(*data)[*location+1] = byte((i >> 48) & 0xff)
	(*data)[*location+2] = byte((i >> 40) & 0xff)
	(*data)[*location+3] = byte((i >> 32) & 0xff)
	(*data)[*location+4] = byte((i >> 24) & 0xff)
	(*data)[*location+5] = byte((i >> 16) & 0xff)
	(*data)[*location+6] = byte((i >> 8) & 0xff)
	(*data)[*location+7] = byte((i) & 0xff)
	*location += 8
}

func (this *Float64) get(data *[]byte, location *int) interface{} {
	var result uint64
	result = uint64(0xff&(*data)[*location])<<56 |
		uint64(0xff&(*data)[*location+1])<<48 |
		uint64(0xff&(*data)[*location+2])<<40 |
		uint64(0xff&(*data)[*location+3])<<32 |
		uint64(0xff&(*data)[*location+4])<<24 |
		uint64(0xff&(*data)[*location+5])<<16 |
		uint64(0xff&(*data)[*location+6])<<8 |
		uint64(0xff&(*data)[*location+7])
	*location += 8
	return math.Float64frombits(result)
}
