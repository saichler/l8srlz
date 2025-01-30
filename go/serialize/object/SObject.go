package object

import (
	"github.com/saichler/shared/go/share/interfaces"
	"reflect"
)

type SObject struct {
	data     []byte
	location int
	typeName string
	registry interfaces.IRegistry
}

type Primitive interface {
	add(interface{}) ([]byte, int)
	get([]byte, int) (interface{}, int)
}

type Complex interface {
	add(interface{}) ([]byte, int)
	get([]byte, int, string, interfaces.IRegistry) (interface{}, int)
}

var primitives = make(map[reflect.Kind]Primitive)
var complex = make(map[reflect.Kind]Complex)

var sizeObjectType = &Int32{}
var stringObjectType = &String{}

var Log interfaces.ILogger

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

func New(data []byte, location int, typeName string, registry interfaces.IRegistry) *SObject {
	obj := &SObject{}
	obj.data = data
	obj.location = location
	obj.registry = registry
	obj.typeName = typeName
	return obj
}

func (obj *SObject) Data() []byte {
	return obj.data
}

func (obj *SObject) Location() int {
	return obj.location
}

func (obj *SObject) Add(any interface{}) error {
	kind := reflect.ValueOf(any).Kind()
	p, pOK := primitives[kind]
	c, cOK := complex[kind]

	if !pOK && !cOK {
		return Log.Error("Did not find any Object for kind", kind.String())
	}

	obj.addKind(kind)
	var b []byte
	var l int

	if pOK {
		b, l = p.add(any)
	} else {
		b, l = c.add(any)
	}
	obj.location += l
	obj.data = append(obj.data, b...)
	return nil
}

func (obj *SObject) Get() (interface{}, error) {
	kind := obj.getKind()
	p, pOK := primitives[kind]
	c, cOK := complex[kind]

	if !pOK && !cOK {
		return nil, Log.Error("Did not find any Object for kind", kind.String())
	}

	var d interface{}
	var l int

	if pOK {
		d, l = p.get(obj.data, obj.location)
	} else {
		d, l = c.get(obj.data, obj.location, obj.typeName, obj.registry)
	}

	obj.location += l
	return d, nil
}

func (obj *SObject) addKind(kind reflect.Kind) {
	b, l := sizeObjectType.add(int32(kind))
	obj.location += l
	obj.data = append(obj.data, b...)
}

func (obj *SObject) getKind() reflect.Kind {
	i, l := sizeObjectType.get(obj.data, obj.location)
	obj.location += l
	return reflect.Kind(i.(int32))
}
