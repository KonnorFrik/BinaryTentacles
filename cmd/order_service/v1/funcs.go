package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/KonnorFrik/BinaryTentacles/cmd/order_service/v1/usecase"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InterceptorLogger adapts slog logger to interceptor logger.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

// ErrorToCode - Map error to codes.Code for logging
func ErrorToCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	stat, ok := status.FromError(err)

	if ok {
		return stat.Code()
	}

	return codes.Internal
}

// RecoveryHandler - Map recovered value to error.
// For use in "recovery.WithRecoveryHandler"
func RecoveryHandler(a any) error {
	return status.Error(codes.Internal, "Something went wrong")
}

// WrapError - wrap usecase error into gRPC error with codes
func WrapError(err error) error {
	var code = codes.Internal
	var msg string

	switch {
	case errors.Is(err, usecase.ErrDoesNotExist):
		code = codes.NotFound
		msg = "object cannot be found"
	case errors.Is(err, usecase.ErrInvalidMarket):
		code = codes.FailedPrecondition
		msg = "market is unavailable"
		// case errors.Is(err, usecase.ErrInvalidData):
		// 	code = codes.InvalidArgument
		// case errors.Is(err, usecase.ErrDbNoAccess):
		// 	// default = Internal
		// case errors.Is(err, usecase.ErrUnknown):
		// default = Internal
	}

	return status.Error(code, msg)
}
