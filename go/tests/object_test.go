package tests

import (
	"github.com/saichler/serializer/go/serialize/object"
	"github.com/saichler/shared/go/tests"
	"strconv"
	"testing"
)

func testType(val interface{}) (interface{}, error) {
	e := object.NewEncode([]byte{}, 0, log)
	err := e.Add(val)
	if err != nil {
		return nil, err
	}
	data := e.Data()
	e = object.NewDecode(data, 0, "", globals.Registry(), log)
	dval, err := e.Get()
	if err != nil {
		return nil, err
	}
	return dval, nil
}

func TestInt(t *testing.T) {
	val := int(-39)
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(int)
	if res != val {
		log.Error("Value do not match:", strconv.Itoa(int(res)), " and ", strconv.Itoa(int(val)))
		t.Fail()
	}
}

func TestInt32(t *testing.T) {
	val := int32(-39)
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(int32)
	if res != val {
		log.Error("Value do not match:", strconv.Itoa(int(res)), " and ", strconv.Itoa(int(val)))
		t.Fail()
	}
}

func TestInt64(t *testing.T) {
	val := int64(-39)
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(int64)
	if res != val {
		log.Error("Value do not match:", strconv.Itoa(int(res)), " and ", strconv.Itoa(int(val)))
		t.Fail()
	}
}

func TestUInt32(t *testing.T) {
	val := uint32(39)
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(uint32)
	if res != val {
		log.Error("Value do not match:", strconv.Itoa(int(res)), " and ", strconv.Itoa(int(val)))
		t.Fail()
	}
}

func TestUInt64(t *testing.T) {
	val := uint64(39)
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(uint64)
	if res != val {
		log.Error("Value do not match:", strconv.Itoa(int(res)), " and ", strconv.Itoa(int(val)))
		t.Fail()
	}
}

func TestFloat64(t *testing.T) {
	val := float64(39.39)
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(float64)
	if res != val {
		log.Error("Value do not match:", strconv.Itoa(int(res)), " and ", strconv.Itoa(int(val)))
		t.Fail()
	}
}

func TestFloat32(t *testing.T) {
	val := float32(39.39)
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(float32)
	if res != val {
		log.Error("Value do not match:", strconv.Itoa(int(res)), " and ", strconv.Itoa(int(val)))
		t.Fail()
	}
}

func TestPbString(t *testing.T) {
	val := "Hello World"
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(string)
	if res != val {
		log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
}

func TestProtoType(t *testing.T) {
	val := &tests.TestProto{}
	val.MyString = "MyString"
	globals.Registry().Register(val)
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(*tests.TestProto)
	if res.MyString != val.MyString {
		log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
}

func TestSliceOfInt32(t *testing.T) {
	val := []int32{1, 2, 3, 4, 5}
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.([]int32)
	if len(val) != len(res) {
		log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for i := 0; i < len(val); i++ {
		if val[i] != res[i] {
			log.Error("int32 Slice Values do not match")
			t.Fail()
		}
	}
}

func TestSliceOfString(t *testing.T) {
	val := []string{"1", "2", "3", "4", "5"}
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.([]string)
	if len(val) != len(res) {
		log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for i := 0; i < len(val); i++ {
		if val[i] != res[i] {
			log.Error("string Slice Values do not match")
			t.Fail()
		}
	}
}

func TestSliceOfProto(t *testing.T) {
	proto1 := &tests.TestProto{}
	proto1.MyString = "UUID-1"

	proto2 := &tests.TestProto{}
	proto2.MyString = "UUID-2"

	globals.Registry().Register(proto1)

	val := []*tests.TestProto{proto1, proto2}
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.([]*tests.TestProto)
	if len(val) != len(res) {
		log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for i := 0; i < len(val); i++ {
		if val[i].MyString != res[i].MyString {
			log.Error("proto Slice Values do not match")
			t.Fail()
		}
	}
}

func TestNilSlice(t *testing.T) {
	var val []*tests.TestProto
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}

	if dval != nil {
		t.Fail()
		log.Error("Excpected nil slice")
	}
}

func TestNilProto(t *testing.T) {
	var val *tests.TestProto
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}

	if dval != nil {
		t.Fail()
		log.Error("Excpected nil proto")
	}
}

func TestSliceOfProtoWithNil(t *testing.T) {
	proto1 := &tests.TestProto{}
	proto1.MyString = "UUID-1"

	proto2 := &tests.TestProto{}
	proto2.MyString = "UUID-2"

	globals.Registry().Register(proto1)

	val := []*tests.TestProto{proto1, nil, proto2}
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.([]*tests.TestProto)
	if len(val) != len(res) {
		log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for i := 0; i < len(val); i++ {
		if val[i] == nil && res[i] != nil {
			log.Error("nil proto Slice Values do not match")
			t.Fail()
		} else if val[i] != nil && val[i].MyString != res[i].MyString {
			log.Error("proto Slice Values do not match")
			t.Fail()
		}
	}
}

func TestMapOfString2Int32(t *testing.T) {
	val := make(map[string]int32)
	val["1"] = 1
	val["2"] = 2
	val["3"] = 3
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(map[string]int32)
	if len(val) != len(res) {
		log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for k, v := range res {
		if val[k] != v {
			log.Error("map[string]int32 Values do not match")
			t.Fail()
		}
	}
}

func TestMapOfInt322String(t *testing.T) {
	val := make(map[int32]string)
	val[1] = "1"
	val[2] = "2"
	val[3] = "3"
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(map[int32]string)
	if len(val) != len(res) {
		log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for k, v := range res {
		if val[k] != v {
			log.Error("map[int32]string Values do not match")
			t.Fail()
		}
	}
}

func TestMapOfString2Proto(t *testing.T) {
	proto1 := &tests.TestProto{}
	proto1.MyString = "UUID-1"

	proto2 := &tests.TestProto{}
	proto2.MyString = "UUID-2"

	globals.Registry().Register(proto1)

	val := make(map[string]*tests.TestProto)

	val[proto1.MyString] = proto1
	val[proto2.MyString] = proto2

	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(map[string]*tests.TestProto)
	if len(val) != len(res) {
		log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for k, v := range res {
		if val[k].MyString != v.MyString {
			log.Error("map[string]proto Values do not match")
			t.Fail()
		}
	}
}

func TestMapOfString2ProtoWithNil(t *testing.T) {
	proto1 := &tests.TestProto{}
	proto1.MyString = "UUID-1"

	proto2 := &tests.TestProto{}
	proto2.MyString = "UUID-2"

	globals.Registry().Register(proto1)

	val := make(map[string]*tests.TestProto)

	val[proto1.MyString] = proto1
	val[proto2.MyString] = proto2
	val["Uuid3"] = nil

	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(map[string]*tests.TestProto)
	if len(val) != len(res) {
		log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for k, v := range val {
		if v == nil && res[k] != nil {
			log.Error("expected nil")
			t.Fail()
		} else if v != nil && res[k].MyString != v.MyString {
			log.Error("map[string]proto Values do not match")
			t.Fail()
		}
	}
}

func TestBool(t *testing.T) {
	val := true
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(bool)
	if res != val {
		log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}

	val = false
	dval, err = testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res = dval.(bool)
	if res != val {
		log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
}
