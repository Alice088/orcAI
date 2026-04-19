package logging

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func MustLogger() zerolog.Logger {
	logRotator := &lumberjack.Logger{
		Filename:   "./logs/whitebox.log",
		MaxSize:    10,
		MaxBackups: 2,
		MaxAge:     28,
		Compress:   true,
	}

	return zerolog.New(logRotator).With().Timestamp().Logger()
}
