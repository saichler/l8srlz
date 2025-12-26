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

// Package serializers provides implementations of the ifs.ISerializer interface
// for various serialization formats. Currently supports binary Protocol Buffers
// serialization using the L8S (Layer 8 Serialization) format.
package serializers

import (
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
)

// ProtoBuffBinary implements the ifs.ISerializer interface for binary
// serialization using the L8S format. It wraps the object package's
// serialization capabilities with a standard interface.
//
// This serializer is suitable for high-performance inter-service communication
// where compact binary representation is preferred over human-readable formats.
type ProtoBuffBinary struct{}

// Mode returns the serializer mode, which is BINARY for this implementation.
func (s *ProtoBuffBinary) Mode() ifs.SerializerMode {
	return ifs.BINARY
}

// Marshal serializes any Go value to binary format using the L8S encoder.
// The resources parameter provides access to the type registry for complex types.
//
// Returns the serialized byte slice and nil error on success.
func (s *ProtoBuffBinary) Marshal(any interface{}, resources ifs.IResources) ([]byte, error) {
	obj := object.NewEncode()
	obj.Add(any)
	return obj.Data(), nil
}

// Unmarshal deserializes binary data back to the original Go value.
// Uses the registry from resources to resolve type information for
// Protocol Buffers messages and other complex types.
//
// Returns the deserialized value and nil error on success.
func (s *ProtoBuffBinary) Unmarshal(data []byte, resources ifs.IResources) (interface{}, error) {
	obj := object.NewDecode(data, 0, resources.Registry())
	return obj.Get()
}
