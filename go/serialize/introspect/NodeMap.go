package introspect

import (
	"github.com/saichler/serializer/go/types"
	"github.com/saichler/shared/go/share/maps"
	"reflect"
)

var _node *types.Node
var nodeType = reflect.TypeOf(_node)

type NodeMap struct {
	impl *maps.SyncMap
}

func NewIntrospectNodeMap() *NodeMap {
	m := &NodeMap{}
	m.impl = maps.NewSyncMap()
	return m
}

func (this *NodeMap) Put(key string, value *types.Node) bool {
	return this.impl.Put(key, value)
}

func (this *NodeMap) Get(key string) (*types.Node, bool) {
	value, ok := this.impl.Get(key)
	if value != nil {
		return value.(*types.Node), ok
	}
	return nil, ok
}

func (this *NodeMap) Contains(key string) bool {
	return this.impl.Contains(key)
}

func (this *NodeMap) NodesList(filter func(v interface{}) bool) []*types.Node {
	return this.impl.ValuesAsList(nodeType, filter).([]*types.Node)
}

func (this *NodeMap) Iterate(do func(k, v interface{})) {
	this.impl.Iterate(do)
}
