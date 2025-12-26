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
	"bytes"
	"encoding/base64"
	"math"
	"testing"

	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/testtypes"
)

// TestObj_PrimitiveTypes provides comprehensive testing of primitive type
// serialization including edge cases for min/max values, zero values,
// and special cases like unicode strings.
func TestObj_PrimitiveTypes(t *testing.T) {
	tests := []struct {
		name string
		val  interface{}
	}{
		{"int32", int32(42)},
		{"int32_negative", int32(-42)},
		{"int32_zero", int32(0)},
		{"int32_max", int32(math.MaxInt32)},
		{"int32_min", int32(math.MinInt32)},
		{"int64", int64(9223372036854775807)},
		{"int64_negative", int64(-9223372036854775807)},
		{"int64_zero", int64(0)},
		{"uint32", uint32(42)},
		{"uint32_max", uint32(math.MaxUint32)},
		{"uint64", uint64(18446744073709551615)},
		{"uint64_zero", uint64(0)},
		{"float32", float32(3.14159)},
		{"float32_negative", float32(-3.14159)},
		{"float32_zero", float32(0.0)},
		{"float64", float64(3.141592653589793)},
		{"float64_negative", float64(-3.141592653589793)},
		{"float64_zero", float64(0.0)},
		{"bool_true", true},
		{"bool_false", false},
		{"string_empty", ""},
		{"string_simple", "Hello, World!"},
		{"string_unicode", "Hello, 世界! 🌍"},
		{"string_long", string(make([]byte, 10000))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize
			obj := object.NewEncode()
			err := obj.Add(tt.val)
			if err != nil {
				t.Fatalf("Failed to serialize %v: %v", tt.val, err)
			}

			data := obj.Data()
			if len(data) == 0 {
				t.Fatal("Serialized data is empty")
			}

			// Deserialize
			decoded := object.NewDecode(data, 0, globals.Registry())
			result, err := decoded.Get()
			if err != nil {
				t.Fatalf("Failed to deserialize: %v", err)
			}

			// Verify
			if result != tt.val {
				t.Errorf("Value mismatch: expected %v (%T), got %v (%T)", tt.val, tt.val, result, result)
			}
		})
	}
}

// TestObj_Slices tests serialization of various slice types including
// int32, string, float64, bool, and byte slices with empty and populated data.
func TestObj_Slices(t *testing.T) {
	tests := []struct {
		name string
		val  interface{}
	}{
		{"slice_int32_empty", []int32{}},
		{"slice_int32", []int32{1, 2, 3, 4, 5}},
		{"slice_int32_negative", []int32{-1, -2, -3}},
		{"slice_string_empty", []string{}},
		{"slice_string", []string{"a", "b", "c"}},
		{"slice_string_with_empty", []string{"", "hello", ""}},
		{"slice_float64", []float64{1.1, 2.2, 3.3}},
		{"slice_bool", []bool{true, false, true}},
		{"slice_bytes", []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}},
		{"slice_bytes_empty", []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize
			obj := object.NewEncode()
			err := obj.Add(tt.val)
			if err != nil {
				t.Fatalf("Failed to serialize %v: %v", tt.val, err)
			}

			data := obj.Data()

			// Deserialize
			decoded := object.NewDecode(data, 0, globals.Registry())
			result, err := decoded.Get()
			if err != nil {
				t.Fatalf("Failed to deserialize: %v", err)
			}

			// Verify based on type
			switch expected := tt.val.(type) {
			case []int32:
				// Handle empty slices that deserialize to nil
				if len(expected) == 0 && result == nil {
					// This is expected behavior - empty slices serialize to nil
					return
				}
				actual, ok := result.([]int32)
				if !ok {
					t.Errorf("Type assertion failed: expected []int32, got %T", result)
					return
				}
				if !compareInt32Slices(expected, actual) {
					t.Errorf("Slice mismatch: expected %v, got %v", expected, actual)
				}
			case []string:
				// Handle empty slices that deserialize to nil
				if len(expected) == 0 && result == nil {
					return
				}
				actual, ok := result.([]string)
				if !ok {
					t.Errorf("Type assertion failed: expected []string, got %T", result)
					return
				}
				if !compareStringSlices(expected, actual) {
					t.Errorf("Slice mismatch: expected %v, got %v", expected, actual)
				}
			case []float64:
				// Handle empty slices that deserialize to nil
				if len(expected) == 0 && result == nil {
					return
				}
				actual, ok := result.([]float64)
				if !ok {
					t.Errorf("Type assertion failed: expected []float64, got %T", result)
					return
				}
				if !compareFloat64Slices(expected, actual) {
					t.Errorf("Slice mismatch: expected %v, got %v", expected, actual)
				}
			case []bool:
				// Handle empty slices that deserialize to nil
				if len(expected) == 0 && result == nil {
					return
				}
				actual, ok := result.([]bool)
				if !ok {
					t.Errorf("Type assertion failed: expected []bool, got %T", result)
					return
				}
				if !compareBoolSlices(expected, actual) {
					t.Errorf("Slice mismatch: expected %v, got %v", expected, actual)
				}
			case []byte:
				// Handle empty slices that deserialize to nil
				if len(expected) == 0 && result == nil {
					return
				}
				actual, ok := result.([]byte)
				if !ok {
					t.Errorf("Type assertion failed: expected []byte, got %T", result)
					return
				}
				if !bytes.Equal(expected, actual) {
					t.Errorf("Byte slice mismatch: expected %v, got %v", expected, actual)
				}
			}
		})
	}
}

// TestObj_Maps tests serialization of various map types including
// string-to-int32, int32-to-string, and string-to-float64/bool maps.
func TestObj_Maps(t *testing.T) {
	tests := []struct {
		name string
		val  interface{}
	}{
		{"map_string_int32_empty", map[string]int32{}},
		{"map_string_int32", map[string]int32{"one": 1, "two": 2, "three": 3}},
		{"map_int32_string", map[int32]string{1: "one", 2: "two", 3: "three"}},
		{"map_string_float64", map[string]float64{"pi": 3.14159, "e": 2.71828}},
		{"map_string_bool", map[string]bool{"true": true, "false": false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize
			obj := object.NewEncode()
			err := obj.Add(tt.val)
			if err != nil {
				t.Fatalf("Failed to serialize %v: %v", tt.val, err)
			}

			data := obj.Data()

			// Deserialize
			decoded := object.NewDecode(data, 0, globals.Registry())
			result, err := decoded.Get()
			if err != nil {
				t.Fatalf("Failed to deserialize: %v", err)
			}

			// Verify based on type
			switch expected := tt.val.(type) {
			case map[string]int32:
				// Handle empty maps that deserialize to nil
				if len(expected) == 0 && result == nil {
					return
				}
				actual, ok := result.(map[string]int32)
				if !ok {
					t.Errorf("Type assertion failed: expected map[string]int32, got %T", result)
					return
				}
				if !compareStringInt32Maps(expected, actual) {
					t.Errorf("Map mismatch: expected %v, got %v", expected, actual)
				}
			case map[int32]string:
				// Handle empty maps that deserialize to nil
				if len(expected) == 0 && result == nil {
					return
				}
				actual, ok := result.(map[int32]string)
				if !ok {
					t.Errorf("Type assertion failed: expected map[int32]string, got %T", result)
					return
				}
				if !compareInt32StringMaps(expected, actual) {
					t.Errorf("Map mismatch: expected %v, got %v", expected, actual)
				}
			case map[string]float64:
				// Handle empty maps that deserialize to nil
				if len(expected) == 0 && result == nil {
					return
				}
				actual, ok := result.(map[string]float64)
				if !ok {
					t.Errorf("Type assertion failed: expected map[string]float64, got %T", result)
					return
				}
				if !compareStringFloat64Maps(expected, actual) {
					t.Errorf("Map mismatch: expected %v, got %v", expected, actual)
				}
			case map[string]bool:
				// Handle empty maps that deserialize to nil
				if len(expected) == 0 && result == nil {
					return
				}
				actual, ok := result.(map[string]bool)
				if !ok {
					t.Errorf("Type assertion failed: expected map[string]bool, got %T", result)
					return
				}
				if !compareStringBoolMaps(expected, actual) {
					t.Errorf("Map mismatch: expected %v, got %v", expected, actual)
				}
			}
		})
	}
}

// TestObj_NilValues tests proper handling of nil slices, maps, and pointers.
func TestObj_NilValues(t *testing.T) {
	tests := []struct {
		name string
		val  interface{}
	}{
		{"nil_slice", []int32(nil)},
		{"nil_map", map[string]int32(nil)},
		{"nil_pointer", (*testtypes.TestProto)(nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize
			obj := object.NewEncode()
			err := obj.Add(tt.val)
			if err != nil {
				t.Fatalf("Failed to serialize nil value: %v", err)
			}

			data := obj.Data()

			// Deserialize
			decoded := object.NewDecode(data, 0, globals.Registry())
			result, err := decoded.Get()
			if err != nil {
				t.Fatalf("Failed to deserialize nil value: %v", err)
			}

			// Verify nil
			if result != nil {
				t.Errorf("Expected nil, got %v (%T)", result, result)
			}
		})
	}
}

// TestObj_ProtoStructs tests Protocol Buffers message serialization
// with populated fields.
func TestObj_ProtoStructs(t *testing.T) {
	// Create test proto
	proto := &testtypes.TestProto{
		MyString: "test-proto-123",
		MyInt32:  42,
		MyBool:   true,
	}
	
	// Register the type
	globals.Registry().Register(proto)

	// Serialize
	obj := object.NewEncode()
	err := obj.Add(proto)
	if err != nil {
		t.Fatalf("Failed to serialize proto: %v", err)
	}

	data := obj.Data()

	// Deserialize
	decoded := object.NewDecode(data, 0, globals.Registry())
	result, err := decoded.Get()
	if err != nil {
		t.Fatalf("Failed to deserialize proto: %v", err)
	}

	// Verify
	resultProto, ok := result.(*testtypes.TestProto)
	if !ok {
		t.Fatalf("Expected *testtypes.TestProto, got %T", result)
	}

	if resultProto.MyString != proto.MyString {
		t.Errorf("MyString mismatch: expected %s, got %s", proto.MyString, resultProto.MyString)
	}
	if resultProto.MyInt32 != proto.MyInt32 {
		t.Errorf("MyInt32 mismatch: expected %d, got %d", proto.MyInt32, resultProto.MyInt32)
	}
	if resultProto.MyBool != proto.MyBool {
		t.Errorf("MyBool mismatch: expected %t, got %t", proto.MyBool, resultProto.MyBool)
	}
}

// TestObj_EmptyProtoStruct tests serialization of empty Protocol Buffers messages.
func TestObj_EmptyProtoStruct(t *testing.T) {
	// Create empty test proto
	proto := &testtypes.TestProto{}
	
	// Register the type
	globals.Registry().Register(proto)

	// Serialize
	obj := object.NewEncode()
	err := obj.Add(proto)
	if err != nil {
		t.Fatalf("Failed to serialize empty proto: %v", err)
	}

	data := obj.Data()

	// Deserialize
	decoded := object.NewDecode(data, 0, globals.Registry())
	result, err := decoded.Get()
	if err != nil {
		t.Fatalf("Failed to deserialize empty proto: %v", err)
	}

	// Verify
	resultProto, ok := result.(*testtypes.TestProto)
	if !ok {
		t.Fatalf("Expected *testtypes.TestProto, got %T", result)
	}

	if resultProto.MyString != proto.MyString {
		t.Errorf("MyString mismatch: expected %s, got %s", proto.MyString, resultProto.MyString)
	}
}

// TestObj_Base64 tests Base64 encoding and decoding round-trip.
func TestObj_Base64(t *testing.T) {
	testValue := "Hello, Base64!"
	
	// Serialize
	obj := object.NewEncode()
	err := obj.Add(testValue)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	// Get Base64 string
	base64Str := obj.Base64()
	if base64Str == "" {
		t.Fatal("Base64 string is empty")
	}

	// Verify it's valid Base64
	_, err = base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		t.Fatalf("Invalid Base64 string: %v", err)
	}

	// Decode from Base64
	data, err := object.FromBase64(base64Str)
	if err != nil {
		t.Fatalf("Failed to decode Base64: %v", err)
	}

	// Deserialize
	decoded := object.NewDecode(data, 0, globals.Registry())
	result, err := decoded.Get()
	if err != nil {
		t.Fatalf("Failed to deserialize from Base64: %v", err)
	}

	// Verify
	if result != testValue {
		t.Errorf("Value mismatch: expected %s, got %s", testValue, result)
	}
}

// TestObj_UtilityFunctions tests the DataOf and ElemOf convenience functions.
func TestObj_UtilityFunctions(t *testing.T) {
	testValue := int32(12345)
	
	// Test DataOf
	data, err := object.DataOf(testValue)
	if err != nil {
		t.Fatalf("DataOf failed: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("DataOf returned empty data")
	}

	// Test ElemOf
	result, err := object.ElemOf(data, globals.Registry())
	if err != nil {
		t.Fatalf("ElemOf failed: %v", err)
	}

	// Verify
	if result != testValue {
		t.Errorf("Value mismatch: expected %v, got %v", testValue, result)
	}
}

// TestObj_UtilityFunctions_Nil tests DataOf and ElemOf with nil inputs.
func TestObj_UtilityFunctions_Nil(t *testing.T) {
	// Test DataOf with nil
	data, err := object.DataOf(nil)
	if err != nil {
		t.Fatalf("DataOf with nil failed: %v", err)
	}
	if data != nil {
		t.Error("Expected nil data for nil input")
	}

	// Test ElemOf with nil
	result, err := object.ElemOf(nil, globals.Registry())
	if err != nil {
		t.Fatalf("ElemOf with nil failed: %v", err)
	}
	if result != nil {
		t.Error("Expected nil result for nil input")
	}
}

// TestObj_LargeData tests serialization of large slices (100,000 elements)
// to verify buffer expansion and performance under load.
func TestObj_LargeData(t *testing.T) {
	// Create large slice
	largeSlice := make([]int32, 100000)
	for i := 0; i < len(largeSlice); i++ {
		largeSlice[i] = int32(i)
	}

	// Serialize
	obj := object.NewEncode()
	err := obj.Add(largeSlice)
	if err != nil {
		t.Fatalf("Failed to serialize large data: %v", err)
	}

	data := obj.Data()
	if len(data) == 0 {
		t.Fatal("Serialized data is empty")
	}

	// Deserialize
	decoded := object.NewDecode(data, 0, globals.Registry())
	result, err := decoded.Get()
	if err != nil {
		t.Fatalf("Failed to deserialize large data: %v", err)
	}

	// Verify
	resultSlice := result.([]int32)
	if len(resultSlice) != len(largeSlice) {
		t.Errorf("Length mismatch: expected %d, got %d", len(largeSlice), len(resultSlice))
	}

	// Spot check some values
	for i := 0; i < len(largeSlice); i += 1000 {
		if resultSlice[i] != largeSlice[i] {
			t.Errorf("Value mismatch at index %d: expected %d, got %d", i, largeSlice[i], resultSlice[i])
		}
	}
}

// TestObj_BufferExpansion tests the exponential buffer growth strategy
// by adding many items to force multiple buffer expansions.
func TestObj_BufferExpansion(t *testing.T) {
	obj := object.NewEncode()
	
	// Add many small items to force buffer expansion
	for i := 0; i < 1000; i++ {
		err := obj.Add(int32(i))
		if err != nil {
			t.Fatalf("Failed to add item %d: %v", i, err)
		}
	}

	data := obj.Data()
	if len(data) == 0 {
		t.Fatal("No data after multiple additions")
	}

	// Verify we can still serialize properly
	finalObj := object.NewEncode()
	err := finalObj.Add("test after expansion")
	if err != nil {
		t.Fatalf("Failed to serialize after buffer expansion: %v", err)
	}
}

// compareInt32Slices compares two int32 slices for equality.
func compareInt32Slices(a, b []int32) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// compareStringSlices compares two string slices for equality.
func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// compareFloat64Slices compares two float64 slices for equality.
func compareFloat64Slices(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// compareBoolSlices compares two bool slices for equality.
func compareBoolSlices(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// compareStringInt32Maps compares two string-to-int32 maps for equality.
func compareStringInt32Maps(a, b map[string]int32) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// compareInt32StringMaps compares two int32-to-string maps for equality.
func compareInt32StringMaps(a, b map[int32]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// compareStringFloat64Maps compares two string-to-float64 maps for equality.
func compareStringFloat64Maps(a, b map[string]float64) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// compareStringBoolMaps compares two string-to-bool maps for equality.
func compareStringBoolMaps(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}