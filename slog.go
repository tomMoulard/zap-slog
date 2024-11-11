// Package zapslog provides a zapcore.Core implementation that forwards logs to
// slog.Logger.
package zapslog

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// WrapCore wraps a slog.Logger as a zapcore.Core.
// The returned zap.Option can be passed to zap.New.
// Example:
//
//	logger := slog.Default()
//	zapLogger, _ := zap.NewProduction(zapslog.WrapCore(logger))
func WrapCore(logger *slog.Logger) zap.Option {
	return zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return &zapSlogCore{logger: logger}
	})
}

// zapSlogCore is a zapcore.Core implementation that forwards logs to
// slog.Logger.
type zapSlogCore struct {
	logger *slog.Logger
}

func (c *zapSlogCore) Enabled(level zapcore.Level) bool {
	return c.logger.Enabled(context.Background(), zapCoreLevelToSlogLevel(level))
}

func fieldToAttr(field zapcore.Field) slog.Attr {
	switch field.Type {
	case zapcore.StringType:
		return slog.String(field.Key, field.String)
	case zapcore.Int64Type:
		return slog.Int64(field.Key, field.Integer)
	case zapcore.Int32Type:
		return slog.Int(field.Key, int(field.Integer))
	case zapcore.Uint64Type:
		return slog.Uint64(field.Key, uint64(field.Integer))
	case zapcore.Float64Type:
		return slog.Float64(field.Key, math.Float64frombits(uint64(field.Integer)))
	case zapcore.BoolType:
		return slog.Bool(field.Key, field.Integer == 1)
	case zapcore.TimeType:
		if field.Interface != nil {
			loc, ok := field.Interface.(*time.Location)
			if ok {
				return slog.Time(field.Key, time.Unix(0, field.Integer).In(loc))
			}
		}

		return slog.Time(field.Key, time.Unix(0, field.Integer))
	case zapcore.DurationType:
		return slog.Duration(field.Key, time.Duration(field.Integer))
	default:
		return slog.Any(field.Key, field.Interface)
	}
}

func fieldToAttrs(fields []zapcore.Field) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fields))
	for _, field := range fields {
		attrs = append(attrs, fieldToAttr(field))
	}

	return attrs
}

func (c *zapSlogCore) With(fields []zapcore.Field) zapcore.Core {
	handler := c.logger.Handler().WithAttrs(fieldToAttrs(fields))

	return &zapSlogCore{logger: slog.New(handler)}
}

func (c *zapSlogCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}

	return ce
}

func (c *zapSlogCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	attrs := make([]slog.Attr, 0, len(fields)+3)
	attrs = append(attrs, slog.String("name", entry.LoggerName))
	attrs = append(attrs, fieldToAttrs(fields)...)
	attrs = append(attrs, slog.String("stack", entry.Stack))

	// https://pkg.go.dev/log/slog#hdr-Writing_a_handler
	r := slog.NewRecord(entry.Time, zapCoreLevelToSlogLevel(entry.Level), entry.Message, entry.Caller.PC)

	err := c.logger.Handler().WithAttrs(attrs).Handle(context.Background(), r)
	if err != nil {
		return fmt.Errorf("failed to write log: %w", err)
	}

	return nil
}

func (c *zapSlogCore) Sync() error {
	return nil
}

// zapCoreLevelToSlogLevel converts a zapcore.Level to a slog.Level.
// unsupported levels are converted to slog.LevelDebug.
func zapCoreLevelToSlogLevel(level zapcore.Level) slog.Level {
	switch level {
	case zapcore.DebugLevel:
		return slog.LevelDebug
	case zapcore.InfoLevel:
		return slog.LevelInfo
	case zapcore.WarnLevel:
		return slog.LevelWarn
	case zapcore.ErrorLevel:
		return slog.LevelError
	case zapcore.DPanicLevel:
		return slog.LevelError
	case zapcore.PanicLevel:
		return slog.LevelError
	case zapcore.FatalLevel:
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}
