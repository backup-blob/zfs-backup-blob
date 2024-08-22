package driver

import (
	"fmt"
	"github.com/backup-blob/zfs-backup-blob/internal/domain"
	"github.com/rs/zerolog"
	"io"
	"strings"
	"time"
)

type log struct {
	logger zerolog.Logger
}

func NewLog(writer io.Writer, logLevel domain.LogLevel) domain.LogDriver {
	return &log{
		logger: initLog(writer, logLevel),
	}
}

func initLog(w io.Writer, logLevel domain.LogLevel) zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: w, TimeFormat: time.RFC3339, NoColor: true}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}

	return zerolog.New(output).Level(logLevel.ToZeroLevel()).With().Timestamp().Logger()
}

func (l *log) Debugf(format string, v ...interface{}) {
	l.logger.Debug().Msgf(format, v...)
}

func (l *log) Infof(format string, v ...interface{}) {
	l.logger.Info().Msgf(format, v...)
}
