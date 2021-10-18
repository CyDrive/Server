// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: rpc/file_stream.proto

package rpc

import (
	models "github.com/CyDrive/models"
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

type SendFileChunkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskId   int64            `protobuf:"varint,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`      // first
	FileInfo *models.FileInfo `protobuf:"bytes,2,opt,name=file_info,json=fileInfo,proto3" json:"file_info,omitempty"` // none of first
	Error    string           `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`                       // last
	Data     []byte           `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`                         // all
}

func (x *SendFileChunkRequest) Reset() {
	*x = SendFileChunkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_file_stream_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendFileChunkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendFileChunkRequest) ProtoMessage() {}

func (x *SendFileChunkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_file_stream_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendFileChunkRequest.ProtoReflect.Descriptor instead.
func (*SendFileChunkRequest) Descriptor() ([]byte, []int) {
	return file_rpc_file_stream_proto_rawDescGZIP(), []int{0}
}

func (x *SendFileChunkRequest) GetTaskId() int64 {
	if x != nil {
		return x.TaskId
	}
	return 0
}

func (x *SendFileChunkRequest) GetFileInfo() *models.FileInfo {
	if x != nil {
		return x.FileInfo
	}
	return nil
}

func (x *SendFileChunkRequest) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *SendFileChunkRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type SendFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendFileResponse) Reset() {
	*x = SendFileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_file_stream_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendFileResponse) ProtoMessage() {}

func (x *SendFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_file_stream_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendFileResponse.ProtoReflect.Descriptor instead.
func (*SendFileResponse) Descriptor() ([]byte, []int) {
	return file_rpc_file_stream_proto_rawDescGZIP(), []int{1}
}

type RecvFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskId int64  `protobuf:"varint,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Error  string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *RecvFileRequest) Reset() {
	*x = RecvFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_file_stream_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecvFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecvFileRequest) ProtoMessage() {}

func (x *RecvFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_file_stream_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecvFileRequest.ProtoReflect.Descriptor instead.
func (*RecvFileRequest) Descriptor() ([]byte, []int) {
	return file_rpc_file_stream_proto_rawDescGZIP(), []int{2}
}

func (x *RecvFileRequest) GetTaskId() int64 {
	if x != nil {
		return x.TaskId
	}
	return 0
}

func (x *RecvFileRequest) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type RecvFileChunkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskId int64  `protobuf:"varint,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Data   []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *RecvFileChunkResponse) Reset() {
	*x = RecvFileChunkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_file_stream_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecvFileChunkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecvFileChunkResponse) ProtoMessage() {}

func (x *RecvFileChunkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_file_stream_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecvFileChunkResponse.ProtoReflect.Descriptor instead.
func (*RecvFileChunkResponse) Descriptor() ([]byte, []int) {
	return file_rpc_file_stream_proto_rawDescGZIP(), []int{3}
}

func (x *RecvFileChunkResponse) GetTaskId() int64 {
	if x != nil {
		return x.TaskId
	}
	return 0
}

func (x *RecvFileChunkResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_rpc_file_stream_proto protoreflect.FileDescriptor

var file_rpc_file_stream_proto_rawDesc = []byte{
	0x0a, 0x15, 0x72, 0x70, 0x63, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x72, 0x70, 0x63, 0x1a, 0x16, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x88, 0x01, 0x0a, 0x14, 0x53, 0x65, 0x6e, 0x64, 0x46, 0x69, 0x6c,
	0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a,
	0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x12, 0x2d, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69,
	0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x66, 0x69, 0x6c,
	0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x12, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x40, 0x0a, 0x0f, 0x52, 0x65, 0x63, 0x76, 0x46, 0x69, 0x6c, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x44, 0x0a, 0x15, 0x52, 0x65, 0x63, 0x76, 0x46, 0x69, 0x6c,
	0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17,
	0x0a, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0x90, 0x01, 0x0a, 0x0a,
	0x46, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x40, 0x0a, 0x08, 0x53, 0x65,
	0x6e, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x19, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x65, 0x6e,
	0x64, 0x46, 0x69, 0x6c, 0x65, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x15, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x46, 0x69, 0x6c, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x12, 0x40, 0x0a, 0x08,
	0x52, 0x65, 0x63, 0x76, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x14, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x52,
	0x65, 0x63, 0x76, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a,
	0x2e, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x63, 0x76, 0x46, 0x69, 0x6c, 0x65, 0x43, 0x68, 0x75,
	0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0x18,
	0x5a, 0x16, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x43, 0x79, 0x44,
	0x72, 0x69, 0x76, 0x65, 0x2f, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_file_stream_proto_rawDescOnce sync.Once
	file_rpc_file_stream_proto_rawDescData = file_rpc_file_stream_proto_rawDesc
)

func file_rpc_file_stream_proto_rawDescGZIP() []byte {
	file_rpc_file_stream_proto_rawDescOnce.Do(func() {
		file_rpc_file_stream_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_file_stream_proto_rawDescData)
	})
	return file_rpc_file_stream_proto_rawDescData
}

var file_rpc_file_stream_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_rpc_file_stream_proto_goTypes = []interface{}{
	(*SendFileChunkRequest)(nil),  // 0: rpc.SendFileChunkRequest
	(*SendFileResponse)(nil),      // 1: rpc.SendFileResponse
	(*RecvFileRequest)(nil),       // 2: rpc.RecvFileRequest
	(*RecvFileChunkResponse)(nil), // 3: rpc.RecvFileChunkResponse
	(*models.FileInfo)(nil),       // 4: models.FileInfo
}
var file_rpc_file_stream_proto_depIdxs = []int32{
	4, // 0: rpc.SendFileChunkRequest.file_info:type_name -> models.FileInfo
	0, // 1: rpc.FileStream.SendFile:input_type -> rpc.SendFileChunkRequest
	2, // 2: rpc.FileStream.RecvFile:input_type -> rpc.RecvFileRequest
	1, // 3: rpc.FileStream.SendFile:output_type -> rpc.SendFileResponse
	3, // 4: rpc.FileStream.RecvFile:output_type -> rpc.RecvFileChunkResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rpc_file_stream_proto_init() }
func file_rpc_file_stream_proto_init() {
	if File_rpc_file_stream_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_file_stream_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendFileChunkRequest); i {
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
		file_rpc_file_stream_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendFileResponse); i {
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
		file_rpc_file_stream_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecvFileRequest); i {
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
		file_rpc_file_stream_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecvFileChunkResponse); i {
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
			RawDescriptor: file_rpc_file_stream_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_file_stream_proto_goTypes,
		DependencyIndexes: file_rpc_file_stream_proto_depIdxs,
		MessageInfos:      file_rpc_file_stream_proto_msgTypes,
	}.Build()
	File_rpc_file_stream_proto = out.File
	file_rpc_file_stream_proto_rawDesc = nil
	file_rpc_file_stream_proto_goTypes = nil
	file_rpc_file_stream_proto_depIdxs = nil
}