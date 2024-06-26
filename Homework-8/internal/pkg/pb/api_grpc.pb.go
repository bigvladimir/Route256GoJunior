// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: api.proto

package pb

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
	PvzManager_TakeOrderFromCourier_FullMethodName   = "/server.PvzManager/TakeOrderFromCourier"
	PvzManager_ReturnOrderToCourier_FullMethodName   = "/server.PvzManager/ReturnOrderToCourier"
	PvzManager_GiveOrderToCustomer_FullMethodName    = "/server.PvzManager/GiveOrderToCustomer"
	PvzManager_TakeRefundFromCustomer_FullMethodName = "/server.PvzManager/TakeRefundFromCustomer"
	PvzManager_GetRefundList_FullMethodName          = "/server.PvzManager/GetRefundList"
	PvzManager_GetCustomerOrderList_FullMethodName   = "/server.PvzManager/GetCustomerOrderList"
	PvzManager_GetPvzByID_FullMethodName             = "/server.PvzManager/GetPvzByID"
	PvzManager_AddPvz_FullMethodName                 = "/server.PvzManager/AddPvz"
	PvzManager_ModifyPvz_FullMethodName              = "/server.PvzManager/ModifyPvz"
	PvzManager_UpdatePvz_FullMethodName              = "/server.PvzManager/UpdatePvz"
	PvzManager_DeletePvz_FullMethodName              = "/server.PvzManager/DeletePvz"
)

// PvzManagerClient is the client API for PvzManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PvzManagerClient interface {
	TakeOrderFromCourier(ctx context.Context, in *OrderInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ReturnOrderToCourier(ctx context.Context, in *OrderIdentifier, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GiveOrderToCustomer(ctx context.Context, in *OrderRequestForCustomer, opts ...grpc.CallOption) (*emptypb.Empty, error)
	TakeRefundFromCustomer(ctx context.Context, in *OrderRequestForCustomer, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetRefundList(ctx context.Context, in *RefundListRequest, opts ...grpc.CallOption) (*OrderList, error)
	GetCustomerOrderList(ctx context.Context, in *CustomerOrderListRequest, opts ...grpc.CallOption) (*OrderList, error)
	GetPvzByID(ctx context.Context, in *PvzIdentifier, opts ...grpc.CallOption) (*Pvz, error)
	AddPvz(ctx context.Context, in *PvzInfo, opts ...grpc.CallOption) (*PvzIdentifier, error)
	ModifyPvz(ctx context.Context, in *Pvz, opts ...grpc.CallOption) (*PvzIdentifier, error)
	UpdatePvz(ctx context.Context, in *Pvz, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeletePvz(ctx context.Context, in *PvzIdentifier, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type pvzManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewPvzManagerClient(cc grpc.ClientConnInterface) PvzManagerClient {
	return &pvzManagerClient{cc}
}

func (c *pvzManagerClient) TakeOrderFromCourier(ctx context.Context, in *OrderInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PvzManager_TakeOrderFromCourier_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pvzManagerClient) ReturnOrderToCourier(ctx context.Context, in *OrderIdentifier, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PvzManager_ReturnOrderToCourier_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pvzManagerClient) GiveOrderToCustomer(ctx context.Context, in *OrderRequestForCustomer, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PvzManager_GiveOrderToCustomer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pvzManagerClient) TakeRefundFromCustomer(ctx context.Context, in *OrderRequestForCustomer, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PvzManager_TakeRefundFromCustomer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pvzManagerClient) GetRefundList(ctx context.Context, in *RefundListRequest, opts ...grpc.CallOption) (*OrderList, error) {
	out := new(OrderList)
	err := c.cc.Invoke(ctx, PvzManager_GetRefundList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pvzManagerClient) GetCustomerOrderList(ctx context.Context, in *CustomerOrderListRequest, opts ...grpc.CallOption) (*OrderList, error) {
	out := new(OrderList)
	err := c.cc.Invoke(ctx, PvzManager_GetCustomerOrderList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pvzManagerClient) GetPvzByID(ctx context.Context, in *PvzIdentifier, opts ...grpc.CallOption) (*Pvz, error) {
	out := new(Pvz)
	err := c.cc.Invoke(ctx, PvzManager_GetPvzByID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pvzManagerClient) AddPvz(ctx context.Context, in *PvzInfo, opts ...grpc.CallOption) (*PvzIdentifier, error) {
	out := new(PvzIdentifier)
	err := c.cc.Invoke(ctx, PvzManager_AddPvz_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pvzManagerClient) ModifyPvz(ctx context.Context, in *Pvz, opts ...grpc.CallOption) (*PvzIdentifier, error) {
	out := new(PvzIdentifier)
	err := c.cc.Invoke(ctx, PvzManager_ModifyPvz_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pvzManagerClient) UpdatePvz(ctx context.Context, in *Pvz, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PvzManager_UpdatePvz_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pvzManagerClient) DeletePvz(ctx context.Context, in *PvzIdentifier, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PvzManager_DeletePvz_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PvzManagerServer is the server API for PvzManager service.
// All implementations must embed UnimplementedPvzManagerServer
// for forward compatibility
type PvzManagerServer interface {
	TakeOrderFromCourier(context.Context, *OrderInfo) (*emptypb.Empty, error)
	ReturnOrderToCourier(context.Context, *OrderIdentifier) (*emptypb.Empty, error)
	GiveOrderToCustomer(context.Context, *OrderRequestForCustomer) (*emptypb.Empty, error)
	TakeRefundFromCustomer(context.Context, *OrderRequestForCustomer) (*emptypb.Empty, error)
	GetRefundList(context.Context, *RefundListRequest) (*OrderList, error)
	GetCustomerOrderList(context.Context, *CustomerOrderListRequest) (*OrderList, error)
	GetPvzByID(context.Context, *PvzIdentifier) (*Pvz, error)
	AddPvz(context.Context, *PvzInfo) (*PvzIdentifier, error)
	ModifyPvz(context.Context, *Pvz) (*PvzIdentifier, error)
	UpdatePvz(context.Context, *Pvz) (*emptypb.Empty, error)
	DeletePvz(context.Context, *PvzIdentifier) (*emptypb.Empty, error)
	mustEmbedUnimplementedPvzManagerServer()
}

// UnimplementedPvzManagerServer must be embedded to have forward compatible implementations.
type UnimplementedPvzManagerServer struct {
}

func (UnimplementedPvzManagerServer) TakeOrderFromCourier(context.Context, *OrderInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TakeOrderFromCourier not implemented")
}
func (UnimplementedPvzManagerServer) ReturnOrderToCourier(context.Context, *OrderIdentifier) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReturnOrderToCourier not implemented")
}
func (UnimplementedPvzManagerServer) GiveOrderToCustomer(context.Context, *OrderRequestForCustomer) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GiveOrderToCustomer not implemented")
}
func (UnimplementedPvzManagerServer) TakeRefundFromCustomer(context.Context, *OrderRequestForCustomer) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TakeRefundFromCustomer not implemented")
}
func (UnimplementedPvzManagerServer) GetRefundList(context.Context, *RefundListRequest) (*OrderList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRefundList not implemented")
}
func (UnimplementedPvzManagerServer) GetCustomerOrderList(context.Context, *CustomerOrderListRequest) (*OrderList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCustomerOrderList not implemented")
}
func (UnimplementedPvzManagerServer) GetPvzByID(context.Context, *PvzIdentifier) (*Pvz, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPvzByID not implemented")
}
func (UnimplementedPvzManagerServer) AddPvz(context.Context, *PvzInfo) (*PvzIdentifier, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPvz not implemented")
}
func (UnimplementedPvzManagerServer) ModifyPvz(context.Context, *Pvz) (*PvzIdentifier, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifyPvz not implemented")
}
func (UnimplementedPvzManagerServer) UpdatePvz(context.Context, *Pvz) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePvz not implemented")
}
func (UnimplementedPvzManagerServer) DeletePvz(context.Context, *PvzIdentifier) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePvz not implemented")
}
func (UnimplementedPvzManagerServer) mustEmbedUnimplementedPvzManagerServer() {}

// UnsafePvzManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PvzManagerServer will
// result in compilation errors.
type UnsafePvzManagerServer interface {
	mustEmbedUnimplementedPvzManagerServer()
}

func RegisterPvzManagerServer(s grpc.ServiceRegistrar, srv PvzManagerServer) {
	s.RegisterService(&PvzManager_ServiceDesc, srv)
}

func _PvzManager_TakeOrderFromCourier_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PvzManagerServer).TakeOrderFromCourier(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PvzManager_TakeOrderFromCourier_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PvzManagerServer).TakeOrderFromCourier(ctx, req.(*OrderInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _PvzManager_ReturnOrderToCourier_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderIdentifier)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PvzManagerServer).ReturnOrderToCourier(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PvzManager_ReturnOrderToCourier_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PvzManagerServer).ReturnOrderToCourier(ctx, req.(*OrderIdentifier))
	}
	return interceptor(ctx, in, info, handler)
}

func _PvzManager_GiveOrderToCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderRequestForCustomer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PvzManagerServer).GiveOrderToCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PvzManager_GiveOrderToCustomer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PvzManagerServer).GiveOrderToCustomer(ctx, req.(*OrderRequestForCustomer))
	}
	return interceptor(ctx, in, info, handler)
}

func _PvzManager_TakeRefundFromCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderRequestForCustomer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PvzManagerServer).TakeRefundFromCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PvzManager_TakeRefundFromCustomer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PvzManagerServer).TakeRefundFromCustomer(ctx, req.(*OrderRequestForCustomer))
	}
	return interceptor(ctx, in, info, handler)
}

func _PvzManager_GetRefundList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefundListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PvzManagerServer).GetRefundList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PvzManager_GetRefundList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PvzManagerServer).GetRefundList(ctx, req.(*RefundListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PvzManager_GetCustomerOrderList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CustomerOrderListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PvzManagerServer).GetCustomerOrderList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PvzManager_GetCustomerOrderList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PvzManagerServer).GetCustomerOrderList(ctx, req.(*CustomerOrderListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PvzManager_GetPvzByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PvzIdentifier)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PvzManagerServer).GetPvzByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PvzManager_GetPvzByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PvzManagerServer).GetPvzByID(ctx, req.(*PvzIdentifier))
	}
	return interceptor(ctx, in, info, handler)
}

func _PvzManager_AddPvz_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PvzInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PvzManagerServer).AddPvz(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PvzManager_AddPvz_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PvzManagerServer).AddPvz(ctx, req.(*PvzInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _PvzManager_ModifyPvz_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Pvz)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PvzManagerServer).ModifyPvz(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PvzManager_ModifyPvz_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PvzManagerServer).ModifyPvz(ctx, req.(*Pvz))
	}
	return interceptor(ctx, in, info, handler)
}

func _PvzManager_UpdatePvz_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Pvz)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PvzManagerServer).UpdatePvz(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PvzManager_UpdatePvz_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PvzManagerServer).UpdatePvz(ctx, req.(*Pvz))
	}
	return interceptor(ctx, in, info, handler)
}

func _PvzManager_DeletePvz_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PvzIdentifier)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PvzManagerServer).DeletePvz(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PvzManager_DeletePvz_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PvzManagerServer).DeletePvz(ctx, req.(*PvzIdentifier))
	}
	return interceptor(ctx, in, info, handler)
}

// PvzManager_ServiceDesc is the grpc.ServiceDesc for PvzManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PvzManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "server.PvzManager",
	HandlerType: (*PvzManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TakeOrderFromCourier",
			Handler:    _PvzManager_TakeOrderFromCourier_Handler,
		},
		{
			MethodName: "ReturnOrderToCourier",
			Handler:    _PvzManager_ReturnOrderToCourier_Handler,
		},
		{
			MethodName: "GiveOrderToCustomer",
			Handler:    _PvzManager_GiveOrderToCustomer_Handler,
		},
		{
			MethodName: "TakeRefundFromCustomer",
			Handler:    _PvzManager_TakeRefundFromCustomer_Handler,
		},
		{
			MethodName: "GetRefundList",
			Handler:    _PvzManager_GetRefundList_Handler,
		},
		{
			MethodName: "GetCustomerOrderList",
			Handler:    _PvzManager_GetCustomerOrderList_Handler,
		},
		{
			MethodName: "GetPvzByID",
			Handler:    _PvzManager_GetPvzByID_Handler,
		},
		{
			MethodName: "AddPvz",
			Handler:    _PvzManager_AddPvz_Handler,
		},
		{
			MethodName: "ModifyPvz",
			Handler:    _PvzManager_ModifyPvz_Handler,
		},
		{
			MethodName: "UpdatePvz",
			Handler:    _PvzManager_UpdatePvz_Handler,
		},
		{
			MethodName: "DeletePvz",
			Handler:    _PvzManager_DeletePvz_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
