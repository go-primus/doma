package http

import (
	"fmt"
	"log/slog"

	"github.com/go-resty/resty/v2"
)

var _ resty.Logger = (*logWrap)(nil)

type logWrap struct {
}

// Debugf implements resty.Logger.
func (l *logWrap) Debugf(format string, v ...interface{}) {
	slog.Debug(fmt.Sprintf(format, v...))
}

// Errorf implements resty.Logger.
func (l *logWrap) Errorf(format string, v ...interface{}) {
	slog.Error(fmt.Sprintf(format, v...))
}

// Warnf implements resty.Logger.
func (l *logWrap) Warnf(format string, v ...interface{}) {
	slog.Warn(fmt.Sprintf(format, v...))
}
