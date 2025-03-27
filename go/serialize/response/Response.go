package response

import (
	"errors"
	"github.com/saichler/serializer/go/serialize/object"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
	"google.golang.org/protobuf/proto"
)

type Response struct {
	list []interface{}
	err  error
}

func (this *Response) AddProto(pb proto.Message) {
	if this.list == nil {
		this.list = make([]interface{}, 0)
	}
	this.list = append(this.list, pb)
}

func (this *Response) SetError(err error) {
	this.err = err
}

func (this *Response) List() []interface{} {
	return this.list
}

func (this *Response) Error() error {
	return this.err
}

func (this *Response) ToProto() *types.Response {
	response := &types.Response{}
	if this.list != nil {
		response.List = make([][]byte, len(this.list))
		for i, pb := range this.list {
			obj := object.NewEncode([]byte{}, 0)
			obj.Add(pb)
			response.List[i] = obj.Data()
		}
	}
	if this.err != nil {
		response.ErrMsg = this.err.Error()
	}
	return response
}

func FromProto(response *types.Response, resourcs common.IResources) *Response {
	resp := &Response{}
	if response.List != nil {
		resp.list = make([]interface{}, len(response.List))
		for i, data := range response.List {
			obj := object.NewDecode(data, 0, "", resourcs.Registry())
			pb, _ := obj.Get()
			resp.list[i] = pb
		}
	}
	if response.ErrMsg != "" {
		resp.err = errors.New(response.ErrMsg)
	}
	return resp
}

func NewMl(pbs []interface{}) *Response {
	resp := &Response{}
	if pbs != nil {
		resp.list = make([]interface{}, len(pbs))
		for i, pb := range pbs {
			resp.list[i] = pb
		}
	}
	return resp
}

func NewSl(pb interface{}) *Response {
	resp := &Response{}
	if pb != nil {
		resp.list = make([]interface{}, 1)
		resp.list[0] = pb
	}
	return resp
}

func NewErr(errMessage string) *Response {
	return &Response{err: errors.New(errMessage)}
}
