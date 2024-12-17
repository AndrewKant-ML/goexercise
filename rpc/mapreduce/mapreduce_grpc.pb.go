// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: mapreduce.proto

package mapreduce

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
	MapperService_Map_FullMethodName = "/rpc.mapreduce.MapperService/Map"
)

// MapperServiceClient is the client API for MapperService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MapperServiceClient interface {
	Map(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[Number, Status], error)
}

type mapperServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMapperServiceClient(cc grpc.ClientConnInterface) MapperServiceClient {
	return &mapperServiceClient{cc}
}

func (c *mapperServiceClient) Map(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[Number, Status], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &MapperService_ServiceDesc.Streams[0], MapperService_Map_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[Number, Status]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MapperService_MapClient = grpc.ClientStreamingClient[Number, Status]

// MapperServiceServer is the server API for MapperService service.
// All implementations must embed UnimplementedMapperServiceServer
// for forward compatibility.
type MapperServiceServer interface {
	Map(grpc.ClientStreamingServer[Number, Status]) error
	mustEmbedUnimplementedMapperServiceServer()
}

// UnimplementedMapperServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMapperServiceServer struct{}

func (UnimplementedMapperServiceServer) Map(grpc.ClientStreamingServer[Number, Status]) error {
	return status.Errorf(codes.Unimplemented, "method Map not implemented")
}
func (UnimplementedMapperServiceServer) mustEmbedUnimplementedMapperServiceServer() {}
func (UnimplementedMapperServiceServer) testEmbeddedByValue()                       {}

// UnsafeMapperServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MapperServiceServer will
// result in compilation errors.
type UnsafeMapperServiceServer interface {
	mustEmbedUnimplementedMapperServiceServer()
}

func RegisterMapperServiceServer(s grpc.ServiceRegistrar, srv MapperServiceServer) {
	// If the following call pancis, it indicates UnimplementedMapperServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MapperService_ServiceDesc, srv)
}

func _MapperService_Map_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MapperServiceServer).Map(&grpc.GenericServerStream[Number, Status]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MapperService_MapServer = grpc.ClientStreamingServer[Number, Status]

// MapperService_ServiceDesc is the grpc.ServiceDesc for MapperService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MapperService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.mapreduce.MapperService",
	HandlerType: (*MapperServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Map",
			Handler:       _MapperService_Map_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "mapreduce.proto",
}

const (
	ReducerService_Reduce_FullMethodName = "/rpc.mapreduce.ReducerService/Reduce"
)

// ReducerServiceClient is the client API for ReducerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReducerServiceClient interface {
	Reduce(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[Number, Status], error)
}

type reducerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReducerServiceClient(cc grpc.ClientConnInterface) ReducerServiceClient {
	return &reducerServiceClient{cc}
}

func (c *reducerServiceClient) Reduce(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[Number, Status], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ReducerService_ServiceDesc.Streams[0], ReducerService_Reduce_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[Number, Status]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ReducerService_ReduceClient = grpc.ClientStreamingClient[Number, Status]

// ReducerServiceServer is the server API for ReducerService service.
// All implementations must embed UnimplementedReducerServiceServer
// for forward compatibility.
type ReducerServiceServer interface {
	Reduce(grpc.ClientStreamingServer[Number, Status]) error
	mustEmbedUnimplementedReducerServiceServer()
}

// UnimplementedReducerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedReducerServiceServer struct{}

func (UnimplementedReducerServiceServer) Reduce(grpc.ClientStreamingServer[Number, Status]) error {
	return status.Errorf(codes.Unimplemented, "method Reduce not implemented")
}
func (UnimplementedReducerServiceServer) mustEmbedUnimplementedReducerServiceServer() {}
func (UnimplementedReducerServiceServer) testEmbeddedByValue()                        {}

// UnsafeReducerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReducerServiceServer will
// result in compilation errors.
type UnsafeReducerServiceServer interface {
	mustEmbedUnimplementedReducerServiceServer()
}

func RegisterReducerServiceServer(s grpc.ServiceRegistrar, srv ReducerServiceServer) {
	// If the following call pancis, it indicates UnimplementedReducerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ReducerService_ServiceDesc, srv)
}

func _ReducerService_Reduce_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ReducerServiceServer).Reduce(&grpc.GenericServerStream[Number, Status]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ReducerService_ReduceServer = grpc.ClientStreamingServer[Number, Status]

// ReducerService_ServiceDesc is the grpc.ServiceDesc for ReducerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReducerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.mapreduce.ReducerService",
	HandlerType: (*ReducerServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Reduce",
			Handler:       _ReducerService_Reduce_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "mapreduce.proto",
}