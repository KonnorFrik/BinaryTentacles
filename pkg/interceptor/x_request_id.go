package interceptor

import (
	"context"
	"log/slog"

	"github.com/KonnorFrik/BinaryTentacles/pkg/logging"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type requestIDKey struct{}

const RequestIDHeader = "x-request-id"

var (
	RequestIDKey = requestIDKey{}
	logger       = logging.Default()
)

func UnaryServerXRequestId(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (
	response any,
	err error,
) {
	id, exist, err := getRequestUUID(ctx)

	if exist {
		logger.LogAttrs(
			ctx,
			slog.LevelInfo,
			"[Interceptor/X-Request-ID]",
			slog.String("UUID", id),
		)

		return handler(ctx, req)
	}

	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"[Interceptor/X-Request-ID]",
			slog.String("ID error", err.Error()),
		)
	}

	id, err = newUUID()

	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"[Interceptor/X-Request-ID]",
			slog.String("ID Create error", err.Error()),
		)
	}

	ctx = context.WithValue(ctx, RequestIDKey, id)
	ctx = metadata.AppendToOutgoingContext(ctx, RequestIDHeader, id)
	logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"[Interceptor/X-Request-ID]",
		slog.String("UUID", id),
	)
	return handler(ctx, req)
}

type serverStreamWrapper struct {
	grpc.ServerStream
	ctx context.Context
}

func (sw *serverStreamWrapper) Context() context.Context {
	return sw.ctx
}

func StreamServerXRequestId(
	srv any,
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	ctx := ss.Context()
	id, exist, err := getRequestUUID(ctx)

	if exist {
		logger.LogAttrs(
			ctx,
			slog.LevelInfo,
			"[X-REQUEST-ID]",
			slog.String("UUID", id),
		)

		return handler(srv, ss)
	}

	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"[X-REQUEST-ID]",
			slog.String("Create-UUID Error", err.Error()),
		)
	}

	id, err = newUUID()

	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"[X-REQUEST-ID]",
			slog.String("ID Create error", err.Error()),
		)
	}

	newCtx := context.WithValue(ctx, RequestIDKey, id)
	wrappedStream := &serverStreamWrapper{
		ServerStream: ss,
		ctx:          newCtx,
	}

	return handler(srv, wrappedStream)
}

func newUUID() (
	string,
	error,
) {
	result, err := uuid.NewRandom()

	if err != nil {
		return uuid.UUID{}.String(), err
	}

	return result.String(), nil
}

func getRequestUUID(
	ctx context.Context,
) (
	id string,
	exist bool,
	err error,
) {
	if mData, ok := metadata.FromIncomingContext(ctx); ok {
		values := mData.Get(RequestIDHeader)

		if len(values) > 0 {
			err = uuid.Validate(values[0])
			return values[0], err == nil, err
		}
	}

	return "", false, nil
}
