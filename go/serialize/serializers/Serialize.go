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
package serializers

import "github.com/saichler/l8types/go/ifs"

// Default is the default serializer instance used throughout the Layer 8 ecosystem.
// It uses ProtoBuffBinary for efficient binary serialization.
//
// Usage:
//
//	data, err := serializers.Default.Marshal(myObject, resources)
//	result, err := serializers.Default.Unmarshal(data, resources)
var Default ifs.ISerializer = &ProtoBuffBinary{}

/*
import (
	"github.com/saichler/l8utils/go/utils/common"
	"google.golang.org/protobuf/proto"
	"reflect"
	"sync"
)

var empty = make([]byte, 0)
var pbMtx = &sync.Mutex{}
var UseProtoBuiltInSerializer = true

var Default ifs.Serializer

func (r *type_registry2.StructRegistryImpl) Marshal(any interface{}) ([]byte, error) {
	if any == nil {
		return empty, nil
	}

	pb, ok := any.(proto.Message)
	if ok && UseProtoBuiltInSerializer {
		pbMtx.Lock()
		defer pbMtx.Unlock()
		return proto.Marshal(pb)
	}

	val := reflect.ValueOf(any)
	if !val.IsValid() {
		return empty, nil
	}

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return empty, nil
		}
		val = val.Elem()
	}

	_, ser, _ := r.types.Get(val.Type().Name())
	if ser != nil {
		b, _ := ser.Add(any, r)
		return b, nil
	}
	return nil, ifs.Error("serializer not found in struct registry for ", val.Type().Name())
}

func (r *type_registry2.StructRegistryImpl) Unmarshal(name string, b []byte) (interface{}, error) {
	ins, ser, err := r.NewInstance(name)
	if err == nil && UseProtoBuiltInSerializer {
		pb := ins.(proto.Message)
		pbMtx.Lock()
		defer pbMtx.Unlock()
		err = proto.Unmarshal(b, pb)
		if err != nil {
			return nil, err
		}
		return pb, nil
	}

	if err == nil && ser != nil {
		dins, _ := ser.Get(b, 0, r)
		return dins, nil
	}
	return nil, ifs.Error("serializer not found in struct registry for ", name)
}*/
