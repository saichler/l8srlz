package object

type String struct{}

func (this *String) add(any interface{}, data *[]byte, location *int) {
	str := any.(string)
	sizeObjectType.add(len(str), data, location)
	copy(*data, str)
	*location += len(str)
}

func (this *String) get(data *[]byte, location *int) interface{} {
	l := sizeObjectType.get(data, location)
	size := int(l.(int32))
	s := string((*data)[*location : *location+int(size)])
	*location += size
	return s
}
