package object

import (
	"encoding/base64"
	"errors"
	"github.com/saichler/l8types/go/ifs"
	"reflect"
)

var MessageSerializer ifs.ISerializer
var TransactionSerializer ifs.ISerializer

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

var primitives = make(map[reflect.Kind]Primitive)
var complex = make(map[reflect.Kind]Complex)

var sizeObjectType = &Int32{}
var stringObjectType = &String{}

func init() {
	primitives[reflect.Int] = &Int{}
	primitives[reflect.Uint32] = &UInt32{}
	primitives[reflect.Uint64] = &UInt64{}
	primitives[reflect.Int32] = &Int32{}
	primitives[reflect.Int64] = &Int64{}
	primitives[reflect.Float32] = &Float32{}
	primitives[reflect.Float64] = &Float64{}
	primitives[reflect.String] = &String{}
	primitives[reflect.Bool] = &Bool{}

	complex[reflect.Ptr] = &Struct{}
	complex[reflect.Slice] = &Slice{}
	complex[reflect.Map] = &Map{}
}

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
	kind := reflect.ValueOf(any).Kind()
	if kind == reflect.Invalid {
		kind = reflect.Ptr
	}
	p, pOK := primitives[kind]
	c, cOK := complex[kind]

	if !pOK && !cOK {
		panic(kind.String())
		return errors.New("Did not find any Object for kind " + kind.String())
	}

	this.addKind(kind)

	var e error

	if pOK {
		p.add(any, this.data, this.location)
	} else {
		e = c.add(any, this.data, this.location)
	}
	return e
}

func (this *Object) Get() (interface{}, error) {
	kind := this.getKind()
	p, pOK := primitives[kind]
	c, cOK := complex[kind]

	if !pOK && !cOK {
		return nil, errors.New("Did not find any Object for kind " + kind.String())
	}

	var d interface{}
	var e error

	if pOK {
		d = p.get(this.data, this.location)
	} else {
		d, e = c.get(this.data, this.location, this.registry)
	}
	return d, e
}

func (this *Object) addKind(kind reflect.Kind) {
	sizeObjectType.add(int32(kind), this.data, this.location)
}

func (this *Object) getKind() reflect.Kind {
	i := sizeObjectType.get(this.data, this.location)
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
		tmp := make([]byte, *location+need+512)
		copy(tmp, *data)
		*data = tmp
	}
}
