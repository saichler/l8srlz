package object

import (
	"errors"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
)

type MObjects struct {
	objects []*MObject
}

type MObject struct {
	element interface{}
	key     interface{}
	error   error
}

func New(err string, elem interface{}) *MObjects {
	result := &MObjects{}
	result.objects = make([]*MObject, 1)
	result.objects[0] = &MObject{}
	if err != "" {
		result.objects[0].error = errors.New(err)
	}
	result.objects[0].element = elem
	return result
}

func NewError(err string) *MObjects {
	return New(err, nil)
}

func (this *MObjects) Elements() []interface{} {
	result := make([]interface{}, len(this.objects))
	for i, o := range this.objects {
		result[i] = o.element
	}
	return result
}

func (this *MObjects) Element() interface{} {
	return this.objects[0].element
}

func (this *MObjects) Keys() []interface{} {
	result := make([]interface{}, len(this.objects))
	for i, o := range this.objects {
		result[i] = o.key
	}
	return result
}

func (this *MObjects) Key() interface{} {
	return this.objects[0].key
}

func (this *MObjects) Errors() []error {
	result := make([]error, len(this.objects))
	for i, o := range this.objects {
		result[i] = o.error
	}
	return result
}

func (this *MObjects) Error() error {
	return this.objects[0].error
}

func (this *MObjects) Serialize() *types.MObjects {
	result := &types.MObjects{}
	result.Objects = make([]*types.MObject, len(this.objects))
	for i, o := range this.objects {
		mo := &types.MObject{}
		mo.Key = DataOf(o.key)
		mo.Data = DataOf(o.element)
		if o.error != nil {
			mo.ErrorMessage = o.error.Error()
		}
		result.Objects[i] = mo
	}
	return result
}

func (this *MObjects) Deserialize(objs *types.MObjects, r common.IRegistry) {
	this.objects = make([]*MObject, len(objs.Objects))
	for i, o := range objs.Objects {
		this.objects[i] = &MObject{}
		this.objects[i].element = ElemOf(o.Data, r)
		this.objects[i].key = ElemOf(o.Data, r)
		if o.ErrorMessage != "" {
			this.objects[i].error = errors.New(o.ErrorMessage)
		}
	}
}
