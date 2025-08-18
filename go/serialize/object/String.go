package object

func addString(any interface{}, data *[]byte, location *int) {
	str := any.(string)
	addInt32(int32(len(str)), data, location)
	checkAndEnlarge(data, location, len(str))
	copy((*data)[*location:*location+len(str)], str)
	*location += len(str)
}

func getString(data *[]byte, location *int) interface{} {
	l := getInt32(data, location)
	size := int(l.(int32))
	s := string((*data)[*location : *location+size])
	*location += size
	return s
}
