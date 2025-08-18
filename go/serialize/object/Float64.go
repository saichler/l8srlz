package object

import (
	"encoding/binary"
	"math"
)

type Float64 struct{}

func (this *Float64) add(any interface{}, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 8)
	f := any.(float64)
	i := math.Float64bits(f)
	loc := *location
	binary.BigEndian.PutUint64((*data)[loc:loc+8], i)
	*location += 8
}

func (this *Float64) get(data *[]byte, location *int) interface{} {
	loc := *location
	result := binary.BigEndian.Uint64((*data)[loc : loc+8])
	*location += 8
	return math.Float64frombits(result)
}
