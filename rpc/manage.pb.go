// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: rpc/manage.proto

package rpc

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

type StorageNodeType int32

const (
	StorageNodeType_Public  StorageNodeType = 0
	StorageNodeType_Private StorageNodeType = 1
)

// Enum value maps for StorageNodeType.
var (
	StorageNodeType_name = map[int32]string{
		0: "Public",
		1: "Private",
	}
	StorageNodeType_value = map[string]int32{
		"Public":  0,
		"Private": 1,
	}
)

func (x StorageNodeType) Enum() *StorageNodeType {
	p := new(StorageNodeType)
	*p = x
	return p
}

func (x StorageNodeType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StorageNodeType) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc_manage_proto_enumTypes[0].Descriptor()
}

func (StorageNodeType) Type() protoreflect.EnumType {
	return &file_rpc_manage_proto_enumTypes[0]
}

func (x StorageNodeType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StorageNodeType.Descriptor instead.
func (StorageNodeType) EnumDescriptor() ([]byte, []int) {
	return file_rpc_manage_proto_rawDescGZIP(), []int{0}
}

type JoinClusterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Capacity int64           `protobuf:"varint,1,opt,name=capacity,proto3" json:"capacity,omitempty"`
	Usage    int64           `protobuf:"varint,2,opt,name=usage,proto3" json:"usage,omitempty"`
	Type     StorageNodeType `protobuf:"varint,3,opt,name=type,proto3,enum=rpc.StorageNodeType" json:"type,omitempty"`
}

func (x *JoinClusterRequest) Reset() {
	*x = JoinClusterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_manage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinClusterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinClusterRequest) ProtoMessage() {}

func (x *JoinClusterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_manage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinClusterRequest.ProtoReflect.Descriptor instead.
func (*JoinClusterRequest) Descriptor() ([]byte, []int) {
	return file_rpc_manage_proto_rawDescGZIP(), []int{0}
}

func (x *JoinClusterRequest) GetCapacity() int64 {
	if x != nil {
		return x.Capacity
	}
	return 0
}

func (x *JoinClusterRequest) GetUsage() int64 {
	if x != nil {
		return x.Usage
	}
	return 0
}

func (x *JoinClusterRequest) GetType() StorageNodeType {
	if x != nil {
		return x.Type
	}
	return StorageNodeType_Public
}

type JoinClusterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *JoinClusterResponse) Reset() {
	*x = JoinClusterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_manage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinClusterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinClusterResponse) ProtoMessage() {}

func (x *JoinClusterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_manage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinClusterResponse.ProtoReflect.Descriptor instead.
func (*JoinClusterResponse) Descriptor() ([]byte, []int) {
	return file_rpc_manage_proto_rawDescGZIP(), []int{1}
}

func (x *JoinClusterResponse) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type HeartBeatsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	StorageUsage    int64 `protobuf:"varint,2,opt,name=storage_usage,json=storageUsage,proto3" json:"storage_usage,omitempty"`
	CpuUsagePercent int32 `protobuf:"varint,3,opt,name=cpu_usage_percent,json=cpuUsagePercent,proto3" json:"cpu_usage_percent,omitempty"`
	TaskNum         int32 `protobuf:"varint,4,opt,name=task_num,json=taskNum,proto3" json:"task_num,omitempty"`
}

func (x *HeartBeatsRequest) Reset() {
	*x = HeartBeatsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_manage_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartBeatsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartBeatsRequest) ProtoMessage() {}

func (x *HeartBeatsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_manage_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartBeatsRequest.ProtoReflect.Descriptor instead.
func (*HeartBeatsRequest) Descriptor() ([]byte, []int) {
	return file_rpc_manage_proto_rawDescGZIP(), []int{2}
}

func (x *HeartBeatsRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *HeartBeatsRequest) GetStorageUsage() int64 {
	if x != nil {
		return x.StorageUsage
	}
	return 0
}

func (x *HeartBeatsRequest) GetCpuUsagePercent() int32 {
	if x != nil {
		return x.CpuUsagePercent
	}
	return 0
}

func (x *HeartBeatsRequest) GetTaskNum() int32 {
	if x != nil {
		return x.TaskNum
	}
	return 0
}

type HeartBeatsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HeartBeatsResponse) Reset() {
	*x = HeartBeatsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_manage_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartBeatsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartBeatsResponse) ProtoMessage() {}

func (x *HeartBeatsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_manage_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartBeatsResponse.ProtoReflect.Descriptor instead.
func (*HeartBeatsResponse) Descriptor() ([]byte, []int) {
	return file_rpc_manage_proto_rawDescGZIP(), []int{3}
}

var File_rpc_manage_proto protoreflect.FileDescriptor

var file_rpc_manage_proto_rawDesc = []byte{
	0x0a, 0x10, 0x72, 0x70, 0x63, 0x2f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x03, 0x72, 0x70, 0x63, 0x22, 0x70, 0x0a, 0x12, 0x4a, 0x6f, 0x69, 0x6e, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a,
	0x08, 0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x08, 0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x75, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x28, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e,
	0x72, 0x70, 0x63, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x25, 0x0a, 0x13, 0x4a, 0x6f, 0x69,
	0x6e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64,
	0x22, 0x8f, 0x01, 0x0a, 0x11, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x73,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2a, 0x0a, 0x11, 0x63,
	0x70, 0x75, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0f, 0x63, 0x70, 0x75, 0x55, 0x73, 0x61, 0x67, 0x65,
	0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x61, 0x73, 0x6b, 0x5f,
	0x6e, 0x75, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x4e,
	0x75, 0x6d, 0x22, 0x14, 0x0a, 0x12, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2a, 0x2a, 0x0a, 0x0f, 0x53, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x50,
	0x75, 0x62, 0x6c, 0x69, 0x63, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x72, 0x69, 0x76, 0x61,
	0x74, 0x65, 0x10, 0x01, 0x32, 0x91, 0x01, 0x0a, 0x06, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x12,
	0x42, 0x0a, 0x0b, 0x4a, 0x6f, 0x69, 0x6e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x17,
	0x2e, 0x72, 0x70, 0x63, 0x2e, 0x4a, 0x6f, 0x69, 0x6e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x4a, 0x6f,
	0x69, 0x6e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x43, 0x0a, 0x0a, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74,
	0x73, 0x12, 0x16, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61,
	0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x72, 0x70, 0x63, 0x2e,
	0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x18, 0x5a, 0x16, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x43, 0x79, 0x44, 0x72, 0x69, 0x76, 0x65, 0x2f, 0x72,
	0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_manage_proto_rawDescOnce sync.Once
	file_rpc_manage_proto_rawDescData = file_rpc_manage_proto_rawDesc
)

func file_rpc_manage_proto_rawDescGZIP() []byte {
	file_rpc_manage_proto_rawDescOnce.Do(func() {
		file_rpc_manage_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_manage_proto_rawDescData)
	})
	return file_rpc_manage_proto_rawDescData
}

var file_rpc_manage_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_rpc_manage_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_rpc_manage_proto_goTypes = []interface{}{
	(StorageNodeType)(0),        // 0: rpc.StorageNodeType
	(*JoinClusterRequest)(nil),  // 1: rpc.JoinClusterRequest
	(*JoinClusterResponse)(nil), // 2: rpc.JoinClusterResponse
	(*HeartBeatsRequest)(nil),   // 3: rpc.HeartBeatsRequest
	(*HeartBeatsResponse)(nil),  // 4: rpc.HeartBeatsResponse
}
var file_rpc_manage_proto_depIdxs = []int32{
	0, // 0: rpc.JoinClusterRequest.type:type_name -> rpc.StorageNodeType
	1, // 1: rpc.Manage.JoinCluster:input_type -> rpc.JoinClusterRequest
	3, // 2: rpc.Manage.HeartBeats:input_type -> rpc.HeartBeatsRequest
	2, // 3: rpc.Manage.JoinCluster:output_type -> rpc.JoinClusterResponse
	4, // 4: rpc.Manage.HeartBeats:output_type -> rpc.HeartBeatsResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rpc_manage_proto_init() }
func file_rpc_manage_proto_init() {
	if File_rpc_manage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_manage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinClusterRequest); i {
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
		file_rpc_manage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinClusterResponse); i {
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
		file_rpc_manage_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeartBeatsRequest); i {
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
		file_rpc_manage_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeartBeatsResponse); i {
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
			RawDescriptor: file_rpc_manage_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_manage_proto_goTypes,
		DependencyIndexes: file_rpc_manage_proto_depIdxs,
		EnumInfos:         file_rpc_manage_proto_enumTypes,
		MessageInfos:      file_rpc_manage_proto_msgTypes,
	}.Build()
	File_rpc_manage_proto = out.File
	file_rpc_manage_proto_rawDesc = nil
	file_rpc_manage_proto_goTypes = nil
	file_rpc_manage_proto_depIdxs = nil
}
