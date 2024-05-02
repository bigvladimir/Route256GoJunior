package grpcserver

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	orders_dto "homework/internal/app/orders/dto"
	pvz_dto "homework/internal/app/pvz/dto"
	"homework/internal/pkg/pb"
)

// orders structs begin

type orderIdentifierModel struct {
	orderID int
	pvzID   int
}

func (m *orderIdentifierModel) mapFromProto(req *pb.OrderIdentifier) {
	m.orderID = int(req.GetOrderID())
	m.pvzID = int(req.GetPvzID())
}

func (m *orderIdentifierModel) mapToProto() *pb.OrderIdentifier {
	return &pb.OrderIdentifier{
		OrderID: int64(m.orderID),
		PvzID:   int32(m.pvzID),
	}
}

type orderBodyModel struct {
	customerID      int
	storageLastTime time.Time
	packageType     string
	weight          float64
	price           int
}

func (m *orderBodyModel) mapFromProto(req *pb.OrderBody) {
	m.customerID = int(req.GetCustomerID())
	m.storageLastTime = req.GetStorageLastTime().AsTime()
	m.packageType = req.GetPackageType()
	m.weight = float64(req.GetWeight())
	m.price = int(req.GetPrice())
}

func (m *orderBodyModel) mapToProto() *pb.OrderBody {
	return &pb.OrderBody{
		CustomerID:      int32(m.customerID),
		StorageLastTime: timestamppb.New(m.storageLastTime),
		PackageType:     m.packageType,
		Weight:          float32(m.weight),
		Price:           int32(m.price),
	}
}

type orderServiceModel struct {
	isCompleted  bool
	completeTime time.Time
	isRefunded   bool
	arrivalTime  time.Time
}

func (m *orderServiceModel) mapToProto() *pb.OrderService {
	return &pb.OrderService{
		IsCompleted:  m.isCompleted,
		CompleteTime: timestamppb.New(m.completeTime),
		IsRefunded:   m.isRefunded,
		ArrivalTime:  timestamppb.New(m.arrivalTime),
	}
}

type orderInfoModel struct {
	identifier orderIdentifierModel
	body       orderBodyModel
}

func (m *orderInfoModel) mapFromProto(req *pb.OrderInfo) {
	m.identifier.mapFromProto(req.GetIdentifier())
	m.body.mapFromProto(req.GetBody())
}

func (m *orderInfoModel) mapToDTO() orders_dto.OrderInput {
	return orders_dto.OrderInput{
		OrderID: m.identifier.orderID,
		PvzID:   m.identifier.pvzID,

		CustomerID:      m.body.customerID,
		StorageLastTime: m.body.storageLastTime,
		PackageType:     m.body.packageType,
		Weight:          m.body.weight,
		Price:           m.body.price,
	}
}

type orderModel struct {
	identifier orderIdentifierModel
	body       orderBodyModel
	service    orderServiceModel
}

func (m *orderModel) mapFromDTO(d orders_dto.Order) {
	m.identifier = orderIdentifierModel{
		orderID: d.OrderID,
		pvzID:   d.PvzID,
	}
	m.body = orderBodyModel{
		customerID:      d.CustomerID,
		storageLastTime: d.StorageLastTime,
		packageType:     d.PackageType,
		weight:          d.Weight,
		price:           d.Price,
	}
	m.service = orderServiceModel{
		isCompleted:  d.IsCompleted,
		completeTime: d.CompleteTime,
		isRefunded:   d.IsRefunded,
		arrivalTime:  d.ArrivalTime,
	}
}

func (m *orderModel) mapToProto() *pb.Order {
	return &pb.Order{
		Identifier: m.identifier.mapToProto(),
		Body:       m.body.mapToProto(),
		Service:    m.service.mapToProto(),
	}
}

type orderListModel struct {
	orders []*orderModel
}

func (m *orderListModel) mapFromDTO(d []orders_dto.Order) {
	m.orders = make([]*orderModel, 0, len(d))
	for i := 0; i < len(m.orders); i++ {
		m.orders[i].mapFromDTO(d[i])
	}
}

func (m *orderListModel) mapToProto() *pb.OrderList {
	var protoOrderList *pb.OrderList = &pb.OrderList{}
	protoOrderList.Orders = make([]*pb.Order, 0, len(m.orders))
	for i := 0; i < len(protoOrderList.Orders); i++ {
		protoOrderList.Orders[i] = m.orders[i].mapToProto()
	}
	return protoOrderList
}

type orderRequestForCustomerModel struct {
	identifier orderIdentifierModel
	customerID int
}

func (m *orderRequestForCustomerModel) mapFromProto(req *pb.OrderRequestForCustomer) {
	m.identifier.mapFromProto(req.GetIdentifier())
	m.customerID = int(req.GetCustomerID())
}

type refundListRequestModel struct {
	pvzID    int
	pageNum  int
	pageSize int
}

func (m *refundListRequestModel) mapFromProto(req *pb.RefundListRequest) {
	m.pvzID = int(req.GetPvzID())
	m.pageNum = int(req.GetPageNum())
	m.pageSize = int(req.GetPageSize())
}

type customerOrderListRequestModel struct {
	pvzID      int
	customerID int
	limit      int
	isInStock  bool
}

func (m *customerOrderListRequestModel) mapFromProto(req *pb.CustomerOrderListRequest) {
	m.pvzID = int(req.GetPvzID())
	m.customerID = int(req.GetCustomerID())
	m.limit = int(req.GetLimit())
	m.isInStock = req.GetIsInStock()
}

// orders structs end

// pvz structs begin

type pvzIdentifierModel struct {
	pvzID int64
}

func (m *pvzIdentifierModel) mapFromProto(req *pb.PvzIdentifier) {
	m.pvzID = int64(req.GetPvzID())
}

func (m *pvzIdentifierModel) mapToProto() *pb.PvzIdentifier {
	return &pb.PvzIdentifier{
		PvzID: int32(m.pvzID),
	}
}

type pvzInfoModel struct {
	name     string
	adress   string
	contacts string
}

func (m *pvzInfoModel) mapToProto() *pb.PvzInfo {
	return &pb.PvzInfo{
		Name:     m.name,
		Adress:   m.adress,
		Contacts: m.contacts,
	}
}

func (m *pvzInfoModel) mapFromProto(req *pb.PvzInfo) {
	m.name = req.GetName()
	m.adress = req.GetAdress()
	m.contacts = req.GetContacts()
}

func (m *pvzInfoModel) mapToDTO() pvz_dto.PvzInput {
	return pvz_dto.PvzInput{
		Name:     m.name,
		Adress:   m.adress,
		Contacts: m.contacts,
	}
}

type pvzModel struct {
	identifier pvzIdentifierModel
	info       pvzInfoModel
}

func (m *pvzModel) mapToDTO() pvz_dto.Pvz {
	return pvz_dto.Pvz{
		ID: m.identifier.pvzID,

		Name:     m.info.name,
		Adress:   m.info.adress,
		Contacts: m.info.contacts,
	}
}

func (m *pvzModel) mapFromDTO(d pvz_dto.Pvz) {
	m.identifier.pvzID = d.ID

	m.info.name = d.Name
	m.info.adress = d.Adress
	m.info.contacts = d.Contacts
}

func (m *pvzModel) mapToProto() *pb.Pvz {
	return &pb.Pvz{
		Identifier: m.identifier.mapToProto(),
		Info:       m.info.mapToProto(),
	}
}

func (m *pvzModel) mapFromProto(req *pb.Pvz) {
	m.identifier.mapFromProto(req.GetIdentifier())
	m.info.mapFromProto(req.GetInfo())
}

// pvz structs end
