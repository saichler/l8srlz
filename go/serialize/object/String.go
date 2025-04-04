package object

type String struct{}

func (this *String) add(any interface{}) ([]byte, int) {
	str := any.(string)
	obj := NewEncode()
	obj.appendBytes(sizeObjectType.add(int32(len(str))))
	obj.appendBytes([]byte(str), len(str))
	return obj.Data(), obj.Location()
}

func (this *String) get(data []byte, location int) (interface{}, int) {
	l, _ := sizeObjectType.get(data, location)
	size := l.(int32)
	location += 4
	s := string(data[location : location+int(size)])
	return s, len(s) + 4
}
