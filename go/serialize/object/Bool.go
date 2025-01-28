package object

type Bool struct{}

func (this *Bool) add(any interface{}) ([]byte, int) {
	b := any.(bool)
	if b {
		return []byte{1}, 1
	}
	return []byte{0}, 1
}

func (this *Bool) get(data []byte, location int) (interface{}, int) {
	b := data[location]
	if b == 1 {
		return true, 1
	}
	return false, 1
}
