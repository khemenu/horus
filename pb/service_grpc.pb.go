//

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.0
// source: service.proto

package pb

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

const (
	Horus_NewOrg_FullMethodName     = "/khepri.horus.Horus/NewOrg"
	Horus_ListOrgs_FullMethodName   = "/khepri.horus.Horus/ListOrgs"
	Horus_UpdateOrg_FullMethodName  = "/khepri.horus.Horus/UpdateOrg"
	Horus_InviteUser_FullMethodName = "/khepri.horus.Horus/InviteUser"
	Horus_JoinOrg_FullMethodName    = "/khepri.horus.Horus/JoinOrg"
	Horus_LeaveOrg_FullMethodName   = "/khepri.horus.Horus/LeaveOrg"
)

// HorusClient is the client API for Horus service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HorusClient interface {
	// Creates organization.
	NewOrg(ctx context.Context, in *NewOrgReq, opts ...grpc.CallOption) (*NewOrgRes, error)
	// Lists organizations the user belongs to.
	ListOrgs(ctx context.Context, in *ListOrgsReq, opts ...grpc.CallOption) (*ListOrgsRes, error)
	// Updates orgnataion info.
	UpdateOrg(ctx context.Context, in *UpdateOrgReq, opts ...grpc.CallOption) (*UpdateOrgRes, error)
	// Invites a user to the organization.
	InviteUser(ctx context.Context, in *InviteUserReq, opts ...grpc.CallOption) (*InviteUserRes, error)
	// Joins an organization.
	JoinOrg(ctx context.Context, in *JoinOrgReq, opts ...grpc.CallOption) (*JoinOrgRes, error)
	// Leaves an organization.
	LeaveOrg(ctx context.Context, in *LeaveOrgReq, opts ...grpc.CallOption) (*LeaveOrgRes, error)
}

type horusClient struct {
	cc grpc.ClientConnInterface
}

func NewHorusClient(cc grpc.ClientConnInterface) HorusClient {
	return &horusClient{cc}
}

func (c *horusClient) NewOrg(ctx context.Context, in *NewOrgReq, opts ...grpc.CallOption) (*NewOrgRes, error) {
	out := new(NewOrgRes)
	err := c.cc.Invoke(ctx, Horus_NewOrg_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *horusClient) ListOrgs(ctx context.Context, in *ListOrgsReq, opts ...grpc.CallOption) (*ListOrgsRes, error) {
	out := new(ListOrgsRes)
	err := c.cc.Invoke(ctx, Horus_ListOrgs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *horusClient) UpdateOrg(ctx context.Context, in *UpdateOrgReq, opts ...grpc.CallOption) (*UpdateOrgRes, error) {
	out := new(UpdateOrgRes)
	err := c.cc.Invoke(ctx, Horus_UpdateOrg_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *horusClient) InviteUser(ctx context.Context, in *InviteUserReq, opts ...grpc.CallOption) (*InviteUserRes, error) {
	out := new(InviteUserRes)
	err := c.cc.Invoke(ctx, Horus_InviteUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *horusClient) JoinOrg(ctx context.Context, in *JoinOrgReq, opts ...grpc.CallOption) (*JoinOrgRes, error) {
	out := new(JoinOrgRes)
	err := c.cc.Invoke(ctx, Horus_JoinOrg_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *horusClient) LeaveOrg(ctx context.Context, in *LeaveOrgReq, opts ...grpc.CallOption) (*LeaveOrgRes, error) {
	out := new(LeaveOrgRes)
	err := c.cc.Invoke(ctx, Horus_LeaveOrg_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HorusServer is the server API for Horus service.
// All implementations must embed UnimplementedHorusServer
// for forward compatibility
type HorusServer interface {
	// Creates organization.
	NewOrg(context.Context, *NewOrgReq) (*NewOrgRes, error)
	// Lists organizations the user belongs to.
	ListOrgs(context.Context, *ListOrgsReq) (*ListOrgsRes, error)
	// Updates orgnataion info.
	UpdateOrg(context.Context, *UpdateOrgReq) (*UpdateOrgRes, error)
	// Invites a user to the organization.
	InviteUser(context.Context, *InviteUserReq) (*InviteUserRes, error)
	// Joins an organization.
	JoinOrg(context.Context, *JoinOrgReq) (*JoinOrgRes, error)
	// Leaves an organization.
	LeaveOrg(context.Context, *LeaveOrgReq) (*LeaveOrgRes, error)
	mustEmbedUnimplementedHorusServer()
}

// UnimplementedHorusServer must be embedded to have forward compatible implementations.
type UnimplementedHorusServer struct {
}

func (UnimplementedHorusServer) NewOrg(context.Context, *NewOrgReq) (*NewOrgRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewOrg not implemented")
}
func (UnimplementedHorusServer) ListOrgs(context.Context, *ListOrgsReq) (*ListOrgsRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListOrgs not implemented")
}
func (UnimplementedHorusServer) UpdateOrg(context.Context, *UpdateOrgReq) (*UpdateOrgRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateOrg not implemented")
}
func (UnimplementedHorusServer) InviteUser(context.Context, *InviteUserReq) (*InviteUserRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InviteUser not implemented")
}
func (UnimplementedHorusServer) JoinOrg(context.Context, *JoinOrgReq) (*JoinOrgRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinOrg not implemented")
}
func (UnimplementedHorusServer) LeaveOrg(context.Context, *LeaveOrgReq) (*LeaveOrgRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveOrg not implemented")
}
func (UnimplementedHorusServer) mustEmbedUnimplementedHorusServer() {}

// UnsafeHorusServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HorusServer will
// result in compilation errors.
type UnsafeHorusServer interface {
	mustEmbedUnimplementedHorusServer()
}

func RegisterHorusServer(s grpc.ServiceRegistrar, srv HorusServer) {
	s.RegisterService(&Horus_ServiceDesc, srv)
}

func _Horus_NewOrg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewOrgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HorusServer).NewOrg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Horus_NewOrg_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HorusServer).NewOrg(ctx, req.(*NewOrgReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Horus_ListOrgs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListOrgsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HorusServer).ListOrgs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Horus_ListOrgs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HorusServer).ListOrgs(ctx, req.(*ListOrgsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Horus_UpdateOrg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateOrgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HorusServer).UpdateOrg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Horus_UpdateOrg_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HorusServer).UpdateOrg(ctx, req.(*UpdateOrgReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Horus_InviteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InviteUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HorusServer).InviteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Horus_InviteUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HorusServer).InviteUser(ctx, req.(*InviteUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Horus_JoinOrg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinOrgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HorusServer).JoinOrg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Horus_JoinOrg_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HorusServer).JoinOrg(ctx, req.(*JoinOrgReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Horus_LeaveOrg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveOrgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HorusServer).LeaveOrg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Horus_LeaveOrg_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HorusServer).LeaveOrg(ctx, req.(*LeaveOrgReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Horus_ServiceDesc is the grpc.ServiceDesc for Horus service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Horus_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "khepri.horus.Horus",
	HandlerType: (*HorusServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewOrg",
			Handler:    _Horus_NewOrg_Handler,
		},
		{
			MethodName: "ListOrgs",
			Handler:    _Horus_ListOrgs_Handler,
		},
		{
			MethodName: "UpdateOrg",
			Handler:    _Horus_UpdateOrg_Handler,
		},
		{
			MethodName: "InviteUser",
			Handler:    _Horus_InviteUser_Handler,
		},
		{
			MethodName: "JoinOrg",
			Handler:    _Horus_JoinOrg_Handler,
		},
		{
			MethodName: "LeaveOrg",
			Handler:    _Horus_LeaveOrg_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
