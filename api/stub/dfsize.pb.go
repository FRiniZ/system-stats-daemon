// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: dfsize.proto

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

type Dfsize struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorMsg string `protobuf:"bytes,1,opt,name=ErrorMsg,proto3" json:"ErrorMsg,omitempty"`
	Name     string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Used     int32  `protobuf:"varint,3,opt,name=Used,proto3" json:"Used,omitempty"`
	Use      int32  `protobuf:"varint,4,opt,name=Use,proto3" json:"Use,omitempty"`
}

func (x *Dfsize) Reset() {
	*x = Dfsize{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dfsize_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Dfsize) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Dfsize) ProtoMessage() {}

func (x *Dfsize) ProtoReflect() protoreflect.Message {
	mi := &file_dfsize_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Dfsize.ProtoReflect.Descriptor instead.
func (*Dfsize) Descriptor() ([]byte, []int) {
	return file_dfsize_proto_rawDescGZIP(), []int{0}
}

func (x *Dfsize) GetErrorMsg() string {
	if x != nil {
		return x.ErrorMsg
	}
	return ""
}

func (x *Dfsize) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Dfsize) GetUsed() int32 {
	if x != nil {
		return x.Used
	}
	return 0
}

func (x *Dfsize) GetUse() int32 {
	if x != nil {
		return x.Use
	}
	return 0
}

var File_dfsize_proto protoreflect.FileDescriptor

var file_dfsize_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x64, 0x66, 0x73, 0x69, 0x7a, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03,
	0x61, 0x70, 0x69, 0x22, 0x5e, 0x0a, 0x06, 0x64, 0x66, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x55, 0x73, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x55, 0x73, 0x65,
	0x64, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x73, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x55, 0x73, 0x65, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x73, 0x74, 0x75, 0x62, 0x2f, 0x3b, 0x61,
	0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_dfsize_proto_rawDescOnce sync.Once
	file_dfsize_proto_rawDescData = file_dfsize_proto_rawDesc
)

func file_dfsize_proto_rawDescGZIP() []byte {
	file_dfsize_proto_rawDescOnce.Do(func() {
		file_dfsize_proto_rawDescData = protoimpl.X.CompressGZIP(file_dfsize_proto_rawDescData)
	})
	return file_dfsize_proto_rawDescData
}

var file_dfsize_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_dfsize_proto_goTypes = []interface{}{
	(*Dfsize)(nil), // 0: api.dfsize
}
var file_dfsize_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_dfsize_proto_init() }
func file_dfsize_proto_init() {
	if File_dfsize_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dfsize_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Dfsize); i {
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
			RawDescriptor: file_dfsize_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_dfsize_proto_goTypes,
		DependencyIndexes: file_dfsize_proto_depIdxs,
		MessageInfos:      file_dfsize_proto_msgTypes,
	}.Build()
	File_dfsize_proto = out.File
	file_dfsize_proto_rawDesc = nil
	file_dfsize_proto_goTypes = nil
	file_dfsize_proto_depIdxs = nil
}
