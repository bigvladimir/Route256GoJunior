package grpcserver

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pvz_errors "homework/internal/app/pvz/errors"
	"homework/internal/pkg/pb"
)

// GetPvzByID возвращает запись ПВЗ по id, если она существует
func (s *Server) GetPvzByID(ctx context.Context, req *pb.PvzIdentifier) (*pb.Pvz, error) {
	spanCtx, span := s.tracer.Start(ctx, "GetPvzByID")
	defer span.End()

	var pvzIdentifier pvzIdentifierModel
	pvzIdentifier.mapFromProto(req)

	if err := pvzIdentifier.validate(); err != nil {
		return &pb.Pvz{}, status.Errorf(codes.InvalidArgument, "pvzIdentifierModel.validate: %v", err)
	}
	pvz, err := s.service.GetPvzByID(spanCtx, pvzIdentifier.pvzID)
	if err != nil {
		if errors.Is(err, pvz_errors.ErrNotFound) {
			return &pb.Pvz{}, status.Error(codes.NotFound, codes.NotFound.String())
		}
		return &pb.Pvz{}, status.Errorf(codes.Internal, "GetPvzByID: %v", err)
	}
	var protoPvz pvzModel
	protoPvz.mapFromDTO(pvz)

	return protoPvz.mapToProto(), nil
}

// AddPvz добавляет запись ПВЗ без указания id
func (s *Server) AddPvz(ctx context.Context, req *pb.PvzInfo) (*pb.PvzIdentifier, error) {
	spanCtx, span := s.tracer.Start(ctx, "AddPvz")
	defer span.End()

	var pvzInfo pvzInfoModel
	pvzInfo.mapFromProto(req)

	if err := pvzInfo.validate(); err != nil {
		return &pb.PvzIdentifier{}, status.Errorf(codes.InvalidArgument, "pvzInfoModel.validate: %v", err)
	}
	pvzID, err := s.service.AddPvz(spanCtx, pvzInfo.mapToDTO())
	if err != nil {
		return &pb.PvzIdentifier{}, status.Errorf(codes.Internal, "AddPvz: %v", err)
	}
	protoPvzIdentifier := pvzIdentifierModel{pvzID: pvzID}

	return protoPvzIdentifier.mapToProto(), nil
}

// ModifyPvz обновляет запись ПВЗ по id, если не находит, что обновить, то вставляет новую
func (s *Server) ModifyPvz(ctx context.Context, req *pb.Pvz) (*pb.PvzIdentifier, error) {
	spanCtx, span := s.tracer.Start(ctx, "ModifyPvz")
	defer span.End()

	var pvz pvzModel
	pvz.mapFromProto(req)

	if err := pvz.validate(); err != nil {
		return &pb.PvzIdentifier{}, status.Errorf(codes.InvalidArgument, "pvzModel.validate: %v", err)
	}
	pvzID, err := s.service.ModifyPvz(spanCtx, pvz.mapToDTO())
	if err != nil {
		return &pb.PvzIdentifier{}, status.Errorf(codes.Internal, "ModifyPvz: %v", err)
	}
	protoPvzIdentifier := pvzIdentifierModel{pvzID: pvzID}

	return protoPvzIdentifier.mapToProto(), nil
}

// UpdatePvz обновляет запись ПВЗ по id, если она существует
func (s *Server) UpdatePvz(ctx context.Context, req *pb.Pvz) (*emptypb.Empty, error) {
	spanCtx, span := s.tracer.Start(ctx, "UpdatePvz")
	defer span.End()

	var pvz pvzModel
	pvz.mapFromProto(req)

	if err := pvz.validate(); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "pvzModel.validate: %v", err)
	}
	if err := s.service.UpdatePvz(spanCtx, pvz.mapToDTO()); err != nil {
		if errors.Is(err, pvz_errors.ErrNotFound) {
			return &emptypb.Empty{}, status.Error(codes.NotFound, codes.NotFound.String())
		}
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "UpdatePvz: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// DeletePvz удаляет запись ПВЗ по id, если она существует
func (s *Server) DeletePvz(ctx context.Context, req *pb.PvzIdentifier) (*emptypb.Empty, error) {
	spanCtx, span := s.tracer.Start(ctx, "DeletePvz")
	defer span.End()

	var pvzIdentifier pvzIdentifierModel
	pvzIdentifier.mapFromProto(req)

	if err := pvzIdentifier.validate(); err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "pvzIdentifier.validate: %v", err)
	}
	if err := s.service.DeletePvz(spanCtx, pvzIdentifier.pvzID); err != nil {
		if errors.Is(err, pvz_errors.ErrNotFound) {
			return &emptypb.Empty{}, status.Error(codes.NotFound, codes.NotFound.String())
		}
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "DeletePvz: %v", err)
	}

	return &emptypb.Empty{}, nil
}
