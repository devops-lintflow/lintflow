// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.26.1
// source: lint/lint.proto

package lint

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	LintProto_SendLint_FullMethodName = "/lint.LintProto/SendLint"
)

// LintProtoClient is the client API for LintProto service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LintProtoClient interface {
	SendLint(ctx context.Context, in *LintRequest, opts ...grpc.CallOption) (*LintReply, error)
}

type lintProtoClient struct {
	cc grpc.ClientConnInterface
}

func NewLintProtoClient(cc grpc.ClientConnInterface) LintProtoClient {
	return &lintProtoClient{cc}
}

func (c *lintProtoClient) SendLint(ctx context.Context, in *LintRequest, opts ...grpc.CallOption) (*LintReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LintReply)
	err := c.cc.Invoke(ctx, LintProto_SendLint_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LintProtoServer is the server API for LintProto service.
// All implementations must embed UnimplementedLintProtoServer
// for forward compatibility
type LintProtoServer interface {
	SendLint(context.Context, *LintRequest) (*LintReply, error)
	mustEmbedUnimplementedLintProtoServer()
}

// UnimplementedLintProtoServer must be embedded to have forward compatible implementations.
type UnimplementedLintProtoServer struct {
}

func (UnimplementedLintProtoServer) SendLint(context.Context, *LintRequest) (*LintReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendLint not implemented")
}
func (UnimplementedLintProtoServer) mustEmbedUnimplementedLintProtoServer() {}

// UnsafeLintProtoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LintProtoServer will
// result in compilation errors.
type UnsafeLintProtoServer interface {
	mustEmbedUnimplementedLintProtoServer()
}

func RegisterLintProtoServer(s grpc.ServiceRegistrar, srv LintProtoServer) {
	s.RegisterService(&LintProto_ServiceDesc, srv)
}

func _LintProto_SendLint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LintRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LintProtoServer).SendLint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LintProto_SendLint_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LintProtoServer).SendLint(ctx, req.(*LintRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LintProto_ServiceDesc is the grpc.ServiceDesc for LintProto service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LintProto_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lint.LintProto",
	HandlerType: (*LintProtoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendLint",
			Handler:    _LintProto_SendLint_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lint/lint.proto",
}
