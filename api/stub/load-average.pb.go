// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: load-average.proto

package api

import (
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

type Loadaverage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorMsg string  `protobuf:"bytes,1,opt,name=ErrorMsg,proto3" json:"ErrorMsg,omitempty"`
	L1       float32 `protobuf:"fixed32,2,opt,name=L1,proto3" json:"L1,omitempty"`
	L2       float32 `protobuf:"fixed32,3,opt,name=L2,proto3" json:"L2,omitempty"`
	L3       float32 `protobuf:"fixed32,4,opt,name=L3,proto3" json:"L3,omitempty"`
}

func (x *Loadaverage) Reset() {
	*x = Loadaverage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_load_average_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Loadaverage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Loadaverage) ProtoMessage() {}

func (x *Loadaverage) ProtoReflect() protoreflect.Message {
	mi := &file_load_average_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Loadaverage.ProtoReflect.Descriptor instead.
func (*Loadaverage) Descriptor() ([]byte, []int) {
	return file_load_average_proto_rawDescGZIP(), []int{0}
}

func (x *Loadaverage) GetErrorMsg() string {
	if x != nil {
		return x.ErrorMsg
	}
	return ""
}

func (x *Loadaverage) GetL1() float32 {
	if x != nil {
		return x.L1
	}
	return 0
}

func (x *Loadaverage) GetL2() float32 {
	if x != nil {
		return x.L2
	}
	return 0
}

func (x *Loadaverage) GetL3() float32 {
	if x != nil {
		return x.L3
	}
	return 0
}

var File_load_average_proto protoreflect.FileDescriptor

var file_load_average_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6c, 0x6f, 0x61, 0x64, 0x2d, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x22, 0x59, 0x0a, 0x0b, 0x6c, 0x6f, 0x61,
	0x64, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x4d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x4d, 0x73, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x4c, 0x31, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02,
	0x52, 0x02, 0x4c, 0x31, 0x12, 0x0e, 0x0a, 0x02, 0x4c, 0x32, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02,
	0x52, 0x02, 0x4c, 0x32, 0x12, 0x0e, 0x0a, 0x02, 0x4c, 0x33, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02,
	0x52, 0x02, 0x4c, 0x33, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x73, 0x74, 0x75, 0x62, 0x2f, 0x3b,
	0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_load_average_proto_rawDescOnce sync.Once
	file_load_average_proto_rawDescData = file_load_average_proto_rawDesc
)

func file_load_average_proto_rawDescGZIP() []byte {
	file_load_average_proto_rawDescOnce.Do(func() {
		file_load_average_proto_rawDescData = protoimpl.X.CompressGZIP(file_load_average_proto_rawDescData)
	})
	return file_load_average_proto_rawDescData
}

var file_load_average_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_load_average_proto_goTypes = []interface{}{
	(*Loadaverage)(nil), // 0: api.loadaverage
}
var file_load_average_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_load_average_proto_init() }
func file_load_average_proto_init() {
	if File_load_average_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_load_average_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Loadaverage); i {
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
			RawDescriptor: file_load_average_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_load_average_proto_goTypes,
		DependencyIndexes: file_load_average_proto_depIdxs,
		MessageInfos:      file_load_average_proto_msgTypes,
	}.Build()
	File_load_average_proto = out.File
	file_load_average_proto_rawDesc = nil
	file_load_average_proto_goTypes = nil
	file_load_average_proto_depIdxs = nil
}
