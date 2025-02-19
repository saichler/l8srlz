package object

import (
	"github.com/saichler/shared/go/share/interfaces"
	"reflect"
)

type Object struct {
	data     []byte
	location int
	typeName string
	registry interfaces.IRegistry
	log      interfaces.ILogger
}

type Primitive interface {
	add(interface{}) ([]byte, int)
	get([]byte, int) (interface{}, int)
}

type Complex interface {
	add(interface{}, interfaces.ILogger) ([]byte, int)
	get([]byte, int, string, interfaces.IRegistry, interfaces.ILogger) (interface{}, int)
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

func NewEncode(data []byte, location int, log interfaces.ILogger) *Object {
	return NewDecode(data, location, "", nil, log)
}

func NewDecode(data []byte, location int, typeName string, registry interfaces.IRegistry, log interfaces.ILogger) *Object {
	obj := &Object{}
	obj.data = data
	obj.location = location
	obj.registry = registry
	obj.log = log
	obj.typeName = typeName
	return obj
}

func (obj *Object) Data() []byte {
	return obj.data
}

func (obj *Object) Location() int {
	return obj.location
}

func (obj *Object) Add(any interface{}) error {
	kind := reflect.ValueOf(any).Kind()
	p, pOK := primitives[kind]
	c, cOK := complex[kind]

	if !pOK && !cOK {
		return obj.log.Error("Did not find any Object for kind", kind.String())
	}

	obj.addKind(kind)
	var b []byte
	var l int

	if pOK {
		b, l = p.add(any)
	} else {
		b, l = c.add(any, obj.log)
	}
	obj.location += l
	obj.data = append(obj.data, b...)
	return nil
}

func (obj *Object) Get() (interface{}, error) {
	kind := obj.getKind()
	p, pOK := primitives[kind]
	c, cOK := complex[kind]

	if !pOK && !cOK {
		return nil, obj.log.Error("Did not find any Object for kind", kind.String())
	}

	var d interface{}
	var l int

	if pOK {
		d, l = p.get(obj.data, obj.location)
	} else {
		d, l = c.get(obj.data, obj.location, obj.typeName, obj.registry, obj.log)
	}

	obj.location += l
	return d, nil
}

func (obj *Object) addKind(kind reflect.Kind) {
	b, l := sizeObjectType.add(int32(kind))
	obj.location += l
	obj.data = append(obj.data, b...)
}

func (obj *Object) getKind() reflect.Kind {
	i, l := sizeObjectType.get(obj.data, obj.location)
	obj.location += l
	return reflect.Kind(i.(int32))
}
