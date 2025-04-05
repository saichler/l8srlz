package object

type Bool struct{}

func (this *Bool) add(any interface{}, data *[]byte, location *int) {
	b := any.(bool)
	if b {
		(*data)[*location] = 1
	}
	*location++
}

func (this *Bool) get(data *[]byte, location *int) interface{} {
	b := (*data)[*location]
	*location++
	if b == 1 {
		return true
	}
	return false
}
