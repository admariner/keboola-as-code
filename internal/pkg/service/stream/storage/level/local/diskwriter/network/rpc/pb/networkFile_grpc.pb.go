// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.2
// source: networkFile.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	NetworkFile_Open_FullMethodName            = "/pb.NetworkFile/Open"
	NetworkFile_KeepAliveStream_FullMethodName = "/pb.NetworkFile/KeepAliveStream"
	NetworkFile_Write_FullMethodName           = "/pb.NetworkFile/Write"
	NetworkFile_Sync_FullMethodName            = "/pb.NetworkFile/Sync"
	NetworkFile_Close_FullMethodName           = "/pb.NetworkFile/Close"
)

// NetworkFileClient is the client API for NetworkFile service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NetworkFileClient interface {
	Open(ctx context.Context, in *OpenRequest, opts ...grpc.CallOption) (*OpenResponse, error)
	KeepAliveStream(ctx context.Context, in *KeepAliveStreamRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[KeepAliveStreamResponse], error)
	Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error)
	Sync(ctx context.Context, in *SyncRequest, opts ...grpc.CallOption) (*SyncResponse, error)
	Close(ctx context.Context, in *CloseRequest, opts ...grpc.CallOption) (*CloseResponse, error)
}

type networkFileClient struct {
	cc grpc.ClientConnInterface
}

func NewNetworkFileClient(cc grpc.ClientConnInterface) NetworkFileClient {
	return &networkFileClient{cc}
}

func (c *networkFileClient) Open(ctx context.Context, in *OpenRequest, opts ...grpc.CallOption) (*OpenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpenResponse)
	err := c.cc.Invoke(ctx, NetworkFile_Open_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkFileClient) KeepAliveStream(ctx context.Context, in *KeepAliveStreamRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[KeepAliveStreamResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &NetworkFile_ServiceDesc.Streams[0], NetworkFile_KeepAliveStream_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[KeepAliveStreamRequest, KeepAliveStreamResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type NetworkFile_KeepAliveStreamClient = grpc.ServerStreamingClient[KeepAliveStreamResponse]

func (c *networkFileClient) Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WriteResponse)
	err := c.cc.Invoke(ctx, NetworkFile_Write_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkFileClient) Sync(ctx context.Context, in *SyncRequest, opts ...grpc.CallOption) (*SyncResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SyncResponse)
	err := c.cc.Invoke(ctx, NetworkFile_Sync_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkFileClient) Close(ctx context.Context, in *CloseRequest, opts ...grpc.CallOption) (*CloseResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CloseResponse)
	err := c.cc.Invoke(ctx, NetworkFile_Close_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetworkFileServer is the server API for NetworkFile service.
// All implementations must embed UnimplementedNetworkFileServer
// for forward compatibility.
type NetworkFileServer interface {
	Open(context.Context, *OpenRequest) (*OpenResponse, error)
	KeepAliveStream(*KeepAliveStreamRequest, grpc.ServerStreamingServer[KeepAliveStreamResponse]) error
	Write(context.Context, *WriteRequest) (*WriteResponse, error)
	Sync(context.Context, *SyncRequest) (*SyncResponse, error)
	Close(context.Context, *CloseRequest) (*CloseResponse, error)
	mustEmbedUnimplementedNetworkFileServer()
}

// UnimplementedNetworkFileServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNetworkFileServer struct{}

func (UnimplementedNetworkFileServer) Open(context.Context, *OpenRequest) (*OpenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Open not implemented")
}
func (UnimplementedNetworkFileServer) KeepAliveStream(*KeepAliveStreamRequest, grpc.ServerStreamingServer[KeepAliveStreamResponse]) error {
	return status.Errorf(codes.Unimplemented, "method KeepAliveStream not implemented")
}
func (UnimplementedNetworkFileServer) Write(context.Context, *WriteRequest) (*WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}
func (UnimplementedNetworkFileServer) Sync(context.Context, *SyncRequest) (*SyncResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sync not implemented")
}
func (UnimplementedNetworkFileServer) Close(context.Context, *CloseRequest) (*CloseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Close not implemented")
}
func (UnimplementedNetworkFileServer) mustEmbedUnimplementedNetworkFileServer() {}
func (UnimplementedNetworkFileServer) testEmbeddedByValue()                     {}

// UnsafeNetworkFileServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NetworkFileServer will
// result in compilation errors.
type UnsafeNetworkFileServer interface {
	mustEmbedUnimplementedNetworkFileServer()
}

func RegisterNetworkFileServer(s grpc.ServiceRegistrar, srv NetworkFileServer) {
	// If the following call pancis, it indicates UnimplementedNetworkFileServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&NetworkFile_ServiceDesc, srv)
}

func _NetworkFile_Open_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkFileServer).Open(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NetworkFile_Open_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkFileServer).Open(ctx, req.(*OpenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetworkFile_KeepAliveStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(KeepAliveStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NetworkFileServer).KeepAliveStream(m, &grpc.GenericServerStream[KeepAliveStreamRequest, KeepAliveStreamResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type NetworkFile_KeepAliveStreamServer = grpc.ServerStreamingServer[KeepAliveStreamResponse]

func _NetworkFile_Write_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkFileServer).Write(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NetworkFile_Write_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkFileServer).Write(ctx, req.(*WriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetworkFile_Sync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkFileServer).Sync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NetworkFile_Sync_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkFileServer).Sync(ctx, req.(*SyncRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetworkFile_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkFileServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NetworkFile_Close_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkFileServer).Close(ctx, req.(*CloseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NetworkFile_ServiceDesc is the grpc.ServiceDesc for NetworkFile service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NetworkFile_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.NetworkFile",
	HandlerType: (*NetworkFileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Open",
			Handler:    _NetworkFile_Open_Handler,
		},
		{
			MethodName: "Write",
			Handler:    _NetworkFile_Write_Handler,
		},
		{
			MethodName: "Sync",
			Handler:    _NetworkFile_Sync_Handler,
		},
		{
			MethodName: "Close",
			Handler:    _NetworkFile_Close_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "KeepAliveStream",
			Handler:       _NetworkFile_KeepAliveStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "networkFile.proto",
}
