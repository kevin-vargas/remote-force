// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: remote-force.proto

package pbs

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

// RemoteClient is the client API for Remote service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RemoteClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	Execute(ctx context.Context, in *ExecuteRequest, opts ...grpc.CallOption) (*ExecuteResponse, error)
	Commands(ctx context.Context, in *CommandsRequest, opts ...grpc.CallOption) (*CommandsResponse, error)
}

type remoteClient struct {
	cc grpc.ClientConnInterface
}

func NewRemoteClient(cc grpc.ClientConnInterface) RemoteClient {
	return &remoteClient{cc}
}

func (c *remoteClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/V1.Remote/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *remoteClient) Execute(ctx context.Context, in *ExecuteRequest, opts ...grpc.CallOption) (*ExecuteResponse, error) {
	out := new(ExecuteResponse)
	err := c.cc.Invoke(ctx, "/V1.Remote/Execute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *remoteClient) Commands(ctx context.Context, in *CommandsRequest, opts ...grpc.CallOption) (*CommandsResponse, error) {
	out := new(CommandsResponse)
	err := c.cc.Invoke(ctx, "/V1.Remote/Commands", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RemoteServer is the server API for Remote service.
// All implementations must embed UnimplementedRemoteServer
// for forward compatibility
type RemoteServer interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Execute(context.Context, *ExecuteRequest) (*ExecuteResponse, error)
	Commands(context.Context, *CommandsRequest) (*CommandsResponse, error)
	mustEmbedUnimplementedRemoteServer()
}

// UnimplementedRemoteServer must be embedded to have forward compatible implementations.
type UnimplementedRemoteServer struct {
}

func (UnimplementedRemoteServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedRemoteServer) Execute(context.Context, *ExecuteRequest) (*ExecuteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Execute not implemented")
}
func (UnimplementedRemoteServer) Commands(context.Context, *CommandsRequest) (*CommandsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Commands not implemented")
}
func (UnimplementedRemoteServer) mustEmbedUnimplementedRemoteServer() {}

// UnsafeRemoteServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RemoteServer will
// result in compilation errors.
type UnsafeRemoteServer interface {
	mustEmbedUnimplementedRemoteServer()
}

func RegisterRemoteServer(s grpc.ServiceRegistrar, srv RemoteServer) {
	s.RegisterService(&Remote_ServiceDesc, srv)
}

func _Remote_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoteServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/V1.Remote/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoteServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Remote_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoteServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/V1.Remote/Execute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoteServer).Execute(ctx, req.(*ExecuteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Remote_Commands_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoteServer).Commands(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/V1.Remote/Commands",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoteServer).Commands(ctx, req.(*CommandsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Remote_ServiceDesc is the grpc.ServiceDesc for Remote service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Remote_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "V1.Remote",
	HandlerType: (*RemoteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _Remote_Login_Handler,
		},
		{
			MethodName: "Execute",
			Handler:    _Remote_Execute_Handler,
		},
		{
			MethodName: "Commands",
			Handler:    _Remote_Commands_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "remote-force.proto",
}
