package database

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormLogger implements gorm logger interface using zap
type GormLogger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	SkipErrRecordNotFound     bool
	IgnoreRecordNotFoundError bool
}

// NewGormLogger creates a new GORM logger using zap
func NewGormLogger(zapLogger *zap.Logger) logger.Interface {
	return &GormLogger{
		ZapLogger:                 zapLogger,
		LogLevel:                  logger.Info,
		SlowThreshold:             200 * time.Millisecond,
		SkipErrRecordNotFound:     true,
		IgnoreRecordNotFoundError: true,
	}
}

// LogMode sets log level
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info logs info messages
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.ZapLogger.Sugar().Infof(msg, data...)
	}
}

// Warn logs warning messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.ZapLogger.Sugar().Warnf(msg, data...)
	}
}

// Error logs error messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.ZapLogger.Sugar().Errorf(msg, data...)
	}
}

// Trace logs SQL queries
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
	}

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		fields = append(fields, zap.Error(err))
		l.ZapLogger.Error("database query error", fields...)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0:
		fields = append(fields, zap.Duration("threshold", l.SlowThreshold))
		l.ZapLogger.Warn("slow SQL query", fields...)
	case l.LogLevel == logger.Info:
		l.ZapLogger.Info("database query", fields...)
	}
}
