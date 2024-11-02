// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: api/job.proto

package api

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
	JobService_Create_FullMethodName   = "/api.JobService/Create"
	JobService_Claim_FullMethodName    = "/api.JobService/Claim"
	JobService_Release_FullMethodName  = "/api.JobService/Release"
	JobService_Complete_FullMethodName = "/api.JobService/Complete"
	JobService_Listen_FullMethodName   = "/api.JobService/Listen"
)

// JobServiceClient is the client API for JobService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type JobServiceClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	Claim(ctx context.Context, in *ClaimRequest, opts ...grpc.CallOption) (*ClaimResponse, error)
	Release(ctx context.Context, in *ReleaseRequest, opts ...grpc.CallOption) (*ReleaseResponse, error)
	Complete(ctx context.Context, in *CompleteRequest, opts ...grpc.CallOption) (*CompleteResponse, error)
	Listen(ctx context.Context, in *ListenRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Job], error)
}

type jobServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewJobServiceClient(cc grpc.ClientConnInterface) JobServiceClient {
	return &jobServiceClient{cc}
}

func (c *jobServiceClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, JobService_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobServiceClient) Claim(ctx context.Context, in *ClaimRequest, opts ...grpc.CallOption) (*ClaimResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClaimResponse)
	err := c.cc.Invoke(ctx, JobService_Claim_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobServiceClient) Release(ctx context.Context, in *ReleaseRequest, opts ...grpc.CallOption) (*ReleaseResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReleaseResponse)
	err := c.cc.Invoke(ctx, JobService_Release_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobServiceClient) Complete(ctx context.Context, in *CompleteRequest, opts ...grpc.CallOption) (*CompleteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CompleteResponse)
	err := c.cc.Invoke(ctx, JobService_Complete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobServiceClient) Listen(ctx context.Context, in *ListenRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Job], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &JobService_ServiceDesc.Streams[0], JobService_Listen_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ListenRequest, Job]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type JobService_ListenClient = grpc.ServerStreamingClient[Job]

// JobServiceServer is the server API for JobService service.
// All implementations must embed UnimplementedJobServiceServer
// for forward compatibility.
type JobServiceServer interface {
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Claim(context.Context, *ClaimRequest) (*ClaimResponse, error)
	Release(context.Context, *ReleaseRequest) (*ReleaseResponse, error)
	Complete(context.Context, *CompleteRequest) (*CompleteResponse, error)
	Listen(*ListenRequest, grpc.ServerStreamingServer[Job]) error
	mustEmbedUnimplementedJobServiceServer()
}

// UnimplementedJobServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedJobServiceServer struct{}

func (UnimplementedJobServiceServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedJobServiceServer) Claim(context.Context, *ClaimRequest) (*ClaimResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Claim not implemented")
}
func (UnimplementedJobServiceServer) Release(context.Context, *ReleaseRequest) (*ReleaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Release not implemented")
}
func (UnimplementedJobServiceServer) Complete(context.Context, *CompleteRequest) (*CompleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Complete not implemented")
}
func (UnimplementedJobServiceServer) Listen(*ListenRequest, grpc.ServerStreamingServer[Job]) error {
	return status.Errorf(codes.Unimplemented, "method Listen not implemented")
}
func (UnimplementedJobServiceServer) mustEmbedUnimplementedJobServiceServer() {}
func (UnimplementedJobServiceServer) testEmbeddedByValue()                    {}

// UnsafeJobServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to JobServiceServer will
// result in compilation errors.
type UnsafeJobServiceServer interface {
	mustEmbedUnimplementedJobServiceServer()
}

func RegisterJobServiceServer(s grpc.ServiceRegistrar, srv JobServiceServer) {
	// If the following call pancis, it indicates UnimplementedJobServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&JobService_ServiceDesc, srv)
}

func _JobService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: JobService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobService_Claim_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClaimRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobServiceServer).Claim(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: JobService_Claim_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).Claim(ctx, req.(*ClaimRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobService_Release_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReleaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobServiceServer).Release(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: JobService_Release_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).Release(ctx, req.(*ReleaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobService_Complete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobServiceServer).Complete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: JobService_Complete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobServiceServer).Complete(ctx, req.(*CompleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobService_Listen_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListenRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(JobServiceServer).Listen(m, &grpc.GenericServerStream[ListenRequest, Job]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type JobService_ListenServer = grpc.ServerStreamingServer[Job]

// JobService_ServiceDesc is the grpc.ServiceDesc for JobService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var JobService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.JobService",
	HandlerType: (*JobServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _JobService_Create_Handler,
		},
		{
			MethodName: "Claim",
			Handler:    _JobService_Claim_Handler,
		},
		{
			MethodName: "Release",
			Handler:    _JobService_Release_Handler,
		},
		{
			MethodName: "Complete",
			Handler:    _JobService_Complete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Listen",
			Handler:       _JobService_Listen_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/job.proto",
}
