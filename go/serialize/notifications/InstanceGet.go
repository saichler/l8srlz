package instance

import (
	"reflect"
)

func (inst *Instance) getMap(parent reflect.Value) []reflect.Value {
	result := make([]reflect.Value, 0)
	if inst.parent.key != nil {
		myValue := parent.MapIndex(reflect.ValueOf(inst.parent.key))
		if !myValue.IsValid() {
			return result
		}
		if myValue.Kind() == reflect.Ptr {
			if myValue.IsNil() {
				return result
			}
			myValue = myValue.Elem()
		}
		myValue = myValue.FieldByName(inst.node.FieldName)
		result = append(result, myValue)
	} else {
		keys := parent.MapKeys()
		for _, key := range keys {
			value := parent.MapIndex(key)
			if value.Kind() == reflect.Ptr {
				value = value.Elem()
			}
			myValue := value.FieldByName(inst.node.FieldName)
			result = append(result, myValue)
		}
	}
	return result
}

func (inst *Instance) getSlice(parent reflect.Value) []reflect.Value {
	result := make([]reflect.Value, 0)
	if inst.parent.key != nil {
		myValue := parent.Index(inst.parent.key.(int))
		if !myValue.IsValid() {
			return result
		}
		if myValue.Kind() == reflect.Ptr {
			if myValue.IsNil() {
				return result
			}
			myValue = myValue.Elem()
		}
		myValue = myValue.FieldByName(inst.node.FieldName)
		result = append(result, myValue)
	} else {
		for i := 0; i < parent.Len(); i++ {
			value := parent.Index(i)
			if value.Kind() == reflect.Interface {
				value = value.Elem()
			}
			if value.Kind() == reflect.Ptr {
				if value.IsNil() {
					continue
				}
				value = value.Elem()
			}

			myValue := value.FieldByName(inst.node.FieldName)
			result = append(result, myValue)
		}
	}
	return result
}

func (inst *Instance) GetValue(any reflect.Value) []reflect.Value {
	if !any.IsValid() {
		return []reflect.Value{}
	}
	if any.Kind() == reflect.Ptr && any.IsNil() {
		return []reflect.Value{}
	}
	if inst.parent == nil {
		return []reflect.Value{any}
	}

	parents := inst.parent.GetValue(any)
	results := make([]reflect.Value, 0)

	for _, parent := range parents {
		if parent.Kind() == reflect.Ptr {
			parent = parent.Elem()
		}
		if parent.Kind() == reflect.Map {
			mapItems := inst.getMap(parent)
			results = append(results, mapItems...)
		} else if parent.Kind() == reflect.Slice {
			sliceItems := inst.getSlice(parent)
			results = append(results, sliceItems...)
		} else {
			value := parent.FieldByName(inst.node.FieldName)
			results = append(results, value)
		}
	}
	return results
}

func (inst *Instance) Get(any interface{}) (interface{}, error) {
	if any == nil {
		if inst == nil {
			panic("nil inst")
		}
		if inst.introspect == nil {
			panic("nil introspect")
		}
		if inst.introspect.Registry() == nil {
			panic("nil registry")
		}
		info, err := inst.introspect.Registry().Info(inst.node.TypeName)
		if err != nil {
			return nil, err
		}
		n, err := info.NewInstance()
		if inst.key != nil {
			inst.SetPrimaryKey(inst.node, n, inst.key)
		}
		return n, nil
	}
	values := inst.GetValue(reflect.ValueOf(any))
	if !values[0].IsValid() {
		return nil, nil
	}
	if values[0].Kind() == reflect.Ptr && values[0].IsNil() {
		return nil, nil
	}
	return values[0].Interface(), nil
}

func (inst *Instance) GetAsValues(any interface{}) []reflect.Value {
	_interface, _ := inst.Get(any)
	if _interface == nil {
		return []reflect.Value{reflect.ValueOf(inst)}
	}
	value := reflect.ValueOf(_interface)
	if value.Kind() == reflect.Map {
		result := make([]reflect.Value, value.Len())
		keys := value.MapKeys()
		for i, key := range keys {
			item := value.MapIndex(key)
			if item.Kind() == reflect.Interface {
				item = item.Elem()
			}
			result[i] = item
		}
		return result
	} else if value.Kind() == reflect.Slice {
		result := make([]reflect.Value, value.Len())
		for i := 0; i < len(result); i++ {
			item := value.Index(i)
			if item.Kind() == reflect.Interface {
				item = item.Elem()
			}
			result[i] = item
		}
		return result
	} else {
		return []reflect.Value{value}
	}
}
