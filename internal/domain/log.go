package domain

import "github.com/rs/zerolog"

type LogLevel int

const (
	UnknownLevel LogLevel = iota
	DebugLevel
	InfoLevel
)

func StringToLevel(s string) LogLevel {
	switch s {
	case "debug":
		return DebugLevel
	default:
		return UnknownLevel
	}
}

func (l LogLevel) ToZeroLevel() zerolog.Level {
	switch l {
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	default:
		return zerolog.Disabled
	}
}

type LogDriver interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
}

type LogRepository interface {
	LogDriver
}
