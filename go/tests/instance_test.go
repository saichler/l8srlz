package tests

import (
	"fmt"
	"github.com/saichler/serializer/go/serialize/common"
	"github.com/saichler/serializer/go/serialize/introspect"
	instance "github.com/saichler/serializer/go/serialize/notifications"
	"github.com/saichler/serializer/go/types"
	"github.com/saichler/shared/go/share/registry"
	"github.com/saichler/shared/go/tests"
	"testing"
)

var _introspect common.IIntrospect

func instanceOf(id string, root interface{}, t *testing.T) (interface{}, bool) {
	ins, err := instance.InstanceOf(id, _introspect)
	if err != nil {
		t.Fail()
		fmt.Println("failed with id: ", id, err)
		return nil, false
	}

	v, err := ins.Get(root)
	if err != nil {
		t.Fail()
		fmt.Println("failed with get: ", id, err)
		return nil, false
	}
	return v, true
}

func TestInstance(t *testing.T) {
	_introspect = introspect.NewIntrospect(registry.NewRegistry(), log)
	node, err := _introspect.Inspect(&tests.TestProto{})
	if err != nil {
		fmt.Println("1", err)
		t.Fail()
		return
	}
	_introspect.AddDecorator(types.DecoratorType_Primary, []string{"MyString"}, node)

	id := "testproto<{24}Hello>"
	v, ok := instanceOf(id, nil, t)
	if !ok {
		return
	}

	mytest := v.(*tests.TestProto)
	if mytest.MyString != "Hello" {
		t.Fail()
		fmt.Println("Expected Hello but got ", mytest.MyString)
	}

	mytest.MyFloat64 = 128.128
	id = "testproto.myfloat64"
	v, ok = instanceOf(id, mytest, t)
	if !ok {
		return
	}

	f := v.(float64)
	if f != mytest.MyFloat64 {
		t.Fail()
		fmt.Println("float64 failed:", mytest.MyFloat64, "!=", f)
		return
	}

	mytest.MySingle = &tests.TestProtoSub{MyString: "Hello"}

	id = "testproto.mysingle.mystring"
	v, ok = instanceOf(id, mytest, t)
	if !ok {
		return
	}
	s := v.(string)
	if s != mytest.MySingle.MyString {
		t.Fail()
		fmt.Println("sum model string failed:", mytest.MySingle.MyString, "!=", f)
		return
	}

	/*
		myInstsnce:=model.MyTestModel{
			MyString: "Hello",
			MySingle: &model.MyTestSubModelSingle{MyString: "World"},
		}

		instance,_:=instance.InstanceOf("mytestmodel.mysingle.mystring",introspect.DefaultIntrospect)

		//Getting a value
		v,_:=instance.Get(myInstsnce)
		//Creating another instance
		myOtherInstance:=model.MyTestModel{}
		//Setting the value we fetched from the original instance
		instance.Set(myOtherInstance,"Metadata")

	*/
}
