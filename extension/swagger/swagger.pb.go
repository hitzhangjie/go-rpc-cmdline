// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.3
// source: swagger.proto

package swagger

import (
	reflect "reflect"
	sync "sync"

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

type SwaggerRule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title       string `protobuf:"bytes,50103,opt,name=title,proto3" json:"title,omitempty"`
	Method      string `protobuf:"bytes,50104,opt,name=method,proto3" json:"method,omitempty"`
	Description string `protobuf:"bytes,50105,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *SwaggerRule) Reset() {
	*x = SwaggerRule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swagger_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SwaggerRule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SwaggerRule) ProtoMessage() {}

func (x *SwaggerRule) ProtoReflect() protoreflect.Message {
	mi := &file_swagger_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SwaggerRule.ProtoReflect.Descriptor instead.
func (*SwaggerRule) Descriptor() ([]byte, []int) {
	return file_swagger_proto_rawDescGZIP(), []int{0}
}

func (x *SwaggerRule) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *SwaggerRule) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *SwaggerRule) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

var file_swagger_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptor.MethodOptions)(nil),
		ExtensionType: (*SwaggerRule)(nil),
		Field:         50101,
		Name:          "swagger.swagger",
		Tag:           "bytes,50101,opt,name=swagger",
		Filename:      "swagger.proto",
	},
}

// Extension fields to descriptor.MethodOptions.
var (
	// optional swagger.SwaggerRule swagger = 50101;
	E_Swagger = &file_swagger_proto_extTypes[0]
)

var File_swagger_proto protoreflect.FileDescriptor

var file_swagger_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x77, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x73, 0x77, 0x61, 0x67, 0x67, 0x65, 0x72, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x63, 0x0a, 0x0b, 0x53, 0x77,
	0x61, 0x67, 0x67, 0x65, 0x72, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x18, 0xb7, 0x87, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x12, 0x18, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0xb8, 0x87, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x22, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0xb9, 0x87, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x3a,
	0x50, 0x0a, 0x07, 0x73, 0x77, 0x61, 0x67, 0x67, 0x65, 0x72, 0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xb5, 0x87, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x77, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x53, 0x77, 0x61,
	0x67, 0x67, 0x65, 0x72, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x07, 0x73, 0x77, 0x61, 0x67, 0x67, 0x65,
	0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_swagger_proto_rawDescOnce sync.Once
	file_swagger_proto_rawDescData = file_swagger_proto_rawDesc
)

func file_swagger_proto_rawDescGZIP() []byte {
	file_swagger_proto_rawDescOnce.Do(func() {
		file_swagger_proto_rawDescData = protoimpl.X.CompressGZIP(file_swagger_proto_rawDescData)
	})
	return file_swagger_proto_rawDescData
}

var file_swagger_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_swagger_proto_goTypes = []interface{}{
	(*SwaggerRule)(nil),              // 0: swagger.SwaggerRule
	(*descriptor.MethodOptions)(nil), // 1: google.protobuf.MethodOptions
}
var file_swagger_proto_depIdxs = []int32{
	1, // 0: swagger.swagger:extendee -> google.protobuf.MethodOptions
	0, // 1: swagger.swagger:type_name -> swagger.SwaggerRule
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	1, // [1:2] is the sub-list for extension type_name
	0, // [0:1] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_swagger_proto_init() }
func file_swagger_proto_init() {
	if File_swagger_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_swagger_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SwaggerRule); i {
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
			RawDescriptor: file_swagger_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_swagger_proto_goTypes,
		DependencyIndexes: file_swagger_proto_depIdxs,
		MessageInfos:      file_swagger_proto_msgTypes,
		ExtensionInfos:    file_swagger_proto_extTypes,
	}.Build()
	File_swagger_proto = out.File
	file_swagger_proto_rawDesc = nil
	file_swagger_proto_goTypes = nil
	file_swagger_proto_depIdxs = nil
}
