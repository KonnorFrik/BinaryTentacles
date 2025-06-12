/*
Simple gRPC user_auth server implemented user_auth/v1
*/
package main

import (
	"context"
	"errors"
	"net"

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
	laddr = ":9999"
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
	// case errors.Is(err, usecase.ErrDoesNotExist):
	// 	code = codes.NotFound
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
