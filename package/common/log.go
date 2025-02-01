package common

import (
	"os"

	"log/slog"
)

type HttpLog struct {
	*slog.Logger
}
type GrpcLog struct {
	*slog.Logger
}

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	})))
}

func NewHttpLog() HttpLog {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	}).WithAttrs([]slog.Attr{
		slog.String("proto", "http"),
	})
	logger := slog.New(logHandler)

	return HttpLog{logger}
}

func NewGrpcLog() GrpcLog {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	}).WithAttrs([]slog.Attr{
		slog.String("proto", "grpc"),
	})
	logger := slog.New(logHandler)

	return GrpcLog{logger}
}
