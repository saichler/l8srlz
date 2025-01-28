package introspect

import (
	"github.com/saichler/serializer/go/serialize/common"
	"github.com/saichler/serializer/go/types"
	"reflect"
)

func (this *Introspect) addAttribute(node *types.Node, _type reflect.Type, _fieldName string) *types.Node {
	this.registry.RegisterType(_type)
	if node != nil && node.Attributes == nil {
		node.Attributes = make(map[string]*types.Node)
	}

	subNode := &types.Node{}
	subNode.TypeName = _type.Name()
	subNode.Parent = node
	subNode.FieldName = _fieldName

	if node != nil {
		node.Attributes[subNode.FieldName] = subNode
	}
	return subNode
}

func (this *Introspect) addNode(_type reflect.Type, _parent *types.Node, _fieldName string) (*types.Node, bool) {
	exist, ok := this.typeToNode.Get(_type.Name())
	if ok && !common.IsLeaf(exist) {
		clone := this.cloner.Clone(exist).(*types.Node)
		clone.Parent = _parent
		clone.FieldName = _fieldName
		clone.CachedKey = ""
		nodePath := NodeKey(clone)
		this.pathToNode.Put(nodePath, clone)
		return clone, true
	}

	node := this.addAttribute(_parent, _type, _fieldName)
	nodePath := NodeKey(node)
	_, ok = this.pathToNode.Get(nodePath)
	if ok {
		return nil, false
	}
	this.pathToNode.Put(nodePath, node)
	if _type.Kind() == reflect.Struct {
		this.typeToNode.Put(node.TypeName, node)
	}
	return node, false
}

func (this *Introspect) inspectStruct(_type reflect.Type, _parent *types.Node, _fieldName string) *types.Node {
	localNode, isClone := this.addNode(_type, _parent, _fieldName)
	if isClone {
		return localNode
	}
	this.registry.RegisterType(_type)
	for index := 0; index < _type.NumField(); index++ {
		field := _type.Field(index)
		if common.IgnoreName(field.Name) {
			continue
		}
		if field.Type.Kind() == reflect.Slice {
			this.inspectSlice(field.Type, localNode, field.Name)
		} else if field.Type.Kind() == reflect.Map {
			this.inspectMap(field.Type, localNode, field.Name)
		} else if field.Type.Kind() == reflect.Ptr {
			subnode := this.inspectPtr(field.Type.Elem(), localNode, field.Name)
			this.typeToNode.Put(subnode.TypeName, subnode)
		} else {
			this.addNode(field.Type, localNode, field.Name)
		}
	}
	this.addTableView(localNode)
	return localNode
}

func (this *Introspect) inspectPtr(_type reflect.Type, _parent *types.Node, _fieldName string) *types.Node {
	switch _type.Kind() {
	case reflect.Struct:
		return this.inspectStruct(_type, _parent, _fieldName)
	}
	panic("unknown ptr kind " + _type.Kind().String())
}

func (this *Introspect) inspectMap(_type reflect.Type, _parent *types.Node, _fieldName string) *types.Node {
	if _type.Elem().Kind() == reflect.Ptr && _type.Elem().Elem().Kind() == reflect.Struct {
		subNode := this.inspectStruct(_type.Elem().Elem(), _parent, _fieldName)
		subNode.IsMap = true
		_parent.Attributes[_fieldName] = subNode
		return subNode
	} else {
		subNode, _ := this.addNode(_type.Elem(), _parent, _fieldName)
		subNode.IsMap = true
		return subNode
	}
}

func (this *Introspect) inspectSlice(_type reflect.Type, _parent *types.Node, _fieldName string) *types.Node {
	if _type.Elem().Kind() == reflect.Ptr && _type.Elem().Elem().Kind() == reflect.Struct {
		subNode := this.inspectStruct(_type.Elem().Elem(), _parent, _fieldName)
		subNode.IsSlice = true
		_parent.Attributes[_fieldName] = subNode
		return subNode
	} else {
		subNode, _ := this.addNode(_type.Elem(), _parent, _fieldName)
		subNode.IsSlice = true
		return subNode
	}
}

func (this *Introspect) printDo(key, val interface{}) {
	node := val.(*types.Node)
	this.log.Debug(key, "-", node.TypeName, ", map=", node.IsMap, ", slice=", node.IsSlice, ", leaf=", common.IsLeaf(node))
}
