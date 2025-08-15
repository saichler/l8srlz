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
	"reflect"
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

func TestNewError(t *testing.T) {
	elems := object.NewError("test error")
	if elems.Error().Error() != "test error" {
		Log.Fail(t, "Expected 'test error', got:", elems.Error().Error())
	}
}

func TestNewNotify(t *testing.T) {
	testData := "test notification"
	elems := object.NewNotify(testData)
	if !elems.Notification() {
		Log.Fail(t, "Expected notification to be true")
	}
	if elems.Element() != testData {
		Log.Fail(t, "Expected element to match test data")
	}
}

func TestNewReplicasRequest(t *testing.T) {
	original := object.New(nil, "test")
	replica := object.NewReplicasRequest(original)
	if !replica.ReplicasRequest() {
		Log.Fail(t, "Expected replicasRequest to be true")
	}
}

func TestNewWithSlice(t *testing.T) {
	slice := []string{"item1", "item2", "item3"}
	elems := object.New(nil, slice)
	
	elements := elems.Elements()
	if len(elements) != 3 {
		Log.Fail(t, "Expected 3 elements, got:", len(elements))
	}
	
	keys := elems.Keys()
	for i, key := range keys {
		if key != i {
			Log.Fail(t, "Expected key to be index", i, "got:", key)
		}
	}
}

func TestNewWithMap(t *testing.T) {
	testMap := map[string]int{"key1": 1, "key2": 2, "key3": 3}
	elems := object.New(nil, testMap)
	
	elements := elems.Elements()
	if len(elements) != 3 {
		Log.Fail(t, "Expected 3 elements, got:", len(elements))
	}
	
	keys := elems.Keys()
	if len(keys) != 3 {
		Log.Fail(t, "Expected 3 keys, got:", len(keys))
	}
}

func TestNewWithFunction(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			Log.Fail(t, "Expected panic when passing function")
		}
	}()
	object.New(nil, func() {})
}

func TestAdd(t *testing.T) {
	elems := &object.Elements{}
	elems.Add("element1", "key1", nil)
	elems.Add("element2", "key2", errors.New("test error"))
	
	if len(elems.Elements()) != 2 {
		Log.Fail(t, "Expected 2 elements")
	}
	
	if elems.Elements()[0] != "element1" {
		Log.Fail(t, "First element mismatch")
	}
	
	if elems.Keys()[0] != "key1" {
		Log.Fail(t, "First key mismatch")
	}
	
	errors := elems.Errors()
	if errors[0] != nil {
		Log.Fail(t, "First error should be nil")
	}
	if errors[1].Error() != "test error" {
		Log.Fail(t, "Second error mismatch")
	}
}

func TestElementAndKey(t *testing.T) {
	elems := object.New(nil, "single element")
	if elems.Element() != "single element" {
		Log.Fail(t, "Element mismatch")
	}
	
	elems.Add("test", "test key", nil)
	if elems.Key() != nil {
		Log.Fail(t, "Expected first element key to be nil")
	}
}

func TestErrorHandling(t *testing.T) {
	elems := object.New(errors.New("test error"), "data")
	if elems.Error().Error() != "test error" {
		Log.Fail(t, "Error mismatch")
	}
	
	emptyElems := &object.Elements{}
	if emptyElems.Error() != nil {
		Log.Fail(t, "Empty elements should have nil error")
	}
}

func TestSerializeDeserialize(t *testing.T) {
	res, _ := CreateResources(25000, 2, ifs.Info_Level)
	res.Registry().Register(&testtypes.TestProto{})
	
	original := object.New(errors.New("serialize error"), "test data")
	original.Add("second element", "second key", nil)
	
	data, err := original.Serialize()
	if err != nil {
		Log.Fail(t, "Serialize failed:", err)
		return
	}
	
	deserialized := &object.Elements{}
	err = deserialized.Deserialize(data, res.Registry())
	if err != nil {
		Log.Fail(t, "Deserialize failed:", err)
		return
	}
	
	if len(deserialized.Elements()) != len(original.Elements()) {
		Log.Fail(t, "Element count mismatch after deserialization")
	}
	
	if deserialized.Error().Error() != original.Error().Error() {
		Log.Fail(t, "Error mismatch after deserialization")
	}
}

func TestAsListError(t *testing.T) {
	res, _ := CreateResources(25000, 2, ifs.Info_Level)
	
	emptyElems := &object.Elements{}
	_, err := emptyElems.AsList(res.Registry())
	if err == nil {
		Log.Fail(t, "Expected error for empty elements")
	}
	if err.Error() != "elements is empty" {
		Log.Fail(t, "Expected 'elements is empty' error, got:", err.Error())
	}
}

func TestAsListWithoutRegistration(t *testing.T) {
	res, _ := CreateResources(25000, 2, ifs.Info_Level)
	res.Registry().Register(&testtypes.TestProto{})
	
	elem := CreateTestModelInstance(2)
	elems := object.New(nil, elem)
	
	result, err := elems.AsList(res.Registry())
	if err != nil {
		Log.Fail(t, "AsList failed:", err)
		return
	}
	
	if !reflect.DeepEqual(result, elem) {
		Log.Fail(t, "Result should be original element when list type not registered")
	}
}

func TestAppend(t *testing.T) {
	elem1 := object.New(nil, "first")
	elem2 := object.New(nil, "second")
	
	elem1.Append(elem2)
	
	if len(elem1.Elements()) != 2 {
		Log.Fail(t, "Expected 2 elements after append, got:", len(elem1.Elements()))
	}
}

func TestAppendBugFix(t *testing.T) {
	elem1 := object.New(nil, "first")
	elem2 := object.New(nil, "second")
	elem2.Add("third", "key3", nil)
	
	initialCount := len(elem1.Elements())
	elem1.Append(elem2)
	finalCount := len(elem1.Elements())
	
	expectedCount := initialCount + len(elem2.Elements())
	if finalCount != expectedCount {
		Log.Fail(t, "Append method has a bug. Expected", expectedCount, "elements, got:", finalCount)
	}
}
