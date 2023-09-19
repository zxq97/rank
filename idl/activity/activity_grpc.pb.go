// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.18.0
// source: idl/activity/activity.proto

package activity

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Activity_DeltaScore_FullMethodName = "/activity.Activity/DeltaScore"
	Activity_DelUser_FullMethodName    = "/activity.Activity/DelUser"
)

// ActivityClient is the client API for Activity service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ActivityClient interface {
	// 修改分数
	DeltaScore(ctx context.Context, in *DeltaScoreRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 作弊或其他行为 踢榜
	DelUser(ctx context.Context, in *DelUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type activityClient struct {
	cc grpc.ClientConnInterface
}

func NewActivityClient(cc grpc.ClientConnInterface) ActivityClient {
	return &activityClient{cc}
}

func (c *activityClient) DeltaScore(ctx context.Context, in *DeltaScoreRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Activity_DeltaScore_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *activityClient) DelUser(ctx context.Context, in *DelUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Activity_DelUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ActivityServer is the server API for Activity service.
// All implementations must embed UnimplementedActivityServer
// for forward compatibility
type ActivityServer interface {
	// 修改分数
	DeltaScore(context.Context, *DeltaScoreRequest) (*emptypb.Empty, error)
	// 作弊或其他行为 踢榜
	DelUser(context.Context, *DelUserRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedActivityServer()
}

// UnimplementedActivityServer must be embedded to have forward compatible implementations.
type UnimplementedActivityServer struct {
}

func (UnimplementedActivityServer) DeltaScore(context.Context, *DeltaScoreRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeltaScore not implemented")
}
func (UnimplementedActivityServer) DelUser(context.Context, *DelUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelUser not implemented")
}
func (UnimplementedActivityServer) mustEmbedUnimplementedActivityServer() {}

// UnsafeActivityServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ActivityServer will
// result in compilation errors.
type UnsafeActivityServer interface {
	mustEmbedUnimplementedActivityServer()
}

func RegisterActivityServer(s grpc.ServiceRegistrar, srv ActivityServer) {
	s.RegisterService(&Activity_ServiceDesc, srv)
}

func _Activity_DeltaScore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeltaScoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActivityServer).DeltaScore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Activity_DeltaScore_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActivityServer).DeltaScore(ctx, req.(*DeltaScoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Activity_DelUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActivityServer).DelUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Activity_DelUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActivityServer).DelUser(ctx, req.(*DelUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Activity_ServiceDesc is the grpc.ServiceDesc for Activity service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Activity_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "activity.Activity",
	HandlerType: (*ActivityServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeltaScore",
			Handler:    _Activity_DeltaScore_Handler,
		},
		{
			MethodName: "DelUser",
			Handler:    _Activity_DelUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "idl/activity/activity.proto",
}
