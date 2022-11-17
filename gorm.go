package gorm

import (
	"context"
	"github.com/olegshishkin/go-logger"
	"github.com/pkg/errors"
	gorm "gorm.io/gorm/logger"
	"log"
	"time"
)

var (
	dummyErr = errors.New("error")
)

type wrapper struct {
	log logger.Logger
}

// FromLogger transforms Logger logger to GORM logger.
func FromLogger(l logger.Logger) gorm.Interface {
	return &wrapper{l}
}

func (l *wrapper) LogMode(ll gorm.LogLevel) gorm.Interface {
	err := l.log.LogLevel(logLevel(ll))
	if err != nil {
		log.Printf("[WARN] Log level hasn't been set. Error: %v\n", err)
		return l
	}
	return l
}

func (l *wrapper) Info(_ context.Context, msg string, args ...interface{}) {
	l.log.Info(msg, args)
}

func (l *wrapper) Warn(_ context.Context, msg string, args ...interface{}) {
	l.log.Warn(msg, args)
}

func (l *wrapper) Error(_ context.Context, msg string, args ...interface{}) {
	l.log.Error(dummyErr, msg, args)
}

func (l *wrapper) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	dur := time.Since(begin).Microseconds()
	if err != nil {
		l.log.Error(err, "sql: (%s), duration: %dμs, rows: %d", sql, dur, rows)
		return
	}
	l.log.Trace("sql: (%s), duration: %dμs, rows: %d", sql, dur, rows)
}

// logLevel transforms GORM log level to Logger log level.
func logLevel(ll gorm.LogLevel) logger.Level {
	switch ll {
	case gorm.Silent:
		return logger.Fatal
	case gorm.Info:
		return logger.Info
	case gorm.Warn:
		return logger.Warn
	case gorm.Error:
		return logger.Error
	default:
		return logger.Fatal
	}
}
