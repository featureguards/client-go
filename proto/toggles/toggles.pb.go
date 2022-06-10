// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: toggles.proto

package toggles

import (
	feature_toggle "github.com/featureguards/featureguards-go/v2/proto/feature_toggle"
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

type FetchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Platform feature_toggle.Platform_Type `protobuf:"varint,1,opt,name=platform,proto3,enum=feature_toggle.Platform_Type" json:"platform,omitempty"`
	Version  int64                        `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *FetchRequest) Reset() {
	*x = FetchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_toggles_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FetchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchRequest) ProtoMessage() {}

func (x *FetchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_toggles_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchRequest.ProtoReflect.Descriptor instead.
func (*FetchRequest) Descriptor() ([]byte, []int) {
	return file_toggles_proto_rawDescGZIP(), []int{0}
}

func (x *FetchRequest) GetPlatform() feature_toggle.Platform_Type {
	if x != nil {
		return x.Platform
	}
	return feature_toggle.Platform_Type(0)
}

func (x *FetchRequest) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

type FetchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FeatureToggles []*feature_toggle.FeatureToggle `protobuf:"bytes,1,rep,name=feature_toggles,json=featureToggles,proto3" json:"feature_toggles,omitempty"`
	Version        int64                           `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *FetchResponse) Reset() {
	*x = FetchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_toggles_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FetchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchResponse) ProtoMessage() {}

func (x *FetchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_toggles_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchResponse.ProtoReflect.Descriptor instead.
func (*FetchResponse) Descriptor() ([]byte, []int) {
	return file_toggles_proto_rawDescGZIP(), []int{1}
}

func (x *FetchResponse) GetFeatureToggles() []*feature_toggle.FeatureToggle {
	if x != nil {
		return x.FeatureToggles
	}
	return nil
}

func (x *FetchResponse) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

type ListenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Platform feature_toggle.Platform_Type `protobuf:"varint,1,opt,name=platform,proto3,enum=feature_toggle.Platform_Type" json:"platform,omitempty"`
	Version  int64                        `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *ListenRequest) Reset() {
	*x = ListenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_toggles_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListenRequest) ProtoMessage() {}

func (x *ListenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_toggles_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListenRequest.ProtoReflect.Descriptor instead.
func (*ListenRequest) Descriptor() ([]byte, []int) {
	return file_toggles_proto_rawDescGZIP(), []int{2}
}

func (x *ListenRequest) GetPlatform() feature_toggle.Platform_Type {
	if x != nil {
		return x.Platform
	}
	return feature_toggle.Platform_Type(0)
}

func (x *ListenRequest) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

type ListenPayload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FeatureToggles []*feature_toggle.FeatureToggle `protobuf:"bytes,1,rep,name=feature_toggles,json=featureToggles,proto3" json:"feature_toggles,omitempty"`
	Version        int64                           `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *ListenPayload) Reset() {
	*x = ListenPayload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_toggles_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListenPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListenPayload) ProtoMessage() {}

func (x *ListenPayload) ProtoReflect() protoreflect.Message {
	mi := &file_toggles_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListenPayload.ProtoReflect.Descriptor instead.
func (*ListenPayload) Descriptor() ([]byte, []int) {
	return file_toggles_proto_rawDescGZIP(), []int{3}
}

func (x *ListenPayload) GetFeatureToggles() []*feature_toggle.FeatureToggle {
	if x != nil {
		return x.FeatureToggles
	}
	return nil
}

func (x *ListenPayload) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

var File_toggles_proto protoreflect.FileDescriptor

var file_toggles_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x73, 0x1a, 0x1b, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64,
	0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x63, 0x0a, 0x0c, 0x46, 0x65, 0x74, 0x63, 0x68, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x39, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x5f, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x2e, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x71, 0x0a, 0x0d, 0x46, 0x65,
	0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x0f, 0x66,
	0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x74,
	0x6f, 0x67, 0x67, 0x6c, 0x65, 0x2e, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x54, 0x6f, 0x67,
	0x67, 0x6c, 0x65, 0x52, 0x0e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x54, 0x6f, 0x67, 0x67,
	0x6c, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x64, 0x0a,
	0x0d, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x39,
	0x0a, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x1d, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x74, 0x6f, 0x67, 0x67, 0x6c,
	0x65, 0x2e, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x22, 0x71, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x50, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x12, 0x46, 0x0a, 0x0f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f,
	0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e,
	0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x2e, 0x46,
	0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x52, 0x0e, 0x66, 0x65,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x32, 0x81, 0x01, 0x0a, 0x07, 0x54, 0x6f, 0x67, 0x67, 0x6c,
	0x65, 0x73, 0x12, 0x38, 0x0a, 0x05, 0x46, 0x65, 0x74, 0x63, 0x68, 0x12, 0x15, 0x2e, 0x74, 0x6f,
	0x67, 0x67, 0x6c, 0x65, 0x73, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x73, 0x2e, 0x46, 0x65, 0x74,
	0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3c, 0x0a, 0x06,
	0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x12, 0x16, 0x2e, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x73,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x50,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x00, 0x30, 0x01, 0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x67, 0x75, 0x61, 0x72, 0x64, 0x73, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x67, 0x75,
	0x61, 0x72, 0x64, 0x73, 0x2d, 0x67, 0x6f, 0x2f, 0x76, 0x32, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_toggles_proto_rawDescOnce sync.Once
	file_toggles_proto_rawDescData = file_toggles_proto_rawDesc
)

func file_toggles_proto_rawDescGZIP() []byte {
	file_toggles_proto_rawDescOnce.Do(func() {
		file_toggles_proto_rawDescData = protoimpl.X.CompressGZIP(file_toggles_proto_rawDescData)
	})
	return file_toggles_proto_rawDescData
}

var file_toggles_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_toggles_proto_goTypes = []interface{}{
	(*FetchRequest)(nil),                 // 0: toggles.FetchRequest
	(*FetchResponse)(nil),                // 1: toggles.FetchResponse
	(*ListenRequest)(nil),                // 2: toggles.ListenRequest
	(*ListenPayload)(nil),                // 3: toggles.ListenPayload
	(feature_toggle.Platform_Type)(0),    // 4: feature_toggle.Platform.Type
	(*feature_toggle.FeatureToggle)(nil), // 5: feature_toggle.FeatureToggle
}
var file_toggles_proto_depIdxs = []int32{
	4, // 0: toggles.FetchRequest.platform:type_name -> feature_toggle.Platform.Type
	5, // 1: toggles.FetchResponse.feature_toggles:type_name -> feature_toggle.FeatureToggle
	4, // 2: toggles.ListenRequest.platform:type_name -> feature_toggle.Platform.Type
	5, // 3: toggles.ListenPayload.feature_toggles:type_name -> feature_toggle.FeatureToggle
	0, // 4: toggles.Toggles.Fetch:input_type -> toggles.FetchRequest
	2, // 5: toggles.Toggles.Listen:input_type -> toggles.ListenRequest
	1, // 6: toggles.Toggles.Fetch:output_type -> toggles.FetchResponse
	3, // 7: toggles.Toggles.Listen:output_type -> toggles.ListenPayload
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_toggles_proto_init() }
func file_toggles_proto_init() {
	if File_toggles_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_toggles_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FetchRequest); i {
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
		file_toggles_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FetchResponse); i {
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
		file_toggles_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListenRequest); i {
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
		file_toggles_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListenPayload); i {
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
			RawDescriptor: file_toggles_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_toggles_proto_goTypes,
		DependencyIndexes: file_toggles_proto_depIdxs,
		MessageInfos:      file_toggles_proto_msgTypes,
	}.Build()
	File_toggles_proto = out.File
	file_toggles_proto_rawDesc = nil
	file_toggles_proto_goTypes = nil
	file_toggles_proto_depIdxs = nil
}
