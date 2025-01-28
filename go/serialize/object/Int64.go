package object

type Int64 struct{}

func (this *Int64) add(any interface{}) ([]byte, int) {
	i := any.(int64)
	data := make([]byte, 8)
	data[0] = byte((i >> 56) & 0xff)
	data[1] = byte((i >> 48) & 0xff)
	data[2] = byte((i >> 40) & 0xff)
	data[3] = byte((i >> 32) & 0xff)
	data[4] = byte((i >> 24) & 0xff)
	data[5] = byte((i >> 16) & 0xff)
	data[6] = byte((i >> 8) & 0xff)
	data[7] = byte((i) & 0xff)
	return data, 8
}

func (this *Int64) get(data []byte, location int) (interface{}, int) {
	var result int64
	result = int64(0xff&data[location])<<56 |
		int64(0xff&data[location+1])<<48 |
		int64(0xff&data[location+2])<<40 |
		int64(0xff&data[location+3])<<32 |
		int64(0xff&data[location+4])<<24 |
		int64(0xff&data[location+5])<<16 |
		int64(0xff&data[location+6])<<8 |
		int64(0xff&data[location+7])
	return result, 8
}
