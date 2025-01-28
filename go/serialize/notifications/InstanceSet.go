package instance

import (
	"errors"
	"github.com/saichler/serializer/go/types"
	"reflect"
)

func (inst *Instance) Set(any interface{}, value interface{}) (interface{}, error) {
	if inst == nil {
		return nil, errors.New("Instance is nil, cannot instantiate")
	}
	if inst.parent == nil {
		if any == nil {
			info, err := inst.introspect.Registry().Info(inst.node.TypeName)
			if err != nil {
				return nil, err
			}
			newAny, err := info.NewInstance()
			if err != nil {
				return nil, err
			}
			any = newAny
		}
		if inst.key != nil {
			inst.SetPrimaryKey(inst.node, any, inst.key)
		}
		return any, nil
	}
	parent, err := inst.parent.Set(any, value)
	if err != nil {
		return nil, err
	}
	parentValue := reflect.ValueOf(parent)
	if parentValue.Kind() == reflect.Ptr {
		parentValue = parentValue.Elem()
	}
	myValue := parentValue.FieldByName(inst.node.FieldName)
	info, err := inst.introspect.Registry().Info(inst.node.TypeName)
	if err != nil {
		return nil, err
	}
	typ := info.Type()
	if inst.node.IsMap {
		return inst.mapSet(myValue)
	} else if inst.node.IsSlice {
		return inst.sliceSet(myValue)
	} else if inst.introspect.Kind(inst.node) == reflect.Struct {
		if !myValue.IsValid() || myValue.IsNil() {
			myValue.Set(reflect.New(typ))
		}
		return myValue.Interface(), err
	} else {
		myValue.Set(reflect.ValueOf(value))
		return value, err
	}
}

func (inst *Instance) sliceSet(myValue reflect.Value) (interface{}, error) {
	index := inst.key.(int)
	info, err := inst.introspect.Registry().Info(inst.node.TypeName)
	if err != nil {
		return nil, err
	}
	typ := info.Type()
	if !myValue.IsValid() || myValue.IsNil() {
		myValue.Set(reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(typ)), index+1, index+1))
	}

	if index >= myValue.Len() {
		newSlice := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(typ)), index+1, index+1)
		for i := 0; i < myValue.Len(); i++ {
			newSlice.Index(i).Set(myValue.Index(i))
		}
		myValue.Set(newSlice)
	}

	sliceValue := reflect.ValueOf(inst.value)
	if inst.introspect.Kind(inst.node) == reflect.Struct && inst.value == nil {
		sliceValue = reflect.New(typ)
	}
	myValue.Index(index).Set(sliceValue)
	return sliceValue.Interface(), err
}

func (inst *Instance) mapSet(myValue reflect.Value) (interface{}, error) {
	info, err := inst.introspect.Registry().Info(inst.node.TypeName)
	if err != nil {
		return nil, err
	}
	typ := info.Type()
	info, err = inst.introspect.Registry().Info(inst.node.KeyTypeName)
	if err != nil {
		return nil, err
	}
	typKey := info.Type()
	if !myValue.IsValid() || myValue.IsNil() {
		myValue.Set(reflect.MakeMap(reflect.MapOf(typKey, reflect.PtrTo(typ))))
	}
	mapKey := reflect.ValueOf(inst.key)
	oldMapValue := myValue.MapIndex(mapKey)
	mapValue := reflect.ValueOf(inst.value)
	if inst.introspect.Kind(inst.node) == reflect.Struct && inst.value == nil {
		if oldMapValue.IsValid() && !oldMapValue.IsNil() {
			mapValue = oldMapValue
		} else {
			mapValue = reflect.New(typ)
		}
	}
	myValue.SetMapIndex(mapKey, mapValue)
	return mapValue.Interface(), err
}

func (inst *Instance) SetPrimaryKey(node *types.Node, any interface{}, anyKey interface{}) {
	if anyKey == nil {
		return
	}
	var fieldsValues []interface{}
	if reflect.ValueOf(anyKey).Kind() == reflect.Slice {
		fieldsValues = anyKey.([]interface{})
	} else {
		fieldsValues = []interface{}{anyKey}
	}
	value := reflect.ValueOf(any)
	if !value.IsValid() {
		return
	}
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return
		}
		value = value.Elem()
	}

	f := inst.introspect.DecoratorOf(types.DecoratorType_Primary, node)
	if f != nil {
		fields := f.([]string)
		for i, attr := range fields {
			fld := value.FieldByName(attr)
			fld.Set(reflect.ValueOf(fieldsValues[i]))
		}
	}
}
