package tests

import (
	"errors"
	"fmt"
	"github.com/saichler/l8srlz/go/serialize/object"
	. "github.com/saichler/l8test/go/infra/t_resources"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestElements(t *testing.T) {
	res, _ := CreateResources(25000, 2, ifs.Info_Level)
	res.Registry().Register(&testtypes.TestProto{})
	elem := CreateTestModelInstance(2)
	elems := object.New(nil, elem)
	data, e := elems.Serialize()
	if e != nil {
		Log.Fail(t, "Failed:", e.Error())
		return
	}
	elems = &object.Elements{}
	e = elems.Deserialize(data, res.Registry())
	if e != nil {
		Log.Fail(t, e.Error())
		return
	}

	elems = object.New(errors.New("Hello Error"), nil)
	data, e = elems.Serialize()
	if e != nil {
		Log.Fail(t, "Failed:", e.Error())
		return
	}
	elems = &object.Elements{}
	e = elems.Deserialize(data, res.Registry())
	if e != nil {
		Log.Fail(t, e.Error())
		return
	}

}

func TestElementsList(t *testing.T) {
	res, _ := CreateResources(25000, 2, ifs.Info_Level)
	res.Registry().Register(&testtypes.TestProto{})
	res.Registry().Register(&testtypes.TestProtoList{})
	elemList := []*testtypes.TestProto{CreateTestModelInstance(2), CreateTestModelInstance(3)}
	elems := object.New(nil, elemList)
	list, err := elems.AsList(res.Registry())
	if err != nil {
		Log.Fail(t, "Failed:", err)
		return
	}
	json, err := protojson.Marshal(list.(proto.Message))
	if err != nil {
		Log.Fail(t, "Failed:", err)
		return
	}
	fmt.Println(string(json))
}
