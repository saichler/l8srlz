package introspect

import (
	"github.com/saichler/serializer/go/types"
	"github.com/saichler/shared/go/share/string_utils"
)

func (this *Introspect) AddDecorator(decoratorType types.DecoratorType, any interface{}, node *types.Node) {
	s := string_utils.New()
	s.TypesPrefix = true
	str := s.StringOf(any)
	if node.Decorators == nil {
		node.Decorators = make(map[int32]string)
	}
	node.Decorators[int32(decoratorType)] = str
}

func (this *Introspect) DecoratorOf(decoratorType types.DecoratorType, node *types.Node) interface{} {
	decValue := node.Decorators[int32(decoratorType)]
	v := string_utils.InstanceOf(decValue, this.registry)
	return v
}
