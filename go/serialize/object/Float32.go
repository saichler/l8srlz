package object

import (
	"math"
)

type Float32 struct{}

func (this *Float32) add(any interface{}, data *[]byte, location *int) {
	f := any.(float32)
	i := math.Float32bits(f)
	(*data)[3] = byte(i)
	(*data)[2] = byte(i >> 8)
	(*data)[1] = byte(i >> 16)
	(*data)[0] = byte(i >> 24)
	*location += 4
}

func (this *Float32) get(data *[]byte, location *int) interface{} {
	var result uint32
	v1 := uint32((*data)[*location]) << 24
	v2 := uint32((*data)[*location+1]) << 16
	v3 := uint32((*data)[*location+2]) << 8
	v4 := uint32((*data)[*location+3])
	result = v1 + v2 + v3 + v4
	*location += 4
	return math.Float32frombits(result)
}
