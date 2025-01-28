package object

type UInt32 struct{}

func (this *UInt32) add(any interface{}) ([]byte, int) {
	i := any.(uint32)
	bytes := make([]byte, 4)
	bytes[3] = byte(i)
	bytes[2] = byte(i >> 8)
	bytes[1] = byte(i >> 16)
	bytes[0] = byte(i >> 24)
	return bytes, 4
}

func (this *UInt32) get(data []byte, location int) (interface{}, int) {
	var result uint32
	v1 := uint32(data[location]) << 24
	v2 := uint32(data[location+1]) << 16
	v3 := uint32(data[location+2]) << 8
	v4 := uint32(data[location+3])
	result = v1 + v2 + v3 + v4
	return result, 4
}
