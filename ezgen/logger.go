package ezgen

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var DefaultLog = New(
	200*time.Millisecond,
	false,
	false,
	logger.Info,
)

// New initialize logger
func New(slowThreshold time.Duration, ignoreRecordNotFoundError bool, parameterizedQueries bool, logLevel logger.LogLevel) logger.Interface {
	return &DbLog{
		SlowThreshold:             slowThreshold,
		IgnoreRecordNotFoundError: ignoreRecordNotFoundError,
		ParameterizedQueries:      parameterizedQueries,
		LogLevel:                  logLevel,
	}
}

type DbLog struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	LogLevel                  logger.LogLevel
}

// LogMode log mode
func (l *DbLog) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l *DbLog) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		logc.Infow(ctx, "[GROM]",
			logc.Field("info", fmt.Sprintf(msg, data...)),
			logc.Field("file", utils.FileWithLineNum()),
		)
	}
}

// Warn print warn messages
func (l *DbLog) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		logc.Infow(ctx, "[GROM]",
			logc.Field("warn", fmt.Sprintf(msg, data...)),
			logc.Field("file", utils.FileWithLineNum()),
		)
	}
}

// Error print error messages
func (l *DbLog) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		logc.Errorw(ctx, "[GROM]",
			logc.Field("err", fmt.Sprintf(msg, data...)),
			logc.Field("file", utils.FileWithLineNum()),
		)
	}
}

// Trace print sql message
//
//nolint:cyclop
func (l *DbLog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			logc.Errorw(ctx, "[GROM]",
				logc.Field("err", err),
				logc.Field("file", utils.FileWithLineNum()),
				logc.Field("rows", "-"),
				logc.Field("duration", float64(elapsed.Nanoseconds())/1e6),
				logc.Field("sql", sql),
			)
		} else {
			logc.Errorw(ctx, "[GROM]",
				logc.Field("Err", err),
				logc.Field("File", utils.FileWithLineNum()),
				logc.Field("Rows", rows),
				logc.Field("duration", float64(elapsed.Nanoseconds())/1e6),
				logc.Field("sql", sql),
			)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			logc.Sloww(ctx, "[GROM]",
				logc.Field("file", utils.FileWithLineNum()),
				logc.Field("slowLog", slowLog),
				logc.Field("rows", "-"),
				logc.Field("duration", float64(elapsed.Nanoseconds())/1e6),
				logc.Field("sql", sql),
			)
		} else {
			logc.Sloww(ctx, "[GROM]",
				logc.Field("file", utils.FileWithLineNum()),
				logc.Field("slowLog", slowLog),
				logc.Field("rows", rows),
				logc.Field("duration", float64(elapsed.Nanoseconds())/1e6),
				logc.Field("sql", sql),
			)
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			logc.Infow(ctx, "[GROM]",
				logc.Field("file", utils.FileWithLineNum()),
				logc.Field("rows", "-"),
				logc.Field("duration", float64(elapsed.Nanoseconds())/1e6),
				logc.Field("sql", sql),
			)
		} else {
			logc.Infow(ctx, "[GROM]",
				logc.Field("file", utils.FileWithLineNum()),
				logc.Field("rows", rows),
				logc.Field("duration", float64(elapsed.Nanoseconds())/1e6),
				logc.Field("sql", sql),
			)
		}
	}
}
