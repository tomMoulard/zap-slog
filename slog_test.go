package zapslog_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	handler := slog.NewTextHandler(
		io.MultiWriter(&b, testLogWriter{t}),
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)
	loggerSlog := slog.New(handler)

	loggerZap, err := zap.NewProduction(zapslog.WrapCore(loggerSlog))
	require.NoError(t, err)

	loggerZap.Debug("debug level")
	loggerZap.Info("info level")
	loggerZap.Warn("warn level")
	loggerZap.Error("error level")

	err = loggerZap.Sync()
	require.NoError(t, err)

	bs := b.String()

	assert.Contains(t, bs, `level=DEBUG msg="debug level"`)
	assert.Contains(t, bs, `level=INFO msg="info level"`)
	assert.Contains(t, bs, `level=WARN msg="warn level"`)
	assert.Contains(t, bs, `level=ERROR msg="error level"`)
}

func BenchmarkWrapCore(b *testing.B) {
	loggerSlog := slog.New(noopSlogHandler{})

	loggerZap, err := zap.NewProduction(zapslog.WrapCore(loggerSlog))
	require.NoError(b, err)

	b.ResetTimer()

	for range b.N {
		loggerZap.Info("hello world")
	}
}

func BenchmarkZap(b *testing.B) {
	loggerZap, err := zap.NewProduction(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewNopCore()
	}))
	require.NoError(b, err)

	b.ResetTimer()

	for range b.N {
		loggerZap.Info("hello world")
	}
}

type noopSlogHandler struct{}

func (noopSlogHandler) Enabled(context.Context, slog.Level) bool  { return true }
func (noopSlogHandler) Handle(context.Context, slog.Record) error { return nil }
func (h noopSlogHandler) WithAttrs([]slog.Attr) slog.Handler      { return h }
func (h noopSlogHandler) WithGroup(string) slog.Handler           { return h }

type testLogWriter struct{ t *testing.T }

func (w testLogWriter) Write(p []byte) (int, error) {
	w.t.Log(strings.TrimSuffix(string(p), "\n"))

	return len(p), nil
}
