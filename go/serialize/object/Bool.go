package object

func addBool(any interface{}, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 1)
	b := any.(bool)
	if b {
		(*data)[*location] = 1
	}
	*location++
}

func getBool(data *[]byte, location *int) interface{} {
	b := (*data)[*location]
	*location++
	if b == 1 {
		return true
	}
	return false
}
