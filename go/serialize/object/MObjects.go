package object

import (
	"errors"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
	"reflect"
)

type MObjects struct {
	objects []*MObject
}

type MObject struct {
	element interface{}
	key     interface{}
	error   error
}

func New(err error, any interface{}) *MObjects {
	result := &MObjects{}
	result.objects = make([]*MObject, 1)
	result.objects[0] = &MObject{}
	if err != nil {
		result.objects[0].error = err
	}

	v := reflect.ValueOf(any)
	if v.IsValid() {
		if v.Kind() == reflect.Slice {
			for i := 0; i < v.Len(); i++ {
				result.Add(v.Index(i).Interface(), i, nil)
			}
		} else if v.Kind() == reflect.Map {
			keys := v.MapKeys()
			for _, key := range keys {
				result.Add(v.MapIndex(key).Interface(), key.Interface(), nil)
			}
		} else {
			result.Add(v.Interface(), nil, nil)
		}
	}

	return result
}

func NewError(err string) *MObjects {
	return New(errors.New(err), nil)
}

func (this *MObjects) Add(elem interface{}, key interface{}, err error) {
	mobject := &MObject{element: elem, key: key, error: err}
	if this.objects == nil {
		this.objects = make([]*MObject, 0)
	}
	this.objects = append(this.objects, mobject)
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

func (this *MObjects) Serialize() (*types.MObjects, error) {
	result := &types.MObjects{}
	result.Objects = make([]*types.MObject, len(this.objects))
	var err error
	for i, o := range this.objects {
		mo := &types.MObject{}
		mo.Key, err = DataOf(o.key)
		if err != nil {
			return nil, err
		}
		mo.Data, err = DataOf(o.element)
		if err != nil {
			return nil, err
		}
		if o.error != nil {
			mo.ErrorMessage = o.error.Error()
		}
		result.Objects[i] = mo
	}
	return result, nil
}

func (this *MObjects) Deserialize(objs *types.MObjects, r common.IRegistry) error {
	this.objects = make([]*MObject, len(objs.Objects))
	var err error
	for i, o := range objs.Objects {
		this.objects[i] = &MObject{}
		this.objects[i].element, err = ElemOf(o.Data, r)
		if err != nil {
			return err
		}
		this.objects[i].key, err = ElemOf(o.Data, r)
		if err != nil {
			return err
		}
		if o.ErrorMessage != "" {
			this.objects[i].error = errors.New(o.ErrorMessage)
		}
	}
	return nil
}
