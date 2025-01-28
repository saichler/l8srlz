package common

import (
	"github.com/saichler/serializer/go/types"
	"github.com/saichler/shared/go/share/interfaces"
	"reflect"
)

type IIntrospect interface {
	Inspect(interface{}) (*types.Node, error)
	Node(string) (*types.Node, bool)
	NodeByType(p reflect.Type) (*types.Node, bool)
	NodeByTypeName(string) (*types.Node, bool)
	NodeByValue(interface{}) (*types.Node, bool)
	Nodes(bool, bool) []*types.Node
	Print()
	Registry() interfaces.IRegistry
	Kind(*types.Node) reflect.Kind
	Clone(interface{}) interface{}
	AddDecorator(types.DecoratorType, interface{}, *types.Node)
	DecoratorOf(types.DecoratorType, *types.Node) interface{}
	TableView(string) (*types.TableView, bool)
	TableViews() []*types.TableView
}
