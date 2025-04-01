package object

import (
	"errors"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
	"reflect"
)

type Elements struct {
	elements []*Element
}

type Element struct {
	element interface{}
	key     interface{}
	error   error
}

func New(err error, any interface{}) *Elements {
	result := &Elements{}
	result.elements = make([]*Element, 1)
	result.elements[0] = &Element{}
	if err != nil {
		result.elements[0].error = err
	}

	v := reflect.ValueOf(any)
	if v.IsValid() {
		if v.Kind() == reflect.Slice {
			result.elements = nil
			for i := 0; i < v.Len(); i++ {
				result.Add(v.Index(i).Interface(), i, nil)
			}
		} else if v.Kind() == reflect.Map {
			result.elements = nil
			keys := v.MapKeys()
			for _, key := range keys {
				result.Add(v.MapIndex(key).Interface(), key.Interface(), nil)
			}
		} else {
			result.elements[0].element = any
		}
	}

	return result
}

func NewError(err string) *Elements {
	return New(errors.New(err), nil)
}

func (this *Elements) Query() common.IQuery {
	return nil
}

func (this *Elements) Add(elem interface{}, key interface{}, err error) {
	mobject := &Element{element: elem, key: key, error: err}
	if this.elements == nil {
		this.elements = make([]*Element, 0)
	}
	this.elements = append(this.elements, mobject)
}

func (this *Elements) Elements() []interface{} {
	result := make([]interface{}, len(this.elements))
	for i, o := range this.elements {
		result[i] = o.element
	}
	return result
}

func (this *Elements) Element() interface{} {
	return this.elements[0].element
}

func (this *Elements) Keys() []interface{} {
	result := make([]interface{}, len(this.elements))
	for i, o := range this.elements {
		result[i] = o.key
	}
	return result
}

func (this *Elements) Key() interface{} {
	return this.elements[0].key
}

func (this *Elements) Errors() []error {
	result := make([]error, len(this.elements))
	for i, o := range this.elements {
		result[i] = o.error
	}
	return result
}

func (this *Elements) Error() error {
	return this.elements[0].error
}

func (this *Elements) Serialize() (*types.Elements, error) {
	result := &types.Elements{}
	result.ElementList = make([]*types.Element, len(this.elements))
	var err error
	for i, o := range this.elements {
		mo := &types.Element{}
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
		result.ElementList[i] = mo
	}
	return result, nil
}

func (this *Elements) Deserialize(objs *types.Elements, r common.IRegistry) error {
	this.elements = make([]*Element, len(objs.ElementList))
	var err error
	for i, o := range objs.ElementList {
		this.elements[i] = &Element{}
		this.elements[i].element, err = ElemOf(o.Data, r)
		if err != nil {
			return err
		}
		this.elements[i].key, err = ElemOf(o.Key, r)
		if err != nil {
			return err
		}
		if o.ErrorMessage != "" {
			this.elements[i].error = errors.New(o.ErrorMessage)
		}
	}
	return nil
}
