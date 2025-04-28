package object

import (
	"errors"
	"github.com/saichler/gsql/go/gsql/interpreter"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
	"reflect"
)

type Elements struct {
	elements        []*Element
	query           common.IQuery
	pquery          *types.Query
	notification    bool
	replicasRequest bool
}

type Element struct {
	element interface{}
	key     interface{}
	error   error
}

func NewQuery(gsql string, resources common.IResources) (*Elements, error) {
	q, e := interpreter.NewQuery(gsql, resources)
	if e != nil {
		return nil, e
	}
	elems := &Elements{pquery: q.Query()}
	return elems, nil
}

func NewNotify(any interface{}) *Elements {
	elems := New(nil, any)
	elems.notification = true
	return elems
}

func NewReplicasRequest(elems common.IElements) *Elements {
	c := clone(elems)
	c.replicasRequest = true
	return c
}

func clone(e common.IElements) *Elements {
	old := e.(*Elements)
	c := &Elements{}
	c.elements = old.elements
	c.query = old.query
	c.pquery = old.pquery
	c.notification = old.notification
	c.replicasRequest = old.replicasRequest
	return c
}

func New(err error, any interface{}) *Elements {

	if reflect.ValueOf(any).Kind() == reflect.Func {
		panic("any is a function, this is probably a mistake")
	}

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

func (this *Elements) Query(resources common.IResources) (common.IQuery, error) {
	var err error
	if this.query == nil && this.pquery != nil {
		this.query, err = interpreter.NewFromQuery(this.pquery, resources)
		if err != nil {
			return nil, err
		}
	}
	return this.query, nil
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

func (this *Elements) Serialize() ([]byte, error) {
	obj := NewEncode()
	obj.Add(len(this.elements))
	var err error

	for _, o := range this.elements {
		err = obj.Add(o.element)
		if err != nil {
			return nil, err
		}
		err = obj.Add(o.key)
		if err != nil {
			return nil, err
		}
		if o.error != nil {
			err = obj.Add(o.error.Error())
		} else {
			err = obj.Add("")
		}
		if err != nil {
			return nil, err
		}
	}
	obj.Add(this.pquery)
	return obj.Data(), nil
}

func (this *Elements) Deserialize(data []byte, r common.IRegistry) error {
	location := 0
	obj := NewDecode(data, location, r)
	s, err := obj.Get()
	if err != nil {
		return err
	}
	size := s.(int)
	this.elements = make([]*Element, size)
	var eMsg interface{}
	for i := 0; i < size; i++ {
		elem := &Element{}
		elem.element, err = obj.Get()
		if err != nil {
			return err
		}
		elem.key, err = obj.Get()
		if err != nil {
			return err
		}
		eMsg, err = obj.Get()
		if err != nil {
			return err
		}
		errMsg := eMsg.(string)
		if errMsg != "" {
			elem.error = errors.New(errMsg)
		}
		this.elements[i] = elem
	}
	pq, err := obj.Get()
	if err != nil {
		return err
	}
	this.pquery, _ = pq.(*types.Query)
	return nil
}

func (this *Elements) Notification() bool {
	return this.notification
}

func (this *Elements) ReplicasRequest() bool {
	return this.replicasRequest
}
