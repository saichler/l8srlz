package object

import (
	"encoding/binary"
	"math"
)

func addFloat64(f float64, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 8)
	i := math.Float64bits(f)
	loc := *location
	binary.BigEndian.PutUint64((*data)[loc:loc+8], i)
	*location += 8
}

func getFloat64(data *[]byte, location *int) float64 {
	loc := *location
	result := binary.BigEndian.Uint64((*data)[loc : loc+8])
	*location += 8
	return math.Float64frombits(result)
}
