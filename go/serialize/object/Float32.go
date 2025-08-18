package object

import (
	"encoding/binary"
	"math"
)

func addFloat32(any interface{}, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 4)
	f := any.(float32)
	i := math.Float32bits(f)
	loc := *location
	binary.BigEndian.PutUint32((*data)[loc:loc+4], i)
	*location += 4
}

func getFloat32(data *[]byte, location *int) interface{} {
	loc := *location
	result := binary.BigEndian.Uint32((*data)[loc : loc+4])
	*location += 4
	return math.Float32frombits(result)
}
