package zapslog_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	zapslog "github.com/tommoulard/zap-slog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ExampleWrapCore() {
	logger, _ := zap.NewProduction(zapslog.WrapCore(slog.Default()))
	logger = logger.Named("example")
	logger.Info("hello world")
}

func TestWrapCore(t *testing.T) {
	t.Parallel()

	var b bytes.Buffer
	loggerSlog := slog.New(slog.NewTextHandler(&b, &slog.HandlerOptions{Level: slog.LevelDebug}))

	loggerZap, err := zap.NewProduction(zapslog.WrapCore(loggerSlog))
	require.NoError(t, err)

	loggerZap.Info("hello world")

	err = loggerZap.Sync()
	require.NoError(t, err)

	bs := b.String()
	t.Log(bs)

	require.Contains(t, bs, "hello world")
}

func BenchmarkWrapCore(b *testing.B) {
	loggerSlog := slog.New(noopSlogHandler{})

	loggerZap, err := zap.NewProduction(zapslog.WrapCore(loggerSlog))
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		loggerZap.Info("hello world")
	}
}

func BenchmarkZap(b *testing.B) {
	loggerZap, err := zap.NewProduction(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewNopCore()
	}))
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		loggerZap.Info("hello world")
	}
}

type noopSlogHandler struct{}

func (noopSlogHandler) Enabled(context.Context, slog.Level) bool  { return true }
func (noopSlogHandler) Handle(context.Context, slog.Record) error { return nil }
func (h noopSlogHandler) WithAttrs([]slog.Attr) slog.Handler      { return h }
func (h noopSlogHandler) WithGroup(string) slog.Handler           { return h }
