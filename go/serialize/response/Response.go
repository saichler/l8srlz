package response

import (
	"errors"
	"github.com/saichler/serializer/go/serialize/object"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
)

type Response struct {
	elems []interface{}
	err   error
}

func NewError(errMessage string) *Response {
	return &Response{err: errors.New(errMessage)}
}

func New(err error, elems ...interface{}) *Response {
	return &Response{err: err, elems: elems}
}

func (this *Response) Add(elem interface{}) {
	if this.elems == nil {
		this.elems = make([]interface{}, 0)
	}
	this.elems = append(this.elems, elem)
}

func (this *Response) SetError(err error) {
	this.err = err
}

func (this *Response) Elems() []interface{} {
	return this.elems
}

func (this *Response) Elem() interface{} {
	if this == nil {
		return nil
	}
	if this.elems == nil || len(this.elems) == 0 {
		return nil
	}
	return this.elems[0]
}

func (this *Response) Err() error {
	if this == nil {
		return nil
	}
	return this.err
}

func (this *Response) ToProto() *types.Response {
	response := &types.Response{}
	if this.elems != nil {
		response.Elems = make([][]byte, len(this.elems))
		for i, elem := range this.elems {
			obj := object.NewEncode([]byte{}, 0)
			obj.Add(elem)
			response.Elems[i] = obj.Data()
		}
	}
	if this.err != nil {
		response.ErrMsg = this.err.Error()
	}
	return response
}

func FromProto(response *types.Response, resourcs common.IResources) *Response {
	resp := &Response{}
	if response.Elems != nil {
		resp.elems = make([]interface{}, len(response.Elems))
		for i, data := range response.Elems {
			obj := object.NewDecode(data, 0, "", resourcs.Registry())
			pb, _ := obj.Get()
			resp.elems[i] = pb
		}
	}
	if response.ErrMsg != "" {
		resp.err = errors.New(response.ErrMsg)
	}
	return resp
}
