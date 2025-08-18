package object

import (
	"encoding/base64"
	"errors"
	"github.com/saichler/l8types/go/ifs"
	"reflect"
)

type Object struct {
	data     *[]byte
	location *int
	registry ifs.IRegistry
}

type Primitive interface {
	add(interface{}, *[]byte, *int)
	get(*[]byte, *int) interface{}
}

type Complex interface {
	add(interface{}, *[]byte, *int) error
	get(*[]byte, *int, ifs.IRegistry) (interface{}, error)
}

var _Int = &Int{}
var _UInt32 = &UInt32{}
var _UInt64 = &UInt64{}
var _Int32 = &Int32{}
var _Int64 = &Int64{}
var _Float32 = &Float32{}
var _Float64 = &Float64{}
var _String = &String{}
var _Bool = &Bool{}
var _Struct = &Struct{}
var _Slice = &Slice{}
var _Map = &Map{}

func NewEncode() *Object {
	obj := &Object{}
	data := make([]byte, 1024)
	location := 0
	obj.data = &data
	obj.location = &location
	return obj
}

func NewDecode(data []byte, location int, registry ifs.IRegistry) *Object {
	obj := &Object{}
	obj.data = &data
	obj.location = &location
	obj.registry = registry
	return obj
}

func newDecode(data *[]byte, location *int, registry ifs.IRegistry) *Object {
	obj := &Object{}
	obj.data = data
	obj.location = location
	obj.registry = registry
	return obj
}

func (this *Object) Data() []byte {
	return (*this.data)[0:*this.location]
}

func (this *Object) Location() int {
	return *this.location
}

func (this *Object) Add(any interface{}) error {

	switch v := any.(type) {
	case int:
		this.addKind(reflect.Int)
		_Int.add(v, this.data, this.location)
		return nil
	case uint32:
		this.addKind(reflect.Uint32)
		_UInt32.add(v, this.data, this.location)
		return nil
	case uint64:
		this.addKind(reflect.Uint64)
		_UInt64.add(v, this.data, this.location)
		return nil
	case int32:
		this.addKind(reflect.Int32)
		_Int32.add(v, this.data, this.location)
		return nil
	case int64:
		this.addKind(reflect.Int64)
		_Int64.add(v, this.data, this.location)
		return nil
	case float32:
		this.addKind(reflect.Float32)
		_Float32.add(v, this.data, this.location)
		return nil
	case float64:
		this.addKind(reflect.Float64)
		_Float64.add(v, this.data, this.location)
		return nil
	case string:
		this.addKind(reflect.String)
		_String.add(v, this.data, this.location)
		return nil
	case bool:
		this.addKind(reflect.Bool)
		_Bool.add(v, this.data, this.location)
		return nil
	case Slice:
		this.addKind(reflect.Slice)
		return _Slice.add(v, this.data, this.location)
	case Map:
		this.addKind(reflect.Map)
		return _Map.add(v, this.data, this.location)
	default:
		kind := reflect.ValueOf(any).Kind()
		switch kind {
		case reflect.Ptr:
			this.addKind(reflect.Ptr)
			return _Struct.add(v, this.data, this.location)
		case reflect.Slice:
			this.addKind(reflect.Slice)
			return _Slice.add(v, this.data, this.location)
		case reflect.Map:
			this.addKind(reflect.Map)
			return _Map.add(v, this.data, this.location)
		}
	}

	kind := reflect.ValueOf(any).Kind()
	//panic("Did not find any Object for kind " + kind.String())
	return errors.New("Did not find any Object for kind " + kind.String())
}

func (this *Object) Get() (interface{}, error) {
	kind := this.getKind()
	switch kind {
	case reflect.Int:
		return _Int.get(this.data, this.location), nil
	case reflect.Uint32:
		return _UInt32.get(this.data, this.location), nil
	case reflect.Uint64:
		return _UInt64.get(this.data, this.location), nil
	case reflect.Int32:
		return _Int32.get(this.data, this.location), nil
	case reflect.Int64:
		return _Int64.get(this.data, this.location), nil
	case reflect.Float32:
		return _Float32.get(this.data, this.location), nil
	case reflect.Float64:
		return _Float64.get(this.data, this.location), nil
	case reflect.String:
		return _String.get(this.data, this.location), nil
	case reflect.Bool:
		return _Bool.get(this.data, this.location), nil
	case reflect.Slice:
		return _Slice.get(this.data, this.location, this.registry)
	case reflect.Map:
		return _Map.get(this.data, this.location, this.registry)
	case reflect.Ptr:
		return _Struct.get(this.data, this.location, this.registry)
	}
	return nil, errors.New("Did not find any Object for kind " + kind.String())
}

func (this *Object) addKind(kind reflect.Kind) {
	_Int32.add(int32(kind), this.data, this.location)
}

func (this *Object) getKind() reflect.Kind {
	i := _Int32.get(this.data, this.location)
	return reflect.Kind(i.(int32))
}

func (this *Object) Base64() string {
	return base64.StdEncoding.EncodeToString(this.Data())
}

/*
func (this *Object) appendBytes(data []byte, l int) {
	if this.location+len(data) > len(this.data) {
		newData := make([]byte, this.location+len(data)+512)
		copy(newData[0:len(this.data)], this.data)
		this.data = newData
	}
	copy(this.data[this.location:this.location+l], data)
	this.location += l
}*/

func FromBase64(b64 string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(b64)

}

func DataOf(elem interface{}) ([]byte, error) {
	if elem == nil {
		return nil, nil
	}
	obj := NewEncode()
	err := obj.Add(elem)
	return obj.Data(), err
}

func ElemOf(data []byte, r ifs.IRegistry) (interface{}, error) {
	if data == nil {
		return nil, nil
	}
	location := 0
	obj := NewDecode(data, location, r)
	return obj.Get()
}

func checkAndEnlarge(data *[]byte, location *int, need int) {
	if *location+need > len(*data) {
		// Exponential growth with minimum threshold
		newCap := len(*data) * 2
		if newCap < *location+need+512 {
			newCap = *location + need + 512
		}
		tmp := make([]byte, newCap)
		copy(tmp, *data)
		*data = tmp
	}
}
