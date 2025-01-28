package common

import (
	"github.com/saichler/serializer/go/types"
	"github.com/saichler/shared/go/share/string_utils"
	"reflect"
	"strings"
)

func ValueAndType(any interface{}) (reflect.Value, reflect.Type) {
	v := reflect.ValueOf(any)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	return v, t
}

func IsLeaf(node *types.Node) bool {
	if node.Attributes == nil || len(node.Attributes) == 0 {
		return true
	}
	return false
}

func IsRoot(node *types.Node) bool {
	if node.Parent == nil {
		return true
	}
	return false
}

func IgnoreName(fieldName string) bool {
	if fieldName == "DoNotCompare" {
		return true
	}
	if fieldName == "DoNotCopy" {
		return true
	}
	if len(fieldName) > 3 && fieldName[0:3] == "XXX" {
		return true
	}
	if fieldName[0:1] == strings.ToLower(fieldName[0:1]) {
		return true
	}
	return false
}

func NodeKey(instanceId string) string {
	buff := string_utils.New()
	open := false
	for _, c := range instanceId {
		if c == '<' {
			open = true
		} else if c == '>' {
			open = false
		} else if !open {
			buff.Add(string(c))
		}
	}
	return buff.String()
}
