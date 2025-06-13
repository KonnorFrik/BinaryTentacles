/*
Simple gRPC user_auth server implemented user_auth/v1
*/
package main

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/KonnorFrik/BinaryTentacles/cmd/order_service/v1/usecase"
	pb "github.com/KonnorFrik/BinaryTentacles/internal/generated/order_service/v1"
	"github.com/KonnorFrik/BinaryTentacles/pkg/logging"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedOrderServiceServer
}

var (
	logger = logging.New()
)

const (
	laddr = ":8888"
)

func main() {
	listener, err := net.Listen("tcp", laddr)

	if err != nil {
		logger.Error("[Server/Listen]", "error", err)
		return
	}

	userServer := &server{}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logger.UnaryServerInterceptor,
		),
	)
	pb.RegisterOrderServiceServer(grpcServer, userServer)
	logger.Info("Listen at", "local address", laddr)

	err = grpcServer.Serve(listener)

	if err != nil {
		logger.Error("Serve", "error", err)
		return
	}
}

// WrapError - wrap usecase error into gRPC error with codes
func WrapError(err error) error {
	var code = codes.Internal
	var msg string

	switch {
	case errors.Is(err, usecase.ErrDoesNotExist):
		code = codes.NotFound
		msg = "Object cannot be found"
		// case errors.Is(err, usecase.ErrAlreadyExist):
		// 	code = codes.AlreadyExists
		// case errors.Is(err, usecase.ErrInvalidData):
		// 	code = codes.InvalidArgument
		// case errors.Is(err, usecase.ErrDbNoAccess):
		// 	// default = Internal
		// case errors.Is(err, usecase.ErrUnknown):
		// default = Internal
	}

	return status.Error(code, msg)
}

func (s *server) Create(
	ctx context.Context,
	req *pb.CreateRequest,
) (
	*pb.CreateResponse,
	error,
) {
	order, err := usecase.Create(ctx, req)

	if err != nil {
		return nil, WrapError(err)
	}

	var response pb.CreateResponse
	order.ToGrpcCreateResponse(&response)
	return &response, status.Error(codes.OK, "ok")
}

func (s *server) OrderStatus(
	ctx context.Context,
	req *pb.OrderStatusRequest,
) (
	*pb.OrderStatusResponse,
	error,
) {
	order, err := usecase.OrderStatus(ctx, req)

	if err != nil {
		return nil, WrapError(err)
	}

	var response pb.OrderStatusResponse
	response.Status = order.GetStatus()
	return &response, status.Error(codes.OK, "ok")
}

func (s *server) OrderUpdates(
	req *pb.OrderUpdatesRequest,
	stream grpc.ServerStreamingServer[pb.OrderUpdatesResponse],
) error {

	order, err := usecase.OrderById(stream.Context(), req.GetOrderId())

	if err != nil {
		return WrapError(err)
	}

	for {
		var resp = new(pb.OrderUpdatesResponse)
		resp.Status = order.GetStatus()

		if e := stream.Send(resp); e != nil {
			// TODO: catch a closed by a client connection
			return e
		}

		time.Sleep(time.Millisecond * time.Duration(req.GetDelayMs()))
	}

	// return nil
}
