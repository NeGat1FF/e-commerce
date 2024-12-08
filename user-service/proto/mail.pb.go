// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: proto/mail.proto

package proto

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

type NotificationType int32

const (
	NotificationType_EMAIL_CONFIRMATION NotificationType = 0
	NotificationType_PASSWORD_RESET     NotificationType = 1
)

// Enum value maps for NotificationType.
var (
	NotificationType_name = map[int32]string{
		0: "EMAIL_CONFIRMATION",
		1: "PASSWORD_RESET",
	}
	NotificationType_value = map[string]int32{
		"EMAIL_CONFIRMATION": 0,
		"PASSWORD_RESET":     1,
	}
)

func (x NotificationType) Enum() *NotificationType {
	p := new(NotificationType)
	*p = x
	return p
}

func (x NotificationType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NotificationType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_mail_proto_enumTypes[0].Descriptor()
}

func (NotificationType) Type() protoreflect.EnumType {
	return &file_proto_mail_proto_enumTypes[0]
}

func (x NotificationType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NotificationType.Descriptor instead.
func (NotificationType) EnumDescriptor() ([]byte, []int) {
	return file_proto_mail_proto_rawDescGZIP(), []int{0}
}

type MailRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	To   []string          `protobuf:"bytes,1,rep,name=to,proto3" json:"to,omitempty"`
	Type NotificationType  `protobuf:"varint,2,opt,name=type,proto3,enum=proto.NotificationType" json:"type,omitempty"`
	Data map[string]string `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *MailRequest) Reset() {
	*x = MailRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mail_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MailRequest) ProtoMessage() {}

func (x *MailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mail_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MailRequest.ProtoReflect.Descriptor instead.
func (*MailRequest) Descriptor() ([]byte, []int) {
	return file_proto_mail_proto_rawDescGZIP(), []int{0}
}

func (x *MailRequest) GetTo() []string {
	if x != nil {
		return x.To
	}
	return nil
}

func (x *MailRequest) GetType() NotificationType {
	if x != nil {
		return x.Type
	}
	return NotificationType_EMAIL_CONFIRMATION
}

func (x *MailRequest) GetData() map[string]string {
	if x != nil {
		return x.Data
	}
	return nil
}

type MailResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *MailResponse) Reset() {
	*x = MailResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_mail_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MailResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MailResponse) ProtoMessage() {}

func (x *MailResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_mail_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MailResponse.ProtoReflect.Descriptor instead.
func (*MailResponse) Descriptor() ([]byte, []int) {
	return file_proto_mail_proto_rawDescGZIP(), []int{1}
}

func (x *MailResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_proto_mail_proto protoreflect.FileDescriptor

var file_proto_mail_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb5, 0x01, 0x0a, 0x0b, 0x4d, 0x61,
	0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x6f, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x2b, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x30, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x61, 0x69,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x37, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0x28, 0x0a, 0x0c, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2a, 0x3e, 0x0a, 0x10, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x16, 0x0a, 0x12, 0x45, 0x4d, 0x41, 0x49, 0x4c, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x52, 0x4d,
	0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x50, 0x41, 0x53, 0x53, 0x57,
	0x4f, 0x52, 0x44, 0x5f, 0x52, 0x45, 0x53, 0x45, 0x54, 0x10, 0x01, 0x32, 0x44, 0x0a, 0x0b, 0x4d,
	0x61, 0x69, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x53, 0x65,
	0x6e, 0x64, 0x4d, 0x61, 0x69, 0x6c, 0x12, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d,
	0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4d, 0x61, 0x69, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x28, 0x5a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x4e, 0x65, 0x47, 0x61, 0x74, 0x31, 0x46, 0x46, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2d, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_proto_mail_proto_rawDescOnce sync.Once
	file_proto_mail_proto_rawDescData = file_proto_mail_proto_rawDesc
)

func file_proto_mail_proto_rawDescGZIP() []byte {
	file_proto_mail_proto_rawDescOnce.Do(func() {
		file_proto_mail_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_mail_proto_rawDescData)
	})
	return file_proto_mail_proto_rawDescData
}

var file_proto_mail_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_mail_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_mail_proto_goTypes = []interface{}{
	(NotificationType)(0), // 0: proto.NotificationType
	(*MailRequest)(nil),   // 1: proto.MailRequest
	(*MailResponse)(nil),  // 2: proto.MailResponse
	nil,                   // 3: proto.MailRequest.DataEntry
}
var file_proto_mail_proto_depIdxs = []int32{
	0, // 0: proto.MailRequest.type:type_name -> proto.NotificationType
	3, // 1: proto.MailRequest.data:type_name -> proto.MailRequest.DataEntry
	1, // 2: proto.MailService.SendMail:input_type -> proto.MailRequest
	2, // 3: proto.MailService.SendMail:output_type -> proto.MailResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_mail_proto_init() }
func file_proto_mail_proto_init() {
	if File_proto_mail_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_mail_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MailRequest); i {
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
		file_proto_mail_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MailResponse); i {
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
			RawDescriptor: file_proto_mail_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_mail_proto_goTypes,
		DependencyIndexes: file_proto_mail_proto_depIdxs,
		EnumInfos:         file_proto_mail_proto_enumTypes,
		MessageInfos:      file_proto_mail_proto_msgTypes,
	}.Build()
	File_proto_mail_proto = out.File
	file_proto_mail_proto_rawDesc = nil
	file_proto_mail_proto_goTypes = nil
	file_proto_mail_proto_depIdxs = nil
}