package object

type String struct{}

func (this *String) add(any interface{}) ([]byte, int) {
	str := any.(string)
	s, _ := sizeObjectType.add(int32(len(str)))
	s = append(s, []byte(str)...)
	return s, len(s)
}

func (this *String) get(data []byte, location int) (interface{}, int) {
	l, _ := sizeObjectType.get(data, location)
	size := l.(int32)
	location += 4
	s := string(data[location : location+int(size)])
	return s, len(s) + 4
}
