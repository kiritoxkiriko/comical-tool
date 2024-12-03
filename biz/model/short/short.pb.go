// idl/short.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.12
// source: short.proto

package short

import (
	_ "github.com/kiritoxkiriko/comical-tool/biz/model/api"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ShortReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code string `protobuf:"bytes,1,opt,name=Code,proto3" json:"Code,omitempty" path:"code"`
}

func (x *ShortReq) Reset() {
	*x = ShortReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_short_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShortReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortReq) ProtoMessage() {}

func (x *ShortReq) ProtoReflect() protoreflect.Message {
	mi := &file_short_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortReq.ProtoReflect.Descriptor instead.
func (*ShortReq) Descriptor() ([]byte, []int) {
	return file_short_proto_rawDescGZIP(), []int{0}
}

func (x *ShortReq) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type ShortResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RespBody string `protobuf:"bytes,1,opt,name=RespBody,proto3" json:"RespBody,omitempty" raw_body:"resp_body"`
}

func (x *ShortResp) Reset() {
	*x = ShortResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_short_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShortResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortResp) ProtoMessage() {}

func (x *ShortResp) ProtoReflect() protoreflect.Message {
	mi := &file_short_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortResp.ProtoReflect.Descriptor instead.
func (*ShortResp) Descriptor() ([]byte, []int) {
	return file_short_proto_rawDescGZIP(), []int{1}
}

func (x *ShortResp) GetRespBody() string {
	if x != nil {
		return x.RespBody
	}
	return ""
}

var File_short_proto protoreflect.FileDescriptor

var file_short_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x73,
	0x68, 0x6f, 0x72, 0x74, 0x1a, 0x09, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x28, 0x0a, 0x08, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x12, 0x1c, 0x0a, 0x04, 0x43,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xd2, 0xbb, 0x18, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x52, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x36, 0x0a, 0x09, 0x53, 0x68, 0x6f,
	0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x12, 0x29, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x42, 0x6f,
	0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0d, 0xaa, 0xbb, 0x18, 0x09, 0x72, 0x65,
	0x73, 0x70, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x52, 0x08, 0x52, 0x65, 0x73, 0x70, 0x42, 0x6f, 0x64,
	0x79, 0x32, 0x42, 0x0a, 0x0c, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x32, 0x0a, 0x05, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x12, 0x0f, 0x2e, 0x73, 0x68, 0x6f,
	0x72, 0x74, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x10, 0x2e, 0x73, 0x68,
	0x6f, 0x72, 0x74, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x22, 0x06, 0xca,
	0xc1, 0x18, 0x02, 0x2f, 0x73, 0x42, 0x37, 0x5a, 0x35, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x69, 0x72, 0x69, 0x74, 0x6f, 0x78, 0x6b, 0x69, 0x72, 0x69, 0x6b,
	0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x69, 0x63, 0x61, 0x6c, 0x2d, 0x74, 0x6f, 0x6f, 0x6c, 0x2f, 0x62,
	0x69, 0x7a, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_short_proto_rawDescOnce sync.Once
	file_short_proto_rawDescData = file_short_proto_rawDesc
)

func file_short_proto_rawDescGZIP() []byte {
	file_short_proto_rawDescOnce.Do(func() {
		file_short_proto_rawDescData = protoimpl.X.CompressGZIP(file_short_proto_rawDescData)
	})
	return file_short_proto_rawDescData
}

var file_short_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_short_proto_goTypes = []interface{}{
	(*ShortReq)(nil),  // 0: short.ShortReq
	(*ShortResp)(nil), // 1: short.ShortResp
}
var file_short_proto_depIdxs = []int32{
	0, // 0: short.ShortService.Short:input_type -> short.ShortReq
	1, // 1: short.ShortService.Short:output_type -> short.ShortResp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_short_proto_init() }
func file_short_proto_init() {
	if File_short_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_short_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShortReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_short_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShortResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_short_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_short_proto_goTypes,
		DependencyIndexes: file_short_proto_depIdxs,
		MessageInfos:      file_short_proto_msgTypes,
	}.Build()
	File_short_proto = out.File
	file_short_proto_rawDesc = nil
	file_short_proto_goTypes = nil
	file_short_proto_depIdxs = nil
}
