/*
Implement logging interceptor for use in gRPC
*/
package logging

import (
	"context"
	"os"

	"log/slog"

	"google.golang.org/grpc"
)

type Logger struct {
	*slog.Logger
}

// New - Create new 'Logger'
func New() *Logger {
	l := &Logger{
		Logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{})),
	}
	return l
}

var defaultLogger = New()

func Default() *Logger {
	return defaultLogger
}

func (l *Logger) UnaryServerInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	l.LogAttrs(
		ctx,
		slog.LevelInfo,
		"[Server]",
		slog.String("method", info.FullMethod),
	)
	res, err := handler(ctx, req)

	if err != nil {
		l.LogAttrs(
			ctx,
			slog.LevelError,
			"After handler",
			slog.String("method", info.FullMethod),
			slog.String("error", err.Error()),
		)
	}

	return res, err
}

func (l *Logger) UnaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	l.LogAttrs(
		ctx,
		slog.LevelInfo,
		"[Client]",
		slog.String("method", method),
	)
	err := invoker(ctx, method, req, reply, cc, opts...)

	if err != nil {
		l.LogAttrs(
			ctx,
			slog.LevelError,
			"After handler",
			slog.String("method", method),
			slog.String("error", err.Error()),
		)
	}

	return nil
}
