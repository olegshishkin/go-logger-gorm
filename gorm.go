package gorm

import (
	"context"
	"github.com/olegshishkin/go-logger"
	"github.com/pkg/errors"
	gorm "gorm.io/gorm/logger"
	"log"
	"time"
)

var errDummy = errors.New("error")

type Wrapper struct {
	log logger.Logger
}

// FromLogger transforms Logger logger to GORM logger.
func FromLogger(l logger.Logger) *Wrapper {
	return &Wrapper{l}
}

func (w *Wrapper) LogMode(ll gorm.LogLevel) gorm.Interface {
	err := w.log.SetLevel(logLevel(ll))
	if err != nil {
		log.Printf("[WARN] Log level hasn't been set. Error: %v\n", err)

		return w
	}

	return w
}

func (w *Wrapper) Info(_ context.Context, msg string, args ...interface{}) {
	w.log.Info(msg, args...)
}

func (w *Wrapper) Warn(_ context.Context, msg string, args ...interface{}) {
	w.log.Warn(msg, args...)
}

func (w *Wrapper) Error(_ context.Context, msg string, args ...interface{}) {
	w.log.Error(errDummy, msg, args...)
}

func (w *Wrapper) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	latency := time.Since(begin).Milliseconds()

	if err != nil {
		w.log.Error(err, "sql: (%s), latency: %dms, rows: %d", sql, latency, rows)

		return
	}

	w.log.Trace("sql: (%s), latency: %dms, rows: %d", sql, latency, rows)
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
