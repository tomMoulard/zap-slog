package zapslog_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

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

func TestWrapCoreFields(t *testing.T) {
	t.Parallel()

	var b bytes.Buffer
	handler := slog.NewTextHandler(
		io.MultiWriter(&b, testLogWriter{t}),
		&slog.HandlerOptions{Level: slog.LevelInfo},
	)
	loggerSlog := slog.New(handler)

	loggerZap, err := zap.NewProduction(zapslog.WrapCore(loggerSlog))
	require.NoError(t, err)

	loggerZap.Info("fields",
		zap.String("str", "hello"),
		zap.Int64("i64", -64),
		zap.Int32("i32", -32),
		zap.Int16("i16", -16),
		zap.Int8("i8", -8),
		zap.Uint64("u64", 64),
		zap.Uint32("u32", 32),
		zap.Uint16("u16", 16),
		zap.Uint8("u8", 8),
		zap.Uintptr("uptr", 128),
		zap.Float64("f64", 6.28),
		zap.Float32("f32", 3.14),
		zap.Bool("flag", true),
		zap.Time("ts", time.Date(2024, 1, 15, 12, 30, 15, 0, time.UTC)),
		zap.Duration("dur", time.Second),
		zap.Stringp("strp", nil),
		zap.Int64p("i64p", nil),
		zap.Int32p("i32p", nil),
		zap.Int16p("i16p", nil),
		zap.Int8p("i8p", nil),
		zap.Uint64p("u64p", nil),
		zap.Uint32p("u32p", nil),
		zap.Uint16p("u16p", nil),
		zap.Uint8p("u8p", nil),
		zap.Uintptrp("uptrp", nil),
		zap.Float64p("f64p", nil),
		zap.Float32p("f32p", nil),
		zap.Boolp("flagp", nil),
		zap.Timep("tsp", nil),
		zap.Durationp("durp", nil),
	)

	err = loggerZap.Sync()
	require.NoError(t, err)

	bs := b.String()

	assert.Contains(t, bs, `str=hello`)
	assert.Contains(t, bs, `i64=-64`)
	assert.Contains(t, bs, `i32=-32`)
	assert.Contains(t, bs, `i16=-16`)
	assert.Contains(t, bs, `i8=-8`)
	assert.Contains(t, bs, `u64=64`)
	assert.Contains(t, bs, `u32=32`)
	assert.Contains(t, bs, `u16=16`)
	assert.Contains(t, bs, `u8=8`)
	assert.Contains(t, bs, `uptr=128`)
	assert.Contains(t, bs, `f64=6.28`)
	assert.Contains(t, bs, `f32=3.14`)
	assert.Contains(t, bs, `flag=true`)
	assert.Contains(t, bs, `ts=2024-01-15T12:30:15`)
	assert.Contains(t, bs, `dur=1s`)
	assert.Contains(t, bs, `strp=<nil>`)
	assert.Contains(t, bs, `i64p=<nil>`)
	assert.Contains(t, bs, `i32p=<nil>`)
	assert.Contains(t, bs, `i16p=<nil>`)
	assert.Contains(t, bs, `i8p=<nil>`)
	assert.Contains(t, bs, `u64p=<nil>`)
	assert.Contains(t, bs, `u32p=<nil>`)
	assert.Contains(t, bs, `u16p=<nil>`)
	assert.Contains(t, bs, `u8p=<nil>`)
	assert.Contains(t, bs, `uptrp=<nil>`)
	assert.Contains(t, bs, `f64p=<nil>`)
	assert.Contains(t, bs, `f32p=<nil>`)
	assert.Contains(t, bs, `flagp=<nil>`)
	assert.Contains(t, bs, `tsp=<nil>`)
	assert.Contains(t, bs, `durp=<nil>`)
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

func BenchmarkWrapCoreNoCaller(b *testing.B) {
	loggerSlog := slog.New(noopSlogHandler{})

	loggerZap, err := zap.NewProduction(zapslog.WrapCore(loggerSlog), zap.WithCaller(false))
	require.NoError(b, err)

	b.ResetTimer()

	for range b.N {
		loggerZap.Info("hello world")
	}
}

func BenchmarkWrapCoreFields(b *testing.B) {
	loggerSlog := slog.New(noopSlogHandler{})

	loggerZap, err := zap.NewProduction(zapslog.WrapCore(loggerSlog), zap.WithCaller(false))
	require.NoError(b, err)

	b.ResetTimer()

	for range b.N {
		loggerZap.Info("hello world",
			zap.String("key", "value"),
			zap.Int("count", 42),
			zap.Bool("flag", true),
		)
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

func BenchmarkSlog(b *testing.B) {
	loggerSlog := slog.New(noopSlogHandler{})

	b.ResetTimer()

	for range b.N {
		loggerSlog.Info("hello world")
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
