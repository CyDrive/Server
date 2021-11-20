// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package rpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ManageClient is the client API for Manage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ManageClient interface {
	JoinCluster(ctx context.Context, in *JoinClusterRequest, opts ...grpc.CallOption) (*JoinClusterResponse, error)
	HeartBeats(ctx context.Context, in *HeartBeatsRequest, opts ...grpc.CallOption) (*HeartBeatsResponse, error)
	ReportFileInfos(ctx context.Context, in *ReportFileInfosRequest, opts ...grpc.CallOption) (*ReportFileInfosResponse, error)
	Notifier(ctx context.Context, in *ConnectNotifierRequest, opts ...grpc.CallOption) (Manage_NotifierClient, error)
}

type manageClient struct {
	cc grpc.ClientConnInterface
}

func NewManageClient(cc grpc.ClientConnInterface) ManageClient {
	return &manageClient{cc}
}

func (c *manageClient) JoinCluster(ctx context.Context, in *JoinClusterRequest, opts ...grpc.CallOption) (*JoinClusterResponse, error) {
	out := new(JoinClusterResponse)
	err := c.cc.Invoke(ctx, "/rpc.Manage/JoinCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *manageClient) HeartBeats(ctx context.Context, in *HeartBeatsRequest, opts ...grpc.CallOption) (*HeartBeatsResponse, error) {
	out := new(HeartBeatsResponse)
	err := c.cc.Invoke(ctx, "/rpc.Manage/HeartBeats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *manageClient) ReportFileInfos(ctx context.Context, in *ReportFileInfosRequest, opts ...grpc.CallOption) (*ReportFileInfosResponse, error) {
	out := new(ReportFileInfosResponse)
	err := c.cc.Invoke(ctx, "/rpc.Manage/ReportFileInfos", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *manageClient) Notifier(ctx context.Context, in *ConnectNotifierRequest, opts ...grpc.CallOption) (Manage_NotifierClient, error) {
	stream, err := c.cc.NewStream(ctx, &Manage_ServiceDesc.Streams[0], "/rpc.Manage/Notifier", opts...)
	if err != nil {
		return nil, err
	}
	x := &manageNotifierClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Manage_NotifierClient interface {
	Recv() (*Notify, error)
	grpc.ClientStream
}

type manageNotifierClient struct {
	grpc.ClientStream
}

func (x *manageNotifierClient) Recv() (*Notify, error) {
	m := new(Notify)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ManageServer is the server API for Manage service.
// All implementations must embed UnimplementedManageServer
// for forward compatibility
type ManageServer interface {
	JoinCluster(context.Context, *JoinClusterRequest) (*JoinClusterResponse, error)
	HeartBeats(context.Context, *HeartBeatsRequest) (*HeartBeatsResponse, error)
	ReportFileInfos(context.Context, *ReportFileInfosRequest) (*ReportFileInfosResponse, error)
	Notifier(*ConnectNotifierRequest, Manage_NotifierServer) error
	mustEmbedUnimplementedManageServer()
}

// UnimplementedManageServer must be embedded to have forward compatible implementations.
type UnimplementedManageServer struct {
}

func (UnimplementedManageServer) JoinCluster(context.Context, *JoinClusterRequest) (*JoinClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinCluster not implemented")
}
func (UnimplementedManageServer) HeartBeats(context.Context, *HeartBeatsRequest) (*HeartBeatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeartBeats not implemented")
}
func (UnimplementedManageServer) ReportFileInfos(context.Context, *ReportFileInfosRequest) (*ReportFileInfosResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportFileInfos not implemented")
}
func (UnimplementedManageServer) Notifier(*ConnectNotifierRequest, Manage_NotifierServer) error {
	return status.Errorf(codes.Unimplemented, "method Notifier not implemented")
}
func (UnimplementedManageServer) mustEmbedUnimplementedManageServer() {}

// UnsafeManageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ManageServer will
// result in compilation errors.
type UnsafeManageServer interface {
	mustEmbedUnimplementedManageServer()
}

func RegisterManageServer(s grpc.ServiceRegistrar, srv ManageServer) {
	s.RegisterService(&Manage_ServiceDesc, srv)
}

func _Manage_JoinCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManageServer).JoinCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Manage/JoinCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManageServer).JoinCluster(ctx, req.(*JoinClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Manage_HeartBeats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartBeatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManageServer).HeartBeats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Manage/HeartBeats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManageServer).HeartBeats(ctx, req.(*HeartBeatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Manage_ReportFileInfos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportFileInfosRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManageServer).ReportFileInfos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Manage/ReportFileInfos",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManageServer).ReportFileInfos(ctx, req.(*ReportFileInfosRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Manage_Notifier_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConnectNotifierRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ManageServer).Notifier(m, &manageNotifierServer{stream})
}

type Manage_NotifierServer interface {
	Send(*Notify) error
	grpc.ServerStream
}

type manageNotifierServer struct {
	grpc.ServerStream
}

func (x *manageNotifierServer) Send(m *Notify) error {
	return x.ServerStream.SendMsg(m)
}

// Manage_ServiceDesc is the grpc.ServiceDesc for Manage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Manage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Manage",
	HandlerType: (*ManageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "JoinCluster",
			Handler:    _Manage_JoinCluster_Handler,
		},
		{
			MethodName: "HeartBeats",
			Handler:    _Manage_HeartBeats_Handler,
		},
		{
			MethodName: "ReportFileInfos",
			Handler:    _Manage_ReportFileInfos_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Notifier",
			Handler:       _Manage_Notifier_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "rpc/manage.proto",
}
