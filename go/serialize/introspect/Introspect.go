package introspect

import (
	"github.com/saichler/serializer/go/serialize/common"
	"github.com/saichler/serializer/go/types"
	"github.com/saichler/shared/go/share/deep_clone"
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/maps"
	"github.com/saichler/shared/go/share/string_utils"
	"reflect"
	"strings"
)

type Introspect struct {
	pathToNode *NodeMap
	typeToNode *NodeMap
	registry   interfaces.IRegistry
	log        interfaces.ILogger
	cloner     *deep_clone.Cloner
	tableViews *maps.SyncMap
}

func NewIntrospect(registry interfaces.IRegistry, logger interfaces.ILogger) *Introspect {
	i := &Introspect{}
	i.registry = registry
	i.log = logger
	i.cloner = deep_clone.NewCloner()
	i.pathToNode = NewIntrospectNodeMap()
	i.typeToNode = NewIntrospectNodeMap()
	i.tableViews = maps.NewSyncMap()
	return i
}

func (this *Introspect) Registry() interfaces.IRegistry {
	return this.registry
}

func (this *Introspect) Inspect(any interface{}) (*types.Node, error) {
	if any == nil {
		return nil, this.log.Error("Cannot introspect a nil value")
	}

	_, t := common.ValueAndType(any)
	if t.Kind() == reflect.Slice && t.Kind() == reflect.Map {
		t = t.Elem().Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, this.log.Error("Cannot introspect a value that is not a struct")
	}
	localNode, ok := this.pathToNode.Get(strings.ToLower(t.Name()))
	if ok {
		return localNode, nil
	}
	return this.inspectStruct(t, nil, ""), nil
}

func (this *Introspect) Node(path string) (*types.Node, bool) {
	return this.pathToNode.Get(strings.ToLower(path))
}

func (this *Introspect) NodeByValue(any interface{}) (*types.Node, bool) {
	val := reflect.ValueOf(any)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return this.NodeByType(val.Type())
}

func (this *Introspect) NodeByType(typ reflect.Type) (*types.Node, bool) {
	return this.NodeByTypeName(typ.Name())
}

func (this *Introspect) NodeByTypeName(name string) (*types.Node, bool) {
	return this.typeToNode.Get(name)
}

func (this *Introspect) Nodes(onlyLeafs, onlyRoots bool) []*types.Node {
	filter := func(any interface{}) bool {
		n := any.(*types.Node)
		if onlyLeafs && !common.IsLeaf(n) {
			return false
		}
		if onlyRoots && !common.IsRoot(n) {
			return false
		}
		return true
	}

	return this.pathToNode.NodesList(filter)
}

func (this *Introspect) Print() {
	this.pathToNode.Iterate(this.printDo)
}

func (this *Introspect) Kind(node *types.Node) reflect.Kind {
	info, err := this.registry.Info(node.TypeName)
	if err != nil {
		panic(err.Error())
	}
	return info.Type().Kind()
}

func (this *Introspect) Clone(any interface{}) interface{} {
	return this.cloner.Clone(any)
}

func (this *Introspect) addTableView(node *types.Node) {
	tv := &types.TableView{Table: node, Columns: make([]*types.Node, 0), SubTables: make([]*types.Node, 0)}
	for _, attr := range node.Attributes {
		if common.IsLeaf(attr) {
			tv.Columns = append(tv.Columns, attr)
		} else {
			tv.SubTables = append(tv.SubTables, attr)
		}
	}
	this.tableViews.Put(node.TypeName, tv)
}

func (this *Introspect) TableView(name string) (*types.TableView, bool) {
	tv, ok := this.tableViews.Get(name)
	if !ok {
		return nil, ok
	}
	return tv.(*types.TableView), ok
}

func (this *Introspect) TableViews() []*types.TableView {
	list := this.tableViews.ValuesAsList(reflect.TypeOf(&types.TableView{}), nil)
	return list.([]*types.TableView)
}

func NodeKey(node *types.Node) string {
	if node.CachedKey != "" {
		return node.CachedKey
	}
	if node.Parent == nil {
		return strings.ToLower(node.TypeName)
	}
	buff := string_utils.New()
	buff.Add(NodeKey(node.Parent))
	buff.Add(".")
	buff.Add(strings.ToLower(node.FieldName))
	node.CachedKey = buff.String()
	return node.CachedKey
}
