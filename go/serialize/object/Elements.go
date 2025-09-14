package object

import (
	"errors"
	"reflect"

	"github.com/saichler/gsql/go/gsql/interpreter"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types"
)

type Elements struct {
	elements        []*Element
	query           ifs.IQuery
	pquery          *types.Query
	stats           map[string]int32
	notification    bool
	replicasRequest bool
}

type Element struct {
	element interface{}
	key     interface{}
	error   error
}

func NewQuery(gsql string, resources ifs.IResources) (*Elements, error) {
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

func NewReplicasRequest(elems ifs.IElements) *Elements {
	c := clone(elems)
	c.replicasRequest = true
	return c
}

func clone(e ifs.IElements) *Elements {
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

func NewQueryResult(any interface{}, stats map[string]int32) *Elements {
	elements := New(nil, any)
	elements.stats = stats
	return elements
}

func NewError(err string) *Elements {
	return New(errors.New(err), nil)
}

func (this *Elements) Query(resources ifs.IResources) (ifs.IQuery, error) {
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
	if this.elements == nil || len(this.elements) == 0 {
		return nil
	}
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
	if this.elements == nil || len(this.elements) == 0 {
		return nil
	}
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

	if this.stats == nil || len(this.stats) == 0 {
		this.stats = make(map[string]int32)
		this.stats["Total"] = int32(len(this.elements))
	}
	obj.Add(this.stats)

	obj.Add(this.pquery)
	return obj.Data(), nil
}

func (this *Elements) PQuery() *types.Query {
	return this.pquery
}

func (this *Elements) Deserialize(data []byte, r ifs.IRegistry) error {
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

	st, err := obj.Get()
	if err != nil {
		return err
	}
	this.stats, _ = st.(map[string]int32)

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

func (this *Elements) Append(elements ifs.IElements) {
	if elements == nil {
		return
	}
	if elements.Elements() == nil {
		return
	}
	for _, elem := range elements.Elements() {
		this.Add(elem, nil, nil)
	}
}

func (this *Elements) AsList(r ifs.IRegistry) (interface{}, error) {
	if len(this.elements) == 0 {
		return nil, errors.New("elements is empty")
	}
	if this.elements[0] == nil || this.elements[0].element == nil {
		return nil, errors.New("element is nil")
	}
	listName := reflect.ValueOf(this.elements[0].element).Elem().Type().Name() + "List"
	info, err := r.Info(listName)

	if err != nil {
		return this.elements[0].element, nil
	}

	listItem, err := info.NewInstance()
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(listItem).Elem()
	f := v.FieldByName("List")
	newList := reflect.MakeSlice(f.Type(), len(this.elements), len(this.elements))
	for i := 0; i < len(this.elements); i++ {
		newList.Index(i).Set(reflect.ValueOf(this.elements[i].element))
	}
	f.Set(newList)

	f = v.FieldByName("Stats")
	if f.IsValid() && f.CanSet() {
		f.Set(reflect.ValueOf(this.stats))
	}

	return listItem, nil
}

func (this *Elements) IsFilterMode() bool {
	if this.pquery == nil && (this.elements != nil || len(this.elements) == 1) {
		return true
	}
	return false
}
