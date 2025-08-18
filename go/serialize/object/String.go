package object

func addString(str string, data *[]byte, location *int) {
	addInt32(int32(len(str)), data, location)
	checkAndEnlarge(data, location, len(str))
	copy((*data)[*location:*location+len(str)], str)
	*location += len(str)
}

func getString(data *[]byte, location *int) string {
	l := getInt32(data, location)
	size := int(l)
	s := string((*data)[*location : *location+size])
	*location += size
	return s
}
