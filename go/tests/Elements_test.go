package tests

import (
	"errors"
	. "github.com/saichler/l8test/go/infra/t_resources"
	"github.com/saichler/serializer/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
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
