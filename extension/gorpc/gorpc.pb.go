// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.3
// source: gorpc.proto

package gorpc

import (
	reflect "reflect"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

var file_gorpc_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptor.MethodOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         50001,
		Name:          "gorpc.alias",
		Tag:           "bytes,50001,opt,name=alias",
		Filename:      "gorpc.proto",
	},
}

// Extension fields to descriptor.MethodOptions.
var (
	// optional string alias = 50001;
	E_Alias = &file_gorpc_proto_extTypes[0]
)

var File_gorpc_proto protoreflect.FileDescriptor

var file_gorpc_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x67, 0x6f, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x67,
	0x6f, 0x72, 0x70, 0x63, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x36, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x12,
	0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0xd1, 0x86, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_gorpc_proto_goTypes = []interface{}{
	(*descriptor.MethodOptions)(nil), // 0: google.protobuf.MethodOptions
}
var file_gorpc_proto_depIdxs = []int32{
	0, // 0: gorpc.alias:extendee -> google.protobuf.MethodOptions
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	0, // [0:1] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gorpc_proto_init() }
func file_gorpc_proto_init() {
	if File_gorpc_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gorpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_gorpc_proto_goTypes,
		DependencyIndexes: file_gorpc_proto_depIdxs,
		ExtensionInfos:    file_gorpc_proto_extTypes,
	}.Build()
	File_gorpc_proto = out.File
	file_gorpc_proto_rawDesc = nil
	file_gorpc_proto_goTypes = nil
	file_gorpc_proto_depIdxs = nil
}
