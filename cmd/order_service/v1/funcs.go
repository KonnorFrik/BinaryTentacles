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

func ContextValueAs[T any](ctx context.Context, key any) (T, bool) {
	valueAny := ctx.Value(key)

	if valueAny == nil {
		var empty T
		return empty, false
	}

	valueTyped, ok := valueAny.(T)
	return valueTyped, ok
}

// InterceptorLogger adapts slog logger to interceptor logger.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

// ErrorToCode - Map error to codes.Code for logging.
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
// For use in "recovery.WithRecoveryHandler".
func RecoveryHandler(a any) error {
	slog.LogAttrs(
		nil,
		slog.LevelError,
		"[RECOVERY]",
		slog.Any("recovered with", a),
	)
	return status.Error(codes.Internal, "Something went wrong")
}

// WrapError - wrap usecase error into gRPC error with codes.
func wrapError(err error) error {
	if err == nil {
		return nil
	}

	var code = codes.Internal
	var msg string

	switch {
	case errors.Is(err, usecase.ErrDoesNotExist):
		code = codes.NotFound
		msg = "object cannot be found"
	case errors.Is(err, usecase.ErrMarketUnavailable):
		code = codes.FailedPrecondition
		msg = "market is unavailable"
	case errors.Is(err, usecase.ErrUnknown):
		code = codes.Internal
		msg = "something went wrong"
	case errors.Is(err, usecase.ErrInternal):
		code = codes.Internal
		msg = "something went wrong"
	}

	return status.Error(code, msg)
}
