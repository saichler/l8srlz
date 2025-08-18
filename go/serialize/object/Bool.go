package object

func addBool(b bool, data *[]byte, location *int) {
	checkAndEnlarge(data, location, 1)
	if b {
		(*data)[*location] = 1
	}
	*location++
}

func getBool(data *[]byte, location *int) bool {
	b := (*data)[*location]
	*location++
	if b == 1 {
		return true
	}
	return false
}
