package object

type String struct{}

func (this *String) add(any interface{}, data *[]byte, location *int) {
	str := any.(string)
	_Int32.add(int32(len(str)), data, location)
	checkAndEnlarge(data, location, len(str))
	copy((*data)[*location:*location+len(str)], str)
	*location += len(str)
}

func (this *String) get(data *[]byte, location *int) interface{} {
	l := _Int32.get(data, location)
	size := int(l.(int32))
	s := string((*data)[*location : *location+size])
	*location += size
	return s
}
