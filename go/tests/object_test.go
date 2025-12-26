/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package tests

import (
	. "github.com/saichler/l8test/go/infra/t_resources"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/testtypes"
	"strconv"
	"testing"
)

// testType is a helper function that serializes a value, then deserializes it
// and returns the result for comparison. This provides a convenient round-trip
// test pattern used by the individual type tests.
func testType(val interface{}) (interface{}, error) {
	e := object.NewEncode()
	err := e.Add(val)
	if err != nil {
		return nil, err
	}
	data := e.Data()
	e = object.NewDecode(data, 0, globals.Registry())
	dval, err := e.Get()
	if err != nil {
		return nil, err
	}
	return dval, nil
}

// TestInt tests serialization of platform-dependent int type.
func TestInt(t *testing.T) {
	val := int(-39)
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(int)
	if res != val {
		Log.Error("Value do not match:", strconv.Itoa(int(res)), " and ", strconv.Itoa(int(val)))
		t.Fail()
	}
}

// TestInt32 tests serialization of 32-bit signed integer.
func TestInt32(t *testing.T) {
	val := int32(-39)
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(int32)
	if res != val {
		Log.Error("Value do not match:", strconv.Itoa(int(res)), " and ", strconv.Itoa(int(val)))
		t.Fail()
	}
}

// TestInt64 tests serialization of 64-bit signed integer.
func TestInt64(t *testing.T) {
	val := int64(-39)
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(int64)
	if res != val {
		Log.Error("Value do not match:", strconv.Itoa(int(res)), " and ", strconv.Itoa(int(val)))
		t.Fail()
	}
}

// TestUInt32 tests serialization of 32-bit unsigned integer.
func TestUInt32(t *testing.T) {
	val := uint32(39)
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(uint32)
	if res != val {
		Log.Error("Value do not match:", strconv.Itoa(int(res)), " and ", strconv.Itoa(int(val)))
		t.Fail()
	}
}

// TestUInt64 tests serialization of 64-bit unsigned integer.
func TestUInt64(t *testing.T) {
	val := uint64(39)
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(uint64)
	if res != val {
		Log.Error("Value do not match:", strconv.Itoa(int(res)), " and ", strconv.Itoa(int(val)))
		t.Fail()
	}
}

// TestFloat64 tests serialization of 64-bit floating point number.
func TestFloat64(t *testing.T) {
	val := float64(39.39)
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(float64)
	if res != val {
		Log.Error("Value do not match:", strconv.Itoa(int(res)), " and ", strconv.Itoa(int(val)))
		t.Fail()
	}
}

// TestFloat32 tests serialization of 32-bit floating point number.
func TestFloat32(t *testing.T) {
	val := float32(39.39)
	dval, err := testType(val)
	if err != nil {
		Log.Fail(t, err)
		return
	}
	res := dval.(float32)
	if res != val {
		Log.Error("Value do not match:", strconv.Itoa(int(res)), " and ", strconv.Itoa(int(val)))
		t.Fail()
	}
}

// TestPbString tests serialization of string values.
func TestPbString(t *testing.T) {
	val := "Hello World"
	dval, err := testType(val)
	if err != nil {
		Log.Fail(t, err)
		t.Fail()
		return
	}
	res := dval.(string)
	if res != val {
		Log.Fail(t, "Value do not match:", res, " and ", val)
	}
}

// TestProtoType tests serialization of Protocol Buffers messages with data.
func TestProtoType(t *testing.T) {
	val := CreateTestModelInstance(1)
	globals.Registry().Register(val)
	dval, err := testType(val)
	if err != nil {
		Log.Fail(t, err)
		return
	}
	res := dval.(*testtypes.TestProto)
	if res.MyString != val.MyString {
		Log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
}

// TestEmptyProtoType tests serialization of empty Protocol Buffers messages.
func TestEmptyProtoType(t *testing.T) {
	val := &testtypes.TestProto{}
	globals.Registry().Register(val)
	dval, err := testType(val)
	if err != nil {
		Log.Fail(t, err)
		return
	}
	res := dval.(*testtypes.TestProto)
	if res.MyString != val.MyString {
		Log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
}

// TestSliceOfInt32 tests serialization of int32 slices.
func TestSliceOfInt32(t *testing.T) {
	val := []int32{1, 2, 3, 4, 5}
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.([]int32)
	if len(val) != len(res) {
		Log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for i := 0; i < len(val); i++ {
		if val[i] != res[i] {
			Log.Error("int32 Slice Values do not match")
			t.Fail()
		}
	}
}

// TestSliceOfString tests serialization of string slices.
func TestSliceOfString(t *testing.T) {
	val := []string{"1", "2", "3", "4", "5"}
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.([]string)
	if len(val) != len(res) {
		Log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for i := 0; i < len(val); i++ {
		if val[i] != res[i] {
			Log.Error("string Slice Values do not match")
			t.Fail()
		}
	}
}

// TestSliceOfProto tests serialization of Protocol Buffers message slices.
func TestSliceOfProto(t *testing.T) {
	proto1 := &testtypes.TestProto{}
	proto1.MyString = "UUID-1"

	proto2 := &testtypes.TestProto{}
	proto2.MyString = "UUID-2"

	globals.Registry().Register(proto1)

	val := []*testtypes.TestProto{proto1, proto2}
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.([]*testtypes.TestProto)
	if len(val) != len(res) {
		Log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for i := 0; i < len(val); i++ {
		if val[i].MyString != res[i].MyString {
			Log.Error("proto Slice Values do not match")
			t.Fail()
		}
	}
}

// TestNilSlice tests that nil slices serialize and deserialize correctly.
func TestNilSlice(t *testing.T) {
	var val []*testtypes.TestProto
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}

	if dval != nil {
		t.Fail()
		Log.Error("Excpected nil slice")
	}
}

// TestNilProto tests that nil Protocol Buffers pointers serialize correctly.
func TestNilProto(t *testing.T) {
	var val *testtypes.TestProto
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}

	if dval != nil {
		t.Fail()
		Log.Error("Excpected nil proto")
	}
}

// TestSliceOfProtoWithNil tests slices containing nil Protocol Buffers elements.
func TestSliceOfProtoWithNil(t *testing.T) {
	proto1 := &testtypes.TestProto{}
	proto1.MyString = "UUID-1"

	proto2 := &testtypes.TestProto{}
	proto2.MyString = "UUID-2"

	globals.Registry().Register(proto1)

	val := []*testtypes.TestProto{proto1, nil, proto2}
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.([]*testtypes.TestProto)
	if len(val) != len(res) {
		Log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for i := 0; i < len(val); i++ {
		if val[i] == nil && res[i] != nil {
			Log.Error("nil proto Slice Values do not match")
			t.Fail()
		} else if val[i] != nil && val[i].MyString != res[i].MyString {
			Log.Error("proto Slice Values do not match")
			t.Fail()
		}
	}
}

// TestMapOfString2Int32 tests serialization of string-to-int32 maps.
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
		Log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for k, v := range res {
		if val[k] != v {
			Log.Error("map[string]int32 Values do not match")
			t.Fail()
		}
	}
}

// TestMapOfInt322String tests serialization of int32-to-string maps.
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
		Log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for k, v := range res {
		if val[k] != v {
			Log.Error("map[int32]string Values do not match")
			t.Fail()
		}
	}
}

// TestMapOfString2Proto tests serialization of string-to-Protocol Buffers maps.
func TestMapOfString2Proto(t *testing.T) {
	proto1 := &testtypes.TestProto{}
	proto1.MyString = "UUID-1"

	proto2 := &testtypes.TestProto{}
	proto2.MyString = "UUID-2"

	globals.Registry().Register(proto1)

	val := make(map[string]*testtypes.TestProto)

	val[proto1.MyString] = proto1
	val[proto2.MyString] = proto2

	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(map[string]*testtypes.TestProto)
	if len(val) != len(res) {
		Log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for k, v := range res {
		if val[k].MyString != v.MyString {
			Log.Error("map[string]proto Values do not match")
			t.Fail()
		}
	}
}

// TestMapOfString2ProtoWithNil tests maps containing nil Protocol Buffers values.
func TestMapOfString2ProtoWithNil(t *testing.T) {
	proto1 := &testtypes.TestProto{}
	proto1.MyString = "UUID-1"

	proto2 := &testtypes.TestProto{}
	proto2.MyString = "UUID-2"

	globals.Registry().Register(proto1)

	val := make(map[string]*testtypes.TestProto)

	val[proto1.MyString] = proto1
	val[proto2.MyString] = proto2
	val["Uuid3"] = nil

	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(map[string]*testtypes.TestProto)
	if len(val) != len(res) {
		Log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
	for k, v := range val {
		if v == nil && res[k] != nil {
			Log.Error("expected nil")
			t.Fail()
		} else if v != nil && res[k].MyString != v.MyString {
			Log.Error("map[string]proto Values do not match")
			t.Fail()
		}
	}
}

// TestBool tests serialization of boolean values (true and false).
func TestBool(t *testing.T) {
	val := true
	dval, err := testType(val)
	if err != nil {
		t.Fail()
		return
	}
	res := dval.(bool)
	if res != val {
		Log.Error("Value do not match:", res, " and ", val)
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
		Log.Error("Value do not match:", res, " and ", val)
		t.Fail()
	}
}
